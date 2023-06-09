// Package trie implements a thread safe trie in go
package trie

import (
	"fmt"
	"sort"
)

// lexicalKeys store the keys of a resulting depth search in lexical order.
// this is useful for obtaining a lexically sorted list from a given prefix
// search term. It implements the sort interface for a slice of runes.
type lexicalKeys []rune

func (lk lexicalKeys) Len() int           { return len(lk) }
func (lk lexicalKeys) Swap(i, j int)      { lk[i], lk[j] = lk[j], lk[i] }
func (lk lexicalKeys) Less(i, j int) bool { return lk[i] < lk[j] }

// SafeTrie is the structure that contains the channel used to
// communicate with the trie.
type T struct {
	op chan (func(*trie))
}

// Insert will place the array data at v into the trie at location k
func (t *T) Insert(k string, v []interface{}) (e error) {
	if k == "" {
		e = fmt.Errorf("Insert may not have empty key.")
		return
	}
	if len(v) == 0 {
		e = fmt.Errorf("Insert may not have empty data.")
		return
	}
	t.op <- func(st *trie) {
		curr := st.root
		for _, char := range k {
			if next, ok := curr.children[char]; !ok {
				nc := newNode()
				curr.children[char] = nc
				curr = nc
			} else {
				curr = next
			}
		}
		curr.data = append(curr.data, v...)
		return
	}
	return
}

// Get retreives the data (if any) at location k
func (t *T) Get(k string) (v []interface{}, e error) {
	if k == "" {
		e = fmt.Errorf("Get requires a string to search for")
		return
	}
	ich := make(chan []interface{})
	t.op <- func(st *trie) {
		curr := st.root
		for _, char := range k {
			next, ok := curr.children[char]
			if !ok {
				ich <- nil
				return
			}
			curr = next
		}
		ich <- curr.data
		return
	}
	v = <-ich
	return
}

// Search makes a DFS Search for the term which begins with startAt
// if it is blank, the entire trie will be returned as a lexically sorted
// slice of interfaces.
func (t *T) Search(startAt string) []interface{} {
	rch := make(chan []interface{})
	t.op <- func(st *trie) {
		curr := st.root
		for _, char := range startAt {
			next, ok := curr.children[char]
			if !ok {
				rch <- nil
				return
			}
			curr = next
		}
		rch <- curr.getDataBelow()
		return
	}
	return <-rch
}

// trie contains the locally available trie
type trie struct {
	root *trieNode
}

// trieNode is the structure that contains the node data
// for the trie
type trieNode struct {
	data     []interface{}
	children map[rune]*trieNode
}

// getDataBelow returns all of the data for all descendants of n
func (n *trieNode) getDataBelow() (d []interface{}) {
	if len(n.data) > 0 {
		d = append(d, n.data...)
	}
	tmpKeys := lexicalKeys{}
	for k := range n.children {
		tmpKeys = append(tmpKeys, k)
	}
	sort.Sort(tmpKeys)
	for _, next := range tmpKeys {
		d = append(d, n.children[next].getDataBelow()...)
	}
	return
}

// newNode is the method for creating a new node for the trie
func newNode() (n *trieNode) {
	n = &trieNode{}
	n.children = make(map[rune]*trieNode)
	return
}

// loop is the method that runs the goroutine for the data structure
func (t *T) loop() {
	core := &trie{}
	core.root = newNode()
	for op := range t.op {
		op(core)
	}
}

// Close stops the running go routine
func (t *T) Close() {
	close(t.op)
	return
}

// New creates a new trie
func New() (t *T) {
	t = &T{make(chan func(*trie))}
	go t.loop()
	return
}
