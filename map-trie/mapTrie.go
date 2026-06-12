package maptrie

import "unsafe"

type trieNode[V any] struct {
	value      V
	key        string
	terminal   bool
	childCount int
	children   [256]*trieNode[V]
}

type trieImpl[V any] struct {
	api    MapTrie[V]
	root   *trieNode[V]
	length int
}

func implOf[V any](m *MapTrie[V]) *trieImpl[V] {
	return (*trieImpl[V])(unsafe.Pointer(m))
}

func newTrieImpl[V any]() *trieImpl[V] {
	impl := &trieImpl[V]{}
	impl.root = &trieNode[V]{}
	return impl
}

// New creates empty trie map.
//
// New returns non-nil map with zero live entries.
//
// Example: m := New[int]()
func New[V any]() *MapTrie[V] {
	return &newTrieImpl[V]().api
}
