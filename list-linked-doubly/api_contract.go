package listlinkeddoubly

import "iter"

// ListLinkedDoubly implements the API interface.
//
// ListLinkedDoubly stores live nodes in a doubly linked sequence from head to
// tail.
type ListLinkedDoubly[T any] struct{}

// API defines doubly linked list behavior.
type API[T any] interface {
	// PushFront inserts v at list head in O(1) time.
	//
	// v is prepended value. When list is empty, pushed node becomes both head and
	// tail.
	//
	// Example: list.PushFront(4)
	PushFront(v T)
	// PushBack inserts v at list tail in O(1) time.
	//
	// v is appended value. When list is empty, pushed node becomes both head and
	// tail.
	//
	// Example: list.PushBack(8)
	PushBack(v T)
	// PopFront removes and returns head value.
	//
	// It returns (zero, false) when list is empty. Popping the only live node
	// resets both head and tail to empty.
	//
	// Example: v, ok := list.PopFront()
	PopFront() (T, bool)
	// PopBack removes and returns tail value.
	//
	// It returns (zero, false) when list is empty. Popping the only live node
	// resets both head and tail to empty.
	//
	// Example: v, ok := list.PopBack()
	PopBack() (T, bool)
	// Len returns number of live nodes.
	//
	// Example: n := list.Len()
	Len() int
	// Clear removes all nodes and resets structural fields.
	//
	// Clear is safe on an already-empty list and leaves it ready for future end
	// operations.
	//
	// Example: list.Clear()
	Clear()
	// Clone returns independent list copy with same length and head-to-tail order.
	//
	// Elements are copied with normal Go assignment.
	//
	// Example: cloned := list.Clone()
	Clone() *ListLinkedDoubly[T]
	// CloneWith returns independent list copy using cloneValue for each live node.
	//
	// CloneWith preserves Len() and head-to-tail order. cloneValue receives each
	// live value once from head to tail and never sees reclaimed or free-list
	// nodes. When cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := list.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *ListLinkedDoubly[T]
	// Values yields values from head to tail.
	//
	// Sequence yields each live element once, supports early stop when yield returns false, and yields nothing for an empty list.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range list.Values() { _ = v }
	Values() iter.Seq[T]
}
