package maptreeavl

import "iter"

// MapTreeAvl implements the API interface.
//
// MapTreeAvl stores ordered key-value pairs using AVL balancing on keys.
type MapTreeAvl[K any, V any] struct{}

// API defines AVL-ordered map behavior.
type API[K any, V any] interface {
	// Put inserts or overwrites key with value.
	//
	// key is map key and value is stored payload.
	// Overwriting existing key does not change Len().
	//
	// Example: m.Put(3, "v")
	Put(key K, value V)
	// Get returns value for key.
	//
	// key is lookup key.
	// It returns (zero, false) when key is missing.
	//
	// Example: v, ok := m.Get(3)
	Get(key K) (V, bool)
	// Delete removes key when present.
	//
	// key is key to remove.
	// It returns false when key is missing.
	//
	// Example: ok := m.Delete(3)
	Delete(key K) bool
	// Has reports whether key exists.
	//
	// key is lookup key.
	//
	// Example: ok := m.Has(3)
	Has(key K) bool
	// Min returns smallest key and its value.
	//
	// It returns (zeroK, zeroV, false) when map is empty.
	//
	// Example: k, v, ok := m.Min()
	Min() (K, V, bool)
	// Max returns largest key and its value.
	//
	// It returns (zeroK, zeroV, false) when map is empty.
	//
	// Example: k, v, ok := m.Max()
	Max() (K, V, bool)
	// Len returns number of live key-value pairs.
	//
	// Example: n := m.Len()
	Len() int
	// Clear removes all entries and resets map state.
	//
	// Example: m.Clear()
	Clear()
	// All yields each key-value pair once in ascending key order.
	//
	// Sequence supports early stop when yield returns false and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for k, v := range m.All() { _, _ = k, v }
	All() iter.Seq2[K, V]
}
