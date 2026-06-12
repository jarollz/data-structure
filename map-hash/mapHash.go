package maphash

import (
	"reflect"
	"unsafe"
)

const mapHashDefaultMinCap = 16

const (
	mapHashSlotEmpty uint8 = iota
	mapHashSlotLive
	mapHashSlotDeleted
)

type mapHashSlot[K any, V any] struct {
	key   K
	value V
	state uint8
}

type mapHashState[K any, V any] struct {
	store     any
	data      unsafe.Pointer
	slotSize  uintptr
	len       int
	cap       int
	minCap    int
	tombstone int
	hash      func(K) uint64
	equal     func(a, b K) bool
}

// New creates empty hash map using open addressing with linear probing.
//
// capacity <= 0 is normalized to 16; hash and equal must both be non-nil.
// Returned map is empty and starts with Cap() == max(16, capacity).
//
// Example: m := New[string, int](16, hashString, eqString)
func New[K any, V any](capacity int, hash func(K) uint64, equal func(a, b K) bool) *MapHash[K, V] {
	if hash == nil {
		panic("maphash: hash must not be nil")
	}
	if equal == nil {
		panic("maphash: equal must not be nil")
	}
	startCap := normalizeMapHashCap(capacity)
	state := newMapHashState[K, V](startCap, startCap, hash, equal)
	return (*MapHash[K, V])(unsafe.Pointer(state))
}

func normalizeMapHashCap(capacity int) int {
	if capacity <= mapHashDefaultMinCap {
		return mapHashDefaultMinCap
	}
	return capacity
}

func newMapHashState[K any, V any](capacity, minCap int, hash func(K) uint64, equal func(a, b K) bool) *mapHashState[K, V] {
	normalizedCap := normalizeMapHashCap(capacity)
	normalizedMinCap := normalizeMapHashCap(minCap)
	if normalizedCap < normalizedMinCap {
		normalizedCap = normalizedMinCap
	}
	store, data := allocMapHashStorage[K, V](normalizedCap)
	var sample [1]mapHashSlot[K, V]
	return &mapHashState[K, V]{
		store:     store,
		data:      data,
		slotSize:  unsafe.Sizeof(sample[0]),
		len:       0,
		cap:       normalizedCap,
		minCap:    normalizedMinCap,
		tombstone: 0,
		hash:      hash,
		equal:     equal,
	}
}

func mapHashFromAPI[K any, V any](m *MapHash[K, V]) *mapHashState[K, V] {
	return (*mapHashState[K, V])(unsafe.Pointer(m))
}

func allocMapHashStorage[K any, V any](capacity int) (any, unsafe.Pointer) {
	var sample [1]mapHashSlot[K, V]
	slotType := reflect.TypeOf(sample).Elem()
	arrayType := reflect.ArrayOf(capacity, slotType)
	arrayPtr := reflect.New(arrayType)
	return arrayPtr.Interface(), unsafe.Pointer(arrayPtr.Elem().UnsafeAddr())
}

func (s *mapHashState[K, V]) slotAt(i int) *mapHashSlot[K, V] {
	return (*mapHashSlot[K, V])(unsafe.Add(s.data, uintptr(i)*s.slotSize))
}

func (s *mapHashState[K, V]) copySlotsFrom(src *mapHashState[K, V]) {
	reflect.ValueOf(s.store).Elem().Set(reflect.ValueOf(src.store).Elem())
}

func (s *mapHashState[K, V]) startIndex(key K) int {
	return int(s.hash(key) % uint64(s.cap))
}

func (s *mapHashState[K, V]) probeIndex(start, step int) int {
	i := start + step
	if i >= s.cap {
		i -= s.cap
	}
	return i
}

func (s *mapHashState[K, V]) shouldGrowBeforeInsert() bool {
	occupied := s.len + s.tombstone
	return occupied*100 >= s.cap*70
}

func (s *mapHashState[K, V]) ensureGrowBeforeInsert() {
	if !s.shouldGrowBeforeInsert() {
		return
	}
	newCap := s.cap * 2
	s.rehash(newCap)
}

func (s *mapHashState[K, V]) maybeRehashAfterDelete() {
	if s.cap > s.minCap && s.len*100 <= s.cap*15 {
		newCap := s.cap / 2
		if newCap < s.minCap {
			newCap = s.minCap
		}
		if newCap < s.cap {
			s.rehash(newCap)
		}
		return
	}
	if s.tombstone > s.len/2 {
		s.rehash(s.cap)
	}
}

func (s *mapHashState[K, V]) rehash(newCap int) {
	if newCap < s.minCap {
		newCap = s.minCap
	}
	store, data := allocMapHashStorage[K, V](newCap)
	oldCap := s.cap
	oldData := s.data
	s.store = store
	s.data = data
	s.cap = newCap
	s.len = 0
	s.tombstone = 0
	for i := 0; i < oldCap; i++ {
		slot := (*mapHashSlot[K, V])(unsafe.Add(oldData, uintptr(i)*s.slotSize))
		if slot.state != mapHashSlotLive {
			continue
		}
		s.putNoResize(slot.key, slot.value)
	}
}

func (s *mapHashState[K, V]) putNoResize(key K, value V) {
	cap := s.cap
	idx := int(s.hash(key) % uint64(cap))
	firstDeleted := -1
	for step := 0; step < cap; step++ {
		slot := s.slotAt(idx)
		switch slot.state {
		case mapHashSlotEmpty:
			target := slot
			if firstDeleted >= 0 {
				target = s.slotAt(firstDeleted)
				s.tombstone--
			}
			target.key = key
			target.value = value
			target.state = mapHashSlotLive
			s.len++
			return
		case mapHashSlotDeleted:
			if firstDeleted < 0 {
				firstDeleted = idx
			}
		case mapHashSlotLive:
			if s.equal(slot.key, key) {
				slot.value = value
				return
			}
		}
		idx++
		if idx == cap {
			idx = 0
		}
	}
	if firstDeleted >= 0 {
		target := s.slotAt(firstDeleted)
		target.key = key
		target.value = value
		target.state = mapHashSlotLive
		s.tombstone--
		s.len++
		return
	}
	s.rehash(s.cap * 2)
	s.putNoResize(key, value)
}
