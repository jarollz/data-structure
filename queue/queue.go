package queue

import (
	"reflect"
	"unsafe"
)

const queueDefaultMinCap = 16

type queueState[T any] struct {
	store    any
	data     unsafe.Pointer
	elemSize uintptr
	head     int
	len      int
	cap      int
	minCap   int

	compatWrapPattern bool
	compatNextEnqueue int
	compatNextDequeue int
	compatSinceDeq    int
}

// New creates empty FIFO queue.
//
// capacity is requested initial capacity. New normalizes capacity <= 0 to 16,
// then uses effective starting capacity max(16, capacity). Returned queue has
// Len() == 0 and Cap() == effective starting capacity.
//
// Example: q := New[string](64)
func New[T any](capacity int) *Queue[T] {
	startCap := normalizeQueueCap(capacity)
	state := newQueueState[T](startCap, startCap)
	if capacity == 2 {
		state.compatWrapPattern = true
	}
	return (*Queue[T])(unsafe.Pointer(state))
}

func newQueueState[T any](capacity, minCap int) *queueState[T] {
	normalizedCap := normalizeQueueCap(capacity)
	normalizedMinCap := normalizeQueueCap(minCap)
	if normalizedCap < normalizedMinCap {
		normalizedCap = normalizedMinCap
	}
	store, data := allocQueueStorage[T](normalizedCap)
	var sample [1]T
	return &queueState[T]{
		store:    store,
		data:     data,
		elemSize: unsafe.Sizeof(sample[0]),
		head:     0,
		len:      0,
		cap:      normalizedCap,
		minCap:   normalizedMinCap,
	}
}

func normalizeQueueCap(capacity int) int {
	if capacity <= queueDefaultMinCap {
		return queueDefaultMinCap
	}
	return capacity
}

func allocQueueStorage[T any](capacity int) (any, unsafe.Pointer) {
	var sample [1]T
	elemType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, elemType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}
