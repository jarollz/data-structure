package treegeneral

import "iter"

// Compile-time check: *TreeGeneral[T] satisfies API[T].
var _ API[int] = (*TreeGeneral[int])(nil)

// AddChild implements the API interface.
//
// AddChild appends value as the new last child of parentID.
// It returns (-1, false) for invalid parent IDs, removed parent IDs, empty
// trees, or when ID space is exhausted.
// Example: childID, ok := tree.AddChild(0, 9)
func (s *TreeGeneral[T]) AddChild(parentID int, value T) (childID int, ok bool) {
	state := stateOf(s)
	if state == nil || !state.isLive(parentID) || state.nextID >= maxNodeIDs {
		return -1, false
	}

	pid := parentID
	po := state.nodeOffset(pid)
	pb := state.blockFor(pid)
	if pb == nil {
		return -1, false
	}

	id := state.nextID
	state.nextID++
	cb := state.ensureBlock(id)
	co := state.nodeOffset(id)

	cb.live[co] = true
	cb.parent[co] = pid
	cb.firstChild[co] = nilNodeID
	cb.lastChild[co] = nilNodeID
	cb.nextSibling[co] = nilNodeID
	cb.value[co] = value

	last := pb.lastChild[po]
	if last == nilNodeID {
		pb.firstChild[po] = id
		pb.lastChild[po] = id
		cb.prevSibling[co] = nilNodeID
	} else {
		lb := state.blockFor(last)
		lo := state.nodeOffset(last)
		lb.nextSibling[lo] = id
		cb.prevSibling[co] = last
		pb.lastChild[po] = id
	}

	state.length++
	return id, true
}

// RemoveSubtree implements the API interface.
//
// RemoveSubtree removes nodeID and all descendants.
// It returns false for invalid or already-removed IDs.
// Example: ok := tree.RemoveSubtree(3)
func (s *TreeGeneral[T]) RemoveSubtree(nodeID int) bool {
	state := stateOf(s)
	if state == nil || !state.isLive(nodeID) {
		return false
	}

	id := nodeID
	o := state.nodeOffset(id)
	b := state.blockFor(id)

	if id != 0 {
		parentID := b.parent[o]
		pb := state.blockFor(parentID)
		po := state.nodeOffset(parentID)

		prev := b.prevSibling[o]
		next := b.nextSibling[o]

		if pb.firstChild[po] == id {
			pb.firstChild[po] = next
		}
		if pb.lastChild[po] == id {
			pb.lastChild[po] = prev
		}
		if prev != nilNodeID {
			prevBlock := state.blockFor(prev)
			prevOffset := state.nodeOffset(prev)
			prevBlock.nextSibling[prevOffset] = next
		}
		if next != nilNodeID {
			nextBlock := state.blockFor(next)
			nextOffset := state.nodeOffset(next)
			nextBlock.prevSibling[nextOffset] = prev
		}
	}

	removed := state.removeSubtreeCount(id)
	state.length -= removed
	if state.length < 0 {
		state.length = 0
	}
	return true
}

// Get implements the API interface.
//
// Get returns value stored at nodeID.
// It returns (zero, false) for invalid or removed IDs.
// Example: v, ok := tree.Get(2)
func (s *TreeGeneral[T]) Get(nodeID int) (T, bool) {
	state := stateOf(s)
	if state == nil || !state.isLive(nodeID) {
		var zero T
		return zero, false
	}
	b := state.blockFor(nodeID)
	o := state.nodeOffset(nodeID)
	return b.value[o], true
}

// Parent implements the API interface.
//
// Parent returns parent ID of nodeID.
// It returns (-1, false) for root, invalid IDs, or removed IDs.
// Example: p, ok := tree.Parent(2)
func (s *TreeGeneral[T]) Parent(nodeID int) (int, bool) {
	state := stateOf(s)
	if state == nil || !state.isLive(nodeID) {
		return -1, false
	}
	b := state.blockFor(nodeID)
	o := state.nodeOffset(nodeID)
	parentID := b.parent[o]
	if parentID == nilNodeID {
		return -1, false
	}
	return parentID, true
}

// ChildCount implements the API interface.
//
// ChildCount returns number of live direct children of nodeID.
// It returns -1 for invalid or removed IDs.
// Example: n := tree.ChildCount(0)
func (s *TreeGeneral[T]) ChildCount(nodeID int) int {
	state := stateOf(s)
	if state == nil || !state.isLive(nodeID) {
		return -1
	}
	b := state.blockFor(nodeID)
	o := state.nodeOffset(nodeID)
	count := 0
	for child := b.firstChild[o]; child != nilNodeID; {
		count++
		cb := state.blockFor(child)
		co := state.nodeOffset(child)
		child = cb.nextSibling[co]
	}
	return count
}

// Len implements the API interface.
//
// Len returns number of live nodes.
// Example: n := tree.Len()
func (s *TreeGeneral[T]) Len() int {
	state := stateOf(s)
	if state == nil {
		return 0
	}
	return state.length
}

// Clone implements the API interface.
//
// Clone returns independent tree copy with identical IDs, removed-ID holes,
// links, and next-ID progression. Values use normal assignment copy.
// Example: cloned := tree.Clone()
func (s *TreeGeneral[T]) Clone() *TreeGeneral[T] {
	return s.CloneWith(nil)
}

// CloneWith implements the API interface.
//
// CloneWith returns independent tree copy while preserving IDs, removed-ID
// holes, links, and next-ID progression.
// cloneValue is called once per live node in pre-order when non-nil.
// Example: cloned := tree.CloneWith(func(v int) int { return v * 10 })
func (s *TreeGeneral[T]) CloneWith(cloneValue func(T) T) *TreeGeneral[T] {
	state := stateOf(s)
	if state == nil {
		return nil
	}

	handle := &treeHandle[T]{}
	cloneState := &handle.state
	cloneState.nextID = state.nextID
	cloneState.length = state.length

	usedBlocks := state.nextID >> blockShift
	if (state.nextID & blockMask) != 0 {
		usedBlocks++
	}
	for i := 0; i < usedBlocks; i++ {
		src := state.blocks[i]
		if src == nil {
			continue
		}
		dst := &nodeBlock[T]{}
		dst.live = src.live
		dst.parent = src.parent
		dst.firstChild = src.firstChild
		dst.lastChild = src.lastChild
		dst.nextSibling = src.nextSibling
		dst.prevSibling = src.prevSibling
		dst.value = src.value
		cloneState.blocks[i] = dst
	}

	if cloneValue != nil && cloneState.isLive(0) {
		state.preOrderIDs(0, func(id int) bool {
			srcBlock := state.blockFor(id)
			srcOffset := state.nodeOffset(id)
			dstBlock := cloneState.blockFor(id)
			dstOffset := cloneState.nodeOffset(id)
			dstBlock.value[dstOffset] = cloneValue(srcBlock.value[srcOffset])
			return true
		})
	}

	return &handle.api
}

// PreOrder implements the API interface.
//
// PreOrder yields live values in parent-before-children order.
// Example: for v := range tree.PreOrder() { _ = v }
func (s *TreeGeneral[T]) PreOrder() iter.Seq[T] {
	return func(yield func(T) bool) {
		state := stateOf(s)
		if state == nil || !state.isLive(0) {
			return
		}
		state.preOrderIDs(0, func(id int) bool {
			b := state.blockFor(id)
			o := state.nodeOffset(id)
			return yield(b.value[o])
		})
	}
}
