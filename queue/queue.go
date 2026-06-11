package queue

// New creates empty FIFO queue.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity). Returned queue has
// Len() == 0 and Cap() == effective starting capacity.
//
// Example: q := New[string](64)
func New[T any](capacity int) *Queue[T] {
	return nil
}
