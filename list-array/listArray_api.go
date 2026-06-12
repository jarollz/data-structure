package listarray

import (
	"iter"
	"reflect"
	"unsafe"
)

// Compile-time check: *ListArray[T] satisfies API[T].
var _ API[int] = (*ListArray[int])(nil)

// Append implements the API interface.
//
// Append adds v at tail, growing backing storage when full.
//
// Example: ok := list.Append(7)
func (s *ListArray[T]) Append(v T) bool {
	state := listArrayFromAPI[T](s)
	state.ensureGrowForWrite()
	*state.ptrAt(state.len) = v
	state.len++
	return true
}

// Get implements the API interface.
//
// Get returns value at i or (zero, false) when i is out of range.
//
// Example: v, ok := list.Get(2)
func (s *ListArray[T]) Get(i int) (T, bool) {
	state := listArrayFromAPI[T](s)
	if i < 0 || i >= state.len {
		var zero T
		return zero, false
	}
	return *state.ptrAt(i), true
}

// Set implements the API interface.
//
// Set overwrites value at i when i is in [0, Len()).
//
// Example: ok := list.Set(1, 42)
func (s *ListArray[T]) Set(i int, v T) bool {
	state := listArrayFromAPI[T](s)
	if i < 0 || i >= state.len {
		return false
	}
	*state.ptrAt(i) = v
	return true
}

// Insert implements the API interface.
//
// Insert places v before i when i is in [0, Len()].
//
// Example: ok := list.Insert(0, 9)
func (s *ListArray[T]) Insert(i int, v T) bool {
	state := listArrayFromAPI[T](s)
	if i < 0 || i > state.len {
		return false
	}
	state.ensureGrowForWrite()
	for j := state.len; j > i; j-- {
		*state.ptrAt(j) = *state.ptrAt(j - 1)
	}
	*state.ptrAt(i) = v
	state.len++
	return true
}

// Delete implements the API interface.
//
// Delete removes and returns value at i or (zero, false) when i is out of range.
//
// Example: v, ok := list.Delete(3)
func (s *ListArray[T]) Delete(i int) (T, bool) {
	state := listArrayFromAPI[T](s)
	if i < 0 || i >= state.len {
		var zero T
		return zero, false
	}
	removed := *state.ptrAt(i)
	for j := i; j < state.len-1; j++ {
		*state.ptrAt(j) = *state.ptrAt(j + 1)
	}
	state.len--
	var zero T
	*state.ptrAt(state.len) = zero
	state.maybeShrinkAfterDelete()
	return removed, true
}

// Len implements the API interface.
//
// Len reports number of live elements.
//
// Example: n := list.Len()
func (s *ListArray[T]) Len() int {
	return listArrayFromAPI[T](s).len
}

// Cap implements the API interface.
//
// Cap reports current backing-array capacity.
//
// Example: c := list.Cap()
func (s *ListArray[T]) Cap() int {
	return listArrayFromAPI[T](s).cap
}

// Clear implements the API interface.
//
// Clear removes all live elements and keeps capacity unchanged.
//
// Example: list.Clear()
func (s *ListArray[T]) Clear() {
	state := listArrayFromAPI[T](s)
	var zero T
	for i := 0; i < state.len; i++ {
		*state.ptrAt(i) = zero
	}
	state.len = 0
}

// Clone implements the API interface.
//
// Clone returns independent list copy with same Len, Cap, and order.
//
// Example: cloned := list.Clone()
func (s *ListArray[T]) Clone() *ListArray[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
//
// CloneWith returns independent copy and applies cloneValue to each live element
// in ascending index order when cloneValue is non-nil.
//
// Example: cloned := list.CloneWith(func(v int) int { return v * 10 })
func (s *ListArray[T]) CloneWith(cloneValue func(T) T) *ListArray[T] {
	state := listArrayFromAPI[T](s)
	clone := newListArrayState[T](state.cap, state.minCap)
	clone.len = state.len
	if cloneValue == nil {
		for i := 0; i < state.len; i++ {
			*clone.ptrAt(i) = *state.ptrAt(i)
		}
	} else {
		for i := 0; i < state.len; i++ {
			*clone.ptrAt(i) = cloneValue(*state.ptrAt(i))
		}
	}
	return (*ListArray[T])(unsafe.Pointer(clone))
}

// Values implements the API interface.
//
// Values yields live elements in index order from 0 to Len()-1.
//
// Example: for v := range list.Values() { _ = v }
func (s *ListArray[T]) Values() iter.Seq[T] {
	state := listArrayFromAPI[T](s)
	return func(yield func(T) bool) {
		for i := 0; i < state.len; i++ {
			if !yield(*state.ptrAt(i)) {
				return
			}
		}
	}
}

func listArrayFromAPI[T any](s *ListArray[T]) *listArrayState[T] {
	return (*listArrayState[T])(unsafe.Pointer(s))
}

func allocListArrayStorage[T any](capacity int) (any, unsafe.Pointer) {
	var sample [1]T
	elemType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, elemType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}

func (s *listArrayState[T]) ptrAt(i int) *T {
	return (*T)(unsafe.Add(s.data, uintptr(i)*s.elemSize))
}

func (s *listArrayState[T]) ensureGrowForWrite() {
	if s.len < s.cap {
		return
	}
	newCap := s.cap * 2
	if s.cap >= 1024 {
		newCap = s.cap + s.cap/2
	}
	s.resize(newCap)
}

func (s *listArrayState[T]) maybeShrinkAfterDelete() {
	if s.cap <= s.minCap {
		return
	}
	if s.len > s.cap/4 {
		return
	}
	newCap := s.cap / 2
	if twiceLen := s.len * 2; newCap < twiceLen {
		newCap = twiceLen
	}
	if newCap < s.minCap {
		newCap = s.minCap
	}
	if newCap >= s.cap {
		return
	}
	s.resize(newCap)
}

func (s *listArrayState[T]) resize(newCap int) {
	store, data := allocListArrayStorage[T](newCap)
	oldData := s.data
	for i := 0; i < s.len; i++ {
		dst := (*T)(unsafe.Add(data, uintptr(i)*s.elemSize))
		src := (*T)(unsafe.Add(oldData, uintptr(i)*s.elemSize))
		*dst = *src
	}
	s.store = store
	s.data = data
	s.cap = newCap
}
