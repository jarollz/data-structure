package listskip

// New creates empty ordered skip list.
//
// maxLevel is configured level bound. New normalizes maxLevel < 1 to 1.
// cmp defines ordering used by Insert, Delete, and Has.
//
// Example: s := New[int](8, cmpInt)
func New[T any](maxLevel int, cmp func(a, b T) int) *ListSkip[T] {
	return nil
}
