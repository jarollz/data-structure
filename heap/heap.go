package heap

// New creates empty binary heap.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity). cmp defines heap
// ordering (min-heap or max-heap behavior) and is preserved by Clone and
// CloneWith. Returned heap has Len() == 0 and Cap() == effective starting
// capacity.
//
// Example: h := New[int](32, cmpInt)
func New[T any](capacity int, cmp func(a, b T) int) *Heap[T] {
	return nil
}
