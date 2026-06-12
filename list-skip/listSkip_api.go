package listskip

import "iter"
import "unsafe"

// Compile-time check: *ListSkip[T] satisfies API[T].
var _ API[int] = (*ListSkip[int])(nil)

// Insert implements the API interface.
//
// Insert adds v when v is missing. Duplicate values are rejected.
//
// Example: ok := list.Insert(10)
func (s *ListSkip[T]) Insert(v T) bool {
	state := listSkipFromAPI[T](s)
	state.findPredecessors(v)
	next := state.nextAt(*state.updatePtrAt(0), 0)
	if next != listSkipNilIndex && state.cmp(*state.valuePtrAt(next), v) == 0 {
		return false
	}

	level := state.nextNodeLevel()
	if level > state.currentLevel {
		for i := state.currentLevel; i < level; i++ {
			*state.updatePtrAt(i) = listSkipHeadIndex
		}
		state.currentLevel = level
	}

	idx := state.allocateNode(v)
	for i := 0; i < level; i++ {
		prev := *state.updatePtrAt(i)
		next = state.nextAt(prev, i)
		state.setNextAt(idx, i, next)
		state.setNextAt(prev, i, idx)
	}
	state.len++
	return true
}

// Delete implements the API interface.
//
// Delete removes v when present and keeps all levels sorted.
//
// Example: ok := list.Delete(10)
func (s *ListSkip[T]) Delete(v T) bool {
	state := listSkipFromAPI[T](s)
	state.findPredecessors(v)

	target := state.nextAt(*state.updatePtrAt(0), 0)
	if target == listSkipNilIndex || state.cmp(*state.valuePtrAt(target), v) != 0 {
		return false
	}

	for i := 0; i < state.currentLevel; i++ {
		prev := *state.updatePtrAt(i)
		if state.nextAt(prev, i) != target {
			continue
		}
		state.setNextAt(prev, i, state.nextAt(target, i))
	}

	state.releaseNode(target)
	state.len--

	for state.currentLevel > 1 && state.nextAt(listSkipHeadIndex, state.currentLevel-1) == listSkipNilIndex {
		state.currentLevel--
	}
	return true
}

// Has implements the API interface.
//
// Has reports whether v exists as live value.
//
// Example: ok := list.Has(10)
func (s *ListSkip[T]) Has(v T) bool {
	state := listSkipFromAPI[T](s)
	idx := listSkipHeadIndex
	for level := state.currentLevel - 1; level >= 0; level-- {
		for {
			next := state.nextAt(idx, level)
			if next == listSkipNilIndex {
				break
			}
			cmp := state.cmp(*state.valuePtrAt(next), v)
			if cmp < 0 {
				idx = next
				continue
			}
			if cmp == 0 {
				return true
			}
			break
		}
	}
	return false
}

// Len implements the API interface.
//
// Len reports number of live values.
//
// Example: n := list.Len()
func (s *ListSkip[T]) Len() int {
	return listSkipFromAPI[T](s).len
}

// Clear implements the API interface.
//
// Clear removes all values and resets structure to empty logical state.
//
// Example: list.Clear()
func (s *ListSkip[T]) Clear() {
	state := listSkipFromAPI[T](s)
	var zero T
	for i := 1; i < state.cap; i++ {
		*state.valuePtrAt(i) = zero
		for level := 0; level < state.maxLevel; level++ {
			state.setNextAt(i, level, listSkipNilIndex)
		}
	}
	state.resetHeadLinks()
	state.resetFreeList()
	state.len = 0
	state.currentLevel = 1
}

