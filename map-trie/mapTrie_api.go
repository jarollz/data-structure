package maptrie

import "iter"

// Compile-time check: *MapTrie[V] satisfies API[V].
var _ API[int] = (*MapTrie[int])(nil)

// Put implements the API interface.
//
// Put inserts key or overwrites existing key value.
//
// Example: m.Put("apple", 7)
func (s *MapTrie[V]) Put(key string, value V) {
	impl := implOf(s)
	node := impl.root
	for i := 0; i < len(key); i++ {
		b := key[i]
		child := node.children[b]
		if child == nil {
			child = &trieNode[V]{}
			node.children[b] = child
			node.childCount++
		}
		node = child
	}
	if !node.terminal {
		impl.length++
		node.terminal = true
		node.key = key
	}
	node.value = value
}

// Get implements the API interface.
//
// Get returns stored value for exact key.
//
// Example: v, ok := m.Get("apple")
func (s *MapTrie[V]) Get(key string) (V, bool) {
	node := implOf(s).root
	for i := 0; i < len(key); i++ {
		b := key[i]
		node = node.children[b]
		if node == nil {
			var zero V
			return zero, false
		}
	}
	if !node.terminal {
		var zero V
		return zero, false
	}
	return node.value, true
}

func deleteAt[V any](node *trieNode[V], key string, depth int, isRoot bool) (bool, bool) {
	if depth == len(key) {
		if !node.terminal {
			return false, false
		}
		node.terminal = false
		node.key = ""
		var zero V
		node.value = zero
		if !isRoot && node.childCount == 0 {
			return true, true
		}
		return true, false
	}

	b := key[depth]
	child := node.children[b]
	if child == nil {
		return false, false
	}

	deleted, pruneChild := deleteAt(child, key, depth+1, false)
	if !deleted {
		return false, false
	}
	if pruneChild {
		node.children[b] = nil
		node.childCount--
	}
	if !isRoot && !node.terminal && node.childCount == 0 {
		return true, true
	}
	return true, false
}

func walkAll[V any](node *trieNode[V], yield func(string, V) bool) bool {
	if node.terminal {
		if !yield(node.key, node.value) {
			return false
		}
	}
	for i := 0; i < 256; i++ {
		child := node.children[i]
		if child == nil {
			continue
		}
		if !walkAll(child, yield) {
			return false
		}
	}
	return true
}

// Delete implements the API interface.
//
// Delete removes live key and prunes dead suffix nodes.
//
// Example: ok := m.Delete("apple")
func (s *MapTrie[V]) Delete(key string) bool {
	impl := implOf(s)
	deleted, _ := deleteAt(impl.root, key, 0, true)
	if deleted {
		impl.length--
	}
	return deleted
}

// Has implements the API interface.
//
// Has reports whether exact key currently exists.
//
// Example: ok := m.Has("apple")
func (s *MapTrie[V]) Has(key string) bool {
	_, ok := s.Get(key)
	return ok
}

// HasPrefix implements the API interface.
//
// HasPrefix reports whether any stored key starts with prefix.
//
// Example: ok := m.HasPrefix("app")
func (s *MapTrie[V]) HasPrefix(prefix string) bool {
	impl := implOf(s)
	if prefix == "" {
		return impl.length != 0
	}
	node := impl.root
	for i := 0; i < len(prefix); i++ {
		b := prefix[i]
		node = node.children[b]
		if node == nil {
			return false
		}
	}
	return node.terminal || node.childCount != 0
}

// Len implements the API interface.
//
// Len returns number of live key-value entries.
//
// Example: n := m.Len()
func (s *MapTrie[V]) Len() int {
	return implOf(s).length
}

// Clear implements the API interface.
//
// Clear resets trie to empty state.
//
// Example: m.Clear()
func (s *MapTrie[V]) Clear() {
	impl := implOf(s)
	impl.root = &trieNode[V]{}
	impl.length = 0
}

// Clone implements the API interface.
//
// Clone returns independent copy using assignment for values.
//
// Example: cloned := m.Clone()
func (s *MapTrie[V]) Clone() *MapTrie[V] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
//
// CloneWith returns independent copy using optional clone hook.
//
// Example: cloned := m.CloneWith(func(v int) int { return v * 10 })
func (s *MapTrie[V]) CloneWith(cloneValue func(V) V) *MapTrie[V] {
	cloned := New[V]()
	if cloneValue == nil {
		for key, value := range s.All() {
			cloned.Put(key, value)
		}
		return cloned
	}
	for key, value := range s.All() {
		cloned.Put(key, cloneValue(value))
	}
	return cloned
}

// All implements the API interface.
//
// All iterates every live entry in ascending byte-lex key order.
//
// Example: for k, v := range m.All() { _, _ = k, v }
func (s *MapTrie[V]) All() iter.Seq2[string, V] {
	impl := implOf(s)
	return func(yield func(string, V) bool) {
		_ = walkAll(impl.root, yield)
	}
}

// WithPrefix implements the API interface.
//
// WithPrefix iterates entries matching prefix in ascending byte-lex order.
//
// Example: for k, v := range m.WithPrefix("app") { _, _ = k, v }
func (s *MapTrie[V]) WithPrefix(prefix string) iter.Seq2[string, V] {
	impl := implOf(s)
	return func(yield func(string, V) bool) {
		node := impl.root
		for i := 0; i < len(prefix); i++ {
			b := prefix[i]
			node = node.children[b]
			if node == nil {
				return
			}
		}
		_ = walkAll(node, yield)
	}
}
