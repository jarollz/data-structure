package heap

import "iter"

// Heap implements the API interface.
//
// Heap stores live values in binary-heap array form ordered by constructor
// comparator.
type Heap[T any] struct{}

// API defines binary-heap behavior.
type API[T any] interface {
	// Push inserts v, restores heap order, and returns true.
	//
	// v is value to insert. Push grows backing storage before writing when heap
	// is full.
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
	// PeekTop does not change Len(), Cap(), heap shape, or internal array order.
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
	// Capacity starts at effective initial capacity and reflects later growth or
	// shrink decisions.
	//
	// Example: c := h.Cap()
	Cap() int
	// Clear removes all elements and resets heap state.
	//
	// Clear is safe on an already-empty heap and leaves it ready for future Push,
	// PopTop, and PeekTop calls with same comparator.
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
	// CloneWith preserves Len(), Cap(), comparator, and internal heap-array
	// order. cloneValue receives each live value once in current internal array
	// order, not sorted order, and never sees unused capacity. When cloneValue is
	// nil, CloneWith uses normal Go assignment.
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
