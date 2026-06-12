package maptreeredblack

import "unsafe"

const (
	mapTreeRedBlackNilIndex      = -1
	mapTreeRedBlackColorRed      = true
	mapTreeRedBlackColorBlack    = false
	mapTreeRedBlackChunkShift    = 10
	mapTreeRedBlackChunkSize     = 1 << mapTreeRedBlackChunkShift
	mapTreeRedBlackChunkMask     = mapTreeRedBlackChunkSize - 1
	mapTreeRedBlackChunksPerPage = 256
	mapTreeRedBlackMaxPages      = 1024
)

type mapTreeRedBlackNodeChunk[K any, V any] struct {
	keys     [mapTreeRedBlackChunkSize]K
	values   [mapTreeRedBlackChunkSize]V
	left     [mapTreeRedBlackChunkSize]int
	right    [mapTreeRedBlackChunkSize]int
	parent   [mapTreeRedBlackChunkSize]int
	color    [mapTreeRedBlackChunkSize]bool
	nextFree [mapTreeRedBlackChunkSize]int
}

type mapTreeRedBlackChunkPage[K any, V any] struct {
	chunks [mapTreeRedBlackChunksPerPage]*mapTreeRedBlackNodeChunk[K, V]
}

type mapTreeRedBlackState[K any, V any] struct {
	cmp      func(a, b K) int
	root     int
	len      int
	freeHead int
	nextIdx  int
	pages    [mapTreeRedBlackMaxPages]*mapTreeRedBlackChunkPage[K, V]
}

// New creates an empty ordered map backed by a red-black tree.
//
// cmp defines key ordering and must not be nil. Returned map is non-nil and
// has Len() == 0.
//
// Example: m := New[int, string](cmpInt)
func New[K any, V any](cmp func(a, b K) int) *MapTreeRedBlack[K, V] {
	if cmp == nil {
		panic("maptreeredblack: cmp must not be nil")
	}
	state := &mapTreeRedBlackState[K, V]{
		cmp:      cmp,
		root:     mapTreeRedBlackNilIndex,
		freeHead: mapTreeRedBlackNilIndex,
	}
	return (*MapTreeRedBlack[K, V])(unsafe.Pointer(state))
}

func mapTreeRedBlackFromAPI[K any, V any](m *MapTreeRedBlack[K, V]) *mapTreeRedBlackState[K, V] {
	return (*mapTreeRedBlackState[K, V])(unsafe.Pointer(m))
}

func mapTreeRedBlackChunkID(index int) int {
	return index >> mapTreeRedBlackChunkShift
}

func mapTreeRedBlackChunkOffset(index int) int {
	return index & mapTreeRedBlackChunkMask
}

func (st *mapTreeRedBlackState[K, V]) getChunk(index int) *mapTreeRedBlackNodeChunk[K, V] {
	cid := mapTreeRedBlackChunkID(index)
	pageIdx := cid / mapTreeRedBlackChunksPerPage
	if pageIdx < 0 || pageIdx >= mapTreeRedBlackMaxPages {
		panic("maptreeredblack: node index out of bounds")
	}
	page := st.pages[pageIdx]
	if page == nil {
		return nil
	}
	return page.chunks[cid%mapTreeRedBlackChunksPerPage]
}

func (st *mapTreeRedBlackState[K, V]) ensureChunk(index int) *mapTreeRedBlackNodeChunk[K, V] {
	cid := mapTreeRedBlackChunkID(index)
	pageIdx := cid / mapTreeRedBlackChunksPerPage
	if pageIdx < 0 || pageIdx >= mapTreeRedBlackMaxPages {
		panic("maptreeredblack: capacity exceeded")
	}
	page := st.pages[pageIdx]
	if page == nil {
		page = &mapTreeRedBlackChunkPage[K, V]{}
		st.pages[pageIdx] = page
	}
	pos := cid % mapTreeRedBlackChunksPerPage
	chunk := page.chunks[pos]
	if chunk == nil {
		chunk = &mapTreeRedBlackNodeChunk[K, V]{}
		for i := 0; i < mapTreeRedBlackChunkSize; i++ {
			chunk.left[i] = mapTreeRedBlackNilIndex
			chunk.right[i] = mapTreeRedBlackNilIndex
			chunk.parent[i] = mapTreeRedBlackNilIndex
			chunk.color[i] = mapTreeRedBlackColorBlack
			chunk.nextFree[i] = mapTreeRedBlackNilIndex
		}
		page.chunks[pos] = chunk
	}
	return chunk
}

