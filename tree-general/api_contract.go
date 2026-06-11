package treegeneral

import "iter"

// TreeGeneral implements the API interface.
//
// TreeGeneral stores hierarchical n-ary nodes with stable integer IDs.
type TreeGeneral[T any] struct{}

// API defines general tree behavior.
type API[T any] interface {
	// AddChild adds new last child of parentID.
	//
	// parentID is existing parent node ID; value is child value.
	// It returns (childID, true) on success, or (-1, false) for invalid parent or empty tree.
	//
	// Example: childID, ok := tree.AddChild(0, "leaf")
	AddChild(parentID int, value T) (childID int, ok bool)
	// RemoveSubtree removes nodeID and all descendants.
	//
	// nodeID is subtree root to remove.
	// It returns false for invalid IDs. Removing root (0) makes tree empty.
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
	// PreOrder yields values parent-before-children.
	//
	// Sequence yields each live node once, supports early stop when yield returns false, and yields nothing when tree is empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range tree.PreOrder() { _ = v }
	PreOrder() iter.Seq[T]
}
