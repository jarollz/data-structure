package listskip

import "iter"

// ListSkip implements the API interface.
//
// ListSkip stores unique values in sorted order according to comparator.
type ListSkip[T any] struct{}

// API defines ordered skip-list behavior.
type API[T any] interface {
	// Insert adds v when v is not already present.
	//
	// v is value to insert according to comparator order.
	// It returns false for duplicate values.
	//
	// Example: ok := list.Insert(10)
	Insert(v T) bool
	// Delete removes v when present.
	//
	// v is value to remove.
	// It returns false when v is missing.
	//
	// Example: ok := list.Delete(10)
	Delete(v T) bool
	// Has reports whether v exists.
	//
	// v is value to test.
	//
	// Example: ok := list.Has(10)
	Has(v T) bool
	// Len returns number of live elements.
	//
	// Example: n := list.Len()
	Len() int
	// Clear removes all elements and resets skip-list state.
	//
	// Example: list.Clear()
	Clear()
	// Clone returns independent skip-list copy with same length and sorted order.
	//
	// Elements are copied with normal Go assignment. Clone also preserves comparator, maxLevel, currentLevel, and deterministic RNG state.
	//
	// Example: cloned := list.Clone()
	Clone() *ListSkip[T]
	// CloneWith returns independent skip-list copy using cloneValue for each live element.
	//
	// cloneValue receives each live value in sorted order. When cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := list.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *ListSkip[T]
	// Values yields values in sorted order via level-0 traversal.
	//
	// Sequence yields each live element once, supports early stop when yield returns false, and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range list.Values() { _ = v }
	Values() iter.Seq[T]
}
