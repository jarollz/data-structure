package maphash

import "iter"

// MapHash implements the API interface.
//
// MapHash stores live key-value pairs using open addressing with linear
// probing and tombstone semantics.
type MapHash[K any, V any] struct{}

// API defines hash map behavior.
type API[K any, V any] interface {
	// Put inserts or overwrites value for key.
	//
	// key is lookup key; value is stored payload.
	// Overwriting existing key does not change Len(). Tombstones must not break
	// reachability of live keys later in the same probe chain.
	//
	// Example: m.Put("id", 7)
	Put(key K, value V)
	// Get returns stored value for key.
	//
	// key is lookup key. Get does not mutate map state.
	// It returns (zero, false) when key is missing.
	//
	// Example: v, ok := m.Get("id")
	Get(key K) (V, bool)
	// Delete removes key when present.
	//
	// key is entry to remove. Delete preserves reachability of live keys later in
	// the same probe chain.
	// It returns false when key is missing.
	//
	// Example: ok := m.Delete("id")
	Delete(key K) bool
	// Has reports whether key exists.
	//
	// key is lookup key. Has does not mutate map state, and tombstones do not
	// stop probe-chain search for later live keys.
	//
	// Example: ok := m.Has("id")
	Has(key K) bool
	// Len returns number of live entries.
	//
	// Example: n := m.Len()
	Len() int
	// Cap returns current table capacity.
	//
	// Capacity starts at effective initial capacity, never drops below minimum
	// capacity, and reflects later growth, shrink, or cleanup rehash decisions.
	//
	// Example: c := m.Cap()
	Cap() int
	// Clear removes all live entries and tombstones.
	//
	// Clear resets capacity to minimum capacity and leaves map ready for future
	// Put, Get, Delete, Has, and LoadFactor calls.
	//
	// Example: m.Clear()
	Clear()
	// LoadFactor returns float64(Len()) / float64(Cap()).
	//
	// Tombstones are excluded from Len() in this ratio. LoadFactor does not
	// mutate map state.
	//
	// Example: lf := m.LoadFactor()
	LoadFactor() float64
	// Clone returns independent map copy with same live entries, capacity, load factor, hash hook, and equal hook.
	//
	// Keys and values are copied with normal Go assignment.
	//
	// Example: cloned := m.Clone()
	Clone() *MapHash[K, V]
	// CloneWith returns independent map copy using cloneKey and cloneValue for each live entry.
	//
	// CloneWith preserves live entries, Len(), Cap(), LoadFactor(), hash hook,
	// and equal hook. cloneKey and cloneValue receive each live key-value pair
	// once and never see empty or tombstone slots. When either hook is nil,
	// CloneWith uses normal Go assignment for that payload type. Cloned keys must
	// remain compatible with hash and equal.
	//
	// Example: cloned := m.CloneWith(func(k int) int { return k }, func(v string) string { return v + "!" })
	CloneWith(cloneKey func(K) K, cloneValue func(V) V) *MapHash[K, V]
	// All yields each live key-value pair exactly once.
	//
	// Iteration order is unspecified and may change after mutation or resize;
	// sequence supports early stop and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for k, v := range m.All() { _, _ = k, v }
	All() iter.Seq2[K, V]
}
