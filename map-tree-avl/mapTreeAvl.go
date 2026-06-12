package maptreeavl

import "unsafe"

const (
	mapTreeAvlNilIndex      = -1
	mapTreeAvlChunkShift    = 10
	mapTreeAvlChunkSize     = 1 << mapTreeAvlChunkShift
	mapTreeAvlChunkMask     = mapTreeAvlChunkSize - 1
	mapTreeAvlChunksPerPage = 256
	mapTreeAvlMaxPages      = 1024
)

type mapTreeAvlNodeChunk[K any, V any] struct {
	keys     [mapTreeAvlChunkSize]K
	values   [mapTreeAvlChunkSize]V
	left     [mapTreeAvlChunkSize]int
	right    [mapTreeAvlChunkSize]int
	height   [mapTreeAvlChunkSize]int
	nextFree [mapTreeAvlChunkSize]int
}

type mapTreeAvlChunkPage[K any, V any] struct {
	chunks [mapTreeAvlChunksPerPage]*mapTreeAvlNodeChunk[K, V]
}

type mapTreeAvlState[K any, V any] struct {
	cmp      func(a, b K) int
	root     int
	len      int
	freeHead int
	nextIdx  int
	pages    [mapTreeAvlMaxPages]*mapTreeAvlChunkPage[K, V]
}

// New creates empty ordered map backed by AVL tree.
//
// cmp defines key ordering used by map operations and clones. Returned map is
// non-nil and has Len() == 0.
//
// Example: m := New[int, string](cmpInt)
func New[K any, V any](cmp func(a, b K) int) *MapTreeAvl[K, V] {
	if cmp == nil {
		panic("maptreeavl: cmp must not be nil")
	}
	state := &mapTreeAvlState[K, V]{
		cmp:      cmp,
		root:     mapTreeAvlNilIndex,
		freeHead: mapTreeAvlNilIndex,
	}
	return (*MapTreeAvl[K, V])(unsafe.Pointer(state))
}

func mapTreeAvlFromAPI[K any, V any](m *MapTreeAvl[K, V]) *mapTreeAvlState[K, V] {
	return (*mapTreeAvlState[K, V])(unsafe.Pointer(m))
}

func mapTreeAvlChunkID(index int) int {
	return index >> mapTreeAvlChunkShift
}

func mapTreeAvlChunkOffset(index int) int {
	return index & mapTreeAvlChunkMask
}

func (st *mapTreeAvlState[K, V]) getChunk(index int) *mapTreeAvlNodeChunk[K, V] {
	cid := mapTreeAvlChunkID(index)
	pageIdx := cid / mapTreeAvlChunksPerPage
	if pageIdx < 0 || pageIdx >= mapTreeAvlMaxPages {
		panic("maptreeavl: node index out of bounds")
	}
	page := st.pages[pageIdx]
	if page == nil {
		return nil
	}
	return page.chunks[cid%mapTreeAvlChunksPerPage]
}

func (st *mapTreeAvlState[K, V]) ensureChunk(index int) *mapTreeAvlNodeChunk[K, V] {
	cid := mapTreeAvlChunkID(index)
	pageIdx := cid / mapTreeAvlChunksPerPage
	if pageIdx < 0 || pageIdx >= mapTreeAvlMaxPages {
		panic("maptreeavl: capacity exceeded")
	}
	page := st.pages[pageIdx]
	if page == nil {
		page = &mapTreeAvlChunkPage[K, V]{}
		st.pages[pageIdx] = page
	}
	pos := cid % mapTreeAvlChunksPerPage
	chunk := page.chunks[pos]
	if chunk == nil {
		chunk = &mapTreeAvlNodeChunk[K, V]{}
		for i := 0; i < mapTreeAvlChunkSize; i++ {
			chunk.left[i] = mapTreeAvlNilIndex
			chunk.right[i] = mapTreeAvlNilIndex
			chunk.nextFree[i] = mapTreeAvlNilIndex
		}
		page.chunks[pos] = chunk
	}
	return chunk
}

