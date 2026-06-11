package maphash

import (
	"math/rand"
	"testing"
)

func TestMapHashSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		m := New[int, string](0, hashInt, eqInt)
		if m == nil {
			t.Fatalf("New(0, hash, eq) = nil, want non-nil")
		}
		if m.Cap() < 16 {
			t.Fatalf("Cap() = %d, want >= 16", m.Cap())
		}
		if m.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", m.Len())
		}
	})

	t.Run("put_get_delete_has_clear_loadfactor", func(t *testing.T) {
		m := New[int, string](4, hashInt, eqInt)
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(1, "aa")
		if v, ok := m.Get(1); !ok || v != "aa" {
			t.Fatalf("Get(1) = (%q, %v), want (aa, true)", v, ok)
		}
		if !m.Has(2) {
			t.Fatalf("Has(2) = false, want true")
		}
		if !m.Delete(2) {
			t.Fatalf("Delete(2) = false, want true")
		}
		if m.Delete(9) {
			t.Fatalf("Delete(9) = true, want false")
		}
		if lf := m.LoadFactor(); lf <= 0 {
			t.Fatalf("LoadFactor() = %f, want > 0", lf)
		}
		m.Clear()
		if m.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", m.Len())
		}
	})

	t.Run("collision_and_tombstone_reachability", func(t *testing.T) {
		constantHash := func(int) uint64 { return 1 }
		m := New[int, string](4, constantHash, eqInt)
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")
		if !m.Delete(2) {
			t.Fatalf("Delete(2) = false, want true")
		}
		if v, ok := m.Get(3); !ok || v != "c" {
			t.Fatalf("Get(3) behind tombstone = (%q, %v), want (c, true)", v, ok)
		}
		m.Put(2, "bb")
		if v, ok := m.Get(2); !ok || v != "bb" {
			t.Fatalf("Get(2) after reinsert = (%q, %v), want (bb, true)", v, ok)
		}
	})

	t.Run("randomized_against_builtin_model", func(t *testing.T) {
		m := New[int, int](8, hashInt, eqInt)
		model := make(map[int]int)
		rng := rand.New(rand.NewSource(10))
		for step := 0; step < 500; step++ {
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
			if m.Len() != len(model) {
				t.Fatalf("Len() = %d, want %d", m.Len(), len(model))
			}
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		m := New[int, string](4, hashInt, eqInt)
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")

		got := collectSeq2(m.All())
		if len(got) != 3 {
			t.Fatalf("All count = %d, want 3", len(got))
		}

		count := 0
		for range m.All() {
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
		t.Log("All order is unspecified; mutation during iteration is not safe by contract")
	})
}

func TestMapHashCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type node struct{ value int }

		a := &node{value: 1}
		b := &node{value: 2}

		m := New[int, *node](8, hashInt, eqInt)
		m.Put(1, a)
		m.Put(2, b)

		cloned := m.Clone()
		if cloned == m {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != m.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), m.Len())
		}
		if cloned.Cap() != m.Cap() {
			t.Fatalf("clone Cap() = %d, want %d", cloned.Cap(), m.Cap())
		}
		if cloned.LoadFactor() != m.LoadFactor() {
			t.Fatalf("clone LoadFactor() = %f, want %f", cloned.LoadFactor(), m.LoadFactor())
		}

		if v, ok := cloned.Get(1); !ok || v != a {
			t.Fatalf("clone Get(1) = (%p, %v), want (%p, true)", v, ok, a)
		}
		if v, ok := cloned.Get(2); !ok || v != b {
			t.Fatalf("clone Get(2) = (%p, %v), want (%p, true)", v, ok, b)
		}

		if !m.Delete(1) {
			t.Fatalf("Delete(1) on original failed")
		}
		if !cloned.Has(1) {
			t.Fatalf("clone lost key 1 after original mutation")
		}

		cloned.Put(3, &node{value: 3})
		if m.Has(3) {
			t.Fatalf("original gained key 3 after clone mutation")
		}
	})

	t.Run("clonewith_key_value_nil_and_tombstone_cases", func(t *testing.T) {
		m := New[int, string](8, hashInt, eqInt)
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")
		if !m.Delete(2) {
			t.Fatalf("Delete(2) to create tombstone failed")
		}

		keyCalls := 0
		valueCalls := 0
		both := m.CloneWith(func(k int) int {
			keyCalls++
			return k + 10
		}, func(v string) string {
			valueCalls++
			return v + "!"
		})
		if keyCalls != 2 || valueCalls != 2 {
			t.Fatalf("both-hook call counts = (%d, %d), want (2, 2)", keyCalls, valueCalls)
		}
		if both.Cap() != m.Cap() {
			t.Fatalf("both-hook clone Cap() = %d, want %d", both.Cap(), m.Cap())
		}
		if both.LoadFactor() != m.LoadFactor() {
			t.Fatalf("both-hook clone LoadFactor() = %f, want %f", both.LoadFactor(), m.LoadFactor())
		}
		if v, ok := both.Get(11); !ok || v != "a!" {
			t.Fatalf("both-hook clone Get(11) = (%q, %v), want (a!, true)", v, ok)
		}
		if both.Has(12) {
			t.Fatalf("both-hook clone should not contain tombstoned key 12")
		}
		if v, ok := both.Get(13); !ok || v != "c!" {
			t.Fatalf("both-hook clone Get(13) = (%q, %v), want (c!, true)", v, ok)
		}

		keyOnlyCalls := 0
		keyOnly := m.CloneWith(func(k int) int {
			keyOnlyCalls++
			return k * 10
		}, nil)
		if keyOnlyCalls != 2 {
			t.Fatalf("keyOnly cloneKey calls = %d, want 2", keyOnlyCalls)
		}
		if v, ok := keyOnly.Get(10); !ok || v != "a" {
			t.Fatalf("keyOnly Get(10) = (%q, %v), want (a, true)", v, ok)
		}

		valueOnlyCalls := 0
		valueOnly := m.CloneWith(nil, func(v string) string {
			valueOnlyCalls++
			return v + "?"
		})
		if valueOnlyCalls != 2 {
			t.Fatalf("valueOnly cloneValue calls = %d, want 2", valueOnlyCalls)
		}
		if v, ok := valueOnly.Get(1); !ok || v != "a?" {
			t.Fatalf("valueOnly Get(1) = (%q, %v), want (a?, true)", v, ok)
		}

		nilClone := m.CloneWith(nil, nil)
		if nilClone.Cap() != m.Cap() {
			t.Fatalf("CloneWith(nil, nil) Cap() = %d, want %d", nilClone.Cap(), m.Cap())
		}
		if nilClone.LoadFactor() != m.LoadFactor() {
			t.Fatalf("CloneWith(nil, nil) LoadFactor() = %f, want %f", nilClone.LoadFactor(), m.LoadFactor())
		}
		if v, ok := nilClone.Get(1); !ok || v != "a" {
			t.Fatalf("CloneWith(nil, nil) Get(1) = (%q, %v), want (a, true)", v, ok)
		}
	})

	t.Run("clone_empty_and_no_tombstone_hook_calls", func(t *testing.T) {
		empty := New[int, int](4, hashInt, eqInt)
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
