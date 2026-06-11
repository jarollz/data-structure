package treeredblack

import "iter"

// Compile-time check: *TreeRedBlack[T] satisfies API[T].
var _ API[int] = (*TreeRedBlack[int])(nil)

// Insert implements the API interface.
// Insert adds v into tree when v is not present.
// v is value to insert.
// It returns false on duplicate.
// Example: ok := tree.Insert(10)
func (s *TreeRedBlack[T]) Insert(v T) bool {
	panic("not implemented")
}

// Delete implements the API interface.
// Delete removes existing value v.
// v is value to remove.
// It returns false when v is missing.
// Example: ok := tree.Delete(10)
func (s *TreeRedBlack[T]) Delete(v T) bool {
	panic("not implemented")
}

// Has implements the API interface.
// Has reports whether value v exists in tree.
// v is lookup value.
// Example: ok := tree.Has(10)
func (s *TreeRedBlack[T]) Has(v T) bool {
	panic("not implemented")
}

// Min implements the API interface.
// Min returns smallest stored value.
// It returns (zero, false) when tree is empty.
// Example: v, ok := tree.Min()
func (s *TreeRedBlack[T]) Min() (T, bool) {
	panic("not implemented")
}

// Max implements the API interface.
// Max returns largest stored value.
// It returns (zero, false) when tree is empty.
// Example: v, ok := tree.Max()
func (s *TreeRedBlack[T]) Max() (T, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := tree.Len()
func (s *TreeRedBlack[T]) Len() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all values and resets tree state.
// Example: tree.Clear()
func (s *TreeRedBlack[T]) Clear() {
	panic("not implemented")
}

// Clone implements the API interface.
// Clone returns independent tree copy with same length, comparator, and ascending in-order sequence.
// Values are copied with normal Go assignment.
// Example: cloned := tree.Clone()
func (s *TreeRedBlack[T]) Clone() *TreeRedBlack[T] {
	panic("not implemented")
}

// CloneWith implements the API interface.
// CloneWith returns independent tree copy using cloneValue for each live value.
// cloneValue receives each live value in ascending in-order traversal; nil means normal Go assignment.
// Example: cloned := tree.CloneWith(func(v int) int { return v * 10 })
func (s *TreeRedBlack[T]) CloneWith(cloneValue func(T) T) *TreeRedBlack[T] {
	panic("not implemented")
}

// InOrder implements the API interface.
// InOrder yields values in ascending sorted order.
// Sequence yields each live value once, supports early stop, and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range tree.InOrder() { _ = v }
func (s *TreeRedBlack[T]) InOrder() iter.Seq[T] {
	panic("not implemented")
}