func (st *mapTreeAvlState[K, V]) allocNode(key K, value V) int {
	idx := st.freeHead
	if idx != mapTreeAvlNilIndex {
		chunk := st.getChunk(idx)
		off := mapTreeAvlChunkOffset(idx)
		st.freeHead = chunk.nextFree[off]
		chunk.nextFree[off] = mapTreeAvlNilIndex
		chunk.left[off] = mapTreeAvlNilIndex
		chunk.right[off] = mapTreeAvlNilIndex
		chunk.height[off] = 1
		chunk.keys[off] = key
		chunk.values[off] = value
		return idx
	}
	idx = st.nextIdx
	st.nextIdx++
	chunk := st.ensureChunk(idx)
	off := mapTreeAvlChunkOffset(idx)
	chunk.left[off] = mapTreeAvlNilIndex
	chunk.right[off] = mapTreeAvlNilIndex
	chunk.height[off] = 1
	chunk.nextFree[off] = mapTreeAvlNilIndex
	chunk.keys[off] = key
	chunk.values[off] = value
	return idx
}

func (st *mapTreeAvlState[K, V]) freeNode(index int) {
	if index == mapTreeAvlNilIndex {
		return
	}
	chunk := st.getChunk(index)
	off := mapTreeAvlChunkOffset(index)
	chunk.left[off] = mapTreeAvlNilIndex
	chunk.right[off] = mapTreeAvlNilIndex
	chunk.height[off] = 0
	var zeroK K
	var zeroV V
	chunk.keys[off] = zeroK
	chunk.values[off] = zeroV
	chunk.nextFree[off] = st.freeHead
	st.freeHead = index
}

func (st *mapTreeAvlState[K, V]) key(index int) K {
	chunk := st.getChunk(index)
	return chunk.keys[mapTreeAvlChunkOffset(index)]
}

func (st *mapTreeAvlState[K, V]) setKey(index int, key K) {
	chunk := st.getChunk(index)
	chunk.keys[mapTreeAvlChunkOffset(index)] = key
}

func (st *mapTreeAvlState[K, V]) value(index int) V {
	chunk := st.getChunk(index)
	return chunk.values[mapTreeAvlChunkOffset(index)]
}

func (st *mapTreeAvlState[K, V]) setValue(index int, value V) {
	chunk := st.getChunk(index)
	chunk.values[mapTreeAvlChunkOffset(index)] = value
}

func (st *mapTreeAvlState[K, V]) left(index int) int {
	if index == mapTreeAvlNilIndex {
		return mapTreeAvlNilIndex
	}
	chunk := st.getChunk(index)
	return chunk.left[mapTreeAvlChunkOffset(index)]
}

func (st *mapTreeAvlState[K, V]) setLeft(index, child int) {
	chunk := st.getChunk(index)
	chunk.left[mapTreeAvlChunkOffset(index)] = child
}

func (st *mapTreeAvlState[K, V]) right(index int) int {
	if index == mapTreeAvlNilIndex {
		return mapTreeAvlNilIndex
	}
	chunk := st.getChunk(index)
	return chunk.right[mapTreeAvlChunkOffset(index)]
}

func (st *mapTreeAvlState[K, V]) setRight(index, child int) {
	chunk := st.getChunk(index)
	chunk.right[mapTreeAvlChunkOffset(index)] = child
}

func (st *mapTreeAvlState[K, V]) heightOf(index int) int {
	if index == mapTreeAvlNilIndex {
		return 0
	}
	chunk := st.getChunk(index)
	return chunk.height[mapTreeAvlChunkOffset(index)]
}

func (st *mapTreeAvlState[K, V]) setHeight(index, height int) {
	chunk := st.getChunk(index)
	chunk.height[mapTreeAvlChunkOffset(index)] = height
}

func mapTreeAvlMax2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (st *mapTreeAvlState[K, V]) updateHeight(index int) {
	h := mapTreeAvlMax2(st.heightOf(st.left(index)), st.heightOf(st.right(index))) + 1
	st.setHeight(index, h)
}

func (st *mapTreeAvlState[K, V]) balance(index int) int {
	return st.heightOf(st.left(index)) - st.heightOf(st.right(index))
}

func (st *mapTreeAvlState[K, V]) rotateRight(y int) int {
	x := st.left(y)
	t2 := st.right(x)
	st.setLeft(y, t2)
	st.setRight(x, y)
	st.updateHeight(y)
	st.updateHeight(x)
	return x
}

func (st *mapTreeAvlState[K, V]) rotateLeft(x int) int {
	y := st.right(x)
	t2 := st.left(y)
	st.setRight(x, t2)
	st.setLeft(y, x)
	st.updateHeight(x)
	st.updateHeight(y)
	return y
}

