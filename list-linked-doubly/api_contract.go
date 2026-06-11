package listlinkeddoubly

import "iter"

// ListLinkedDoubly implements the API interface.
//
// ListLinkedDoubly stores values in a doubly linked sequence from head to tail.
type ListLinkedDoubly[T any] struct{}

// API defines doubly linked list behavior.
type API[T any] interface {
	// PushFront inserts v at list head in O(1) time.
	//
	// v is prepended value.
	//
	// Example: list.PushFront(4)
	PushFront(v T)
	// PushBack inserts v at list tail in O(1) time.
	//
	// v is appended value.
	//
	// Example: list.PushBack(8)
	PushBack(v T)
	// PopFront removes and returns head value.
	//
	// It returns (zero, false) when list is empty.
	//
	// Example: v, ok := list.PopFront()
	PopFront() (T, bool)
	// PopBack removes and returns tail value.
	//
	// It returns (zero, false) when list is empty.
	//
	// Example: v, ok := list.PopBack()
	PopBack() (T, bool)
	// Len returns number of live nodes.
	//
	// Example: n := list.Len()
	Len() int
	// Clear removes all nodes and resets structural fields.
	//
	// Example: list.Clear()
	Clear()
	// Values yields values from head to tail.
	//
	// Sequence yields each live element once, supports early stop when yield returns false, and yields nothing for an empty list.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range list.Values() { _ = v }
	Values() iter.Seq[T]
}