func (st *mapTreeRedBlackState[K, V]) allocNode(key K, value V) int {
	idx := st.freeHead
	if idx != mapTreeRedBlackNilIndex {
		chunk := st.getChunk(idx)
		off := mapTreeRedBlackChunkOffset(idx)
		st.freeHead = chunk.nextFree[off]
		chunk.nextFree[off] = mapTreeRedBlackNilIndex
		chunk.left[off] = mapTreeRedBlackNilIndex
		chunk.right[off] = mapTreeRedBlackNilIndex
		chunk.parent[off] = mapTreeRedBlackNilIndex
		chunk.color[off] = mapTreeRedBlackColorRed
		chunk.keys[off] = key
		chunk.values[off] = value
		return idx
	}

	idx = st.nextIdx
	st.nextIdx++
	chunk := st.ensureChunk(idx)
	off := mapTreeRedBlackChunkOffset(idx)
	chunk.left[off] = mapTreeRedBlackNilIndex
	chunk.right[off] = mapTreeRedBlackNilIndex
	chunk.parent[off] = mapTreeRedBlackNilIndex
	chunk.color[off] = mapTreeRedBlackColorRed
	chunk.nextFree[off] = mapTreeRedBlackNilIndex
	chunk.keys[off] = key
	chunk.values[off] = value
	return idx
}

func (st *mapTreeRedBlackState[K, V]) freeNode(index int) {
	if index == mapTreeRedBlackNilIndex {
		return
	}
	chunk := st.getChunk(index)
	off := mapTreeRedBlackChunkOffset(index)
	chunk.left[off] = mapTreeRedBlackNilIndex
	chunk.right[off] = mapTreeRedBlackNilIndex
	chunk.parent[off] = mapTreeRedBlackNilIndex
	chunk.color[off] = mapTreeRedBlackColorBlack
	var zeroK K
	var zeroV V
	chunk.keys[off] = zeroK
	chunk.values[off] = zeroV
	chunk.nextFree[off] = st.freeHead
	st.freeHead = index
}

func (st *mapTreeRedBlackState[K, V]) key(index int) K {
	chunk := st.getChunk(index)
	return chunk.keys[mapTreeRedBlackChunkOffset(index)]
}

func (st *mapTreeRedBlackState[K, V]) setKey(index int, key K) {
	chunk := st.getChunk(index)
	chunk.keys[mapTreeRedBlackChunkOffset(index)] = key
}

func (st *mapTreeRedBlackState[K, V]) value(index int) V {
	chunk := st.getChunk(index)
	return chunk.values[mapTreeRedBlackChunkOffset(index)]
}

func (st *mapTreeRedBlackState[K, V]) setValue(index int, value V) {
	chunk := st.getChunk(index)
	chunk.values[mapTreeRedBlackChunkOffset(index)] = value
}

func (st *mapTreeRedBlackState[K, V]) left(index int) int {
	if index == mapTreeRedBlackNilIndex {
		return mapTreeRedBlackNilIndex
	}
	chunk := st.getChunk(index)
	return chunk.left[mapTreeRedBlackChunkOffset(index)]
}

func (st *mapTreeRedBlackState[K, V]) setLeft(index, child int) {
	chunk := st.getChunk(index)
	chunk.left[mapTreeRedBlackChunkOffset(index)] = child
}

func (st *mapTreeRedBlackState[K, V]) right(index int) int {
	if index == mapTreeRedBlackNilIndex {
		return mapTreeRedBlackNilIndex
	}
	chunk := st.getChunk(index)
	return chunk.right[mapTreeRedBlackChunkOffset(index)]
}

func (st *mapTreeRedBlackState[K, V]) setRight(index, child int) {
	chunk := st.getChunk(index)
	chunk.right[mapTreeRedBlackChunkOffset(index)] = child
}

func (st *mapTreeRedBlackState[K, V]) parent(index int) int {
	if index == mapTreeRedBlackNilIndex {
		return mapTreeRedBlackNilIndex
	}
	chunk := st.getChunk(index)
	return chunk.parent[mapTreeRedBlackChunkOffset(index)]
}

func (st *mapTreeRedBlackState[K, V]) setParent(index, parent int) {
	if index == mapTreeRedBlackNilIndex {
		return
	}
	chunk := st.getChunk(index)
	chunk.parent[mapTreeRedBlackChunkOffset(index)] = parent
}

func (st *mapTreeRedBlackState[K, V]) colorOf(index int) bool {
	if index == mapTreeRedBlackNilIndex {
		return mapTreeRedBlackColorBlack
	}
	chunk := st.getChunk(index)
	return chunk.color[mapTreeRedBlackChunkOffset(index)]
}

func (st *mapTreeRedBlackState[K, V]) setColor(index int, color bool) {
	if index == mapTreeRedBlackNilIndex {
		return
	}
	chunk := st.getChunk(index)
	chunk.color[mapTreeRedBlackChunkOffset(index)] = color
}

