package treeredblack

import "iter"

// TreeRedBlack implements the API interface.
//
// TreeRedBlack stores unique values ordered by comparator with red-black balancing.
type TreeRedBlack[T any] struct{}

// API defines set-like red-black tree behavior.
type API[T any] interface {
	// Insert adds v when missing and returns true.
	//
	// v is value to insert; duplicate value returns false.
	//
	// Example: ok := tree.Insert(10)
	Insert(v T) bool
	// Delete removes existing v and returns true.
	//
	// v is value to delete; missing value returns false.
	//
	// Example: ok := tree.Delete(10)
	Delete(v T) bool
	// Has reports whether v exists.
	//
	// v is value to test.
	//
	// Example: ok := tree.Has(10)
	Has(v T) bool
	// Min returns minimum value.
	//
	// It returns (zero, false) when tree is empty.
	//
	// Example: v, ok := tree.Min()
	Min() (T, bool)
	// Max returns maximum value.
	//
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
	// Example: tree.Clear()
	Clear()
	// InOrder yields values in ascending sorted order.
	//
	// Sequence yields each live value once, supports early stop when yield returns false, and yields nothing when empty.
	// Mutation during iteration is not safe.
	//
	// Example: for v := range tree.InOrder() { _ = v }
	InOrder() iter.Seq[T]
}
