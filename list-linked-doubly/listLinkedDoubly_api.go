package listlinkeddoubly

import "iter"
import "unsafe"

// Compile-time check: *ListLinkedDoubly[T] satisfies API[T].
var _ API[int] = (*ListLinkedDoubly[int])(nil)

// PushFront implements the API interface.
// PushFront inserts v at list head.
// Example: list.PushFront(4)
func (s *ListLinkedDoubly[T]) PushFront(v T) {
	state := listLinkedDoublyFromAPI[T](s)
	idx := state.allocateNode(v)
	if state.head == listLinkedDoublyNilIndex {
		state.head = idx
		state.tail = idx
		return
	}
	oldHead := state.head
	*state.nextPtrAt(idx) = oldHead
	*state.prevPtrAt(oldHead) = idx
	state.head = idx
}

// PushBack implements the API interface.
// PushBack inserts v at list tail.
// Example: list.PushBack(8)
func (s *ListLinkedDoubly[T]) PushBack(v T) {
	state := listLinkedDoublyFromAPI[T](s)
	idx := state.allocateNode(v)
	if state.tail == listLinkedDoublyNilIndex {
		state.head = idx
		state.tail = idx
		return
	}
	oldTail := state.tail
	*state.prevPtrAt(idx) = oldTail
	*state.nextPtrAt(oldTail) = idx
	state.tail = idx
}

// PopFront implements the API interface.
// PopFront removes and returns head value.
// Example: v, ok := list.PopFront()
func (s *ListLinkedDoubly[T]) PopFront() (T, bool) {
	state := listLinkedDoublyFromAPI[T](s)
	if state.head == listLinkedDoublyNilIndex {
		var zero T
		return zero, false
	}
	idx := state.head
	v := *state.valuePtrAt(idx)
	newHead := *state.nextPtrAt(idx)
	state.head = newHead
	if newHead == listLinkedDoublyNilIndex {
		state.tail = listLinkedDoublyNilIndex
	} else {
		*state.prevPtrAt(newHead) = listLinkedDoublyNilIndex
	}
	state.releaseNode(idx)
	return v, true
}