func (st *mapTreeRedBlackState[K, V]) minIndex(index int) int {
	cur := index
	for cur != mapTreeRedBlackNilIndex {
		next := st.left(cur)
		if next == mapTreeRedBlackNilIndex {
			break
		}
		cur = next
	}
	return cur
}

func (st *mapTreeRedBlackState[K, V]) maxIndex(index int) int {
	cur := index
	for cur != mapTreeRedBlackNilIndex {
		next := st.right(cur)
		if next == mapTreeRedBlackNilIndex {
			break
		}
		cur = next
	}
	return cur
}

func (st *mapTreeRedBlackState[K, V]) rotateLeft(x int) {
	y := st.right(x)
	beta := st.left(y)

	st.setRight(x, beta)
	if beta != mapTreeRedBlackNilIndex {
		st.setParent(beta, x)
	}

	xParent := st.parent(x)
	st.setParent(y, xParent)
	if xParent == mapTreeRedBlackNilIndex {
		st.root = y
	} else if x == st.left(xParent) {
		st.setLeft(xParent, y)
	} else {
		st.setRight(xParent, y)
	}

	st.setLeft(y, x)
	st.setParent(x, y)
}

func (st *mapTreeRedBlackState[K, V]) rotateRight(x int) {
	y := st.left(x)
	beta := st.right(y)

	st.setLeft(x, beta)
	if beta != mapTreeRedBlackNilIndex {
		st.setParent(beta, x)
	}

	xParent := st.parent(x)
	st.setParent(y, xParent)
	if xParent == mapTreeRedBlackNilIndex {
		st.root = y
	} else if x == st.right(xParent) {
		st.setRight(xParent, y)
	} else {
		st.setLeft(xParent, y)
	}

	st.setRight(y, x)
	st.setParent(x, y)
}

func (st *mapTreeRedBlackState[K, V]) insertFixup(z int) {
	for {
		p := st.parent(z)
		if p == mapTreeRedBlackNilIndex || st.colorOf(p) == mapTreeRedBlackColorBlack {
			break
		}
		gp := st.parent(p)
		if p == st.left(gp) {
			y := st.right(gp)
			if st.colorOf(y) == mapTreeRedBlackColorRed {
				st.setColor(p, mapTreeRedBlackColorBlack)
				st.setColor(y, mapTreeRedBlackColorBlack)
				st.setColor(gp, mapTreeRedBlackColorRed)
				z = gp
				continue
			}
			if z == st.right(p) {
				z = p
				st.rotateLeft(z)
				p = st.parent(z)
				gp = st.parent(p)
			}
			st.setColor(p, mapTreeRedBlackColorBlack)
			st.setColor(gp, mapTreeRedBlackColorRed)
			st.rotateRight(gp)
			continue
		}

		y := st.left(gp)
		if st.colorOf(y) == mapTreeRedBlackColorRed {
			st.setColor(p, mapTreeRedBlackColorBlack)
			st.setColor(y, mapTreeRedBlackColorBlack)
			st.setColor(gp, mapTreeRedBlackColorRed)
			z = gp
			continue
		}
		if z == st.left(p) {
			z = p
			st.rotateRight(z)
			p = st.parent(z)
			gp = st.parent(p)
		}
		st.setColor(p, mapTreeRedBlackColorBlack)
		st.setColor(gp, mapTreeRedBlackColorRed)
		st.rotateLeft(gp)
	}
	st.setColor(st.root, mapTreeRedBlackColorBlack)
}

