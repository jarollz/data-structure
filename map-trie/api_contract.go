package maptrie

import "iter"

// MapTrie implements the API interface.
//
// MapTrie stores string-keyed values in trie form with byte-wise exact-key and
// prefix-query semantics.
type MapTrie[V any] struct{}

// API defines trie-map behavior.
type API[V any] interface {
	// Put inserts or overwrites value for key.
	//
	// key is string lookup key; value is stored payload.
	// Overwriting existing key does not change Len(). Empty string keys are stored
	// at root terminal state, and shared prefixes are reused.
	//
	// Example: m.Put("apple", 7)
	Put(key string, value V)
	// Get returns stored value for key.
	//
	// key is exact lookup key. Get does not mutate trie state. A prefix path that
	// exists without terminal value is reported as missing.
	// It returns (zero, false) when key is missing.
	//
	// Example: v, ok := m.Get("apple")
	Get(key string) (V, bool)
	// Delete removes key when present.
	//
	// key is exact key to remove. Deleting a key that is prefix of longer keys
	// clears only that terminal value. Deleting a leaf key prunes dead suffix
	// nodes.
	// It returns false when key is missing.
	//
	// Example: ok := m.Delete("apple")
	Delete(key string) bool
	// Has reports whether key exists.
	//
	// key is exact lookup key. A prefix path that exists without terminal value is
	// reported as absent.
	//
	// Example: ok := m.Has("apple")
	Has(key string) bool
	// HasPrefix reports whether any stored key starts with prefix.
	//
	// prefix is byte-wise string prefix. HasPrefix("") is true if and only if the
	// trie is non-empty.
	//
	// Example: ok := m.HasPrefix("app")
	HasPrefix(prefix string) bool
	// Len returns number of live key-value pairs.
	//
	// Example: n := m.Len()
	Len() int
	// Clear removes all entries and resets trie to empty state.
	//
	// Clear is safe on an already-empty trie and leaves it ready for future exact
	// and prefix queries.
	//
	// Example: m.Clear()
	Clear()
	// Clone returns independent trie copy with same keys, values, exact lookups,
	// prefix behavior, and iteration order.
	//
	// Values are copied with normal Go assignment.
	//
	// Example: cloned := m.Clone()
	Clone() *MapTrie[V]
	// CloneWith returns independent trie copy using cloneValue for each live value.
	//
	// CloneWith preserves Len(), exact lookups, prefix behavior, and ascending
	// byte-lex order. cloneValue receives each live value once in ascending
	// byte-lex key order and never sees non-terminal prefix nodes. When
	// cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := m.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(V) V) *MapTrie[V]
	// All yields each key-value pair once in ascending byte-lex key order.
	//
	// Shorter key sorts before longer descendant. Sequence supports early stop
	// when yield returns false and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for k, v := range m.All() { _, _ = k, v }
	All() iter.Seq2[string, V]
	// WithPrefix yields each key-value pair whose key starts with prefix.
	//
	// WithPrefix("") yields same entries and order as All(). Exact prefix key, if
	// present, is yielded before longer descendants. Results are yielded in
	// ascending byte-lex key order. Sequence supports early stop when yield
	// returns false and yields nothing when no key matches.
	// Mutation during iteration is not safe.
	//
	// Example: for k, v := range m.WithPrefix("app") { _, _ = k, v }
	WithPrefix(prefix string) iter.Seq2[string, V]
}