// PopBack implements the API interface.
// PopBack removes and returns tail value.
// Example: v, ok := list.PopBack()
func (s *ListLinkedDoubly[T]) PopBack() (T, bool) {
	state := listLinkedDoublyFromAPI[T](s)
	if state.tail == listLinkedDoublyNilIndex {
		var zero T
		return zero, false
	}
	idx := state.tail
	v := *state.valuePtrAt(idx)
	newTail := *state.prevPtrAt(idx)
	state.tail = newTail
	if newTail == listLinkedDoublyNilIndex {
		state.head = listLinkedDoublyNilIndex
	} else {
		*state.nextPtrAt(newTail) = listLinkedDoublyNilIndex
	}
	state.releaseNode(idx)
	return v, true
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := list.Len()
func (s *ListLinkedDoubly[T]) Len() int {
	return listLinkedDoublyFromAPI[T](s).len
}

// Clear implements the API interface.
// Clear removes all nodes and resets structural state.
// Example: list.Clear()
func (s *ListLinkedDoubly[T]) Clear() {
	state := listLinkedDoublyFromAPI[T](s)
	if state.len == 0 {
		return
	}
	var zero T
	for i := 0; i < state.cap; i++ {
		*state.valuePtrAt(i) = zero
	}
	state.head = listLinkedDoublyNilIndex
	state.tail = listLinkedDoublyNilIndex
	state.len = 0
	state.resetFreeList(0)
}

// Clone implements the API interface.
// Clone returns independent list copy with same length and head-to-tail order.
// Example: cloned := list.Clone()
func (s *ListLinkedDoubly[T]) Clone() *ListLinkedDoubly[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
// CloneWith returns independent list copy using cloneValue for each live node.
// Example: cloned := list.CloneWith(func(v int) int { return v * 10 })
func (s *ListLinkedDoubly[T]) CloneWith(cloneValue func(T) T) *ListLinkedDoubly[T] {
	state := listLinkedDoublyFromAPI[T](s)
	cloneState := newListLinkedDoublyState[T](state.cap)
	for idx := state.head; idx != listLinkedDoublyNilIndex; idx = *state.nextPtrAt(idx) {
		v := *state.valuePtrAt(idx)
		if cloneValue != nil {
			v = cloneValue(v)
		}
		cloneState.appendKnownValue(v)
	}
	return (*ListLinkedDoubly[T])(unsafe.Pointer(cloneState))
}

// Values implements the API interface.
// Values yields values from head to tail.
// Example: for v := range list.Values() { _ = v }
func (s *ListLinkedDoubly[T]) Values() iter.Seq[T] {
	state := listLinkedDoublyFromAPI[T](s)
	return func(yield func(T) bool) {
		for idx := state.head; idx != listLinkedDoublyNilIndex; idx = *state.nextPtrAt(idx) {
			if !yield(*state.valuePtrAt(idx)) {
				return
			}
		}
	}
}

func listLinkedDoublyFromAPI[T any](s *ListLinkedDoubly[T]) *listLinkedDoublyState[T] {
	return (*listLinkedDoublyState[T])(unsafe.Pointer(s))
}

func (s *listLinkedDoublyState[T]) allocateNode(v T) int {
	if s.free == listLinkedDoublyNilIndex {
		s.grow()
	}
	idx := s.free
	s.free = *s.nextPtrAt(idx)
	*s.valuePtrAt(idx) = v
	*s.prevPtrAt(idx) = listLinkedDoublyNilIndex
	*s.nextPtrAt(idx) = listLinkedDoublyNilIndex
	s.len++
	return idx
}

func (s *listLinkedDoublyState[T]) appendKnownValue(v T) {
	idx := s.allocateNode(v)
	if s.tail == listLinkedDoublyNilIndex {
		s.head = idx
		s.tail = idx
		return
	}
	oldTail := s.tail
	*s.prevPtrAt(idx) = oldTail
	*s.nextPtrAt(oldTail) = idx
	s.tail = idx
}

func (s *listLinkedDoublyState[T]) releaseNode(idx int) {
	var zero T
	*s.valuePtrAt(idx) = zero
	*s.prevPtrAt(idx) = listLinkedDoublyNilIndex
	*s.nextPtrAt(idx) = s.free
	s.free = idx
	s.len--
}

func (s *listLinkedDoublyState[T]) grow() {
	newCap := s.cap * 2
	if newCap < listLinkedDoublyDefaultCap {
		newCap = listLinkedDoublyDefaultCap
	}
	newValueStore, newValueData, newValueSize := allocListLinkedDoublyValueStorage[T](newCap)
	newPrevStore, newPrevData := allocListLinkedDoublyIntStorage(newCap)
	newNextStore, newNextData := allocListLinkedDoublyIntStorage(newCap)

	for i := 0; i < s.cap; i++ {
		dstValue := (*T)(unsafe.Add(newValueData, uintptr(i)*newValueSize))
		srcValue := s.valuePtrAt(i)
		*dstValue = *srcValue
		dstPrev := (*int)(unsafe.Add(newPrevData, uintptr(i)*unsafe.Sizeof(0)))
		*dstPrev = *s.prevPtrAt(i)
		dstNext := (*int)(unsafe.Add(newNextData, uintptr(i)*unsafe.Sizeof(0)))
		*dstNext = *s.nextPtrAt(i)
	}
	for i := s.cap; i < newCap; i++ {
		*(*int)(unsafe.Add(newPrevData, uintptr(i)*unsafe.Sizeof(0))) = listLinkedDoublyNilIndex
		next := i + 1
		if i == newCap-1 {
			next = listLinkedDoublyNilIndex
		}
		*(*int)(unsafe.Add(newNextData, uintptr(i)*unsafe.Sizeof(0))) = next
	}

	oldCap := s.cap
	s.valueStore = newValueStore
	s.valueData = newValueData
	s.valueSize = newValueSize
	s.prevStore = newPrevStore
	s.prevData = newPrevData
	s.nextStore = newNextStore
	s.nextData = newNextData
	s.cap = newCap
	s.free = oldCap
}
