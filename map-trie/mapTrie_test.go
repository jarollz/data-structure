package maptrie

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestMapTrieSpec(t *testing.T) {
	t.Run("new_and_empty_key", func(t *testing.T) {
		m := New[int]()
		if m == nil {
			t.Fatalf("New() = nil, want non-nil")
		}
		if m.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", m.Len())
		}
		if m.HasPrefix("") {
			t.Fatalf("HasPrefix(\"\") = true on empty trie, want false")
		}

		m.Put("", 7)
		if v, ok := m.Get(""); !ok || v != 7 {
			t.Fatalf("Get(\"\") = (%d, %v), want (7, true)", v, ok)
		}
		if !m.Has("") {
			t.Fatalf("Has(\"\") = false, want true")
		}
		if !m.HasPrefix("") {
			t.Fatalf("HasPrefix(\"\") = false after insert, want true")
		}
	})

	t.Run("put_get_delete_has_clear", func(t *testing.T) {
		m := New[int]()
		m.Put("app", 1)
		m.Put("apple", 2)
		m.Put("app", 3)
		if m.Len() != 2 {
			t.Fatalf("Len() = %d, want 2", m.Len())
		}
		if v, ok := m.Get("app"); !ok || v != 3 {
			t.Fatalf("Get(app) = (%d, %v), want (3, true)", v, ok)
		}
		if m.Has("ap") {
			t.Fatalf("Has(ap) = true for prefix-only path, want false")
		}
		if !m.Delete("apple") {
			t.Fatalf("Delete(apple) = false, want true")
		}
		if m.Delete("missing") {
			t.Fatalf("Delete(missing) = true, want false")
		}
		m.Clear()
		if m.Len() != 0 {
			t.Fatalf("Len() after Clear = %d, want 0", m.Len())
		}
		if m.HasPrefix("a") {
			t.Fatalf("HasPrefix(a) = true after Clear, want false")
		}
	})

	t.Run("shared_prefix_delete_pruning", func(t *testing.T) {
		m := New[int]()
		m.Put("app", 1)
		m.Put("apple", 2)
		m.Put("apt", 3)

		if !m.Delete("app") {
			t.Fatalf("Delete(app) = false, want true")
		}
		if m.Has("app") {
			t.Fatalf("Has(app) = true after delete, want false")
		}
		if v, ok := m.Get("apple"); !ok || v != 2 {
			t.Fatalf("Get(apple) = (%d, %v), want (2, true)", v, ok)
		}
		if v, ok := m.Get("apt"); !ok || v != 3 {
			t.Fatalf("Get(apt) = (%d, %v), want (3, true)", v, ok)
		}

		if !m.Delete("apple") {
			t.Fatalf("Delete(apple) = false, want true")
		}
		if m.HasPrefix("app") {
			t.Fatalf("HasPrefix(app) = true after deleting last app* key, want false")
		}
		if !m.HasPrefix("ap") {
			t.Fatalf("HasPrefix(ap) = false while apt remains, want true")
		}
		if !m.Delete("apt") {
			t.Fatalf("Delete(apt) = false, want true")
		}
		if m.HasPrefix("ap") {
			t.Fatalf("HasPrefix(ap) = true after deleting all ap* keys, want false")
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		m := New[int]()
		m.Put("", 0)
		m.Put("bat", 4)
		m.Put("app", 1)
		m.Put("apple", 2)
		m.Put("apt", 3)

		allKeys := collectKeys(m.All())
		wantAll := []string{"", "app", "apple", "apt", "bat"}
		if !reflect.DeepEqual(allKeys, wantAll) {
			t.Fatalf("All() keys = %v, want %v", allKeys, wantAll)
		}

		prefixKeys := collectKeys(m.WithPrefix("app"))
		wantPrefix := []string{"app", "apple"}
		if !reflect.DeepEqual(prefixKeys, wantPrefix) {
			t.Fatalf("WithPrefix(app) = %v, want %v", prefixKeys, wantPrefix)
		}

		emptyKeys := collectKeys(m.WithPrefix("zzz"))
		if len(emptyKeys) != 0 {
			t.Fatalf("WithPrefix(zzz) len = %d, want 0", len(emptyKeys))
		}

		count := 0
		for range m.All() {
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("All early-stop count = %d, want 2", count)
		}

		prefixCount := 0
		for range m.WithPrefix("ap") {
			prefixCount++
			if prefixCount == 2 {
				break
			}
		}
		if prefixCount != 2 {
			t.Fatalf("WithPrefix early-stop count = %d, want 2", prefixCount)
		}
		t.Log("Mutation during iteration is not safe by contract")
	})

	t.Run("randomized_against_builtin_model", func(t *testing.T) {
		m := New[int]()
		model := make(map[string]int)
		rng := rand.New(rand.NewSource(10))
		for step := 0; step < 500; step++ {
			key := randomKey(rng.Intn(125))
			switch rng.Intn(5) {
			case 0:
				value := rng.Intn(1000)
				m.Put(key, value)
				model[key] = value
			case 1:
				deleted := m.Delete(key)
				_, exists := model[key]
				if deleted != exists {
					t.Fatalf("Delete(%q) = %v, exists=%v", key, deleted, exists)
				}
				delete(model, key)
			case 2:
				got, ok := m.Get(key)
				want, exists := model[key]
				if ok != exists || (ok && got != want) {
					t.Fatalf("Get(%q) = (%d, %v), want (%d, %v)", key, got, ok, want, exists)
				}
			case 3:
				_, exists := model[key]
				if m.Has(key) != exists {
					t.Fatalf("Has(%q) mismatch", key)
				}
			case 4:
				prefixLen := rng.Intn(len(key) + 1)
				prefix := key[:prefixLen]
				if m.HasPrefix(prefix) != oracleHasPrefix(model, prefix) {
					t.Fatalf("HasPrefix(%q) mismatch", prefix)
				}
			}
			if m.Len() != len(model) {
				t.Fatalf("Len() = %d, want %d", m.Len(), len(model))
			}
		}

		gotKeys := collectKeys(m.All())
		wantKeys := make([]string, 0, len(model))
		for key := range model {
			wantKeys = append(wantKeys, key)
		}
		sort.Strings(wantKeys)
		if !reflect.DeepEqual(gotKeys, wantKeys) {
			t.Fatalf("All() keys = %v, want %v", gotKeys, wantKeys)
		}
	})
}

func TestMapTrieCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		a := &cloneNode{value: 1}
		b := &cloneNode{value: 2}

		m := New[*cloneNode]()
		m.Put("app", a)
		m.Put("bat", b)

		cloned := m.Clone()
		if cloned == m {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != m.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), m.Len())
		}
		if v, ok := cloned.Get("app"); !ok || v != a {
			t.Fatalf("clone Get(app) = (%p, %v), want (%p, true)", v, ok, a)
		}
		if v, ok := cloned.Get("bat"); !ok || v != b {
			t.Fatalf("clone Get(bat) = (%p, %v), want (%p, true)", v, ok, b)
		}

		if !m.Delete("app") {
			t.Fatalf("Delete(app) on original failed")
		}
		if !cloned.Has("app") {
			t.Fatalf("clone lost app after original mutation")
		}

		cloned.Put("cat", &cloneNode{value: 3})
		if m.Has("cat") {
			t.Fatalf("original gained cat after clone mutation")
		}
	})

	t.Run("clonewith_nil_and_custom_hook_order", func(t *testing.T) {
		m := New[int]()
		m.Put("b", 2)
		m.Put("a", 1)
		m.Put("aa", 3)

		calls := make([]int, 0, 3)
		cloned := m.CloneWith(func(v int) int {
			calls = append(calls, v)
			return v * 10
		})
		if !reflect.DeepEqual(calls, []int{1, 3, 2}) {
			t.Fatalf("CloneWith hook order = %v, want [1 3 2]", calls)
		}
		if v, ok := cloned.Get("a"); !ok || v != 10 {
			t.Fatalf("cloned Get(a) = (%d, %v), want (10, true)", v, ok)
		}
		if v, ok := cloned.Get("aa"); !ok || v != 30 {
			t.Fatalf("cloned Get(aa) = (%d, %v), want (30, true)", v, ok)
		}
		if v, ok := cloned.Get("b"); !ok || v != 20 {
			t.Fatalf("cloned Get(b) = (%d, %v), want (20, true)", v, ok)
		}

		nilClone := m.CloneWith(nil)
		if got := collectPairs(nilClone.All()); !reflect.DeepEqual(got, collectPairs(m.All())) {
			t.Fatalf("CloneWith(nil) pairs = %v, want %v", got, collectPairs(m.All()))
		}
	})

	t.Run("clone_empty_has_no_hook_calls", func(t *testing.T) {
		empty := New[int]()
		calls := 0
		cloned := empty.CloneWith(func(v int) int {
			calls++
			return v
		})
		if calls != 0 {
			t.Fatalf("CloneWith hook calls on empty = %d, want 0", calls)
		}
		if cloned.Len() != 0 {
			t.Fatalf("empty clone Len() = %d, want 0", cloned.Len())
		}
	})
}
