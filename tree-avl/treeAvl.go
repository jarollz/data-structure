package treeavl

import "unsafe"

const (
	nilIndex      = -1
	chunkShift    = 10
	chunkSize     = 1 << chunkShift
	chunkMask     = chunkSize - 1
	chunksPerPage = 256
	maxPages      = 1024
)

type nodeChunk[T any] struct {
	values   [chunkSize]T
	left     [chunkSize]int
	right    [chunkSize]int
	height   [chunkSize]int
	nextFree [chunkSize]int
}

type chunkPage[T any] struct {
	chunks [chunksPerPage]*nodeChunk[T]
}

type treeState[T any] struct {
	cmp      func(a, b T) int
	root     int
	len      int
	freeHead int
	nextIdx  int
	pages    [maxPages]*chunkPage[T]
}

type stateEntry struct {
	owner unsafe.Pointer
	state any
	next  *stateEntry
}

type treeBox[T any] struct {
	tree TreeAvl[T]
	pad  byte
}

var states *stateEntry

// New creates empty set-like AVL tree.
//
// cmp defines value ordering for all future tree operations and clones.
// Returned tree is non-nil and has Len() == 0.
//
// Example: tr := New[int](cmpInt)
func New[T any](cmp func(a, b T) int) *TreeAvl[T] {
	box := &treeBox[T]{}
	tree := &box.tree
	st := ensureState(tree, cmp)
	st.cmp = cmp
	return tree
}

func ensureState[T any](tree *TreeAvl[T], cmp func(a, b T) int) *treeState[T] {
	if tree == nil {
		return nil
	}
	owner := unsafe.Pointer(tree)
	for entry := states; entry != nil; entry = entry.next {
		if entry.owner == owner {
			st, ok := entry.state.(*treeState[T])
			if !ok {
				panic("tree-avl: mismatched tree state type")
			}
			if st.cmp == nil && cmp != nil {
				st.cmp = cmp
			}
			return st
		}
	}
	st := &treeState[T]{
		cmp:      cmp,
		root:     nilIndex,
		freeHead: nilIndex,
	}
	states = &stateEntry{owner: owner, state: st, next: states}
	return st
}

func chunkID(index int) int {
	return index >> chunkShift
}

func chunkOffset(index int) int {
	return index & chunkMask
}

func (st *treeState[T]) getChunk(index int) *nodeChunk[T] {
	cid := chunkID(index)
	pageIdx := cid / chunksPerPage
	if pageIdx < 0 || pageIdx >= maxPages {
		panic("tree-avl: node index out of bounds")
	}
	page := st.pages[pageIdx]
	if page == nil {
		return nil
	}
	return page.chunks[cid%chunksPerPage]
}

func (st *treeState[T]) ensureChunk(index int) *nodeChunk[T] {
	cid := chunkID(index)
	pageIdx := cid / chunksPerPage
	if pageIdx < 0 || pageIdx >= maxPages {
		panic("tree-avl: capacity exceeded")
	}
	page := st.pages[pageIdx]
	if page == nil {
		page = &chunkPage[T]{}
		st.pages[pageIdx] = page
	}
	pos := cid % chunksPerPage
	chunk := page.chunks[pos]
	if chunk == nil {
		chunk = &nodeChunk[T]{}
		for i := 0; i < chunkSize; i++ {
			chunk.left[i] = nilIndex
			chunk.right[i] = nilIndex
			chunk.nextFree[i] = nilIndex
		}
		page.chunks[pos] = chunk
	}
	return chunk
}

func (st *treeState[T]) allocNode(value T) int {
	idx := st.freeHead
	if idx != nilIndex {
		chunk := st.getChunk(idx)
		off := chunkOffset(idx)
		st.freeHead = chunk.nextFree[off]
		chunk.nextFree[off] = nilIndex
		chunk.left[off] = nilIndex
		chunk.right[off] = nilIndex
		chunk.height[off] = 1
		chunk.values[off] = value
		return idx
	}
	idx = st.nextIdx
	st.nextIdx++
	chunk := st.ensureChunk(idx)
	off := chunkOffset(idx)
	chunk.left[off] = nilIndex
	chunk.right[off] = nilIndex
	chunk.height[off] = 1
	chunk.nextFree[off] = nilIndex
	chunk.values[off] = value
	return idx
}

func (st *treeState[T]) freeNode(index int) {
	if index == nilIndex {
		return
	}
	chunk := st.getChunk(index)
	off := chunkOffset(index)
	chunk.left[off] = nilIndex
	chunk.right[off] = nilIndex
	chunk.height[off] = 0
	var zero T
	chunk.values[off] = zero
	chunk.nextFree[off] = st.freeHead
	st.freeHead = index
}

func (st *treeState[T]) value(index int) T {
	chunk := st.getChunk(index)
	return chunk.values[chunkOffset(index)]
}

func (st *treeState[T]) setValue(index int, value T) {
	chunk := st.getChunk(index)
	chunk.values[chunkOffset(index)] = value
}

func (st *treeState[T]) left(index int) int {
	if index == nilIndex {
		return nilIndex
	}
	chunk := st.getChunk(index)
	return chunk.left[chunkOffset(index)]
}

func (st *treeState[T]) setLeft(index, child int) {
	chunk := st.getChunk(index)
	chunk.left[chunkOffset(index)] = child
}

func (st *treeState[T]) right(index int) int {
	if index == nilIndex {
		return nilIndex
	}
	chunk := st.getChunk(index)
	return chunk.right[chunkOffset(index)]
}

