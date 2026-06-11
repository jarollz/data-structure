package treeavl

import "iter"

// TreeAvl implements the API interface.
//
// TreeAvl stores unique comparator-ordered values with AVL balancing.
type TreeAvl[T any] struct{}

// API defines set-like AVL tree behavior.
type API[T any] interface {
	// Insert adds v when missing and returns true.
	//
	// v is value to insert; duplicate value returns false. Successful insert
	// preserves binary-search ordering and AVL balance invariants.
	//
	// Example: ok := tree.Insert(10)
	Insert(v T) bool
	// Delete removes existing v and returns true.
	//
	// v is value to delete; missing value returns false. Successful delete
	// preserves binary-search ordering and AVL balance invariants, including root
	// deletion with zero, one, or two children.
	//
	// Example: ok := tree.Delete(10)
	Delete(v T) bool
	// Has reports whether v exists.
	//
	// v is value to test. Has does not mutate tree state.
	//
	// Example: ok := tree.Has(10)
	Has(v T) bool
	// Min returns minimum value.
	//
	// Min returns smallest value by comparator order without mutating tree state.
	// It returns (zero, false) when tree is empty.
	//
	// Example: v, ok := tree.Min()
	Min() (T, bool)
	// Max returns maximum value.
	//
	// Max returns largest value by comparator order without mutating tree state.
	// It returns (zero, false) when tree is empty.
	//
	// Example: v, ok := tree.Max()
	Max() (T, bool)
	// Len returns number of live nodes.
	//
	// Example: n := tree.Len()
	Len() int
	// Clear removes all values and resets tree state.
	//
	// Clear is safe on an already-empty tree, resets root and length state, and
	// leaves comparator unchanged for future operations.
	//
	// Example: tree.Clear()
	Clear()
	// Clone returns independent tree copy with same length, comparator, lookup
	// results, ascending in-order sequence, and AVL validity.
	//
	// Values are copied with normal Go assignment.
	//
	// Example: cloned := tree.Clone()
	Clone() *TreeAvl[T]
	// CloneWith returns independent tree copy using cloneValue for each live value.
	//
	// CloneWith preserves length, comparator, ascending in-order sequence, and
	// AVL validity. cloneValue receives each live value once in ascending
	// in-order traversal. Cloned values must remain comparator-compatible. When
	// cloneValue is nil, CloneWith uses normal Go assignment.
	//
	// Example: cloned := tree.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *TreeAvl[T]
	// InOrder yields values in ascending sorted order.
	//
	// Sequence yields each live value once, supports early stop when yield returns false, and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range tree.InOrder() { _ = v }
	InOrder() iter.Seq[T]
}
