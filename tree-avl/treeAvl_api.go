package treeavl

import "iter"

// Compile-time check: *TreeAvl[T] satisfies API[T].
var _ API[int] = (*TreeAvl[int])(nil)

// Insert implements the API interface.
// Insert adds v into tree when v is not present.
// v is value to insert.
// It returns false on duplicate.
// Example: ok := tree.Insert(10)
func (s *TreeAvl[T]) Insert(v T) bool {
	panic("not implemented")
}

// Delete implements the API interface.
// Delete removes existing value v.
// v is value to remove.
// It returns false when v is missing.
// Example: ok := tree.Delete(10)
func (s *TreeAvl[T]) Delete(v T) bool {
	panic("not implemented")
}

// Has implements the API interface.
// Has reports whether value v exists in tree.
// v is value to test.
// Example: ok := tree.Has(10)
func (s *TreeAvl[T]) Has(v T) bool {
	panic("not implemented")
}

// Min implements the API interface.
// Min returns smallest stored value.
// It returns (zero, false) when tree is empty.
// Example: v, ok := tree.Min()
func (s *TreeAvl[T]) Min() (T, bool) {
	panic("not implemented")
}

// Max implements the API interface.
// Max returns largest stored value.
// It returns (zero, false) when tree is empty.
// Example: v, ok := tree.Max()
func (s *TreeAvl[T]) Max() (T, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := tree.Len()
func (s *TreeAvl[T]) Len() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all values and resets tree state.
// Example: tree.Clear()
func (s *TreeAvl[T]) Clear() {
	panic("not implemented")
}

// InOrder implements the API interface.
// InOrder yields values in ascending sorted order.
// Sequence yields each live value once, supports early stop, and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range tree.InOrder() { _ = v }
func (s *TreeAvl[T]) InOrder() iter.Seq[T] {
	panic("not implemented")
}
