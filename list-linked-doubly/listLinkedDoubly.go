package listlinkeddoubly

import (
	"reflect"
	"unsafe"
)

const (
	listLinkedDoublyNilIndex   = -1
	listLinkedDoublyDefaultCap = 16
)

type listLinkedDoublyState[T any] struct {
	valueStore any
	valueData  unsafe.Pointer
	valueSize  uintptr

	prevStore any
	prevData  unsafe.Pointer

	nextStore any
	nextData  unsafe.Pointer

	head int
	tail int
	free int
	len  int
	cap  int
}

// New creates empty doubly linked list.
//
// Returned list is non-nil with Len() == 0 and empty head/tail state.
//
// Example: list := New[int]()
func New[T any]() *ListLinkedDoubly[T] {
	state := newListLinkedDoublyState[T](listLinkedDoublyDefaultCap)
	return (*ListLinkedDoubly[T])(unsafe.Pointer(state))
}

func newListLinkedDoublyState[T any](capacity int) *listLinkedDoublyState[T] {
	if capacity < 1 {
		capacity = listLinkedDoublyDefaultCap
	}
	valueStore, valueData, valueSize := allocListLinkedDoublyValueStorage[T](capacity)
	prevStore, prevData := allocListLinkedDoublyIntStorage(capacity)
	nextStore, nextData := allocListLinkedDoublyIntStorage(capacity)
	state := &listLinkedDoublyState[T]{
		valueStore: valueStore,
		valueData:  valueData,
		valueSize:  valueSize,
		prevStore:  prevStore,
		prevData:   prevData,
		nextStore:  nextStore,
		nextData:   nextData,
		head:       listLinkedDoublyNilIndex,
		tail:       listLinkedDoublyNilIndex,
		free:       0,
		len:        0,
		cap:        capacity,
	}
	state.resetFreeList(0)
	return state
}

func allocListLinkedDoublyValueStorage[T any](capacity int) (any, unsafe.Pointer, uintptr) {
	var sample [1]T
	elemType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, elemType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr()), unsafe.Sizeof(sample[0])
}

func allocListLinkedDoublyIntStorage(capacity int) (any, unsafe.Pointer) {
	arrayType := reflect.ArrayOf(capacity, reflect.TypeOf(0))
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}

func (s *listLinkedDoublyState[T]) valuePtrAt(i int) *T {
	return (*T)(unsafe.Add(s.valueData, uintptr(i)*s.valueSize))
}

func (s *listLinkedDoublyState[T]) prevPtrAt(i int) *int {
	return (*int)(unsafe.Add(s.prevData, uintptr(i)*unsafe.Sizeof(0)))
}

func (s *listLinkedDoublyState[T]) nextPtrAt(i int) *int {
	return (*int)(unsafe.Add(s.nextData, uintptr(i)*unsafe.Sizeof(0)))
}

func (s *listLinkedDoublyState[T]) resetFreeList(start int) {
	for i := start; i < s.cap; i++ {
		*s.prevPtrAt(i) = listLinkedDoublyNilIndex
		next := i + 1
		if i == s.cap-1 {
			next = listLinkedDoublyNilIndex
		}
		*s.nextPtrAt(i) = next
	}
	if start < s.cap {
		s.free = start
	} else {
		s.free = listLinkedDoublyNilIndex
	}
}
