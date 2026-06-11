package maptrie

import "iter"

// Compile-time check: *MapTrie[V] satisfies API[V].
var _ API[int] = (*MapTrie[int])(nil)

// Put implements the API interface.
// Put inserts key-value pair or overwrites existing key.
// key is exact string key; value is payload to store. Overwrite keeps Len()
// unchanged. Empty string keys are stored at root terminal state, and shared
// prefixes are reused.
// Example: m.Put("apple", 7)
func (s *MapTrie[V]) Put(key string, value V) {
	panic("not implemented")
}

// Get implements the API interface.
// Get returns value associated with key.
// key is exact lookup key. Get does not mutate trie state. A prefix path that
// exists without terminal value is reported as missing.
// It returns (zero, false) when key is missing.
// Example: v, ok := m.Get("apple")
func (s *MapTrie[V]) Get(key string) (V, bool) {
	panic("not implemented")
}

// Delete implements the API interface.
// Delete removes key when present.
// key identifies exact entry to remove. Deleting a key that is prefix of longer
// keys clears only that terminal value. Deleting a leaf key prunes dead suffix
// nodes.
// It returns false when key does not exist.
// Example: ok := m.Delete("apple")
func (s *MapTrie[V]) Delete(key string) bool {
	panic("not implemented")
}

// Has implements the API interface.
// Has reports whether key exists.
// key is exact lookup key. A prefix path that exists without terminal value is
// reported as absent.
// Example: ok := m.Has("apple")
func (s *MapTrie[V]) Has(key string) bool {
	panic("not implemented")
}

// HasPrefix implements the API interface.
// HasPrefix reports whether any stored key starts with prefix.
// prefix uses byte-wise string prefix semantics. HasPrefix("") is true if and
// only if trie is non-empty.
// Example: ok := m.HasPrefix("app")
func (s *MapTrie[V]) HasPrefix(prefix string) bool {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live key-value pairs.
// Example: n := m.Len()
func (s *MapTrie[V]) Len() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all entries and resets trie to empty state.
// Clear is safe on an already-empty trie and leaves the trie ready for future
// exact and prefix queries.
// Example: m.Clear()
func (s *MapTrie[V]) Clear() {
	panic("not implemented")
}

// Clone implements the API interface.
// Clone returns independent trie copy with same keys, values, exact lookups,
// prefix behavior, and iteration order.
// Values are copied with normal Go assignment.
// Example: cloned := m.Clone()
func (s *MapTrie[V]) Clone() *MapTrie[V] {
	panic("not implemented")
}

// CloneWith implements the API interface.
// CloneWith returns independent trie copy using cloneValue for each live value.
// CloneWith preserves Len(), exact lookups, prefix behavior, and ascending
// byte-lex order. cloneValue receives each live value once in ascending
// byte-lex key order and never sees non-terminal prefix nodes; nil means
// normal Go assignment.
// Example: cloned := m.CloneWith(func(v int) int { return v * 10 })
func (s *MapTrie[V]) CloneWith(cloneValue func(V) V) *MapTrie[V] {
	panic("not implemented")
}

// All implements the API interface.
// All yields each key-value pair exactly once in ascending byte-lex key order.
// Shorter key sorts before longer descendant. Sequence supports early stop and
// yields nothing when trie is empty.
// Mutation during iteration is not safe.
// Example: for k, v := range m.All() { _, _ = k, v }
func (s *MapTrie[V]) All() iter.Seq2[string, V] {
	panic("not implemented")
}

// WithPrefix implements the API interface.
// WithPrefix yields each key-value pair whose key starts with prefix.
// WithPrefix("") yields same entries and order as All(). Exact prefix key, if
// present, is yielded before longer descendants. Results are in ascending
// byte-lex key order; sequence supports early stop and yields nothing when no
// key matches.
// Mutation during iteration is not safe.
// Example: for k, v := range m.WithPrefix("app") { _, _ = k, v }
func (s *MapTrie[V]) WithPrefix(prefix string) iter.Seq2[string, V] {
	panic("not implemented")
}
