// Package stack implements a thread safe stack in go.
package stack

// S is the structure that contains the channel used to
// communicate with the stack.
type S struct {
	op chan (func(*stack))
}

// Len will get return the number of items in the stack.
func (s *S) Len() (i int64) {
	lChan := make(chan int64)
	s.op <- func(curr *stack) {
		lChan <- int64(len(*curr))
	}
	return <-lChan
}

// Pop will perform a pop on the stack, removing the first item
// and returning it's value.
func (s *S) Pop() (v interface{}) {
	vChan := make(chan interface{})
	s.op <- func(curr *stack) {
		old := *curr
		n := len(old)
		if n == 0 {
			vChan <- nil
			return
		}
		item := old[n-1]
		*curr = old[0 : n-1]
		vChan <- item
		return
	}
	return <-vChan
}

// Push will push the value v onto the stack.
func (s *S) Push(v interface{}) {
	s.op <- func(curr *stack) {
		//		*curr = append([]interface{}{v}, *curr...)
		*curr = append(*curr, v)
		return
	}
	return
}

// Close closes the primary channel thus stopping
// the running go-routine.
func (s *S) Close() {
	close(s.op)
}

// New creates a new Safe Stack, this also starts the go-routine
// so once this is called, you need to clean up after yourself
// by using the Destroy method.
func New() (s *S) {
	s = &S{make(chan func(*stack))}
	go s.loop()
	return
}

// We emulate a stack using an interface slice to reduce memory overhead
type stack []interface{}

// loop creates the guarded data structure and listens for
// methods on the op channel. loop terminates when the op
// channel is closed.
func (s *S) loop() {
	st := &stack{}
	for op := range s.op {
		op(st)
	}
}
