package maptreeredblack

import "iter"

// MapTreeRedBlack implements the API interface.
//
// MapTreeRedBlack stores ordered key-value pairs using red-black balancing on
// keys.
type MapTreeRedBlack[K any, V any] struct{}

// API defines red-black ordered map behavior.
type API[K any, V any] interface {
	// Put inserts or overwrites key with value.
	//
	// key is map key and value is stored payload.
	// Overwriting existing key does not change Len(). Successful Put preserves
	// key ordering and red-black invariants.
	//
	// Example: m.Put(3, "v")
	Put(key K, value V)
	// Get returns value for key.
	//
	// key is lookup key. Get does not mutate map state.
	// It returns (zero, false) when key is missing.
	//
	// Example: v, ok := m.Get(3)
	Get(key K) (V, bool)
	// Delete removes key when present.
	//
	// key is key to remove. Successful Delete preserves key ordering and
	// red-black invariants, including root and required fix-up cases.
	// It returns false when key is missing.
	//
	// Example: ok := m.Delete(3)
	Delete(key K) bool
	// Has reports whether key exists.
	//
	// key is lookup key. Has does not mutate map state.
	//
	// Example: ok := m.Has(3)
	Has(key K) bool
	// Min returns smallest key and its value.
	//
	// Min returns smallest key by comparator order and its value without mutating
	// map state. It returns (zeroK, zeroV, false) when map is empty.
	//
	// Example: k, v, ok := m.Min()
	Min() (K, V, bool)
	// Max returns largest key and its value.
	//
	// Max returns largest key by comparator order and its value without mutating
	// map state. It returns (zeroK, zeroV, false) when map is empty.
	//
	// Example: k, v, ok := m.Max()
	Max() (K, V, bool)
	// Len returns number of live key-value pairs.
	//
	// Example: n := m.Len()
	Len() int
	// Clear removes all entries and resets map state.
	//
	// Clear is safe on an already-empty map, resets root and length state, and
	// leaves comparator unchanged for future operations.
	//
	// Example: m.Clear()
	Clear()
	// Clone returns independent map copy with same length, comparator, lookup
	// results, ascending key order, and red-black validity.
	//
	// Keys and values are copied with normal Go assignment.
	//
	// Example: cloned := m.Clone()
	Clone() *MapTreeRedBlack[K, V]
	// CloneWith returns independent map copy using cloneKey and cloneValue for each live entry.
	//
	// CloneWith preserves length, comparator, ascending key order, red-black
	// validity, and lookup results under transformed keys. cloneKey and
	// cloneValue receive each live key-value pair once in ascending key order.
	// Cloned keys must remain comparator-compatible. When either hook is nil,
	// CloneWith uses normal Go assignment for that payload type.
	//
	// Example: cloned := m.CloneWith(func(k int) int { return k }, func(v string) string { return v + "!" })
	CloneWith(cloneKey func(K) K, cloneValue func(V) V) *MapTreeRedBlack[K, V]
	// All yields each key-value pair once in ascending key order.
	//
	// Sequence supports early stop when yield returns false and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for k, v := range m.All() { _, _ = k, v }
	All() iter.Seq2[K, V]
}
