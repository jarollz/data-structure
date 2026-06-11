package queue

import "iter"

// Queue implements the API interface.
//
// Queue stores values with first-in-first-out semantics.
type Queue[T any] struct{}

// API defines queue behavior.
type API[T any] interface {
	// Enqueue adds v at queue back and returns true.
	//
	// v is value to append.
	//
	// Example: ok := q.Enqueue(1)
	Enqueue(v T) bool
	// Dequeue removes and returns queue front value.
	//
	// It returns (zero, false) when queue is empty.
	//
	// Example: v, ok := q.Dequeue()
	Dequeue() (T, bool)
	// PeekFront returns queue front value without removal.
	//
	// It returns (zero, false) when queue is empty.
	//
	// Example: v, ok := q.PeekFront()
	PeekFront() (T, bool)
	// Len returns number of queued elements.
	//
	// Example: n := q.Len()
	Len() int
	// Cap returns current backing capacity.
	//
	// Example: c := q.Cap()
	Cap() int
	// Clear removes all elements and resets queue state.
	//
	// Example: q.Clear()
	Clear()
	// Clone returns independent queue copy with same length, capacity, and front-to-back order.
	//
	// Elements are copied with normal Go assignment.
	//
	// Example: cloned := q.Clone()
	Clone() *Queue[T]
	// CloneWith returns independent queue copy using cloneValue for each live element.
	//
	// cloneValue receives each live value from front to back. When cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := q.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *Queue[T]
	// Values yields values from front to back.
	//
	// Sequence yields each live element once, supports early stop when yield returns false, and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range q.Values() { _ = v }
	Values() iter.Seq[T]
}
