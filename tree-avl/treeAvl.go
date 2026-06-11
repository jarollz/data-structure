package treeavl

// New creates empty set-like AVL tree.
//
// cmp defines value ordering used for search, insert, and delete operations.
// Returned tree has Len() == 0 and preserves cmp for all future operations and
// clones.
//
// Example: tr := New[int](cmpInt)
func New[T any](cmp func(a, b T) int) *TreeAvl[T] {
	return nil
}
