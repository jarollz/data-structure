package listlinkedsingly

import "iter"

// ListLinkedSingly implements the API interface.
//
// ListLinkedSingly stores values in head-to-tail order with singly linked semantics.
type ListLinkedSingly[T any] struct{}

// API defines singly linked list behavior.
type API[T any] interface {
	// PushFront inserts v at list head.
	//
	// v is value to prepend.
	//
	// Example: list.PushFront(5)
	PushFront(v T)
	// PopFront removes and returns head value.
	//
	// It returns (zero, false) when list is empty.
	//
	// Example: v, ok := list.PopFront()
	PopFront() (T, bool)
	// Append adds v at list tail in O(1) time using tracked tail state.
	//
	// Implementation must not rescan from head to find tail.
	//
	// v is value to append.
	//
	// Example: list.Append(9)
	Append(v T)
	// DeleteFirst removes first node that satisfies match.
	//
	// match receives node value and should return true for the first value to remove.
	// It returns false when no matching node exists.
	//
	// Example: ok := list.DeleteFirst(func(v int) bool { return v == 3 })
	DeleteFirst(match func(T) bool) bool
	// Len returns number of live nodes.
	//
	// Example: n := list.Len()
	Len() int
	// Clear removes all nodes and resets structural state.
	//
	// Example: list.Clear()
	Clear()
	// Clone returns independent list copy with same length and head-to-tail order.
	//
	// Elements are copied with normal Go assignment.
	//
	// Example: cloned := list.Clone()
	Clone() *ListLinkedSingly[T]
	// CloneWith returns independent list copy using cloneValue for each live node.
	//
	// cloneValue receives each live value from head to tail. When cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := list.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *ListLinkedSingly[T]
	// Values yields values from head to tail.
	//
	// Sequence yields each element once, supports early stop when yield returns false, and yields nothing for an empty list.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range list.Values() { _ = v }
	Values() iter.Seq[T]
}
