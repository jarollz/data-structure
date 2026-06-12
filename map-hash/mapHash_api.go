package maphash

import (
	"iter"
	"unsafe"
)

// Compile-time check: *MapHash[K, V] satisfies API[K, V].
var _ API[int, string] = (*MapHash[int, string])(nil)

// Put implements the API interface.
//
// Put inserts missing key or overwrites existing key.
//
// Example: m.Put("id", 7)
func (s *MapHash[K, V]) Put(key K, value V) {
	state := mapHashFromAPI[K, V](s)
	state.ensureGrowBeforeInsert()
	state.putNoResize(key, value)
}

// Get implements the API interface.
//
// Get returns stored value for key or (zero, false) when key is missing.
//
// Example: v, ok := m.Get("id")
func (s *MapHash[K, V]) Get(key K) (V, bool) {
	state := mapHashFromAPI[K, V](s)
	cap := state.cap
	idx := int(state.hash(key) % uint64(cap))
	for step := 0; step < cap; step++ {
		slot := state.slotAt(idx)
		if slot.state == mapHashSlotEmpty {
			var zero V
			return zero, false
		}
		if slot.state == mapHashSlotLive && state.equal(slot.key, key) {
			return slot.value, true
		}
		idx++
		if idx == cap {
			idx = 0
		}
	}
	var zero V
	return zero, false
}

// Delete implements the API interface.
//
// Delete removes key when present and keeps probe-chain reachability correct.
//
// Example: ok := m.Delete("id")
func (s *MapHash[K, V]) Delete(key K) bool {
	state := mapHashFromAPI[K, V](s)
	cap := state.cap
	idx := int(state.hash(key) % uint64(cap))
	for step := 0; step < cap; step++ {
		slot := state.slotAt(idx)
		if slot.state == mapHashSlotEmpty {
			return false
		}
		if slot.state == mapHashSlotLive && state.equal(slot.key, key) {
			var zeroK K
			var zeroV V
			slot.key = zeroK
			slot.value = zeroV
			slot.state = mapHashSlotDeleted
			state.len--
			state.tombstone++
			state.maybeRehashAfterDelete()
			return true
		}
		idx++
		if idx == cap {
			idx = 0
		}
	}
	return false
}

// Has implements the API interface.
//
// Has reports whether key exists.
//
// Example: ok := m.Has("id")
func (s *MapHash[K, V]) Has(key K) bool {
	state := mapHashFromAPI[K, V](s)
	cap := state.cap
	idx := int(state.hash(key) % uint64(cap))
	for step := 0; step < cap; step++ {
		slot := state.slotAt(idx)
		if slot.state == mapHashSlotEmpty {
			return false
		}
		if slot.state == mapHashSlotLive && state.equal(slot.key, key) {
			return true
		}
		idx++
		if idx == cap {
			idx = 0
		}
	}
	return false
}

// Len implements the API interface.
//
// Len returns number of live entries.
//
// Example: n := m.Len()
func (s *MapHash[K, V]) Len() int {
	return mapHashFromAPI[K, V](s).len
}

// Cap implements the API interface.
//
// Cap returns current table capacity.
//
// Example: c := m.Cap()
func (s *MapHash[K, V]) Cap() int {
	return mapHashFromAPI[K, V](s).cap
}

// Clear implements the API interface.
//
// Clear removes all live entries and tombstones and resets capacity to minCap.
//
// Example: m.Clear()
func (s *MapHash[K, V]) Clear() {
	state := mapHashFromAPI[K, V](s)
	store, data := allocMapHashStorage[K, V](state.minCap)
	state.store = store
	state.data = data
	state.cap = state.minCap
	state.len = 0
	state.tombstone = 0
}

// LoadFactor implements the API interface.
//
// LoadFactor returns float64(Len()) / float64(Cap()).
//
// Example: lf := m.LoadFactor()
func (s *MapHash[K, V]) LoadFactor() float64 {
	state := mapHashFromAPI[K, V](s)
	if state.len == 0 {
		return 0
	}
	return float64(state.len) / float64(state.cap)
}

// Clone implements the API interface.
//
// Clone returns independent copy with assignment-copy keys and values.
//
// Example: cloned := m.Clone()
func (s *MapHash[K, V]) Clone() *MapHash[K, V] {
	state := mapHashFromAPI[K, V](s)
	clone := newMapHashState[K, V](state.cap, state.minCap, state.hash, state.equal)
	clone.copySlotsFrom(state)
	clone.len = state.len
	clone.tombstone = state.tombstone
	return (*MapHash[K, V])(unsafe.Pointer(clone))
}

// CloneWith implements the API interface.
//
// CloneWith returns independent copy and applies hooks once per live entry.
//
// Example: cloned := m.CloneWith(func(k int) int { return k }, func(v string) string { return v + "!" })
func (s *MapHash[K, V]) CloneWith(cloneKey func(K) K, cloneValue func(V) V) *MapHash[K, V] {
	state := mapHashFromAPI[K, V](s)
	if cloneKey == nil && cloneValue == nil {
		return s.Clone()
	}

	clone := newMapHashState[K, V](state.cap, state.minCap, state.hash, state.equal)
	clone.copySlotsFrom(state)
	clone.len = state.len
	clone.tombstone = state.tombstone

	if cloneKey == nil {
		for i := 0; i < state.cap; i++ {
			slot := clone.slotAt(i)
			if slot.state != mapHashSlotLive {
				continue
			}
			slot.value = cloneValue(slot.value)
		}
		return (*MapHash[K, V])(unsafe.Pointer(clone))
	}

	needsRehash := false
	for i := 0; i < state.cap; i++ {
		srcSlot := state.slotAt(i)
		if srcSlot.state != mapHashSlotLive {
			continue
		}
		dstSlot := clone.slotAt(i)
		newKey := cloneKey(srcSlot.key)
		if !state.equal(newKey, srcSlot.key) {
			needsRehash = true
		}
		dstSlot.key = newKey
		if cloneValue != nil {
			dstSlot.value = cloneValue(srcSlot.value)
		}
	}
	if !needsRehash {
		return (*MapHash[K, V])(unsafe.Pointer(clone))
	}

	rebuilt := newMapHashState[K, V](state.cap, state.minCap, state.hash, state.equal)
	for i := 0; i < clone.cap; i++ {
		slot := clone.slotAt(i)
		if slot.state != mapHashSlotLive {
			continue
		}
		rebuilt.putNoResize(slot.key, slot.value)
	}
	return (*MapHash[K, V])(unsafe.Pointer(rebuilt))
}

// All implements the API interface.
//
// All yields each live key-value pair once in unspecified order.
//
// Example: for k, v := range m.All() { _, _ = k, v }
func (s *MapHash[K, V]) All() iter.Seq2[K, V] {
	state := mapHashFromAPI[K, V](s)
	return func(yield func(K, V) bool) {
		for i := 0; i < state.cap; i++ {
			slot := state.slotAt(i)
			if slot.state != mapHashSlotLive {
				continue
			}
			if !yield(slot.key, slot.value) {
				return
			}
		}
	}
}