func (st *mapTreeAvlState[K, V]) rebalance(index int) int {
	if index == mapTreeAvlNilIndex {
		return mapTreeAvlNilIndex
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

func (st *mapTreeAvlState[K, V]) minIndex(index int) int {
	cur := index
	for cur != mapTreeAvlNilIndex {
		next := st.left(cur)
		if next == mapTreeAvlNilIndex {
			break
		}
		cur = next
	}
	return cur
}

func (st *mapTreeAvlState[K, V]) maxIndex(index int) int {
	cur := index
	for cur != mapTreeAvlNilIndex {
		next := st.right(cur)
		if next == mapTreeAvlNilIndex {
			break
		}
		cur = next
	}
	return cur
}

func (st *mapTreeAvlState[K, V]) clearAll() {
	st.root = mapTreeAvlNilIndex
	st.len = 0
	st.freeHead = mapTreeAvlNilIndex
	var zeroK K
	var zeroV V
	for i := st.nextIdx - 1; i >= 0; i-- {
		chunk := st.getChunk(i)
		off := mapTreeAvlChunkOffset(i)
		chunk.left[off] = mapTreeAvlNilIndex
		chunk.right[off] = mapTreeAvlNilIndex
		chunk.height[off] = 0
		chunk.keys[off] = zeroK
		chunk.values[off] = zeroV
		chunk.nextFree[off] = st.freeHead
		st.freeHead = i
	}
}

func (st *mapTreeAvlState[K, V]) cloneStructureInto(dst *mapTreeAvlState[K, V], index int) int {
	if index == mapTreeAvlNilIndex {
		return mapTreeAvlNilIndex
	}
	newIdx := dst.allocNode(st.key(index), st.value(index))
	left := st.cloneStructureInto(dst, st.left(index))
	right := st.cloneStructureInto(dst, st.right(index))
	dst.setLeft(newIdx, left)
	dst.setRight(newIdx, right)
	dst.setHeight(newIdx, st.heightOf(index))
	return newIdx
}

func (st *mapTreeAvlState[K, V]) applyHookInOrder(index int, cloneKey func(K) K, cloneValue func(V) V) {
	if index == mapTreeAvlNilIndex {
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

func (st *mapTreeAvlState[K, V]) walkInOrder(index int, yield func(K, V) bool) bool {
	if index == mapTreeAvlNilIndex {
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

func (st *mapTreeAvlState[K, V]) insertAt(index int, key K, value V) (int, bool) {
	if index == mapTreeAvlNilIndex {
		return st.allocNode(key, value), true
	}
	cmp := st.cmp(key, st.key(index))
	if cmp < 0 {
		left, inserted := st.insertAt(st.left(index), key, value)
		st.setLeft(index, left)
		if !inserted {
			return index, false
		}
		return st.rebalance(index), true
	}
	if cmp > 0 {
		right, inserted := st.insertAt(st.right(index), key, value)
		st.setRight(index, right)
		if !inserted {
			return index, false
		}
		return st.rebalance(index), true
	}
	st.setValue(index, value)
	return index, false
}

func (st *mapTreeAvlState[K, V]) deleteAt(index int, key K) (int, bool) {
	if index == mapTreeAvlNilIndex {
		return mapTreeAvlNilIndex, false
	}
	cmp := st.cmp(key, st.key(index))
	if cmp < 0 {
		left, ok := st.deleteAt(st.left(index), key)
		if !ok {
			return index, false
		}
		st.setLeft(index, left)
		return st.rebalance(index), true
	}
	if cmp > 0 {
		right, ok := st.deleteAt(st.right(index), key)
		if !ok {
			return index, false
		}
		st.setRight(index, right)
		return st.rebalance(index), true
	}

	left := st.left(index)
	right := st.right(index)
	if left == mapTreeAvlNilIndex || right == mapTreeAvlNilIndex {
		next := left
		if next == mapTreeAvlNilIndex {
			next = right
		}
		st.freeNode(index)
		return next, true
	}

	succ := st.minIndex(right)
	succKey := st.key(succ)
	succValue := st.value(succ)
	st.setKey(index, succKey)
	st.setValue(index, succValue)
	newRight, _ := st.deleteAt(right, succKey)
	st.setRight(index, newRight)
	return st.rebalance(index), true
}
