package heap

import "iter"

// Heap implements the API interface.
//
// Heap stores values in binary-heap array form ordered by constructor comparator.
type Heap[T any] struct{}

// API defines binary-heap behavior.
type API[T any] interface {
	// Push inserts v, restores heap order, and returns true.
	//
	// v is value to insert.
	//
	// Example: ok := h.Push(5)
	Push(v T) bool
	// PopTop removes and returns current heap top.
	//
	// It returns (zero, false) when heap is empty.
	//
	// Example: v, ok := h.PopTop()
	PopTop() (T, bool)
	// PeekTop returns current heap top without removal.
	//
	// It returns (zero, false) when heap is empty.
	//
	// Example: v, ok := h.PeekTop()
	PeekTop() (T, bool)
	// Len returns number of stored elements.
	//
	// Example: n := h.Len()
	Len() int
	// Cap returns current backing capacity.
	//
	// Example: c := h.Cap()
	Cap() int
	// Clear removes all elements and resets heap state.
	//
	// Example: h.Clear()
	Clear()
	// Clone returns independent heap copy with same length, capacity, comparator, and internal array order.
	//
	// Elements are copied with normal Go assignment.
	//
	// Example: cloned := h.Clone()
	Clone() *Heap[T]
	// CloneWith returns independent heap copy using cloneValue for each live element.
	//
	// cloneValue receives each live value in internal array order. When cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := h.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *Heap[T]
	// Values yields each stored element once in internal array order.
	//
	// Yield order is not sorted order. Sequence supports early stop when yield returns false and yields nothing when heap is empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range h.Values() { _ = v }
	Values() iter.Seq[T]
}
