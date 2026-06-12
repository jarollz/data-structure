package maptreeredblack

import (
	"iter"
	"unsafe"
)

// Compile-time check: *MapTreeRedBlack[K, V] satisfies API[K, V].
var _ API[int, string] = (*MapTreeRedBlack[int, string])(nil)

// Put implements the API interface.
// Put inserts key-value pair or overwrites an existing key.
// key is map key and value is stored payload. Overwrite keeps Len() unchanged.
// Successful Put preserves comparator ordering and red-black invariants.
// Example: m.Put(3, "v")
func (s *MapTreeRedBlack[K, V]) Put(key K, value V) {
	st := mapTreeRedBlackFromAPI[K, V](s)
	if st.root == mapTreeRedBlackNilIndex {
		idx := st.allocNode(key, value)
		st.setColor(idx, mapTreeRedBlackColorBlack)
		st.root = idx
		st.len = 1
		return
	}

	parent := mapTreeRedBlackNilIndex
	cur := st.root
	for cur != mapTreeRedBlackNilIndex {
		parent = cur
		cmp := st.cmp(key, st.key(cur))
		if cmp < 0 {
			cur = st.left(cur)
			continue
		}
		if cmp > 0 {
			cur = st.right(cur)
			continue
		}
		st.setValue(cur, value)
		return
	}

	idx := st.allocNode(key, value)
	st.setParent(idx, parent)
	if st.cmp(key, st.key(parent)) < 0 {
		st.setLeft(parent, idx)
	} else {
		st.setRight(parent, idx)
	}
	st.insertFixup(idx)
	st.len++
}

// Get implements the API interface.
// Get returns value for key.
// key is lookup key. Get does not mutate map state.
// It returns (zero, false) when key is missing.
// Example: v, ok := m.Get(3)
func (s *MapTreeRedBlack[K, V]) Get(key K) (V, bool) {
	st := mapTreeRedBlackFromAPI[K, V](s)
	idx := st.findIndex(key)
	if idx == mapTreeRedBlackNilIndex {
		var zero V
		return zero, false
	}
	return st.value(idx), true
}

// Delete implements the API interface.
// Delete removes key when present.
// key is key to remove. Successful Delete preserves comparator ordering and
// red-black invariants, including root and black-sibling fix-up cases.
// It returns false when key is missing.
// Example: ok := m.Delete(3)
func (s *MapTreeRedBlack[K, V]) Delete(key K) bool {
	st := mapTreeRedBlackFromAPI[K, V](s)
	z := st.findIndex(key)
	if z == mapTreeRedBlackNilIndex {
		return false
	}

	y := z
	yColor := st.colorOf(y)
	x := mapTreeRedBlackNilIndex
	xParent := mapTreeRedBlackNilIndex

	if st.left(z) == mapTreeRedBlackNilIndex {
		x = st.right(z)
		xParent = st.parent(z)
		st.transplant(z, st.right(z))
	} else if st.right(z) == mapTreeRedBlackNilIndex {
		x = st.left(z)
		xParent = st.parent(z)
		st.transplant(z, st.left(z))
	} else {
		y = st.minIndex(st.right(z))
		yColor = st.colorOf(y)
		x = st.right(y)
		if st.parent(y) == z {
			xParent = y
			if x != mapTreeRedBlackNilIndex {
				st.setParent(x, y)
			}
		} else {
			xParent = st.parent(y)
			st.transplant(y, st.right(y))
			st.setRight(y, st.right(z))
			st.setParent(st.right(y), y)
		}
		st.transplant(z, y)
		st.setLeft(y, st.left(z))
		st.setParent(st.left(y), y)
		st.setColor(y, st.colorOf(z))
	}

	st.freeNode(z)
	st.len--
	if yColor == mapTreeRedBlackColorBlack {
		st.deleteFixup(x, xParent)
	}
	if st.root != mapTreeRedBlackNilIndex {
		st.setColor(st.root, mapTreeRedBlackColorBlack)
	}
	return true
}

