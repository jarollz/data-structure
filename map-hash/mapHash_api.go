package maphash

import "iter"

// Compile-time check: *MapHash[K, V] satisfies API[K, V].
var _ API[int, string] = (*MapHash[int, string])(nil)

// Put implements the API interface.
// Put inserts key-value pair or overwrites existing key.
// key is lookup key; value is payload to store. Overwrite keeps Len()
// unchanged. Tombstones must not break reachability of live keys later in the
// same probe chain.
// Example: m.Put("id", 7)
func (s *MapHash[K, V]) Put(key K, value V) {
	panic("not implemented")
}

// Get implements the API interface.
// Get returns value associated with key.
// key is lookup key. Get does not mutate map state.
// It returns (zero, false) when key is missing.
// Example: v, ok := m.Get("id")
func (s *MapHash[K, V]) Get(key K) (V, bool) {
	panic("not implemented")
}

// Delete implements the API interface.
// Delete removes key when present.
// key identifies entry to remove. Delete preserves reachability of live keys
// later in the same probe chain.
// It returns false when key does not exist.
// Example: ok := m.Delete("id")
func (s *MapHash[K, V]) Delete(key K) bool {
	panic("not implemented")
}

// Has implements the API interface.
// Has reports whether key exists.
// key is lookup key. Has does not mutate map state, and tombstones do not stop
// probe-chain search for later live keys.
// Example: ok := m.Has("id")
func (s *MapHash[K, V]) Has(key K) bool {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live entries.
// Example: n := m.Len()
func (s *MapHash[K, V]) Len() int {
	panic("not implemented")
}

// Cap implements the API interface.
// Cap returns table capacity.
// Capacity starts at effective initial capacity, never drops below minimum
// capacity, and reflects later growth, shrink, or cleanup rehash decisions.
// Example: c := m.Cap()
func (s *MapHash[K, V]) Cap() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all live entries and tombstones.
// Clear resets capacity to minimum capacity and leaves the map ready for
// future Put, Get, Delete, Has, and LoadFactor calls.
// Example: m.Clear()
func (s *MapHash[K, V]) Clear() {
	panic("not implemented")
}

// LoadFactor implements the API interface.
// LoadFactor returns float64(Len()) / float64(Cap()).
// Tombstones are excluded from Len() in this ratio. LoadFactor does not mutate
// map state.
// Example: lf := m.LoadFactor()
func (s *MapHash[K, V]) LoadFactor() float64 {
	panic("not implemented")
}

// Clone implements the API interface.
// Clone returns independent map copy with same live entries, capacity, load factor, hash hook, and equal hook.
// Keys and values are copied with normal Go assignment.
// Example: cloned := m.Clone()
func (s *MapHash[K, V]) Clone() *MapHash[K, V] {
	panic("not implemented")
}

// CloneWith implements the API interface.
// CloneWith returns independent map copy using cloneKey and cloneValue for
// each live entry.
// CloneWith preserves live entries, Len(), Cap(), LoadFactor(), hash hook, and
// equal hook. cloneKey and cloneValue receive each live key-value pair once
// and never see empty or tombstone slots; nil means normal Go assignment for
// that payload type. Cloned keys must remain compatible with hash and equal.
// Example: cloned := m.CloneWith(func(k int) int { return k }, func(v string) string { return v + "!" })
func (s *MapHash[K, V]) CloneWith(cloneKey func(K) K, cloneValue func(V) V) *MapHash[K, V] {
	panic("not implemented")
}

// All implements the API interface.
// All yields each live key-value pair exactly once in unspecified order.
// Iteration order may change after mutation or resize. Sequence supports early
// stop and yields nothing when map is empty.
// Mutation during iteration is not safe.
// Example: for k, v := range m.All() { _, _ = k, v }
func (s *MapHash[K, V]) All() iter.Seq2[K, V] {
	panic("not implemented")
}
