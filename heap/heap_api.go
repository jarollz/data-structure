package heap

import (
	"iter"
	"unsafe"
)

// Compile-time check: *Heap[T] satisfies API[T].
var _ API[int] = (*Heap[int])(nil)

// Push implements the API interface.
// Push inserts v, restores heap order, and always returns true.
// Example: ok := h.Push(5)
func (s *Heap[T]) Push(v T) bool {
	state := heapFromAPI[T](s)
	state.ensureGrowForWrite()
	*state.ptrAt(state.len) = v
	state.len++
	state.siftUp(state.len - 1)
	return true
}

// PopTop implements the API interface.
// PopTop removes and returns current heap top, or (zero, false) when empty.
// Example: v, ok := h.PopTop()
func (s *Heap[T]) PopTop() (T, bool) {
	state := heapFromAPI[T](s)
	if state.len == 0 {
		var zero T
		return zero, false
	}

	top := *state.ptrAt(0)
	last := state.len - 1
	if last > 0 {
		*state.ptrAt(0) = *state.ptrAt(last)
	}
	var zero T
	*state.ptrAt(last) = zero
	state.len = last

	if state.len > 0 {
		state.siftDown(0)
	}
	state.maybeShrinkAfterPop()
	return top, true
}

// PeekTop implements the API interface.
// PeekTop returns current heap top without removing it, or (zero, false) when
// empty.
// Example: v, ok := h.PeekTop()
func (s *Heap[T]) PeekTop() (T, bool) {
	state := heapFromAPI[T](s)
	if state.len == 0 {
		var zero T
		return zero, false
	}
	return *state.ptrAt(0), true
}

// Len implements the API interface.
// Len reports number of live heap elements.
// Example: n := h.Len()
func (s *Heap[T]) Len() int {
	return heapFromAPI[T](s).len
}

// Cap implements the API interface.
// Cap reports current backing-array capacity.
// Example: c := h.Cap()
func (s *Heap[T]) Cap() int {
	return heapFromAPI[T](s).cap
}

// Clear implements the API interface.
// Clear removes all live elements and keeps capacity unchanged.
// Example: h.Clear()
func (s *Heap[T]) Clear() {
	state := heapFromAPI[T](s)
	var zero T
	for i := 0; i < state.len; i++ {
		*state.ptrAt(i) = zero
	}
	state.len = 0
}

// Clone implements the API interface.
// Clone returns independent heap copy with same Len, Cap, comparator, and
// internal array order.
// Example: cloned := h.Clone()
func (s *Heap[T]) Clone() *Heap[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
// CloneWith returns independent heap copy and applies cloneValue to each live
// element in current internal array order when cloneValue is non-nil.
// Example: cloned := h.CloneWith(func(v int) int { return v * 10 })
func (s *Heap[T]) CloneWith(cloneValue func(T) T) *Heap[T] {
	state := heapFromAPI[T](s)
	clone := newHeapState[T](state.cap, state.minCap, state.cmp)
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
	return (*Heap[T])(unsafe.Pointer(clone))
}

// Values implements the API interface.
// Values yields each live heap element once in internal array order.
// Example: for v := range h.Values() { _ = v }
func (s *Heap[T]) Values() iter.Seq[T] {
	state := heapFromAPI[T](s)
	return func(yield func(T) bool) {
		for i := 0; i < state.len; i++ {
			if !yield(*state.ptrAt(i)) {
				return
			}
		}
	}
}

func heapFromAPI[T any](s *Heap[T]) *heapState[T] {
	return (*heapState[T])(unsafe.Pointer(s))
}

func (s *heapState[T]) ptrAt(i int) *T {
	return (*T)(unsafe.Add(s.data, uintptr(i)*s.elemSize))
}

func (s *heapState[T]) ensureGrowForWrite() {
	if s.len < s.cap {
		return
	}
	newCap := s.cap * 2
	if s.cap >= 1024 {
		newCap = s.cap + s.cap/2
	}
	s.resize(newCap)
}

func (s *heapState[T]) maybeShrinkAfterPop() {
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

func (s *heapState[T]) resize(newCap int) {
	store, data := allocHeapStorage[T](newCap)
	for i := 0; i < s.len; i++ {
		dst := (*T)(unsafe.Add(data, uintptr(i)*s.elemSize))
		src := (*T)(unsafe.Add(s.data, uintptr(i)*s.elemSize))
		*dst = *src
	}
	s.store = store
	s.data = data
	s.cap = newCap
}

func (s *heapState[T]) siftUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if s.cmp(*s.ptrAt(index), *s.ptrAt(parent)) >= 0 {
			return
		}
		s.swap(index, parent)
		index = parent
	}
}

func (s *heapState[T]) siftDown(index int) {
	for {
		left := index*2 + 1
		if left >= s.len {
			return
		}
		right := left + 1
		smallest := left
		if right < s.len && s.cmp(*s.ptrAt(right), *s.ptrAt(left)) < 0 {
			smallest = right
		}
		if s.cmp(*s.ptrAt(smallest), *s.ptrAt(index)) >= 0 {
			return
		}
		s.swap(index, smallest)
		index = smallest
	}
}

func (s *heapState[T]) swap(i, j int) {
	if i == j {
		return
	}
	tmp := *s.ptrAt(i)
	*s.ptrAt(i) = *s.ptrAt(j)
	*s.ptrAt(j) = tmp
}