// Clone implements the API interface.
//
// Clone returns independent copy with same ordering, levels, and RNG state.
//
// Example: cloned := list.Clone()
func (s *ListSkip[T]) Clone() *ListSkip[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
//
// CloneWith returns independent copy and applies cloneValue for each live value
// in sorted order when cloneValue is non-nil.
//
// Example: cloned := list.CloneWith(func(v int) int { return v * 10 })
func (s *ListSkip[T]) CloneWith(cloneValue func(T) T) *ListSkip[T] {
	state := listSkipFromAPI[T](s)
	clone := newListSkipState[T](state.maxLevel, state.cap, state.cmp)

	clone.currentLevel = state.currentLevel
	clone.len = state.len
	clone.rng = state.rng
	clone.free = state.free

	for i := 0; i < state.cap; i++ {
		*clone.freePtrAt(i) = *state.freePtrAt(i)
		for level := 0; level < state.maxLevel; level++ {
			clone.setNextAt(i, level, state.nextAt(i, level))
		}
	}

	if cloneValue == nil {
		for i := 1; i < state.cap; i++ {
			*clone.valuePtrAt(i) = *state.valuePtrAt(i)
		}
		return (*ListSkip[T])(unsafe.Pointer(clone))
	}

	for idx := state.nextAt(listSkipHeadIndex, 0); idx != listSkipNilIndex; idx = state.nextAt(idx, 0) {
		*clone.valuePtrAt(idx) = cloneValue(*state.valuePtrAt(idx))
	}

	return (*ListSkip[T])(unsafe.Pointer(clone))
}

// Values implements the API interface.
//
// Values yields each live value once in sorted order using level-0 links.
//
// Example: for v := range list.Values() { _ = v }
func (s *ListSkip[T]) Values() iter.Seq[T] {
	state := listSkipFromAPI[T](s)
	return func(yield func(T) bool) {
		for idx := state.nextAt(listSkipHeadIndex, 0); idx != listSkipNilIndex; idx = state.nextAt(idx, 0) {
			if !yield(*state.valuePtrAt(idx)) {
				return
			}
		}
	}
}

func listSkipFromAPI[T any](s *ListSkip[T]) *listSkipState[T] {
	return (*listSkipState[T])(unsafe.Pointer(s))
}

func (s *listSkipState[T]) valuePtrAt(i int) *T {
	return (*T)(unsafe.Add(s.valueData, uintptr(i)*s.valueSize))
}

func (s *listSkipState[T]) nextPtrAt(i, level int) *int {
	idx := i*s.maxLevel + level
	return (*int)(unsafe.Add(s.nextData, uintptr(idx)*unsafe.Sizeof(0)))
}

func (s *listSkipState[T]) freePtrAt(i int) *int {
	return (*int)(unsafe.Add(s.freeData, uintptr(i)*unsafe.Sizeof(0)))
}

func (s *listSkipState[T]) updatePtrAt(level int) *int {
	return (*int)(unsafe.Add(s.updateData, uintptr(level)*unsafe.Sizeof(0)))
}

func (s *listSkipState[T]) nextAt(i, level int) int {
	return *s.nextPtrAt(i, level)
}

func (s *listSkipState[T]) setNextAt(i, level, next int) {
	*s.nextPtrAt(i, level) = next
}

func (s *listSkipState[T]) resetHeadLinks() {
	for level := 0; level < s.maxLevel; level++ {
		s.setNextAt(listSkipHeadIndex, level, listSkipNilIndex)
	}
}

func (s *listSkipState[T]) resetFreeList() {
	if s.cap <= 1 {
		s.free = listSkipNilIndex
		return
	}
	for i := 1; i < s.cap; i++ {
		next := i + 1
		if i == s.cap-1 {
			next = listSkipNilIndex
		}
		*s.freePtrAt(i) = next
	}
	*s.freePtrAt(listSkipHeadIndex) = listSkipNilIndex
	s.free = 1
}

func (s *listSkipState[T]) findPredecessors(v T) {
	idx := listSkipHeadIndex
	for level := s.currentLevel - 1; level >= 0; level-- {
		for {
			next := s.nextAt(idx, level)
			if next == listSkipNilIndex {
				break
			}
			if s.cmp(*s.valuePtrAt(next), v) >= 0 {
				break
			}
			idx = next
		}
		*s.updatePtrAt(level) = idx
	}
}

func (s *listSkipState[T]) allocateNode(v T) int {
	if s.free == listSkipNilIndex {
		s.grow()
	}
	idx := s.free
	s.free = *s.freePtrAt(idx)
	*s.valuePtrAt(idx) = v
	for level := 0; level < s.maxLevel; level++ {
		s.setNextAt(idx, level, listSkipNilIndex)
	}
	*s.freePtrAt(idx) = listSkipNilIndex
	return idx
}

func (s *listSkipState[T]) releaseNode(idx int) {
	var zero T
	*s.valuePtrAt(idx) = zero
	for level := 0; level < s.maxLevel; level++ {
		s.setNextAt(idx, level, listSkipNilIndex)
	}
	*s.freePtrAt(idx) = s.free
	s.free = idx
}

func (s *listSkipState[T]) nextNodeLevel() int {
	x := s.rng
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	s.rng = x

	level := 1
	for level < s.maxLevel && (x&1) == 1 {
		level++
		x >>= 1
	}
	return level
}

func (s *listSkipState[T]) grow() {
	newCap := s.cap * 2
	if newCap < 2 {
		newCap = 2
	}

	newValueStore, newValueData, newValueSize := allocListSkipValueStorage[T](newCap)
	newNextStore, newNextData := allocListSkipIntStorage(newCap * s.maxLevel)
	newFreeStore, newFreeData := allocListSkipIntStorage(newCap)

	for i := 0; i < s.cap; i++ {
		dstValue := (*T)(unsafe.Add(newValueData, uintptr(i)*newValueSize))
		srcValue := s.valuePtrAt(i)
		*dstValue = *srcValue

		dstFree := (*int)(unsafe.Add(newFreeData, uintptr(i)*unsafe.Sizeof(0)))
		*dstFree = *s.freePtrAt(i)

		for level := 0; level < s.maxLevel; level++ {
			dstIdx := i*s.maxLevel + level
			dstNext := (*int)(unsafe.Add(newNextData, uintptr(dstIdx)*unsafe.Sizeof(0)))
			*dstNext = s.nextAt(i, level)
		}
	}

	for i := s.cap; i < newCap; i++ {
		next := i + 1
		if i == newCap-1 {
			next = listSkipNilIndex
		}
		*(*int)(unsafe.Add(newFreeData, uintptr(i)*unsafe.Sizeof(0))) = next
		for level := 0; level < s.maxLevel; level++ {
			dstIdx := i*s.maxLevel + level
			*(*int)(unsafe.Add(newNextData, uintptr(dstIdx)*unsafe.Sizeof(0))) = listSkipNilIndex
		}
	}

	oldCap := s.cap
	s.valueStore = newValueStore
	s.valueData = newValueData
	s.valueSize = newValueSize
	s.nextStore = newNextStore
	s.nextData = newNextData
	s.freeStore = newFreeStore
	s.freeData = newFreeData
	s.cap = newCap
	s.free = oldCap
}
