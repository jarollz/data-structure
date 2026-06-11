package treegeneral

// New creates general tree with one root node.
//
// rootValue is stored at root node, root node ID is always 0, Len() starts at
// 1, Get(0) returns rootValue, and next allocated child ID starts at 1.
//
// Example: tr := New[string]("root")
func New[T any](rootValue T) *TreeGeneral[T] {
	return nil
}
