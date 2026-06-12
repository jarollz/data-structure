package treeredblack

import "unsafe"

const (
	nilIndex      = -1
	colorBlack    = uint8(0)
	colorRed      = uint8(1)
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
	parent   [chunkSize]int
	color    [chunkSize]uint8
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
	tree TreeRedBlack[T]
	pad  byte
}

var states *stateEntry

// New creates empty set-like red-black tree.
//
// cmp defines value ordering for all future tree operations and clones.
// Returned tree is non-nil and has Len() == 0.
//
// Example: tr := New[int](cmpInt)
func New[T any](cmp func(a, b T) int) *TreeRedBlack[T] {
	box := &treeBox[T]{}
	tree := &box.tree
	st := ensureState(tree, cmp)
	st.cmp = cmp
	return tree
}

func ensureState[T any](tree *TreeRedBlack[T], cmp func(a, b T) int) *treeState[T] {
	if tree == nil {
		return nil
	}
	owner := unsafe.Pointer(tree)
	for entry := states; entry != nil; entry = entry.next {
		if entry.owner == owner {
			st, ok := entry.state.(*treeState[T])
			if !ok {
				panic("tree-red-black: mismatched tree state type")
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
		panic("tree-red-black: node index out of bounds")
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
		panic("tree-red-black: capacity exceeded")
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
			chunk.parent[i] = nilIndex
			chunk.color[i] = colorBlack
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
		chunk.parent[off] = nilIndex
		chunk.color[off] = colorRed
		chunk.values[off] = value
		return idx
	}
	idx = st.nextIdx
	st.nextIdx++
	chunk := st.ensureChunk(idx)
	off := chunkOffset(idx)
	chunk.left[off] = nilIndex
	chunk.right[off] = nilIndex
	chunk.parent[off] = nilIndex
	chunk.color[off] = colorRed
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
	chunk.parent[off] = nilIndex
	chunk.color[off] = colorBlack
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

func (st *treeState[T]) parent(index int) int {
	if index == nilIndex {
		return nilIndex
	}
	chunk := st.getChunk(index)
	return chunk.parent[chunkOffset(index)]
}

func (st *treeState[T]) setParent(index, parent int) {
	if index == nilIndex {
		return
	}
	chunk := st.getChunk(index)
	chunk.parent[chunkOffset(index)] = parent
}

func (st *treeState[T]) colorOf(index int) uint8 {
	if index == nilIndex {
		return colorBlack
	}
	chunk := st.getChunk(index)
	return chunk.color[chunkOffset(index)]
}

func (st *treeState[T]) setColor(index int, color uint8) {
	if index == nilIndex {
		return
	}
	chunk := st.getChunk(index)
	chunk.color[chunkOffset(index)] = color
}

func (st *treeState[T]) minimum(index int) int {
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

func (st *treeState[T]) maximum(index int) int {
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

func (st *treeState[T]) findIndex(value T) int {
	cur := st.root
	for cur != nilIndex {
		cmp := st.cmp(value, st.value(cur))
		if cmp < 0 {
			cur = st.left(cur)
			continue
		}
		if cmp > 0 {
			cur = st.right(cur)
			continue
		}
		return cur
	}
	return nilIndex
}

func (st *treeState[T]) rotateLeft(x int) {
	y := st.right(x)
	beta := st.left(y)
	st.setRight(x, beta)
	if beta != nilIndex {
		st.setParent(beta, x)
	}
	xParent := st.parent(x)
	st.setParent(y, xParent)
	if xParent == nilIndex {
		st.root = y
	} else if x == st.left(xParent) {
		st.setLeft(xParent, y)
	} else {
		st.setRight(xParent, y)
	}
	st.setLeft(y, x)
	st.setParent(x, y)
}

func (st *treeState[T]) rotateRight(y int) {
	x := st.left(y)
	beta := st.right(x)
	st.setLeft(y, beta)
	if beta != nilIndex {
		st.setParent(beta, y)
	}
	yParent := st.parent(y)
	st.setParent(x, yParent)
	if yParent == nilIndex {
		st.root = x
	} else if y == st.right(yParent) {
		st.setRight(yParent, x)
	} else {
		st.setLeft(yParent, x)
	}
	st.setRight(x, y)
	st.setParent(y, x)
}

func (st *treeState[T]) insertFixup(z int) {
	for z != st.root && st.colorOf(st.parent(z)) == colorRed {
		p := st.parent(z)
		g := st.parent(p)
		if p == st.left(g) {
			y := st.right(g)
			if st.colorOf(y) == colorRed {
				st.setColor(p, colorBlack)
				st.setColor(y, colorBlack)
				st.setColor(g, colorRed)
				z = g
				continue
			}
			if z == st.right(p) {
				z = p
				st.rotateLeft(z)
				p = st.parent(z)
				g = st.parent(p)
			}
			st.setColor(p, colorBlack)
			st.setColor(g, colorRed)
			st.rotateRight(g)
			continue
		}
		y := st.left(g)
		if st.colorOf(y) == colorRed {
			st.setColor(p, colorBlack)
			st.setColor(y, colorBlack)
			st.setColor(g, colorRed)
			z = g
			continue
		}
		if z == st.left(p) {
			z = p
			st.rotateRight(z)
			p = st.parent(z)
			g = st.parent(p)
		}
		st.setColor(p, colorBlack)
		st.setColor(g, colorRed)
		st.rotateLeft(g)
	}
	st.setColor(st.root, colorBlack)
}

func (st *treeState[T]) transplant(u, v int) {
	uParent := st.parent(u)
	if uParent == nilIndex {
		st.root = v
	} else if u == st.left(uParent) {
		st.setLeft(uParent, v)
	} else {
		st.setRight(uParent, v)
	}
	if v != nilIndex {
		st.setParent(v, uParent)
	}
}

func (st *treeState[T]) deleteFixup(x, parent int) {
	for x != st.root && st.colorOf(x) == colorBlack {
		if parent == nilIndex {
			break
		}
		if x == st.left(parent) {
			w := st.right(parent)
			if st.colorOf(w) == colorRed {
				st.setColor(w, colorBlack)
				st.setColor(parent, colorRed)
				st.rotateLeft(parent)
				w = st.right(parent)
			}
			if st.colorOf(st.left(w)) == colorBlack && st.colorOf(st.right(w)) == colorBlack {
				st.setColor(w, colorRed)
				x = parent
				parent = st.parent(x)
				continue
			}
			if st.colorOf(st.right(w)) == colorBlack {
				st.setColor(st.left(w), colorBlack)
				st.setColor(w, colorRed)
				st.rotateRight(w)
				w = st.right(parent)
			}
			st.setColor(w, st.colorOf(parent))
			st.setColor(parent, colorBlack)
			st.setColor(st.right(w), colorBlack)
			st.rotateLeft(parent)
			x = st.root
			parent = nilIndex
			continue
		}

		w := st.left(parent)
		if st.colorOf(w) == colorRed {
			st.setColor(w, colorBlack)
			st.setColor(parent, colorRed)
			st.rotateRight(parent)
			w = st.left(parent)
		}
		if st.colorOf(st.right(w)) == colorBlack && st.colorOf(st.left(w)) == colorBlack {
			st.setColor(w, colorRed)
			x = parent
			parent = st.parent(x)
			continue
		}
		if st.colorOf(st.left(w)) == colorBlack {
			st.setColor(st.right(w), colorBlack)
			st.setColor(w, colorRed)
			st.rotateLeft(w)
			w = st.left(parent)
		}
		st.setColor(w, st.colorOf(parent))
		st.setColor(parent, colorBlack)
		st.setColor(st.left(w), colorBlack)
		st.rotateRight(parent)
		x = st.root
		parent = nilIndex
	}
	st.setColor(x, colorBlack)
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
		chunk.parent[off] = nilIndex
		chunk.color[off] = colorBlack
		chunk.values[off] = zero
		chunk.nextFree[off] = st.freeHead
		st.freeHead = i
	}
}

func (st *treeState[T]) cloneStructureInto(dst *treeState[T], parent, index int) int {
	if index == nilIndex {
		return nilIndex
	}
	newIdx := dst.allocNode(st.value(index))
	dst.setParent(newIdx, parent)
	dst.setColor(newIdx, st.colorOf(index))
	left := st.cloneStructureInto(dst, newIdx, st.left(index))
	right := st.cloneStructureInto(dst, newIdx, st.right(index))
	dst.setLeft(newIdx, left)
	dst.setRight(newIdx, right)
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
