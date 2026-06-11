package maphash

// New creates empty hash map.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity). hash and equal
// define key contract used by probing and key equality checks and must be
// non-nil. Returned map has Len() == 0 and Cap() == effective starting
// capacity.
//
// Example: m := New[string, int](32, hashString, eqString)
func New[K any, V any](capacity int, hash func(K) uint64, equal func(a, b K) bool) *MapHash[K, V] {
	return nil
}
