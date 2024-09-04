//go:build go1.23

package goavl

import (
	"fmt"
)

func ExampleGo123Iterators() {
	// no need to specify a comparator for NewComparable().
	tree := NewComparable[int, int]()
	for _, v := range [...]int{7, 1, 3, 10, 2} {
		valuePtr, inserted := tree.Insert(v, v)
		if !inserted || *valuePtr != v {
			panic("invalid insert result")
		}
	}
	fmt.Println("tree, normal order")
	for k, v := range tree.All() {
		fmt.Printf("k: %d, v: %d\n", k, v)
	}
	// Output: tree, normal order
	// k: 1, v: 1
	// k: 2, v: 2
	// k: 3, v: 3
	// k: 7, v: 7
	// k: 10, v: 10
}

func ExampleGo123MutIterators() {
	// no need to specify a comparator for NewComparable().
	tree := NewComparable[int, int]()
	for _, v := range [...]int{7, 1, 3, 10, 2} {
		valuePtr, inserted := tree.Insert(v, v)
		if !inserted || *valuePtr != v {
			panic("invalid insert result")
		}
	}
	fmt.Println("tree, before modifications")
	for m := range tree.AllMut() {
		fmt.Printf("k: %d, v: %d\n", m.E.Key, *m.E.Value)
		*m.E.Value *= 2
		if m.E.Key == 2 {
			m.Delete()
		}
	}
	fmt.Println("tree, after modifications")
	for m := range tree.AllMut() {
		fmt.Printf("k: %d, v: %d\n", m.E.Key, *m.E.Value)
	}
	// Output: tree, before modifications
	// k: 1, v: 1
	// k: 2, v: 2
	// k: 3, v: 3
	// k: 7, v: 7
	// k: 10, v: 10
	// tree, after modifications
	// k: 1, v: 2
	// k: 3, v: 6
	// k: 7, v: 14
	// k: 10, v: 20
}
