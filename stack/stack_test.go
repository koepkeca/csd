package stack

import (
	"testing"
)

func TestStackCreation(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Errorf("Failed, invalid stack length.")
	}
	s.Close()
	return
}

func TestStackPushLength(t *testing.T) {
	s := New()
	s.Push(14)
	s.Push(42)
	s.Push("testing")
	s.Push([]byte("Viper"))
	len := s.Len()
	if len != 4 {
		t.Errorf("Failed, invalid stack length, got %d expected 4", len)
	}
	s.Close()
	return
}

func TestStackPushOrder(t *testing.T) {
	s := New()
	defer s.Close()
	s.Push(16)
	s.Push(32)
	s.Push(64)
	nv, ok := s.Pop().(int)
	if !ok {
		t.Errorf("Failed, Pop got wrong type")
		return
	}
	if nv != 64 {
		t.Errorf("Failed, got incorrect value order")
		return
	}
	return
}

func TestSizeAfterPop(t *testing.T) {
	s := New()
	s.Push(16)
	s.Push("test")
	s.Push("私は笑い男だ")
	_ = s.Pop()
	_ = s.Pop()
	_ = s.Pop()
	if s.Len() != 0 {
		t.Errorf("Failed, poped through entire stack, yet size is non-zero")
	}
	s.Close()
	return
}

func TestEmptyPop(t *testing.T) {
	s := New()
	v := s.Pop()
	if v != nil {
		t.Errorf("Empty Pop got non-nil value")
	}
	s.Close()
}

func TestEmptyPopWithValues(t *testing.T) {
	s := New()
	s.Push("Thingy")
	_ = s.Pop()
	v := s.Pop()
	if v != nil {
		t.Errorf("Empty stack with values got non-nil value")
	}
	s.Close()
}

func BenchmarkEqualRWWithInt(b *testing.B) {
	s := New()
	write := false
	for i := 0; i < b.N; i++ {
		if s.Pop() == nil || write {
			s.Push(i)
		} else {
			s.Pop()
			write = true
		}
	}
	s.Close()
}

func BenchmarkROnlyWithInt(b *testing.B) {
	s := New()
	nbr := b.N
	for i := 0; i < nbr; i++ {
		s.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Pop()
	}
	s.Close()
}

func BenchmarkWOnlyWithInt(b *testing.B) {
	s := New()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	s.Close()
}
