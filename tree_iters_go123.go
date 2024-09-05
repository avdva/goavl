//go:build go1.23

package goavl

import "iter"

// Mutator allows to modify values and delete keys from a tree in a for-range loop.
// Use mut.E.K and mut.E.V to access keys and modify values.
// Use Delete to delete current key from the tree.
type Mutator[K, V any, Cmp func(a, b K) int] struct {
	E     Entry[K, V]
	t     *Tree[K, V, Cmp]
	it    Iterator[K, V]
	acted bool
}

// Delete deletes the key from the tree.
// The operation is idempotent: second Delete() call is a noop.
func (m *Mutator[K, V, Cmp]) Delete() {
	if !m.acted {
		m.acted = true
		m.it = m.t.DeleteIterator(m.it)
	}
}

// All returns an iterator over the tree's kv pairs.
// It can be used in a for-range loop (Go 1.23+).
func (t *Tree[K, V, Cmp]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		it := t.AscendFromStart()
		for {
			e, ok := it.Next()
			if !ok || !yield(e.Key, *e.Value) {
				break
			}
		}
	}
}

// AllMut returns an iterator over the tree's kv pairs.
// Unlike Mut(), for-range loop will iterate over a sequence of Mutator objects.
// Via a Mutator object one can change a value, or delete current element from the tree.
func (t *Tree[K, V, Cmp]) AllMut() iter.Seq[*Mutator[K, V, Cmp]] {
	return func(yield func(*Mutator[K, V, Cmp]) bool) {
		it := t.AscendFromStart()
		for {
			e, ok := it.Value()
			if !ok {
				break
			}
			m := &Mutator[K, V, Cmp]{
				E:  e,
				t:  t,
				it: it,
			}
			if !yield(m) {
				break
			}
			if !m.acted {
				it.Next()
			} else {
				it = m.it
			}
		}
	}
}
