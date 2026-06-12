package treegeneral

import "unsafe"

const (
	nilNodeID  = -1
	blockShift = 8
	blockSize  = 1 << blockShift
	blockMask  = blockSize - 1
	blockCount = 4096
	maxNodeIDs = blockCount * blockSize
)

type nodeBlock[T any] struct {
	live        [blockSize]bool
	parent      [blockSize]int
	firstChild  [blockSize]int
	lastChild   [blockSize]int
	nextSibling [blockSize]int
	prevSibling [blockSize]int
	value       [blockSize]T
}

type treeState[T any] struct {
	blocks [blockCount]*nodeBlock[T]
	nextID int
	length int
}

type treeHandle[T any] struct {
	api   TreeGeneral[T]
	state treeState[T]
}

func stateOf[T any](tree *TreeGeneral[T]) *treeState[T] {
	if tree == nil {
		return nil
	}
	return &(*treeHandle[T])(unsafe.Pointer(tree)).state
}

func newNodeBlock[T any]() *nodeBlock[T] {
	b := &nodeBlock[T]{}
	for i := 0; i < blockSize; i++ {
		b.parent[i] = nilNodeID
		b.firstChild[i] = nilNodeID
		b.lastChild[i] = nilNodeID
		b.nextSibling[i] = nilNodeID
		b.prevSibling[i] = nilNodeID
	}
	return b
}

func (s *treeState[T]) ensureBlock(id int) *nodeBlock[T] {
	bi := id >> blockShift
	if bi < 0 || bi >= blockCount {
		return nil
	}
	b := s.blocks[bi]
	if b == nil {
		b = newNodeBlock[T]()
		s.blocks[bi] = b
	}
	return b
}

func (s *treeState[T]) blockFor(id int) *nodeBlock[T] {
	bi := id >> blockShift
	if bi < 0 || bi >= blockCount {
		return nil
	}
	return s.blocks[bi]
}

func (s *treeState[T]) nodeOffset(id int) int {
	return id & blockMask
}

func (s *treeState[T]) isLive(id int) bool {
	if s == nil || id < 0 || id >= s.nextID {
		return false
	}
	b := s.blockFor(id)
	if b == nil {
		return false
	}
	return b.live[s.nodeOffset(id)]
}

func (s *treeState[T]) preOrderIDs(nodeID int, visit func(int) bool) bool {
	b := s.blockFor(nodeID)
	if b == nil {
		return true
	}
	o := s.nodeOffset(nodeID)
	if !b.live[o] {
		return true
	}
	if !visit(nodeID) {
		return false
	}
	for child := b.firstChild[o]; child != nilNodeID; {
		cb := s.blockFor(child)
		if cb == nil {
			return false
		}
		co := s.nodeOffset(child)
		next := cb.nextSibling[co]
		if !s.preOrderIDs(child, visit) {
			return false
		}
		child = next
	}
	return true
}

func (s *treeState[T]) removeSubtreeCount(nodeID int) int {
	b := s.blockFor(nodeID)
	o := s.nodeOffset(nodeID)
	removed := 1
	for child := b.firstChild[o]; child != nilNodeID; {
		cb := s.blockFor(child)
		co := s.nodeOffset(child)
		next := cb.nextSibling[co]
		removed += s.removeSubtreeCount(child)
		child = next
	}
	b.live[o] = false
	b.parent[o] = nilNodeID
	b.firstChild[o] = nilNodeID
	b.lastChild[o] = nilNodeID
	b.nextSibling[o] = nilNodeID
	b.prevSibling[o] = nilNodeID
	var zero T
	b.value[o] = zero
	return removed
}

// New creates non-empty general tree with one root node.
//
// rootValue is stored at root ID 0. Returned tree has Len() == 1 and starts
// next allocated ID progression at 1.
//
// Example: tr := New[string]("root")
func New[T any](rootValue T) *TreeGeneral[T] {
	handle := &treeHandle[T]{}
	state := &handle.state
	state.nextID = 1
	state.length = 1
	b := state.ensureBlock(0)
	b.live[0] = true
	b.value[0] = rootValue
	return &handle.api
}
