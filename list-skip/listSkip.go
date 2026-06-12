package listskip

import (
	"reflect"
	"unsafe"
)

const (
	listSkipNilIndex   = -1
	listSkipHeadIndex  = 0
	listSkipDefaultCap = 16
)

type listSkipState[T any] struct {
	cmp func(a, b T) int

	maxLevel     int
	currentLevel int
	len          int
	rng          uint32

	cap int

	valueStore any
	valueData  unsafe.Pointer
	valueSize  uintptr

	nextStore any
	nextData  unsafe.Pointer

	freeStore any
	freeData  unsafe.Pointer

	updateStore any
	updateData  unsafe.Pointer

	free int
}

// New creates empty skip list using comparator-defined sorted order.
//
// New normalizes maxLevel < 1 to 1. Returned list is non-nil, has Len() == 0,
// starts with currentLevel == 1, and starts deterministic xorshift32 state at 1.
//
// Example: s := New[int](8, cmpInt)
func New[T any](maxLevel int, cmp func(a, b T) int) *ListSkip[T] {
	if maxLevel < 1 {
		maxLevel = 1
	}
	state := newListSkipState[T](maxLevel, listSkipDefaultCap, cmp)
	return (*ListSkip[T])(unsafe.Pointer(state))
}

func newListSkipState[T any](maxLevel, capacity int, cmp func(a, b T) int) *listSkipState[T] {
	if maxLevel < 1 {
		maxLevel = 1
	}
	if capacity < 2 {
		capacity = 2
	}
	if capacity < listSkipDefaultCap {
		capacity = listSkipDefaultCap
	}
	valueStore, valueData, valueSize := allocListSkipValueStorage[T](capacity)
	nextStore, nextData := allocListSkipIntStorage(capacity * maxLevel)
	freeStore, freeData := allocListSkipIntStorage(capacity)
	updateStore, updateData := allocListSkipIntStorage(maxLevel)

	state := &listSkipState[T]{
		cmp:          cmp,
		maxLevel:     maxLevel,
		currentLevel: 1,
		len:          0,
		rng:          1,
		cap:          capacity,
		valueStore:   valueStore,
		valueData:    valueData,
		valueSize:    valueSize,
		nextStore:    nextStore,
		nextData:     nextData,
		freeStore:    freeStore,
		freeData:     freeData,
		updateStore:  updateStore,
		updateData:   updateData,
		free:         listSkipNilIndex,
	}

	state.resetHeadLinks()
	state.resetFreeList()
	return state
}

func allocListSkipValueStorage[T any](capacity int) (any, unsafe.Pointer, uintptr) {
	var sample [1]T
	elemType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, elemType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr()), unsafe.Sizeof(sample[0])
}

func allocListSkipIntStorage(length int) (any, unsafe.Pointer) {
	arrayType := reflect.ArrayOf(length, reflect.TypeOf(0))
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}