func (st *treeState[T]) setRight(index, child int) {
	chunk := st.getChunk(index)
	chunk.right[chunkOffset(index)] = child
}

func (st *treeState[T]) heightOf(index int) int {
	if index == nilIndex {
		return 0
	}
	chunk := st.getChunk(index)
	return chunk.height[chunkOffset(index)]
}

func (st *treeState[T]) setHeight(index, height int) {
	chunk := st.getChunk(index)
	chunk.height[chunkOffset(index)] = height
}

func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (st *treeState[T]) updateHeight(index int) {
	h := max2(st.heightOf(st.left(index)), st.heightOf(st.right(index))) + 1
	st.setHeight(index, h)
}

func (st *treeState[T]) balance(index int) int {
	return st.heightOf(st.left(index)) - st.heightOf(st.right(index))
}

func (st *treeState[T]) rotateRight(y int) int {
	x := st.left(y)
	t2 := st.right(x)
	st.setLeft(y, t2)
	st.setRight(x, y)
	st.updateHeight(y)
	st.updateHeight(x)
	return x
}

func (st *treeState[T]) rotateLeft(x int) int {
	y := st.right(x)
	t2 := st.left(y)
	st.setRight(x, t2)
	st.setLeft(y, x)
	st.updateHeight(x)
	st.updateHeight(y)
	return y
}

func (st *treeState[T]) rebalance(index int) int {
	if index == nilIndex {
		return nilIndex
	}
	st.updateHeight(index)
	balance := st.balance(index)
	if balance > 1 {
		left := st.left(index)
		if st.balance(left) < 0 {
			st.setLeft(index, st.rotateLeft(left))
		}
		return st.rotateRight(index)
	}
	if balance < -1 {
		right := st.right(index)
		if st.balance(right) > 0 {
			st.setRight(index, st.rotateRight(right))
		}
		return st.rotateLeft(index)
	}
	return index
}

func (st *treeState[T]) minIndex(index int) int {
	cur := index
	for cur != nilIndex {
		next := st.left(cur)
		if next == nilIndex {
			break
		}
		cur = next
	}
	return cur
}

func (st *treeState[T]) maxIndex(index int) int {
	cur := index
	for cur != nilIndex {
		next := st.right(cur)
		if next == nilIndex {
			break
		}
		cur = next
	}
	return cur
}

func (st *treeState[T]) clearAll() {
	st.root = nilIndex
	st.len = 0
	st.freeHead = nilIndex
	var zero T
	for i := st.nextIdx - 1; i >= 0; i-- {
		chunk := st.getChunk(i)
		off := chunkOffset(i)
		chunk.left[off] = nilIndex
		chunk.right[off] = nilIndex
		chunk.height[off] = 0
		chunk.values[off] = zero
		chunk.nextFree[off] = st.freeHead
		st.freeHead = i
	}
}

func (st *treeState[T]) cloneStructureInto(dst *treeState[T], index int) int {
	if index == nilIndex {
		return nilIndex
	}
	newIdx := dst.allocNode(st.value(index))
	left := st.cloneStructureInto(dst, st.left(index))
	right := st.cloneStructureInto(dst, st.right(index))
	dst.setLeft(newIdx, left)
	dst.setRight(newIdx, right)
	dst.setHeight(newIdx, st.heightOf(index))
	return newIdx
}

func (st *treeState[T]) applyHookInOrder(index int, hook func(T) T) {
	if index == nilIndex {
		return
	}
	st.applyHookInOrder(st.left(index), hook)
	st.setValue(index, hook(st.value(index)))
	st.applyHookInOrder(st.right(index), hook)
}

func (st *treeState[T]) walkInOrder(index int, yield func(T) bool) bool {
	if index == nilIndex {
		return true
	}
	if !st.walkInOrder(st.left(index), yield) {
		return false
	}
	if !yield(st.value(index)) {
		return false
	}
	return st.walkInOrder(st.right(index), yield)
}

func (st *treeState[T]) insertAt(index int, value T) (int, bool) {
	if index == nilIndex {
		return st.allocNode(value), true
	}
	cmp := st.cmp(value, st.value(index))
	if cmp < 0 {
		left, ok := st.insertAt(st.left(index), value)
		if !ok {
			return index, false
		}
		st.setLeft(index, left)
		return st.rebalance(index), true
	}
	if cmp > 0 {
		right, ok := st.insertAt(st.right(index), value)
		if !ok {
			return index, false
		}
		st.setRight(index, right)
		return st.rebalance(index), true
	}
	return index, false
}

func (st *treeState[T]) deleteAt(index int, value T) (int, bool) {
	if index == nilIndex {
		return nilIndex, false
	}
	cmp := st.cmp(value, st.value(index))
	if cmp < 0 {
		left, ok := st.deleteAt(st.left(index), value)
		if !ok {
			return index, false
		}
		st.setLeft(index, left)
		return st.rebalance(index), true
	}
	if cmp > 0 {
		right, ok := st.deleteAt(st.right(index), value)
		if !ok {
			return index, false
		}
		st.setRight(index, right)
		return st.rebalance(index), true
	}

	left := st.left(index)
	right := st.right(index)
	if left == nilIndex || right == nilIndex {
		next := left
		if next == nilIndex {
			next = right
		}
		st.freeNode(index)
		return next, true
	}

	succ := st.minIndex(right)
	succValue := st.value(succ)
	st.setValue(index, succValue)
	newRight, _ := st.deleteAt(right, succValue)
	st.setRight(index, newRight)
	return st.rebalance(index), true
}
