package stack

import "iter"

// Compile-time check: *Stack[T] satisfies API[T].
var _ API[int] = (*Stack[int])(nil)

// Push implements the API interface.
// Push puts v on top and always returns true.
// v is value to push.
// Example: ok := s.Push(1)
func (s *Stack[T]) Push(v T) bool {
	panic("not implemented")
}

// Pop implements the API interface.
// Pop removes and returns top value.
// It returns (zero, false) on empty stack.
// Example: v, ok := s.Pop()
func (s *Stack[T]) Pop() (T, bool) {
	panic("not implemented")
}

// PeekTop implements the API interface.
// PeekTop returns top value without removal.
// It returns (zero, false) on empty stack.
// Example: v, ok := s.PeekTop()
func (s *Stack[T]) PeekTop() (T, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of live elements.
// Example: n := s.Len()
func (s *Stack[T]) Len() int {
	panic("not implemented")
}

// Cap implements the API interface.
// Cap returns backing storage capacity.
// Example: c := s.Cap()
func (s *Stack[T]) Cap() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all elements and resets stack state.
// Example: s.Clear()
func (s *Stack[T]) Clear() {
	panic("not implemented")
}

// Values implements the API interface.
// Values yields elements from top to bottom in current stack order.
// Sequence yields each live element once, supports early stop, and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range s.Values() { _ = v }
func (s *Stack[T]) Values() iter.Seq[T] {
	panic("not implemented")
}
