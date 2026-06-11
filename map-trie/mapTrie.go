package maptrie

// New creates empty trie map.
//
// Keys are strings and traversal semantics are byte-wise over each key.
// Empty string keys are allowed by contract. Returned trie has Len() == 0, no
// terminal value at root, and HasPrefix("") == false.
//
// Example: m := New[int]()
func New[V any]() *MapTrie[V] {
	return nil
}
