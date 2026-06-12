package treeredblack

import "iter"

// Compile-time check: *TreeRedBlack[T] satisfies API[T].
var _ API[int] = (*TreeRedBlack[int])(nil)

// Insert implements the API interface.
// Insert adds v when it is not already present.
// v is value to insert. Insert returns false on duplicate value and preserves
// binary-search ordering and red-black invariants.
// Example: ok := tree.Insert(10)
func (s *TreeRedBlack[T]) Insert(v T) bool {
	st := ensureState(s, nil)
	if st == nil || st.cmp == nil {
		return false
	}
	parent := nilIndex
	cur := st.root
	for cur != nilIndex {
		parent = cur
		cmp := st.cmp(v, st.value(cur))
		if cmp < 0 {
			cur = st.left(cur)
			continue
		}
		if cmp > 0 {
			cur = st.right(cur)
			continue
		}
		return false
	}
	idx := st.allocNode(v)
	st.setParent(idx, parent)
	if parent == nilIndex {
		st.root = idx
	} else if st.cmp(v, st.value(parent)) < 0 {
		st.setLeft(parent, idx)
	} else {
		st.setRight(parent, idx)
	}
	st.len++
	st.insertFixup(idx)
	return true
}

// Delete implements the API interface.
// Delete removes existing v and returns true.
// v is value to delete. Missing value returns false. Successful delete keeps
// binary-search ordering and red-black invariants.
// Example: ok := tree.Delete(10)
func (s *TreeRedBlack[T]) Delete(v T) bool {
	st := ensureState(s, nil)
	if st == nil || st.cmp == nil {
		return false
	}
	z := st.findIndex(v)
	if z == nilIndex {
		return false
	}

	y := z
	yOriginalColor := st.colorOf(y)
	x := nilIndex
	xParent := nilIndex

	if st.left(z) == nilIndex {
		x = st.right(z)
		xParent = st.parent(z)
		st.transplant(z, x)
	} else if st.right(z) == nilIndex {
		x = st.left(z)
		xParent = st.parent(z)
		st.transplant(z, x)
	} else {
		y = st.minimum(st.right(z))
		yOriginalColor = st.colorOf(y)
		x = st.right(y)
		yParent := st.parent(y)
		if yParent == z {
			xParent = y
		} else {
			xParent = yParent
			st.transplant(y, x)
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
	if yOriginalColor == colorBlack {
		st.deleteFixup(x, xParent)
	}
	return true
}

// Has implements the API interface.
// Has reports whether v exists.
// v is lookup value. Has does not mutate tree state.
// Example: ok := tree.Has(10)
func (s *TreeRedBlack[T]) Has(v T) bool {
	st := ensureState(s, nil)
	if st == nil || st.cmp == nil {
		return false
	}
	return st.findIndex(v) != nilIndex
}

// Min implements the API interface.
// Min returns minimum value.
// Min returns smallest stored value by comparator order. Empty tree returns
// (zero, false).
// Example: v, ok := tree.Min()
func (s *TreeRedBlack[T]) Min() (T, bool) {
	st := ensureState(s, nil)
	if st == nil || st.root == nilIndex {
		var zero T
		return zero, false
	}
	idx := st.minimum(st.root)
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
func (s *TreeRedBlack[T]) Max() (T, bool) {
	st := ensureState(s, nil)
	if st == nil || st.root == nilIndex {
		var zero T
		return zero, false
	}
	idx := st.maximum(st.root)
	if idx == nilIndex {
		var zero T
		return zero, false
	}
	return st.value(idx), true
}

// Len implements the API interface.
// Len returns number of live nodes.
// Example: n := tree.Len()
func (s *TreeRedBlack[T]) Len() int {
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
func (s *TreeRedBlack[T]) Clear() {
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
func (s *TreeRedBlack[T]) Clone() *TreeRedBlack[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
// CloneWith returns independent tree copy using cloneValue for each live value.
// CloneWith calls cloneValue once per live value in ascending in-order order.
// Nil cloneValue uses normal Go assignment.
// Example: cloned := tree.CloneWith(func(v int) int { return v * 10 })
func (s *TreeRedBlack[T]) CloneWith(cloneValue func(T) T) *TreeRedBlack[T] {
	st := ensureState(s, nil)
	if st == nil {
		return New[T](nil)
	}
	cloned := New[T](st.cmp)
	dst := ensureState(cloned, st.cmp)
	if st.root == nilIndex {
		return cloned
	}
	dst.root = st.cloneStructureInto(dst, nilIndex, st.root)
	dst.len = st.len
	if cloneValue != nil {
		dst.applyHookInOrder(dst.root, cloneValue)
	}
	return cloned
}

// RootNode implements the API interface.
// RootNode returns read-only view of root node.
// RootNode returns (zero, false) for empty tree.
// Mutation during node traversal is not safe.
// Example: root, ok := tree.RootNode()
func (s *TreeRedBlack[T]) RootNode() (NodeAPI[T], bool) {
	st := ensureState(s, nil)
	if st == nil || st.root == nilIndex {
		var zero NodeAPI[T]
		return zero, false
	}
	return &treeRedBlackNode[T]{state: st, index: st.root}, true
}

// InOrder implements the API interface.
// InOrder yields values in ascending sorted order.
// InOrder yields each live value once and supports early stop.
// Mutation during iteration is not safe.
// Example: for v := range tree.InOrder() { _ = v }
func (s *TreeRedBlack[T]) InOrder() iter.Seq[T] {
	st := ensureState(s, nil)
	return func(yield func(T) bool) {
		if st == nil || st.root == nilIndex {
			return
		}
		st.walkInOrder(st.root, yield)
	}
}

// Value returns node value.
func (n *treeRedBlackNode[T]) Value() T {
	if n == nil || n.state == nil || n.index == nilIndex {
		var zero T
		return zero
	}
	return n.state.value(n.index)
}

// Color returns node color.
func (n *treeRedBlackNode[T]) Color() Color {
	if n == nil || n.state == nil || n.index == nilIndex {
		return ColorBlack
	}
	if n.state.colorOf(n.index) == colorRed {
		return ColorRed
	}
	return ColorBlack
}

// ChildCount returns number of direct child nodes.
func (n *treeRedBlackNode[T]) ChildCount() int {
	if n == nil || n.state == nil || n.index == nilIndex {
		return 0
	}
	count := 0
	if n.state.left(n.index) != nilIndex {
		count++
	}
	if n.state.right(n.index) != nilIndex {
		count++
	}
	return count
}

// Children yields direct child nodes in left-to-right order.
func (n *treeRedBlackNode[T]) Children() iter.Seq[NodeAPI[T]] {
	return func(yield func(NodeAPI[T]) bool) {
		if n == nil || n.state == nil || n.index == nilIndex {
			return
		}
		left := n.state.left(n.index)
		if left != nilIndex {
			if !yield(&treeRedBlackNode[T]{state: n.state, index: left}) {
				return
			}
		}
		right := n.state.right(n.index)
		if right != nilIndex {
			_ = yield(&treeRedBlackNode[T]{state: n.state, index: right})
		}
	}
}
