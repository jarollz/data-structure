package treegeneral

import "iter"

// TreeGeneral implements the API interface.
//
// TreeGeneral stores hierarchical n-ary nodes with stable integer IDs and
// preserved holes after subtree removal.
type TreeGeneral[T any] struct{}

// API defines general tree behavior.
type API[T any] interface {
	// AddChild adds new last child of parentID.
	//
	// parentID is existing parent node ID; value is child value.
	// It returns (childID, true) on success, or (-1, false) for invalid parent,
	// removed parent, or empty tree. Child IDs are stable, monotonically
	// increasing, and never reused.
	//
	// Example: childID, ok := tree.AddChild(0, "leaf")
	AddChild(parentID int, value T) (childID int, ok bool)
	// RemoveSubtree removes nodeID and all descendants.
	//
	// nodeID is subtree root to remove. Invalid or already-removed IDs return
	// false. Removing root (0) makes tree empty. Removed IDs stay as holes, and
	// future child IDs keep increasing past all previously allocated IDs.
	//
	// Example: ok := tree.RemoveSubtree(3)
	RemoveSubtree(nodeID int) bool
	// Get returns value at nodeID.
	//
	// nodeID is node to read.
	// It returns (zero, false) for invalid or removed IDs.
	//
	// Example: v, ok := tree.Get(2)
	Get(nodeID int) (T, bool)
	// Parent returns parent ID of nodeID.
	//
	// nodeID is node to inspect.
	// It returns (-1, false) for root, invalid, or removed IDs.
	//
	// Example: p, ok := tree.Parent(2)
	Parent(nodeID int) (int, bool)
	// ChildCount returns number of live direct children.
	//
	// nodeID is parent candidate.
	// It returns -1 for invalid or removed IDs.
	//
	// Example: n := tree.ChildCount(0)
	ChildCount(nodeID int) int
	// Len returns number of live nodes.
	//
	// Example: n := tree.Len()
	Len() int
	// Clone returns independent tree copy with same live IDs, removed-ID holes,
	// child order, parent links, root ID, and next-ID progression.
	//
	// Node values are copied with normal Go assignment.
	//
	// Example: cloned := tree.Clone()
	Clone() *TreeGeneral[T]
	// CloneWith returns independent tree copy using cloneValue for each live node.
	//
	// CloneWith preserves live IDs, removed-ID holes, child order, parent links,
	// root ID, and next-ID progression. cloneValue receives each live value once
	// in pre-order and never sees removed-ID holes. When cloneValue is nil,
	// CloneWith uses normal Go assignment.
	//
	// Example: cloned := tree.CloneWith(func(v int) int { return v * 10 })
	CloneWith(cloneValue func(T) T) *TreeGeneral[T]
	// PreOrder yields values parent-before-children.
	//
	// Sibling order is preserved. Sequence yields each live node once, supports
	// early stop when yield returns false, and yields nothing when tree is empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range tree.PreOrder() { _ = v }
	PreOrder() iter.Seq[T]
}
