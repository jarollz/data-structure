# STRUCTURE-OVERVIEW

## Purpose
This file is quick cross-folder overview for humans.
Use it to compare module shape at glance and jump into folder docs.
Authoritative contracts live in each folder `SPECS.md` and `api_contract.go`.

## Repo-wide notes
- Implementation uses simple Go types and arrays only.
- `slice` is forbidden in implementation code.
- `map` is forbidden in implementation code.
- Mutation during iteration is not safe.
- Every structure exposes `Clone()` plus `CloneWith(...)`.
- `Clone()` copies payloads with normal Go assignment.
- `CloneWith(...)` applies caller-provided clone hook for each live payload; nil hook means normal Go assignment.

## Module matrix

| Folder | Shape | Constructor | Iterator | Order | Note |
|---|---|---|---|---|---|
| `list-array` | `List[T]` | `New(capacity int)` | `Values() iter.Seq[T]` | Index order | Dynamic array style list. |
| `list-linked-singly` | `List[T]` | `New()` | `Values() iter.Seq[T]` | Head to tail | Singly linked list with tracked tail. |
| `list-linked-doubly` | `List[T]` | `New()` | `Values() iter.Seq[T]` | Head to tail | Doubly linked list with O(1) ends. |
| `list-skip` | `SkipList[T]` | `New(maxLevel, cmp)` | `Values() iter.Seq[T]` | Sorted | Ordered set-like skip list. |
| `queue` | `Queue[T]` | `New(capacity int)` | `Values() iter.Seq[T]` | Front to back | FIFO circular buffer. |
| `stack` | `Stack[T]` | `New(capacity int)` | `Values() iter.Seq[T]` | Top to bottom | LIFO stack. |
| `heap` | `Heap[T]` | `New(capacity, cmp)` | `Values() iter.Seq[T]` | Internal heap order | Priority queue base structure. |
| `tree-general` | `Tree[T]` | `New(rootValue T)` | `PreOrder() iter.Seq[T]` | Pre-order | N-ary tree with stable IDs. |
| `tree-avl` | `Tree[T]` | `New(cmp)` | `InOrder() iter.Seq[T]` | Sorted | AVL set-like tree. |
| `tree-red-black` | `Tree[T]` | `New(cmp)` | `InOrder() iter.Seq[T]` | Sorted | Red-black set-like tree. |
| `map-hash` | `Map[K,V]` | `New(capacity, hash, equal)` | `All() iter.Seq2[K, V]` | Unspecified | Open-addressed hash map. |
| `map-tree-avl` | `Map[K,V]` | `New(cmp)` | `All() iter.Seq2[K, V]` | Ascending key | AVL ordered map. |
| `map-tree-red-black` | `Map[K,V]` | `New(cmp)` | `All() iter.Seq2[K, V]` | Ascending key | Red-black ordered map. |

## Read next
- Repo-wide process constraints: `AGENTS.md`
- Human overview and usage notes: each folder `README.md`
- Authoritative behavior contract: each folder `SPECS.md`
- Public Go signatures and doc comments: each folder `api_contract.go`
