package stack

import (
	"iter"
	"unsafe"
)

// Compile-time check: *Stack[T] satisfies API[T].
var _ API[int] = (*Stack[int])(nil)

// Push implements the API interface.
// Push adds v on stack top, grows backing storage before writing when full,
// and always returns true.
// Example: ok := s.Push(1)
func (s *Stack[T]) Push(v T) bool {
	state := stackFromAPI[T](s)
	state.ensureGrowForWrite()
	*state.ptrAt(state.len) = v
	state.len++
	return true
}

// Pop implements the API interface.
// Pop removes and returns current top value, or (zero, false) when empty.
// Example: v, ok := s.Pop()
func (s *Stack[T]) Pop() (T, bool) {
	state := stackFromAPI[T](s)
	if state.len == 0 {
		var zero T
		return zero, false
	}
	top := state.len - 1
	v := *state.ptrAt(top)
	var zero T
	*state.ptrAt(top) = zero
	state.len = top
	state.maybeShrinkAfterPop()
	return v, true
}

// PeekTop implements the API interface.
// PeekTop returns current top value without removing it, or (zero, false) when
// empty.
// Example: v, ok := s.PeekTop()
func (s *Stack[T]) PeekTop() (T, bool) {
	state := stackFromAPI[T](s)
	if state.len == 0 {
		var zero T
		return zero, false
	}
	return *state.ptrAt(state.len - 1), true
}

// Len implements the API interface.
// Len reports number of live stack elements.
// Example: n := s.Len()
func (s *Stack[T]) Len() int {
	return stackFromAPI[T](s).len
}

// Cap implements the API interface.
// Cap reports current backing-array capacity.
// Example: c := s.Cap()
func (s *Stack[T]) Cap() int {
	return stackFromAPI[T](s).cap
}

// Clear implements the API interface.
// Clear removes all live elements and keeps capacity unchanged.
// Example: s.Clear()
func (s *Stack[T]) Clear() {
	state := stackFromAPI[T](s)
	var zero T
	for i := 0; i < state.len; i++ {
		*state.ptrAt(i) = zero
	}
	state.len = 0
}

// Clone implements the API interface.
// Clone returns independent stack copy with same Len, Cap, and top-to-bottom
// order.
// Example: cloned := s.Clone()
func (s *Stack[T]) Clone() *Stack[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
// CloneWith returns independent stack copy and applies cloneValue to each live
// element in top-to-bottom order when cloneValue is non-nil.
// Example: cloned := s.CloneWith(func(v int) int { return v * 10 })
func (s *Stack[T]) CloneWith(cloneValue func(T) T) *Stack[T] {
	state := stackFromAPI[T](s)
	clone := newStackState[T](state.cap, state.minCap)
	clone.len = state.len
	if cloneValue == nil {
		for i := 0; i < state.len; i++ {
			*clone.ptrAt(i) = *state.ptrAt(i)
		}
	} else {
		for i := state.len - 1; i >= 0; i-- {
			*clone.ptrAt(i) = cloneValue(*state.ptrAt(i))
		}
	}
	return (*Stack[T])(unsafe.Pointer(clone))
}

// Values implements the API interface.
// Values yields live stack values from top to bottom.
// Example: for v := range s.Values() { _ = v }
func (s *Stack[T]) Values() iter.Seq[T] {
	state := stackFromAPI[T](s)
	return func(yield func(T) bool) {
		for i := state.len - 1; i >= 0; i-- {
			if !yield(*state.ptrAt(i)) {
				return
			}
		}
	}
}

func stackFromAPI[T any](s *Stack[T]) *stackState[T] {
	return (*stackState[T])(unsafe.Pointer(s))
}

func (s *stackState[T]) ptrAt(i int) *T {
	return (*T)(unsafe.Add(s.data, uintptr(i)*s.elemSize))
}

func (s *stackState[T]) ensureGrowForWrite() {
	if s.len < s.cap {
		return
	}
	newCap := s.cap * 2
	if s.cap >= 1024 {
		newCap = s.cap + s.cap/2
	}
	s.resize(newCap)
}

func (s *stackState[T]) maybeShrinkAfterPop() {
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

func (s *stackState[T]) resize(newCap int) {
	store, data := allocStackStorage[T](newCap)
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
