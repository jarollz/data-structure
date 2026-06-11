package stack

import "iter"

// Stack implements the API interface.
//
// Stack stores live values with last-in-first-out semantics.
type Stack[T any] struct{}

// API defines stack behavior.
type API[T any] interface {
	// Push adds v on stack top and returns true.
	//
	// v is value to push. Push grows backing storage before writing when stack is
	// full.
	//
	// Example: ok := s.Push(1)
	Push(v T) bool
	// Pop removes and returns current top value.
	//
	// It returns (zero, false) when stack is empty.
	//
	// Example: v, ok := s.Pop()
	Pop() (T, bool)
	// PeekTop returns current top value without removing it.
	//
	// PeekTop does not change Len(), Cap(), or top-to-bottom order. It returns
	// (zero, false) when stack is empty.
	//
	// Example: v, ok := s.PeekTop()
	PeekTop() (T, bool)
	// Len returns number of live elements.
	//
	// Example: n := s.Len()
	Len() int
	// Cap returns current backing capacity.
	//
	// Capacity starts at effective initial capacity and reflects later growth or
	// shrink decisions.
	//
	// Example: c := s.Cap()
	Cap() int
	// Clear removes all elements and resets length state.
	//
	// Clear is safe on an already-empty stack and leaves it ready for future Push
	// calls and empty Pop or PeekTop calls.
	//
	// Example: s.Clear()
	Clear()
	// Clone returns independent stack copy with same length, capacity, and top-to-bottom order.
	//
	// Elements are copied with normal Go assignment.
	//
	// Example: cloned := s.Clone()
	Clone() *Stack[T]
	// CloneWith returns independent stack copy using cloneValue for each live element.
	//
	// CloneWith preserves Len(), Cap(), and top-to-bottom order. cloneValue
	// receives each live value once from top to bottom and never sees unused
	// capacity. When cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := s.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *Stack[T]
	// Values yields values from top to bottom.
	//
	// Sequence yields each live element once, supports early stop when yield returns false, and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range s.Values() { _ = v }
	Values() iter.Seq[T]
}
