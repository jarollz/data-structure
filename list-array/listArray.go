package listarray

// New creates empty array-backed list.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity). Returned list has
// Len() == 0 and Cap() == effective starting capacity.
//
// Example: s := New[int](32)
func New[T any](capacity int) *ListArray[T] {
	return nil
}
