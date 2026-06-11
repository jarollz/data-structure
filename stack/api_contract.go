package stack

import "iter"

// Stack implements the API interface.
//
// Stack stores values with last-in-first-out semantics.
type Stack[T any] struct{}

// API defines stack behavior.
type API[T any] interface {
	// Push adds v on stack top and returns true.
	//
	// v is value to push.
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
	// It returns (zero, false) when stack is empty.
	//
	// Example: v, ok := s.PeekTop()
	PeekTop() (T, bool)
	// Len returns number of live elements.
	//
	// Example: n := s.Len()
	Len() int
	// Cap returns current backing capacity.
	//
	// Example: c := s.Cap()
	Cap() int
	// Clear removes all elements and resets length state.
	//
	// Example: s.Clear()
	Clear()
	// Values yields values from top to bottom.
	//
	// Sequence yields each live element once, supports early stop when yield returns false, and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range s.Values() { _ = v }
	Values() iter.Seq[T]
}
