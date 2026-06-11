package stack

// New creates empty LIFO stack.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity).
//
// Example: s := New[int](16)
func New[T any](capacity int) *Stack[T] {
	return nil
}
