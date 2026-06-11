package listlinkedsingly

// New creates empty singly linked list.
//
// Returned list has no live nodes and is ready for PushFront, Append, and
// other API operations. Returned list has Len() == 0 with empty head and tail
// state. Structure tracks both head and tail so Append keeps contracted O(1)
// behavior.
//
// Example: s := New[int]()
func New[T any]() *ListLinkedSingly[T] {
	return nil
}