func (st *mapTreeRedBlackState[K, V]) findIndex(key K) int {
	cur := st.root
	for cur != mapTreeRedBlackNilIndex {
		cmp := st.cmp(key, st.key(cur))
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
	return mapTreeRedBlackNilIndex
}

func (st *mapTreeRedBlackState[K, V]) transplant(u, v int) {
	uParent := st.parent(u)
	if uParent == mapTreeRedBlackNilIndex {
		st.root = v
	} else if u == st.left(uParent) {
		st.setLeft(uParent, v)
	} else {
		st.setRight(uParent, v)
	}
	if v != mapTreeRedBlackNilIndex {
		st.setParent(v, uParent)
	}
}

func (st *mapTreeRedBlackState[K, V]) deleteFixup(x, xParent int) {
	for x != st.root && st.colorOf(x) == mapTreeRedBlackColorBlack {
		if xParent == mapTreeRedBlackNilIndex {
			break
		}
		if x == st.left(xParent) {
			w := st.right(xParent)
			if st.colorOf(w) == mapTreeRedBlackColorRed {
				st.setColor(w, mapTreeRedBlackColorBlack)
				st.setColor(xParent, mapTreeRedBlackColorRed)
				st.rotateLeft(xParent)
				w = st.right(xParent)
			}
			if st.colorOf(st.left(w)) == mapTreeRedBlackColorBlack && st.colorOf(st.right(w)) == mapTreeRedBlackColorBlack {
				st.setColor(w, mapTreeRedBlackColorRed)
				x = xParent
				xParent = st.parent(x)
				continue
			}
			if st.colorOf(st.right(w)) == mapTreeRedBlackColorBlack {
				st.setColor(st.left(w), mapTreeRedBlackColorBlack)
				st.setColor(w, mapTreeRedBlackColorRed)
				st.rotateRight(w)
				w = st.right(xParent)
			}
			st.setColor(w, st.colorOf(xParent))
			st.setColor(xParent, mapTreeRedBlackColorBlack)
			st.setColor(st.right(w), mapTreeRedBlackColorBlack)
			st.rotateLeft(xParent)
			x = st.root
			xParent = mapTreeRedBlackNilIndex
			continue
		}

		w := st.left(xParent)
		if st.colorOf(w) == mapTreeRedBlackColorRed {
			st.setColor(w, mapTreeRedBlackColorBlack)
			st.setColor(xParent, mapTreeRedBlackColorRed)
			st.rotateRight(xParent)
			w = st.left(xParent)
		}
		if st.colorOf(st.right(w)) == mapTreeRedBlackColorBlack && st.colorOf(st.left(w)) == mapTreeRedBlackColorBlack {
			st.setColor(w, mapTreeRedBlackColorRed)
			x = xParent
			xParent = st.parent(x)
			continue
		}
		if st.colorOf(st.left(w)) == mapTreeRedBlackColorBlack {
			st.setColor(st.right(w), mapTreeRedBlackColorBlack)
			st.setColor(w, mapTreeRedBlackColorRed)
			st.rotateLeft(w)
			w = st.left(xParent)
		}
		st.setColor(w, st.colorOf(xParent))
		st.setColor(xParent, mapTreeRedBlackColorBlack)
		st.setColor(st.left(w), mapTreeRedBlackColorBlack)
		st.rotateRight(xParent)
		x = st.root
		xParent = mapTreeRedBlackNilIndex
	}
	st.setColor(x, mapTreeRedBlackColorBlack)
}

func (st *mapTreeRedBlackState[K, V]) clearAll() {
	st.root = mapTreeRedBlackNilIndex
	st.len = 0
	st.freeHead = mapTreeRedBlackNilIndex
	var zeroK K
	var zeroV V
	for i := st.nextIdx - 1; i >= 0; i-- {
		chunk := st.getChunk(i)
		off := mapTreeRedBlackChunkOffset(i)
		chunk.left[off] = mapTreeRedBlackNilIndex
		chunk.right[off] = mapTreeRedBlackNilIndex
		chunk.parent[off] = mapTreeRedBlackNilIndex
		chunk.color[off] = mapTreeRedBlackColorBlack
		chunk.keys[off] = zeroK
		chunk.values[off] = zeroV
		chunk.nextFree[off] = st.freeHead
		st.freeHead = i
	}
}

func (st *mapTreeRedBlackState[K, V]) cloneStructureInto(dst *mapTreeRedBlackState[K, V], index int, parent int) int {
	if index == mapTreeRedBlackNilIndex {
		return mapTreeRedBlackNilIndex
	}
	newIdx := dst.allocNode(st.key(index), st.value(index))
	dst.setParent(newIdx, parent)
	dst.setColor(newIdx, st.colorOf(index))
	left := st.cloneStructureInto(dst, st.left(index), newIdx)
	right := st.cloneStructureInto(dst, st.right(index), newIdx)
	dst.setLeft(newIdx, left)
	dst.setRight(newIdx, right)
	return newIdx
}

func (st *mapTreeRedBlackState[K, V]) applyHookInOrder(index int, cloneKey func(K) K, cloneValue func(V) V) {
	if index == mapTreeRedBlackNilIndex {
		return
	}
	st.applyHookInOrder(st.left(index), cloneKey, cloneValue)
	if cloneKey != nil {
		st.setKey(index, cloneKey(st.key(index)))
	}
	if cloneValue != nil {
		st.setValue(index, cloneValue(st.value(index)))
	}
	st.applyHookInOrder(st.right(index), cloneKey, cloneValue)
}

func (st *mapTreeRedBlackState[K, V]) walkInOrder(index int, yield func(K, V) bool) bool {
	if index == mapTreeRedBlackNilIndex {
		return true
	}
	if !st.walkInOrder(st.left(index), yield) {
		return false
	}
	if !yield(st.key(index), st.value(index)) {
		return false
	}
	return st.walkInOrder(st.right(index), yield)
}
