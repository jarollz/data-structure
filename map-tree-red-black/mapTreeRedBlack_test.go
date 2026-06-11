package maptreeredblack

import (
	"math/rand"
	"sort"
	"testing"
)

func TestMapTreeRedBlackSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		m := New[int, string](cmpInt)
		if m == nil {
			t.Fatalf("New(cmp) = nil, want non-nil")
		}
		if m.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", m.Len())
		}
	})

	t.Run("put_get_delete_has_min_max_clear", func(t *testing.T) {
		m := New[int, string](cmpInt)
		if _, _, ok := m.Min(); ok {
			t.Fatalf("Min() ok = true on empty map")
		}
		if _, _, ok := m.Max(); ok {
			t.Fatalf("Max() ok = true on empty map")
		}
		m.Put(2, "b")
		m.Put(1, "a")
		m.Put(3, "c")
		m.Put(2, "bb")
		if v, ok := m.Get(2); !ok || v != "bb" {
			t.Fatalf("Get(2) = (%q, %v), want (bb, true)", v, ok)
		}
		if !m.Has(1) {
			t.Fatalf("Has(1) = false, want true")
		}
		if !m.Delete(1) {
			t.Fatalf("Delete(1) = false, want true")
		}
		if m.Delete(9) {
			t.Fatalf("Delete(9) = true, want false")
		}
		m.Clear()
		if m.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", m.Len())
		}
	})

	t.Run("overwrite_len_and_sorted_order", func(t *testing.T) {
		m := New[int, string](cmpInt)
		m.Put(2, "b")
		m.Put(1, "a")
		m.Put(2, "bb")
		if m.Len() != 2 {
			t.Fatalf("Len() = %d, want 2 after overwrite", m.Len())
		}
		keys := make([]int, 0, 2)
		values := make([]string, 0, 2)
		for k, v := range m.All() {
			keys = append(keys, k)
			values = append(values, v)
		}
		if len(keys) != 2 || keys[0] != 1 || keys[1] != 2 || values[1] != "bb" {
			t.Fatalf("All() = keys=%v values=%v, want sorted keys [1 2] with overwritten value", keys, values)
		}
	})

	t.Run("randomized_against_builtin_model", func(t *testing.T) {
		m := New[int, int](cmpInt)
		model := make(map[int]int)
		rng := rand.New(rand.NewSource(12))
		for step := 0; step < 400; step++ {
			key := rng.Intn(200)
			switch rng.Intn(4) {
			case 0:
				value := rng.Intn(1000)
				m.Put(key, value)
				model[key] = value
			case 1:
				deleted := m.Delete(key)
				_, exists := model[key]
				if deleted != exists {
					t.Fatalf("Delete(%d) returned %v, exists=%v", key, deleted, exists)
				}
				delete(model, key)
			case 2:
				got, ok := m.Get(key)
				want, exists := model[key]
				if ok != exists || (ok && got != want) {
					t.Fatalf("Get(%d) = (%d, %v), want (%d, %v)", key, got, ok, want, exists)
				}
			case 3:
				_, exists := model[key]
				if m.Has(key) != exists {
					t.Fatalf("Has(%d) mismatch", key)
				}
			}
			keys := make([]int, 0, len(model))
			for k := range model {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			gotKeys := make([]int, 0, len(keys))
			gotValues := make([]int, 0, len(keys))
			for k, v := range m.All() {
				gotKeys = append(gotKeys, k)
				gotValues = append(gotValues, v)
			}
			if len(gotKeys) != len(keys) {
				t.Fatalf("All len = %d, want %d", len(gotKeys), len(keys))
			}
			for i := range keys {
				if gotKeys[i] != keys[i] || gotValues[i] != model[keys[i]] {
					t.Fatalf("All()[%d] = (%d,%d), want (%d,%d)", i, gotKeys[i], gotValues[i], keys[i], model[keys[i]])
				}
			}
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		m := New[int, string](cmpInt)
		m.Put(3, "c")
		m.Put(1, "a")
		m.Put(2, "b")

		got := collectSeq2(m.All())
		if len(got) != 3 {
			t.Fatalf("All count = %d, want 3", len(got))
		}

		last := -1
		count := 0
		for k := range m.All() {
			if count > 0 && k < last {
				t.Fatalf("All keys not ascending: last=%d now=%d", last, k)
			}
			last = k
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("early-stop count = %d, want 2", count)
		}

		m.Clear()
		emptyCount := 0
		for range m.All() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty All count = %d, want 0", emptyCount)
		}
		t.Log("mutation during iteration is not safe by contract")
	})
}

func TestMapTreeRedBlackCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type node struct{ value int }

		a := &node{value: 1}
		b := &node{value: 2}
		c := &node{value: 3}

		m := New[int, *node](cmpInt)
		m.Put(2, b)
		m.Put(1, a)
		m.Put(3, c)

		cloned := m.Clone()
		if cloned == m {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != m.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), m.Len())
		}
		if k, v, ok := cloned.Min(); !ok || k != 1 || v != a {
			t.Fatalf("clone Min() = (%d, %p, %v), want (1, %p, true)", k, v, ok, a)
		}
		if k, v, ok := cloned.Max(); !ok || k != 3 || v != c {
			t.Fatalf("clone Max() = (%d, %p, %v), want (3, %p, true)", k, v, ok, c)
		}

		keys := make([]int, 0, 3)
		values := make([]*node, 0, 3)
		for k, v := range cloned.All() {
			keys = append(keys, k)
			values = append(values, v)
		}
		if len(keys) != 3 || keys[0] != 1 || keys[1] != 2 || keys[2] != 3 || values[0] != a || values[1] != b || values[2] != c {
			t.Fatalf("clone All() = keys=%v values=%v, want sorted shared refs", keys, values)
		}

		if !m.Delete(1) {
			t.Fatalf("Delete(1) on original failed")
		}
		if !cloned.Has(1) {
			t.Fatalf("clone lost key 1 after original mutation")
		}

		if !cloned.Delete(3) {
			t.Fatalf("Delete(3) on clone failed")
		}
		if !m.Has(3) {
			t.Fatalf("original lost key 3 after clone mutation")
		}
	})

	t.Run("clonewith_nil_key_only_value_only_and_both_hooks", func(t *testing.T) {
		m := New[int, string](cmpInt)
		m.Put(3, "c")
		m.Put(1, "a")
		m.Put(2, "b")

		nilClone := m.CloneWith(nil, nil)
		nilKeys := make([]int, 0, 3)
		for k := range nilClone.All() {
			nilKeys = append(nilKeys, k)
		}
		if len(nilKeys) != 3 || nilKeys[0] != 1 || nilKeys[1] != 2 || nilKeys[2] != 3 {
			t.Fatalf("CloneWith(nil, nil) keys = %v, want [1 2 3]", nilKeys)
		}

		keyCalls := make([]int, 0, 3)
		valueCalls := 0
		both := m.CloneWith(func(k int) int {
			keyCalls = append(keyCalls, k)
			return k + 10
		}, func(v string) string {
			valueCalls++
			return v + "!"
		})
		if len(keyCalls) != 3 || keyCalls[0] != 1 || keyCalls[1] != 2 || keyCalls[2] != 3 || valueCalls != 3 {
			t.Fatalf("both-hook calls = keys=%v values=%d, want [1 2 3] and 3", keyCalls, valueCalls)
		}

		bothKeys := make([]int, 0, 3)
		bothValues := make([]string, 0, 3)
		for k, v := range both.All() {
			bothKeys = append(bothKeys, k)
			bothValues = append(bothValues, v)
		}
		if len(bothKeys) != 3 || bothKeys[0] != 11 || bothKeys[1] != 12 || bothKeys[2] != 13 || bothValues[0] != "a!" || bothValues[1] != "b!" || bothValues[2] != "c!" {
			t.Fatalf("both-hook clone = keys=%v values=%v, want [11 12 13] and [a! b! c!]", bothKeys, bothValues)
		}

		keyOnly := m.CloneWith(func(k int) int { return k + 100 }, nil)
		if v, ok := keyOnly.Get(101); !ok || v != "a" {
			t.Fatalf("keyOnly Get(101) = (%q, %v), want (a, true)", v, ok)
		}

		valueOnly := m.CloneWith(nil, func(v string) string { return v + "?" })
		if v, ok := valueOnly.Get(1); !ok || v != "a?" {
			t.Fatalf("valueOnly Get(1) = (%q, %v), want (a?, true)", v, ok)
		}
	})

	t.Run("clone_empty_hook_calls", func(t *testing.T) {
		empty := New[int, int](cmpInt)
		keyCalls := 0
		valueCalls := 0
		empty.CloneWith(func(k int) int {
			keyCalls++
			return k
		}, func(v int) int {
			valueCalls++
			return v
		})
		if keyCalls != 0 || valueCalls != 0 {
			t.Fatalf("empty hook calls = (%d, %d), want (0, 0)", keyCalls, valueCalls)
		}
	})
}