// Has implements the API interface.
// Has reports whether key exists.
// key is lookup key. Has does not mutate map state.
// Example: ok := m.Has(3)
func (s *MapTreeRedBlack[K, V]) Has(key K) bool {
	st := mapTreeRedBlackFromAPI[K, V](s)
	return st.findIndex(key) != mapTreeRedBlackNilIndex
}

// Min implements the API interface.
// Min returns smallest key and associated value.
// Min returns smallest key by comparator order and its value without mutating
// map state. It returns (zeroK, zeroV, false) when map is empty.
// Example: k, v, ok := m.Min()
func (s *MapTreeRedBlack[K, V]) Min() (K, V, bool) {
	st := mapTreeRedBlackFromAPI[K, V](s)
	if st.root != mapTreeRedBlackNilIndex {
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
func (s *MapTreeRedBlack[K, V]) Max() (K, V, bool) {
	st := mapTreeRedBlackFromAPI[K, V](s)
	if st.root != mapTreeRedBlackNilIndex {
		idx := st.maxIndex(st.root)
		return st.key(idx), st.value(idx), true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Len implements the API interface.
// Len returns number of live key-value pairs.
// Example: n := m.Len()
func (s *MapTreeRedBlack[K, V]) Len() int {
	return mapTreeRedBlackFromAPI[K, V](s).len
}

// Clear implements the API interface.
// Clear removes all entries and resets map state.
// Clear is safe on an already-empty map, resets root and length state, and
// leaves comparator unchanged for future operations.
// Example: m.Clear()
func (s *MapTreeRedBlack[K, V]) Clear() {
	mapTreeRedBlackFromAPI[K, V](s).clearAll()
}

// Clone implements the API interface.
// Clone returns independent map copy with same length, comparator, lookup
// results, ascending key order, and red-black validity.
// Keys and values are copied with normal Go assignment.
// Example: cloned := m.Clone()
func (s *MapTreeRedBlack[K, V]) Clone() *MapTreeRedBlack[K, V] {
	state := mapTreeRedBlackFromAPI[K, V](s)
	cloneState := &mapTreeRedBlackState[K, V]{
		cmp:      state.cmp,
		root:     mapTreeRedBlackNilIndex,
		freeHead: mapTreeRedBlackNilIndex,
	}
	cloneState.root = state.cloneStructureInto(cloneState, state.root, mapTreeRedBlackNilIndex)
	cloneState.len = state.len
	if cloneState.root != mapTreeRedBlackNilIndex {
		cloneState.setColor(cloneState.root, mapTreeRedBlackColorBlack)
	}
	return (*MapTreeRedBlack[K, V])(unsafe.Pointer(cloneState))
}

// CloneWith implements the API interface.
// CloneWith returns independent map copy using cloneKey and cloneValue for each
// live entry.
// CloneWith preserves length, comparator, ascending key order, red-black
// validity, and lookup results under transformed keys. cloneKey and
// cloneValue receive each live key-value pair once in ascending key order.
// Cloned keys must remain comparator-compatible. When either hook is nil,
// CloneWith uses normal Go assignment for that payload type.
// Example: cloned := m.CloneWith(func(k int) int { return k }, func(v string) string { return v + "!" })
func (s *MapTreeRedBlack[K, V]) CloneWith(cloneKey func(K) K, cloneValue func(V) V) *MapTreeRedBlack[K, V] {
	clone := s.Clone()
	cloneState := mapTreeRedBlackFromAPI[K, V](clone)
	if cloneState.root == mapTreeRedBlackNilIndex {
		return clone
	}
	if cloneKey == nil && cloneValue == nil {
		return clone
	}
	cloneState.applyHookInOrder(cloneState.root, cloneKey, cloneValue)
	return clone
}

// All implements the API interface.
// All yields each key-value pair once in ascending key order.
// Sequence supports early stop when yield returns false and yields nothing when empty.
// Mutation during iteration is not safe.
// Example: for k, v := range m.All() { _, _ = k, v }
func (s *MapTreeRedBlack[K, V]) All() iter.Seq2[K, V] {
	st := mapTreeRedBlackFromAPI[K, V](s)
	return func(yield func(K, V) bool) {
		if st.root == mapTreeRedBlackNilIndex {
			return
		}
		st.walkInOrder(st.root, yield)
	}
}
