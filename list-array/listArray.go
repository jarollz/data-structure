package listarray

import "unsafe"

const listArrayDefaultMinCap = 16

type listArrayState[T any] struct {
	store    any
	data     unsafe.Pointer
	elemSize uintptr
	len      int
	cap      int
	minCap   int
}

// New creates empty array-backed list.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity). Returned list has
// Len() == 0 and Cap() == effective starting capacity.
//
// Example: s := New[int](32)
func New[T any](capacity int) *ListArray[T] {
	startCap := normalizeListArrayCap(capacity)
	state := newListArrayState[T](startCap, startCap)
	return (*ListArray[T])(unsafe.Pointer(state))
}

func newListArrayState[T any](capacity, minCap int) *listArrayState[T] {
	normalizedCap := normalizeListArrayCap(capacity)
	normalizedMinCap := normalizeListArrayCap(minCap)
	if normalizedCap < normalizedMinCap {
		normalizedCap = normalizedMinCap
	}
	store, data := allocListArrayStorage[T](normalizedCap)
	var sample [1]T
	return &listArrayState[T]{
		store:    store,
		data:     data,
		elemSize: unsafe.Sizeof(sample[0]),
		len:      0,
		cap:      normalizedCap,
		minCap:   normalizedMinCap,
	}
}

func normalizeListArrayCap(capacity int) int {
	if capacity <= listArrayDefaultMinCap {
		return listArrayDefaultMinCap
	}
	return capacity
}
