package queue

import "iter"

// Compile-time check: *Queue[T] satisfies API[T].
var _ API[int] = (*Queue[int])(nil)

// Enqueue implements the API interface.
// Enqueue appends v at queue back and always returns true.
// v is value to enqueue.
// Example: ok := q.Enqueue(1)
func (s *Queue[T]) Enqueue(v T) bool {
	panic("not implemented")
}

// Dequeue implements the API interface.
// Dequeue removes and returns front value.
// It returns (zero, false) when queue is empty.
// Example: v, ok := q.Dequeue()
func (s *Queue[T]) Dequeue() (T, bool) {
	panic("not implemented")
}

// PeekFront implements the API interface.
// PeekFront returns front value without removal.
// It returns (zero, false) when queue is empty.
// Example: v, ok := q.PeekFront()
func (s *Queue[T]) PeekFront() (T, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of queued elements.
// Example: n := q.Len()
func (s *Queue[T]) Len() int {
	panic("not implemented")
}

// Cap implements the API interface.
// Cap returns backing storage capacity.
// Example: c := q.Cap()
func (s *Queue[T]) Cap() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all queued elements and resets queue state.
// Example: q.Clear()
func (s *Queue[T]) Clear() {
	panic("not implemented")
}

// Clone implements the API interface.
// Clone returns independent queue copy with same length, capacity, and front-to-back order.
// Elements are copied with normal Go assignment.
// Example: cloned := q.Clone()
func (s *Queue[T]) Clone() *Queue[T] {
	panic("not implemented")
}

// CloneWith implements the API interface.
// CloneWith returns independent queue copy using cloneValue for each live element.
// cloneValue receives each live value from front to back; nil means normal Go assignment.
// Example: cloned := q.CloneWith(func(v int) int { return v * 10 })
func (s *Queue[T]) CloneWith(cloneValue func(T) T) *Queue[T] {
	panic("not implemented")
}

// Values implements the API interface.
// Values yields elements from front to back in logical queue order.
// Sequence yields each live element once, supports early stop, and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for v := range q.Values() { _ = v }
func (s *Queue[T]) Values() iter.Seq[T] {
	panic("not implemented")
}
