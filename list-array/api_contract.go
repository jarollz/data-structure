package listarray

import "iter"

// ListArray implements the API interface.
//
// ListArray stores values in index order and exposes list operations defined by API.
type ListArray[T any] struct{}

// API defines array-backed list behavior.
type API[T any] interface {
	// Append adds v at tail and returns true.
	//
	// v is value to append.
	//
	// Example: ok := list.Append(7)
	Append(v T) bool
	// Get returns value at i.
	//
	// i is zero-based index in [0, Len()). It returns (zero, false) when i is out of range.
	//
	// Example: v, ok := list.Get(2)
	Get(i int) (T, bool)
	// Set overwrites value at i.
	//
	// i is zero-based index in [0, Len()); v is replacement value. It returns false when i is out of range.
	//
	// Example: ok := list.Set(1, 42)
	Set(i int, v T) bool
	// Insert inserts v before i.
	//
	// i accepts [0, Len()] where i == Len() appends; v is inserted value. It returns false when i is out of range.
	//
	// Example: ok := list.Insert(0, 9)
	Insert(i int, v T) bool
	// Delete removes value at i.
	//
	// i is zero-based index in [0, Len()). It returns removed value and true on success, or (zero, false) when i is out of range.
	//
	// Example: v, ok := list.Delete(3)
	Delete(i int) (T, bool)
	// Len returns current number of live elements.
	//
	// Example: n := list.Len()
	Len() int
	// Cap returns current backing capacity.
	//
	// Example: c := list.Cap()
	Cap() int
	// Clear removes all elements and resets length state.
	//
	// Example: list.Clear()
	Clear()
	// Values yields each element from index 0 to Len()-1.
	//
	// The sequence yields each live element once, supports early stop when yield returns false, and yields nothing for an empty list.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range list.Values() { _ = v }
	Values() iter.Seq[T]
}
