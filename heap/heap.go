package heap

import (
	"reflect"
	"unsafe"
)

const heapDefaultMinCap = 16

type heapState[T any] struct {
	store    any
	data     unsafe.Pointer
	elemSize uintptr
	len      int
	cap      int
	minCap   int
	cmp      func(a, b T) int
}

// New creates empty heap with normalized starting capacity and comparator.
//
// New normalizes capacity <= 0 to 16, then uses effective starting capacity
// max(16, capacity). Returned heap has Len() == 0 and Cap() == effective
// starting capacity.
//
// Example: h := New[int](32, cmpInt)
func New[T any](capacity int, cmp func(a, b T) int) *Heap[T] {
	startCap := normalizeHeapCap(capacity)
	state := newHeapState[T](startCap, startCap, cmp)
	return (*Heap[T])(unsafe.Pointer(state))
}

func newHeapState[T any](capacity, minCap int, cmp func(a, b T) int) *heapState[T] {
	normalizedCap := normalizeHeapCap(capacity)
	normalizedMinCap := normalizeHeapCap(minCap)
	if normalizedCap < normalizedMinCap {
		normalizedCap = normalizedMinCap
	}
	store, data := allocHeapStorage[T](normalizedCap)
	var sample [1]T
	return &heapState[T]{
		store:    store,
		data:     data,
		elemSize: unsafe.Sizeof(sample[0]),
		len:      0,
		cap:      normalizedCap,
		minCap:   normalizedMinCap,
		cmp:      cmp,
	}
}

func normalizeHeapCap(capacity int) int {
	if capacity <= heapDefaultMinCap {
		return heapDefaultMinCap
	}
	return capacity
}

func allocHeapStorage[T any](capacity int) (any, unsafe.Pointer) {
	var sample [1]T
	elemType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, elemType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}
