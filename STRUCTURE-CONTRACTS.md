# STRUCTURE-CONTRACTS

## Purpose
This file is a compact cross-folder summary.
Use it as quick navigation for API shape, iterator naming, and invariants.
Detailed implementation and test guidance lives in each folder `RULES.md`.

## Global constraints
- Implementation uses only simple types and arrays.
- `slice` is forbidden in implementation code.
- `map` is forbidden in implementation code.
- Mutation during iteration is not safe.

## API and iterator matrix

| Folder | Generic shape | Minimum API summary | Iterator | Iterator order |
|---|---|---|---|---|
| `map-hash` | `Map[K,V]` | `New(capacity,hash,equal)/Put/Get/Delete/Has/Len/Cap/Clear/LoadFactor` | `All() iter.Seq2[K,V]` | Unspecified |
| `map-tree-avl` | `Map[K,V]` | `Put/Get/Delete/Has/Min/Max/Len/Clear` | `All() iter.Seq2[K,V]` | Ascending key |
| `map-tree-red-black` | `Map[K,V]` | `Put/Get/Delete/Has/Min/Max/Len/Clear` | `All() iter.Seq2[K,V]` | Ascending key |
| `heap` | `Heap[T]` | `Push/PopTop/PeekTop/Len/Cap/Clear` | `Values() iter.Seq[T]` | Internal heap order |
| `list-array` | `List[T]` | `Append/Get/Set/Insert/Delete/Len/Cap/Clear` | `Values() iter.Seq[T]` | Index order |
| `list-linked-singly` | `List[T]` | `PushFront/PopFront/Append/DeleteFirst/Len/Clear` | `Values() iter.Seq[T]` | Head to tail |
| `list-linked-doubly` | `List[T]` | `PushFront/PushBack/PopFront/PopBack/Len/Clear` | `Values() iter.Seq[T]` | Head to tail |
| `list-skip` | `SkipList[T]` | `Insert/Delete/Has/Len/Clear` | `Values() iter.Seq[T]` | Sorted |
| `queue` | `Queue[T]` | `Enqueue/Dequeue/PeekFront/Len/Cap/Clear` | `Values() iter.Seq[T]` | Front to back |
| `stack` | `Stack[T]` | `Push/Pop/PeekTop/Len/Cap/Clear` | `Values() iter.Seq[T]` | Top to bottom |
| `tree-general` | `Tree[T]` | `AddChild/RemoveSubtree/Get/Parent/ChildCount/Len` | `PreOrder() iter.Seq[T]` | Pre-order |
| `tree-avl` | `Tree[T]` | `Insert/Delete/Has/Min/Max/Len/Clear` | `InOrder() iter.Seq[T]` | Sorted |
| `tree-red-black` | `Tree[T]` | `Insert/Delete/Has/Min/Max/Len/Clear` | `InOrder() iter.Seq[T]` | Sorted |

Notes:
- Capacity-backed constructors normalize non-positive capacity to `16`.
- Set-like structures (`list-skip`, `tree-avl`, `tree-red-black`) reject duplicates.
- Ordered maps overwrite existing keys on `Put` without changing `Len()`.
- `tree-general` root ID is always `0` while tree is non-empty.

## Invariant summary

| Folder | Core invariant focus |
|---|---|
| `map-hash` | unique keys, probe-chain correctness, tombstone-safe lookup |
| `map-tree-avl` | BST ordering + AVL balance factor bounds |
| `map-tree-red-black` | BST ordering + red-black color/black-height rules |
| `heap` | complete-tree shape + heap property |
| `list-array` | contiguous live range and preserved sequence order |
| `list-linked-singly` | acyclic chain and correct traversal length |
| `list-linked-doubly` | prev/next consistency with valid head/tail |
| `list-skip` | sorted level-0 backbone and ordered upper levels |
| `queue` | FIFO order |
| `stack` | LIFO order |
| `tree-general` | one parent per non-root node, no cycles |
| `tree-avl` | BST ordering + AVL balancing |
| `tree-red-black` | BST ordering + red-black balancing |

## Storage model summary

| Folder | Required representation summary |
|---|---|
| `map-hash` | array of slots with explicit empty/occupied/deleted state |
| `map-tree-avl` | index-based node pool with `-1` nil sentinel |
| `map-tree-red-black` | index-based node pool with `-1` nil sentinel |
| `heap` | contiguous array-backed complete tree |
| `list-array` | contiguous backing storage with live range `[0, Len())` |
| `list-linked-singly` | index-based node pool with head/tail/free-list |
| `list-linked-doubly` | index-based node pool with head/tail/free-list |
| `list-skip` | index-based node pool plus head sentinel and forward links |
| `queue` | circular buffer with head and size |
| `stack` | contiguous backing storage with top at `Len()-1` |
| `tree-general` | stable ID arrays with `parent/firstChild/nextSibling/prevSibling` |
| `tree-avl` | index-based node pool with height/balance metadata |
| `tree-red-black` | index-based node pool with parent/color metadata |

## Auto-resize policy summary

| Folder | Policy |
|---|---|
| `list-array` | Auto-resize capacity. Grow at `Len()==Cap()`. Grow factor `2x` (<1024) else `1.5x`. Shrink at `Len()<=Cap()/4` with `minCap` equal to normalized starting capacity. |
| `queue` | Auto-resize circular buffer with same thresholds as `list-array`. On resize, repack logical front->back order. |
| `stack` | Auto-resize capacity with same thresholds as `list-array`. |
| `heap` | Auto-resize backing array with same thresholds as `list-array`; heap-order invariants preserved after copy. |
| `map-hash` | Internal rehash policy. Grow when probe occupancy `(live+tombstones)/Cap() >= 0.70`; shrink when `LoadFactor()<=0.15` above `minCap`; same-cap cleanup rehash when tombstones are high. No public `Grow()/Shrink()`. |
| `list-linked-singly` | No public resize API. Grow node-pool arrays when free-list is empty; reclaim/reuse nodes on delete/clear. |
| `list-linked-doubly` | No public resize API. Grow node-pool arrays when free-list is empty; reclaim/reuse nodes on delete/clear. |
| `list-skip` | No public resize API. Grow node-pool arrays when free-list is empty; reclaim/reuse nodes on delete/clear; reduce current level when upper levels become empty. |
| `map-tree-avl` | No capacity resize API. Node alloc on insert, reclaim/reuse on delete/clear (free-list recommended). |
| `map-tree-red-black` | No capacity resize API. Node alloc on insert, reclaim/reuse on delete/clear (free-list recommended). |
| `tree-general` | No public resize API. Grow node arrays when full. `RemoveSubtree` removes live nodes but does not reuse public IDs. |
| `tree-avl` | No capacity resize API. Node alloc on insert, reclaim/reuse on delete/clear (free-list recommended). |
| `tree-red-black` | No capacity resize API. Node alloc on insert, reclaim/reuse on delete/clear (free-list recommended). |

Notes:
- Use hysteresis for capacity-backed structures; avoid resize thrash around thresholds.
- Mutation during iteration remains not safe for all folders.

## Where to read next
- Repo-wide process constraints: `AGENTS.md`
- Per-structure implementation and test contracts: each folder `RULES.md`
