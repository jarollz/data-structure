package listlinkedsingly

import "iter"

// Compile-time check: *ListLinkedSingly[T] satisfies API[T].
var _ API[int] = (*ListLinkedSingly[int])(nil)

// PushFront implements the API interface.
// PushFront inserts v at head.
// v is prepended value.
// Example: list.PushFront(5)
func (s *ListLinkedSingly[T]) PushFront(v T) {
	panic("not implemented")
}

// PopFront implements the API interface.
// PopFront removes and returns current head value.
// It returns (zero, false) when list is empty.
// Example: v, ok := list.PopFront()
func (s *ListLinkedSingly[T]) PopFront() (T, bool) {
	panic("not implemented")
}

// Append implements the API interface.
// Append adds v at tail in O(1) time using tracked tail.
// v is appended value.
// Example: list.Append(9)
func (s *ListLinkedSingly[T]) Append(v T) {
	panic("not implemented")
}

// DeleteFirst implements the API interface.
// DeleteFirst removes first node whose value makes match return true.
// match receives each value in head-to-tail order.
// It returns false when no matching node exists.
// Example: ok := list.DeleteFirst(func(v int) bool { return v == 3 })
func (s *ListLinkedSingly[T]) DeleteFirst(match func(T) bool) bool {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := list.Len()
func (s *ListLinkedSingly[T]) Len() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all nodes and resets head, tail, free-list, and length state.
// Example: list.Clear()
func (s *ListLinkedSingly[T]) Clear() {
	panic("not implemented")
}

// Values implements the API interface.
// Values yields values head-to-tail exactly once per live node.
// Sequence supports early stop and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range list.Values() { _ = v }
func (s *ListLinkedSingly[T]) Values() iter.Seq[T] {
	panic("not implemented")
}
