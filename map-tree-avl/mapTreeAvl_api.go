package maptreeavl

import (
	"iter"
	"unsafe"
)

// Compile-time check: *MapTreeAvl[K, V] satisfies API[K, V].
var _ API[int, string] = (*MapTreeAvl[int, string])(nil)

// Put implements the API interface.
// Put inserts key-value pair or overwrites existing key.
// key is lookup key and value is payload to store. Overwrite keeps Len()
// unchanged. Successful Put preserves key ordering and AVL invariants.
// Example: m.Put(3, "v")
func (s *MapTreeAvl[K, V]) Put(key K, value V) {
	st := mapTreeAvlFromAPI[K, V](s)
	root, inserted := st.insertAt(st.root, key, value)
	st.root = root
	if inserted {
		st.len++
	}
}

// Get implements the API interface.
// Get returns value for key.
// key is lookup key. Get does not mutate map state.
// It returns (zero, false) when key is missing.
// Example: v, ok := m.Get(3)
func (s *MapTreeAvl[K, V]) Get(key K) (V, bool) {
	st := mapTreeAvlFromAPI[K, V](s)
	cur := st.root
	for cur != mapTreeAvlNilIndex {
		cmp := st.cmp(key, st.key(cur))
		if cmp < 0 {
			cur = st.left(cur)
			continue
		}
		if cmp > 0 {
			cur = st.right(cur)
			continue
		}
		return st.value(cur), true
	}
	var zero V
	return zero, false
}

// Delete implements the API interface.
// Delete removes key when present.
// key is key to remove. Successful Delete preserves key ordering and AVL
// invariants, including root deletion with zero, one, or two children.
// It returns false when key is missing.
// Example: ok := m.Delete(3)
func (s *MapTreeAvl[K, V]) Delete(key K) bool {
	st := mapTreeAvlFromAPI[K, V](s)
	root, ok := st.deleteAt(st.root, key)
	if !ok {
		return false
	}
	st.root = root
	st.len--
	return true
}

// Has implements the API interface.
// Has reports whether key exists.
// key is lookup key. Has does not mutate map state.
// Example: ok := m.Has(3)
func (s *MapTreeAvl[K, V]) Has(key K) bool {
	_, ok := s.Get(key)
	return ok
}

// Min implements the API interface.
// Min returns smallest key and associated value.
// Min returns smallest key by comparator order and its value without mutating
// map state. It returns (zeroK, zeroV, false) when map is empty.
// Example: k, v, ok := m.Min()
func (s *MapTreeAvl[K, V]) Min() (K, V, bool) {
	st := mapTreeAvlFromAPI[K, V](s)
	if st.root != mapTreeAvlNilIndex {
		idx := st.minIndex(st.root)
		return st.key(idx), st.value(idx), true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Max implements the API interface.
// Max returns largest key and associated value.
// Max returns largest key by comparator order and its value without mutating
// map state. It returns (zeroK, zeroV, false) when map is empty.
// Example: k, v, ok := m.Max()
func (s *MapTreeAvl[K, V]) Max() (K, V, bool) {
	st := mapTreeAvlFromAPI[K, V](s)
	if st.root != mapTreeAvlNilIndex {
		idx := st.maxIndex(st.root)
		return st.key(idx), st.value(idx), true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Len implements the API interface.
// Len returns number of live entries.
// Example: n := m.Len()
func (s *MapTreeAvl[K, V]) Len() int {
	return mapTreeAvlFromAPI[K, V](s).len
}

// Clear implements the API interface.
// Clear removes all entries and resets map state.
// Clear is safe on an already-empty map, resets root and length state, and
// leaves comparator unchanged for future operations.
// Example: m.Clear()
func (s *MapTreeAvl[K, V]) Clear() {
	state := mapTreeAvlFromAPI[K, V](s)
	state.clearAll()
}

// Clone implements the API interface.
// Clone returns independent map copy with same length, comparator, lookup
// results, ascending key order, and AVL validity.
// Keys and values are copied with normal Go assignment.
// Example: cloned := m.Clone()
func (s *MapTreeAvl[K, V]) Clone() *MapTreeAvl[K, V] {
	state := mapTreeAvlFromAPI[K, V](s)
	cloneState := &mapTreeAvlState[K, V]{
		cmp:      state.cmp,
		root:     mapTreeAvlNilIndex,
		freeHead: mapTreeAvlNilIndex,
	}
	cloneState.root = state.cloneStructureInto(cloneState, state.root)
	cloneState.len = state.len
	return (*MapTreeAvl[K, V])(unsafe.Pointer(cloneState))
}

// CloneWith implements the API interface.
// CloneWith returns independent map copy using cloneKey and cloneValue for
// each live entry.
// CloneWith preserves length, comparator, ascending key order, AVL validity,
// and lookup results under transformed keys. cloneKey and cloneValue receive
// each live key-value pair once in ascending key order. Cloned keys must
// remain comparator-compatible; nil means normal Go assignment for that payload
// type.
// Example: cloned := m.CloneWith(func(k int) int { return k }, func(v string) string { return v + "!" })
func (s *MapTreeAvl[K, V]) CloneWith(cloneKey func(K) K, cloneValue func(V) V) *MapTreeAvl[K, V] {
	clone := s.Clone()
	cloneState := mapTreeAvlFromAPI[K, V](clone)
	if cloneState.root == mapTreeAvlNilIndex {
		return clone
	}
	if cloneKey == nil && cloneValue == nil {
		return clone
	}
	cloneState.applyHookInOrder(cloneState.root, cloneKey, cloneValue)
	return clone
}

// All implements the API interface.
// All yields each key-value pair exactly once in ascending key order.
// Sequence supports early stop and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for k, v := range m.All() { _, _ = k, v }
func (s *MapTreeAvl[K, V]) All() iter.Seq2[K, V] {
	st := mapTreeAvlFromAPI[K, V](s)
	return func(yield func(K, V) bool) {
		if st.root == mapTreeAvlNilIndex {
			return
		}
		st.walkInOrder(st.root, yield)
	}
}
