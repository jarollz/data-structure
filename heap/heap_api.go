package heap

import "iter"

// Compile-time check: *Heap[T] satisfies API[T].
var _ API[int] = (*Heap[int])(nil)

// Push implements the API interface.
// Push inserts v, grows backing storage before writing when full, restores heap
// property, and always returns true.
// v is value to insert.
// Example: ok := h.Push(5)
func (s *Heap[T]) Push(v T) bool {
	panic("not implemented")
}

// PopTop implements the API interface.
// PopTop removes and returns heap top element.
// It returns (zero, false) when heap is empty.
// Example: v, ok := h.PopTop()
func (s *Heap[T]) PopTop() (T, bool) {
	panic("not implemented")
}

// PeekTop implements the API interface.
// PeekTop returns heap top without removing it.
// PeekTop does not change Len(), Cap(), heap shape, or internal array order.
// It returns (zero, false) when heap is empty.
// Example: v, ok := h.PeekTop()
func (s *Heap[T]) PeekTop() (T, bool) {
	panic("not implemented")
}

// Len implements the API interface.
// Len returns number of stored elements.
// Example: n := h.Len()
func (s *Heap[T]) Len() int {
	panic("not implemented")
}

// Cap implements the API interface.
// Cap returns backing storage capacity.
// Capacity starts at effective initial capacity and reflects later growth or
// shrink decisions.
// Example: c := h.Cap()
func (s *Heap[T]) Cap() int {
	panic("not implemented")
}

// Clear implements the API interface.
// Clear removes all elements and resets heap length state.
// Clear is safe on an already-empty heap and leaves the heap ready for future
// Push, PopTop, and PeekTop calls with same comparator.
// Example: h.Clear()
func (s *Heap[T]) Clear() {
	panic("not implemented")
}

// Clone implements the API interface.
// Clone returns independent heap copy with same length, capacity, comparator, and internal array order.
// Elements are copied with normal Go assignment.
// Example: cloned := h.Clone()
func (s *Heap[T]) Clone() *Heap[T] {
	panic("not implemented")
}

// CloneWith implements the API interface.
// CloneWith returns independent heap copy using cloneValue for each live
// element.
// CloneWith preserves Len(), Cap(), comparator, and internal heap-array order.
// cloneValue receives each live value once in current internal array order,
// not sorted order, and never sees unused capacity; nil means normal Go
// assignment.
// Example: cloned := h.CloneWith(func(v int) int { return v * 10 })
func (s *Heap[T]) CloneWith(cloneValue func(T) T) *Heap[T] {
	panic("not implemented")
}

// Values implements the API interface.
// Values yields each live element once in internal array order, not sorted order.
// Sequence supports early stop and yields nothing when heap is empty.
// Mutation during iteration is not safe.
// Example: for v := range h.Values() { _ = v }
func (s *Heap[T]) Values() iter.Seq[T] {
	panic("not implemented")
}
