package queue

// New creates empty FIFO queue.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity).
//
// Example: q := New[string](64)
func New[T any](capacity int) *Queue[T] {
	return nil
}
