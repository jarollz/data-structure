package listlinkedsingly

import "iter"

// ListLinkedSingly implements the API interface.
//
// ListLinkedSingly stores live nodes in head-to-tail order with singly linked
// semantics and tracked tail state for O(1) Append.
type ListLinkedSingly[T any] struct{}

// API defines singly linked list behavior.
type API[T any] interface {
	// PushFront inserts v at list head.
	//
	// v is value to prepend. When list is empty, pushed node becomes both head
	// and tail.
	//
	// Example: list.PushFront(5)
	PushFront(v T)
	// PopFront removes and returns head value.
	//
	// It returns (zero, false) when list is empty. Popping the only live node
	// resets both head and tail to empty.
	//
	// Example: v, ok := list.PopFront()
	PopFront() (T, bool)
	// Append adds v at list tail in O(1) time using tracked tail state.
	//
	// Implementation must not rescan from head to find tail. Appending to an
	// empty list updates both head and tail.
	//
	// v is value to append.
	//
	// Example: list.Append(9)
	Append(v T)
	// DeleteFirst removes first node that satisfies match.
	//
	// match receives node value and should return true for the first value to remove.
	// DeleteFirst scans head to tail, removes only the first matching node,
	// preserves order of remaining nodes, updates head or tail when needed, and
	// returns false when no matching node exists.
	//
	// Example: ok := list.DeleteFirst(func(v int) bool { return v == 3 })
	DeleteFirst(match func(T) bool) bool
	// Len returns number of live nodes.
	//
	// Example: n := list.Len()
	Len() int
	// Clear removes all nodes and resets structural state.
	//
	// Clear is safe on an already-empty list and leaves it ready for future
	// PushFront, Append, PopFront, and DeleteFirst calls.
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
	// CloneWith preserves Len() and head-to-tail order. cloneValue receives each
	// live value once from head to tail and never sees reclaimed or free-list
	// nodes. When cloneValue is nil, CloneWith uses normal Go assignment.
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
