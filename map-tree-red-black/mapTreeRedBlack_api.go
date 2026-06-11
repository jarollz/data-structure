package maptreeredblack

import "iter"

// Compile-time check: *MapTreeRedBlack[K, V] satisfies API[K, V].
var _ API[int, string] = (*MapTreeRedBlack[int, string])(nil)

// Put implements the API interface.
// Put inserts key-value pair or overwrites existing key.
// key is lookup key; value is payload to store.
// Overwrite keeps Len() unchanged.
// Example: m.Put(3, "v")
func (s *MapTreeRedBlack[K, V]) Put(key K, value V) {
	panic("not implemented")
}

// Get implements the API interface.
// Get returns value for key.
// key is lookup key.
// It returns (zero, false) when key is missing.
// Example: v, ok := m.Get(3)
func (s *MapTreeRedBlack[K, V]) Get(key K) (V, bool) {
	panic("not implemented")
}

// Delete implements the API interface.
// Delete removes key when present.
// key is key to remove.
// It returns false when key is missing.
// Example: ok := m.Delete(3)
func (s *MapTreeRedBlack[K, V]) Delete(key K) bool {
	panic("not implemented")
}

// Has implements the API interface.
// Has reports whether key exists.
// key is lookup key.
// Example: ok := m.Has(3)
func (s *MapTreeRedBlack[K, V]) Has(key K) bool {
	panic("not implemented")
}

// Min implements the API interface.
// Min returns smallest key and associated value.
// It returns (zeroK, zeroV, false) when map is empty.
// Example: k, v, ok := m.Min()
func (s *MapTreeRedBlack[K, V]) Min() (K, V, bool) {
	panic("not implemented")
}

// Max implements the API interface.
// Max returns largest key and associated value.
// It returns (zeroK, zeroV, false) when map is empty.
// Example: k, v, ok := m.Max()
func (s *MapTreeRedBlack[K, V]) Max() (K, V, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live entries.
// Example: n := m.Len()
func (s *MapTreeRedBlack[K, V]) Len() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all entries and resets map state.
// Example: m.Clear()
func (s *MapTreeRedBlack[K, V]) Clear() {
	panic("not implemented")
}

// All implements the API interface.
// All yields each key-value pair exactly once in ascending key order.
// Sequence supports early stop and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for k, v := range m.All() { _, _ = k, v }
func (s *MapTreeRedBlack[K, V]) All() iter.Seq2[K, V] {
	panic("not implemented")
}
