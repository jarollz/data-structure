package listlinkedsingly

import (
	"reflect"
	"unsafe"
)

const (
	listLinkedSinglyNilIndex   = -1
	listLinkedSinglyDefaultCap = 16
)

type listLinkedSinglyState[T any] struct {
	valueStore any
	valueData  unsafe.Pointer
	valueSize  uintptr

	nextStore any
	nextData  unsafe.Pointer

	head int
	tail int
	free int
	len  int
	cap  int
}

// New creates empty singly linked list.
//
// Returned list is non-nil with Len() == 0 and empty head/tail state.
//
// Example: s := New[int]()
func New[T any]() *ListLinkedSingly[T] {
	state := newListLinkedSinglyState[T](listLinkedSinglyDefaultCap)
	return (*ListLinkedSingly[T])(unsafe.Pointer(state))
}

func newListLinkedSinglyState[T any](capacity int) *listLinkedSinglyState[T] {
	if capacity < 1 {
		capacity = listLinkedSinglyDefaultCap
	}
	valueStore, valueData, valueSize := allocListLinkedSinglyValueStorage[T](capacity)
	nextStore, nextData := allocListLinkedSinglyNextStorage(capacity)
	state := &listLinkedSinglyState[T]{
		valueStore: valueStore,
		valueData:  valueData,
		valueSize:  valueSize,
		nextStore:  nextStore,
		nextData:   nextData,
		head:       listLinkedSinglyNilIndex,
		tail:       listLinkedSinglyNilIndex,
		free:       0,
		len:        0,
		cap:        capacity,
	}
	state.resetFreeList(0)
	return state
}

func allocListLinkedSinglyValueStorage[T any](capacity int) (any, unsafe.Pointer, uintptr) {
	var sample [1]T
	elemType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, elemType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr()), unsafe.Sizeof(sample[0])
}

func allocListLinkedSinglyNextStorage(capacity int) (any, unsafe.Pointer) {
	arrayType := reflect.ArrayOf(capacity, reflect.TypeOf(0))
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}

func (s *listLinkedSinglyState[T]) valuePtrAt(i int) *T {
	return (*T)(unsafe.Add(s.valueData, uintptr(i)*s.valueSize))
}

func (s *listLinkedSinglyState[T]) nextPtrAt(i int) *int {
	return (*int)(unsafe.Add(s.nextData, uintptr(i)*unsafe.Sizeof(0)))
}

func (s *listLinkedSinglyState[T]) resetFreeList(start int) {
	for i := start; i < s.cap; i++ {
		next := i + 1
		if i == s.cap-1 {
			next = listLinkedSinglyNilIndex
		}
		*s.nextPtrAt(i) = next
	}
	if start < s.cap {
		s.free = start
	} else {
		s.free = listLinkedSinglyNilIndex
	}
}
