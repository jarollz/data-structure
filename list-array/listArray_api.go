package listarray

import "iter"

// Compile-time check: *ListArray[T] satisfies API[T].
var _ API[int] = (*ListArray[int])(nil)

// Append implements the API interface.
// Append appends v at list tail and always returns true.
// v is value to append.
// Example: ok := list.Append(7)
func (s *ListArray[T]) Append(v T) bool {
	panic("not implemented")
}

// Get implements the API interface.
// Get returns value at index i when i is in [0, Len()), else (zero, false).
// i is zero-based index to read.
// Example: v, ok := list.Get(2)
func (s *ListArray[T]) Get(i int) (T, bool) {
	panic("not implemented")
}

// Set implements the API interface.
// Set writes v to existing index i and returns true on success.
// i is zero-based target index; v is replacement value.
// It returns false when i is outside [0, Len()).
// Example: ok := list.Set(1, 42)
func (s *ListArray[T]) Set(i int, v T) bool {
	panic("not implemented")
}

// Insert implements the API interface.
// Insert places v before index i and shifts later elements right.
// i accepts [0, Len()] where i == Len() appends; v is inserted value.
// It returns false when i is outside [0, Len()].
// Example: ok := list.Insert(0, 9)
func (s *ListArray[T]) Insert(i int, v T) bool {
	panic("not implemented")
}

// Delete implements the API interface.
// Delete removes and returns element at index i, shifting later elements left.
// i is zero-based index to remove.
// It returns (zero, false) when i is outside [0, Len()).
// Example: v, ok := list.Delete(3)
func (s *ListArray[T]) Delete(i int) (T, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live elements.
// Example: n := list.Len()
func (s *ListArray[T]) Len() int {
	panic("not implemented")
}

// Cap implements the API interface.
// Cap returns backing storage capacity.
// Example: c := list.Cap()
func (s *ListArray[T]) Cap() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all elements and resets list length.
// Example: list.Clear()
func (s *ListArray[T]) Clear() {
	panic("not implemented")
}

// Values implements the API interface.
// Values yields elements in index order from 0 to Len()-1.
// Sequence yields each live element once, supports early stop, and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range list.Values() { _ = v }
func (s *ListArray[T]) Values() iter.Seq[T] {
	panic("not implemented")
}
