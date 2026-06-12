package treeavl

import "iter"

// Compile-time check: *TreeAvl[T] satisfies API[T].
var _ API[int] = (*TreeAvl[int])(nil)

// Insert implements the API interface.
// Insert adds v when it is not already present.
// v is value to insert. Insert returns false on duplicate value and preserves
// AVL balance and binary-search ordering invariants.
// Example: ok := tree.Insert(10)
func (s *TreeAvl[T]) Insert(v T) bool {
	st := ensureState(s, nil)
	if st == nil || st.cmp == nil {
		return false
	}
	root, ok := st.insertAt(st.root, v)
	if !ok {
		return false
	}
	st.root = root
	st.len++
	return true
}

// Delete implements the API interface.
// Delete removes existing v and returns true.
// v is value to delete. Missing value returns false. Successful delete keeps
// AVL balance and binary-search ordering invariants.
// Example: ok := tree.Delete(10)
func (s *TreeAvl[T]) Delete(v T) bool {
	st := ensureState(s, nil)
	if st == nil || st.cmp == nil {
		return false
	}
	root, ok := st.deleteAt(st.root, v)
	if !ok {
		return false
	}
	st.root = root
	st.len--
	return true
}

// Has implements the API interface.
// Has reports whether v exists.
// v is lookup value. Has does not mutate tree state.
// Example: ok := tree.Has(10)
func (s *TreeAvl[T]) Has(v T) bool {
	st := ensureState(s, nil)
	if st == nil || st.cmp == nil {
		return false
	}
	cur := st.root
	for cur != nilIndex {
		cmp := st.cmp(v, st.value(cur))
		if cmp < 0 {
			cur = st.left(cur)
			continue
		}
		if cmp > 0 {
			cur = st.right(cur)
			continue
		}
		return true
	}
	return false
}

// Min implements the API interface.
// Min returns minimum value.
// Min returns smallest stored value by comparator order. Empty tree returns
// (zero, false).
// Example: v, ok := tree.Min()
func (s *TreeAvl[T]) Min() (T, bool) {
	st := ensureState(s, nil)
	if st == nil || st.root == nilIndex {
		var zero T
		return zero, false
	}
	idx := st.minIndex(st.root)
	if idx == nilIndex {
		var zero T
		return zero, false
	}
	return st.value(idx), true
}

// Max implements the API interface.
// Max returns maximum value.
// Max returns largest stored value by comparator order. Empty tree returns
// (zero, false).
// Example: v, ok := tree.Max()
func (s *TreeAvl[T]) Max() (T, bool) {
	st := ensureState(s, nil)
	if st == nil || st.root == nilIndex {
		var zero T
		return zero, false
	}
	idx := st.maxIndex(st.root)
	if idx == nilIndex {
		var zero T
		return zero, false
	}
	return st.value(idx), true
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := tree.Len()
func (s *TreeAvl[T]) Len() int {
	st := ensureState(s, nil)
	if st == nil {
		return 0
	}
	return st.len
}

// Clear implements the API interface.
// Clear removes all values and resets tree state.
// Clear is safe on an already-empty tree and keeps comparator for future use.
// Example: tree.Clear()
func (s *TreeAvl[T]) Clear() {
	st := ensureState(s, nil)
	if st == nil {
		return
	}
	st.clearAll()
}

// Clone implements the API interface.
// Clone returns independent tree copy with same contents and comparator.
// Values are copied with normal Go assignment.
// Example: cloned := tree.Clone()
func (s *TreeAvl[T]) Clone() *TreeAvl[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
// CloneWith returns independent tree copy using cloneValue for each live value.
// CloneWith calls cloneValue once per live value in ascending in-order order.
// Nil cloneValue uses normal Go assignment.
// Example: cloned := tree.CloneWith(func(v int) int { return v * 10 })
func (s *TreeAvl[T]) CloneWith(cloneValue func(T) T) *TreeAvl[T] {
	st := ensureState(s, nil)
	if st == nil {
		return New[T](nil)
	}
	cloned := New[T](st.cmp)
	dst := ensureState(cloned, st.cmp)
	if st.root == nilIndex {
		return cloned
	}
	dst.root = st.cloneStructureInto(dst, st.root)
	dst.len = st.len
	if cloneValue != nil {
		dst.applyHookInOrder(dst.root, cloneValue)
	}
	return cloned
}

// InOrder implements the API interface.
// InOrder yields values in ascending sorted order.
// InOrder yields each live value once and supports early stop.
// Mutation during iteration is not safe.
// Example: for v := range tree.InOrder() { _ = v }
func (s *TreeAvl[T]) InOrder() iter.Seq[T] {
	st := ensureState(s, nil)
	return func(yield func(T) bool) {
		if st == nil || st.root == nilIndex {
			return
		}
		st.walkInOrder(st.root, yield)
	}
}
