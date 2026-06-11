package listskip

import "iter"

// Compile-time check: *ListSkip[T] satisfies API[T].
var _ API[int] = (*ListSkip[int])(nil)

// Insert implements the API interface.
// Insert adds v into sorted structure when v is not already present.
// v is value to insert according to comparator order.
// It returns false for duplicates.
// Example: ok := list.Insert(10)
func (s *ListSkip[T]) Insert(v T) bool {
	panic("not implemented")
}

// Delete implements the API interface.
// Delete removes existing value v.
// v is value to remove.
// It returns false when v is missing.
// Example: ok := list.Delete(10)
func (s *ListSkip[T]) Delete(v T) bool {
	panic("not implemented")
}

// Has implements the API interface.
// Has reports whether v is stored.
// v is lookup value.
// Example: ok := list.Has(10)
func (s *ListSkip[T]) Has(v T) bool {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live elements.
// Example: n := list.Len()
func (s *ListSkip[T]) Len() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all values and resets skip-list state.
// Example: list.Clear()
func (s *ListSkip[T]) Clear() {
	panic("not implemented")
}

// Values implements the API interface.
// Values yields values in ascending sorted order from level 0.
// Sequence yields each live value once, supports early stop, and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range list.Values() { _ = v }
func (s *ListSkip[T]) Values() iter.Seq[T] {
	panic("not implemented")
}
