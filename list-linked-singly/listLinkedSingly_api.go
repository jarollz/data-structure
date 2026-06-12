package listlinkedsingly

import "iter"
import "unsafe"

// Compile-time check: *ListLinkedSingly[T] satisfies API[T].
var _ API[int] = (*ListLinkedSingly[int])(nil)

// PushFront implements the API interface.
// PushFront inserts v at list head.
// Example: list.PushFront(5)
func (s *ListLinkedSingly[T]) PushFront(v T) {
	state := listLinkedSinglyFromAPI[T](s)
	idx := state.allocateNode(v)
	*state.nextPtrAt(idx) = state.head
	state.head = idx
	if state.tail == listLinkedSinglyNilIndex {
		state.tail = idx
	}
}

// PopFront implements the API interface.
// PopFront removes and returns head value.
// Example: v, ok := list.PopFront()
func (s *ListLinkedSingly[T]) PopFront() (T, bool) {
	state := listLinkedSinglyFromAPI[T](s)
	if state.head == listLinkedSinglyNilIndex {
		var zero T
		return zero, false
	}
	idx := state.head
	v := *state.valuePtrAt(idx)
	state.head = *state.nextPtrAt(idx)
	if state.head == listLinkedSinglyNilIndex {
		state.tail = listLinkedSinglyNilIndex
	}
	state.releaseNode(idx)
	return v, true
}

// Append implements the API interface.
// Append adds v at list tail in O(1) time.
// Example: list.Append(9)
func (s *ListLinkedSingly[T]) Append(v T) {
	state := listLinkedSinglyFromAPI[T](s)
	idx := state.allocateNode(v)
	if state.tail == listLinkedSinglyNilIndex {
		state.head = idx
		state.tail = idx
		return
	}
	*state.nextPtrAt(state.tail) = idx
	state.tail = idx
}

// DeleteFirst implements the API interface.
// DeleteFirst removes first node that satisfies match.
// Example: ok := list.DeleteFirst(func(v int) bool { return v == 3 })
func (s *ListLinkedSingly[T]) DeleteFirst(match func(T) bool) bool {
	state := listLinkedSinglyFromAPI[T](s)
	prev := listLinkedSinglyNilIndex
	cur := state.head
	for cur != listLinkedSinglyNilIndex {
		if match(*state.valuePtrAt(cur)) {
			next := *state.nextPtrAt(cur)
			if prev == listLinkedSinglyNilIndex {
				state.head = next
			} else {
				*state.nextPtrAt(prev) = next
			}
			if cur == state.tail {
				state.tail = prev
			}
			state.releaseNode(cur)
			return true
		}
		prev = cur
		cur = *state.nextPtrAt(cur)
	}
	return false
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := list.Len()
func (s *ListLinkedSingly[T]) Len() int {
	return listLinkedSinglyFromAPI[T](s).len
}

// Clear implements the API interface.
// Clear removes all nodes and resets structural state.
// Example: list.Clear()
func (s *ListLinkedSingly[T]) Clear() {
	state := listLinkedSinglyFromAPI[T](s)
	if state.len == 0 {
		return
	}
	var zero T
	for i := 0; i < state.cap; i++ {
		*state.valuePtrAt(i) = zero
	}
	state.head = listLinkedSinglyNilIndex
	state.tail = listLinkedSinglyNilIndex
	state.len = 0
	state.resetFreeList(0)
}

// Clone implements the API interface.
// Clone returns independent list copy with same length and head-to-tail order.
// Example: cloned := list.Clone()
func (s *ListLinkedSingly[T]) Clone() *ListLinkedSingly[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
// CloneWith returns independent list copy using cloneValue for each live node.
// Example: cloned := list.CloneWith(func(v int) int { return v * 10 })
func (s *ListLinkedSingly[T]) CloneWith(cloneValue func(T) T) *ListLinkedSingly[T] {
	state := listLinkedSinglyFromAPI[T](s)
	cloneState := newListLinkedSinglyState[T](state.cap)
	for idx := state.head; idx != listLinkedSinglyNilIndex; idx = *state.nextPtrAt(idx) {
		v := *state.valuePtrAt(idx)
		if cloneValue != nil {
			v = cloneValue(v)
		}
		cloneState.appendKnownValue(v)
	}
	return (*ListLinkedSingly[T])(unsafe.Pointer(cloneState))
}

// Values implements the API interface.
// Values yields values from head to tail.
// Example: for v := range list.Values() { _ = v }
func (s *ListLinkedSingly[T]) Values() iter.Seq[T] {
	state := listLinkedSinglyFromAPI[T](s)
	return func(yield func(T) bool) {
		for idx := state.head; idx != listLinkedSinglyNilIndex; idx = *state.nextPtrAt(idx) {
			if !yield(*state.valuePtrAt(idx)) {
				return
			}
		}
	}
}

func listLinkedSinglyFromAPI[T any](s *ListLinkedSingly[T]) *listLinkedSinglyState[T] {
	return (*listLinkedSinglyState[T])(unsafe.Pointer(s))
}

func (s *listLinkedSinglyState[T]) allocateNode(v T) int {
	if s.free == listLinkedSinglyNilIndex {
		s.grow()
	}
	idx := s.free
	s.free = *s.nextPtrAt(idx)
	*s.valuePtrAt(idx) = v
	*s.nextPtrAt(idx) = listLinkedSinglyNilIndex
	s.len++
	return idx
}

func (s *listLinkedSinglyState[T]) appendKnownValue(v T) {
	idx := s.allocateNode(v)
	if s.tail == listLinkedSinglyNilIndex {
		s.head = idx
		s.tail = idx
		return
	}
	*s.nextPtrAt(s.tail) = idx
	s.tail = idx
}

func (s *listLinkedSinglyState[T]) releaseNode(idx int) {
	var zero T
	*s.valuePtrAt(idx) = zero
	*s.nextPtrAt(idx) = s.free
	s.free = idx
	s.len--
}

func (s *listLinkedSinglyState[T]) grow() {
	newCap := s.cap * 2
	if newCap < listLinkedSinglyDefaultCap {
		newCap = listLinkedSinglyDefaultCap
	}
	newValueStore, newValueData, newValueSize := allocListLinkedSinglyValueStorage[T](newCap)
	newNextStore, newNextData := allocListLinkedSinglyNextStorage(newCap)

	for i := 0; i < s.cap; i++ {
		dstValue := (*T)(unsafe.Add(newValueData, uintptr(i)*newValueSize))
		srcValue := s.valuePtrAt(i)
		*dstValue = *srcValue
		dstNext := (*int)(unsafe.Add(newNextData, uintptr(i)*unsafe.Sizeof(0)))
		*dstNext = *s.nextPtrAt(i)
	}
	for i := s.cap; i < newCap; i++ {
		next := i + 1
		if i == newCap-1 {
			next = listLinkedSinglyNilIndex
		}
		*(*int)(unsafe.Add(newNextData, uintptr(i)*unsafe.Sizeof(0))) = next
	}

	oldCap := s.cap
	s.valueStore = newValueStore
	s.valueData = newValueData
	s.valueSize = newValueSize
	s.nextStore = newNextStore
	s.nextData = newNextData
	s.cap = newCap
	s.free = oldCap
}
