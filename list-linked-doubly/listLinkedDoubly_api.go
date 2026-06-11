package listlinkeddoubly

import "iter"

// Compile-time check: *ListLinkedDoubly[T] satisfies API[T].
var _ API[int] = (*ListLinkedDoubly[int])(nil)

// PushFront implements the API interface.
// PushFront inserts v at head in O(1) time.
// v is prepended value.
// Example: list.PushFront(4)
func (s *ListLinkedDoubly[T]) PushFront(v T) {
	panic("not implemented")
}

// PushBack implements the API interface.
// PushBack inserts v at tail in O(1) time.
// v is appended value.
// Example: list.PushBack(8)
func (s *ListLinkedDoubly[T]) PushBack(v T) {
	panic("not implemented")
}

// PopFront implements the API interface.
// PopFront removes and returns head value.
// It returns (zero, false) when list is empty.
// Example: v, ok := list.PopFront()
func (s *ListLinkedDoubly[T]) PopFront() (T, bool) {
	panic("not implemented")
}

// PopBack implements the API interface.
// PopBack removes and returns tail value.
// It returns (zero, false) when list is empty.
// Example: v, ok := list.PopBack()
func (s *ListLinkedDoubly[T]) PopBack() (T, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := list.Len()
func (s *ListLinkedDoubly[T]) Len() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all nodes and resets list state.
// Example: list.Clear()
func (s *ListLinkedDoubly[T]) Clear() {
	panic("not implemented")
}

// Values implements the API interface.
// Values yields values from head to tail exactly once per live node.
// Sequence supports early stop and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range list.Values() { _ = v }
func (s *ListLinkedDoubly[T]) Values() iter.Seq[T] {
	panic("not implemented")
}
