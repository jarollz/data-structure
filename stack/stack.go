package stack

import (
	"reflect"
	"unsafe"
)

const stackDefaultMinCap = 16

type stackState[T any] struct {
	store    any
	data     unsafe.Pointer
	elemSize uintptr
	len      int
	cap      int
	minCap   int
}

// New creates empty LIFO stack with normalized starting capacity.
//
// New normalizes capacity <= 0 to 16, then uses effective starting capacity
// max(16, capacity). Returned stack has Len() == 0 and Cap() == effective
// starting capacity.
//
// Example: s := New[int](16)
func New[T any](capacity int) *Stack[T] {
	startCap := normalizeStackCap(capacity)
	state := newStackState[T](startCap, startCap)
	return (*Stack[T])(unsafe.Pointer(state))
}

func newStackState[T any](capacity, minCap int) *stackState[T] {
	normalizedCap := normalizeStackCap(capacity)
	normalizedMinCap := normalizeStackCap(minCap)
	if normalizedCap < normalizedMinCap {
		normalizedCap = normalizedMinCap
	}
	store, data := allocStackStorage[T](normalizedCap)
	var sample [1]T
	return &stackState[T]{
		store:    store,
		data:     data,
		elemSize: unsafe.Sizeof(sample[0]),
		len:      0,
		cap:      normalizedCap,
		minCap:   normalizedMinCap,
	}
}

func normalizeStackCap(capacity int) int {
	if capacity <= stackDefaultMinCap {
		return stackDefaultMinCap
	}
	return capacity
}

func allocStackStorage[T any](capacity int) (any, unsafe.Pointer) {
	var sample [1]T
	elemType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, elemType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}
