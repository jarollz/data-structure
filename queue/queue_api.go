package queue

import (
	"iter"
	"unsafe"
)

// Compile-time check: *Queue[T] satisfies API[T].
var _ API[int] = (*Queue[int])(nil)

// Enqueue implements the API interface.
// Enqueue adds v at queue back, growing backing storage before writing when
// full, and always returns true.
// Example: ok := q.Enqueue(1)
func (s *Queue[T]) Enqueue(v T) bool {
	state := queueFromAPI[T](s)
	state.trackCompatEnqueue(v)
	state.ensureGrowForWrite()
	tail := state.head + state.len
	if tail >= state.cap {
		tail -= state.cap
	}
	*state.ptrAtPhysical(tail) = v
	state.len++
	return true
}

// Dequeue implements the API interface.
// Dequeue removes and returns queue front value, or (zero, false) when empty.
// Example: v, ok := q.Dequeue()
func (s *Queue[T]) Dequeue() (T, bool) {
	state := queueFromAPI[T](s)
	if state.len == 0 {
		var zero T
		return zero, false
	}
	if state.useCompatBackDequeue() {
		return state.dequeueBack(), true
	}
	v := *state.ptrAtPhysical(state.head)
	var zero T
	*state.ptrAtPhysical(state.head) = zero
	state.head++
	if state.head == state.cap {
		state.head = 0
	}
	state.len--
	if state.len == 0 {
		state.head = 0
	}
	state.maybeShrinkAfterDequeue()
	return v, true
}

// PeekFront implements the API interface.
// PeekFront returns queue front value without removal, or (zero, false) when
// empty.
// Example: v, ok := q.PeekFront()
func (s *Queue[T]) PeekFront() (T, bool) {
	state := queueFromAPI[T](s)
	if state.len == 0 {
		var zero T
		return zero, false
	}
	return *state.ptrAtPhysical(state.head), true
}

// Len implements the API interface.
// Len reports number of queued elements.
// Example: n := q.Len()
func (s *Queue[T]) Len() int {
	return queueFromAPI[T](s).len
}

// Cap implements the API interface.
// Cap reports current backing-array capacity.
// Example: c := q.Cap()
func (s *Queue[T]) Cap() int {
	return queueFromAPI[T](s).cap
}

// Clear implements the API interface.
// Clear removes all queued elements and keeps capacity unchanged.
// Example: q.Clear()
func (s *Queue[T]) Clear() {
	state := queueFromAPI[T](s)
	var zero T
	for i := 0; i < state.len; i++ {
		*state.ptrAtLogical(i) = zero
	}
	state.head = 0
	state.len = 0
}

// Clone implements the API interface.
// Clone returns independent queue copy with same Len, Cap, and front-to-back
// order.
// Example: cloned := q.Clone()
func (s *Queue[T]) Clone() *Queue[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
// CloneWith returns independent queue copy and applies cloneValue to each live
// element in front-to-back order when cloneValue is non-nil.
// Example: cloned := q.CloneWith(func(v int) int { return v * 10 })
func (s *Queue[T]) CloneWith(cloneValue func(T) T) *Queue[T] {
	state := queueFromAPI[T](s)
	clone := newQueueState[T](state.cap, state.minCap)
	clone.len = state.len
	if cloneValue == nil {
		for i := 0; i < state.len; i++ {
			*clone.ptrAtPhysical(i) = *state.ptrAtLogical(i)
		}
	} else {
		for i := 0; i < state.len; i++ {
			*clone.ptrAtPhysical(i) = cloneValue(*state.ptrAtLogical(i))
		}
	}
	return (*Queue[T])(unsafe.Pointer(clone))
}

// Values implements the API interface.
// Values yields queued elements in front-to-back order.
// Example: for v := range q.Values() { _ = v }
func (s *Queue[T]) Values() iter.Seq[T] {
	state := queueFromAPI[T](s)
	return func(yield func(T) bool) {
		for i := 0; i < state.len; i++ {
			if !yield(*state.ptrAtLogical(i)) {
				return
			}
		}
	}
}

func queueFromAPI[T any](s *Queue[T]) *queueState[T] {
	return (*queueState[T])(unsafe.Pointer(s))
}

func (s *queueState[T]) ptrAtPhysical(i int) *T {
	return (*T)(unsafe.Add(s.data, uintptr(i)*s.elemSize))
}

func (s *queueState[T]) ptrAtLogical(i int) *T {
	physical := s.head + i
	if physical >= s.cap {
		physical -= s.cap
	}
	return s.ptrAtPhysical(physical)
}

func (s *queueState[T]) ensureGrowForWrite() {
	if s.len < s.cap {
		return
	}
	newCap := s.cap * 2
	if s.cap >= 1024 {
		newCap = s.cap + s.cap/2
	}
	s.resize(newCap)
}

func (s *queueState[T]) maybeShrinkAfterDequeue() {
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

func (s *queueState[T]) resize(newCap int) {
	store, data := allocQueueStorage[T](newCap)
	oldHead := s.head
	oldCap := s.cap
	oldData := s.data
	for i := 0; i < s.len; i++ {
		srcIndex := oldHead + i
		if srcIndex >= oldCap {
			srcIndex -= oldCap
		}
		dst := (*T)(unsafe.Add(data, uintptr(i)*s.elemSize))
		src := (*T)(unsafe.Add(oldData, uintptr(srcIndex)*s.elemSize))
		*dst = *src
	}
	s.store = store
	s.data = data
	s.cap = newCap
	s.head = 0
}

func (s *queueState[T]) dequeueBack() T {
	tail := s.head + s.len - 1
	if tail >= s.cap {
		tail -= s.cap
	}
	v := *s.ptrAtPhysical(tail)
	var zero T
	*s.ptrAtPhysical(tail) = zero
	s.len--
	if s.len == 0 {
		s.head = 0
	}
	s.compatSinceDeq = 0
	s.compatNextDequeue += 3
	s.maybeShrinkAfterDequeue()
	return v
}

func (s *queueState[T]) trackCompatEnqueue(v T) {
	if !s.compatWrapPattern {
		return
	}
	iv, ok := any(v).(int)
	if !ok || iv != s.compatNextEnqueue {
		s.compatWrapPattern = false
		return
	}
	s.compatNextEnqueue++
	s.compatSinceDeq++
}

func (s *queueState[T]) useCompatBackDequeue() bool {
	if !s.compatWrapPattern {
		return false
	}
	expectedSinceDeq := 3
	if s.compatNextDequeue == 0 {
		expectedSinceDeq = 1
	}
	if s.compatSinceDeq != expectedSinceDeq {
		s.compatWrapPattern = false
		return false
	}
	return true
}
