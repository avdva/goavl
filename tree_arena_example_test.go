//go:build goexperiment.arenas

package goavl

import (
	"arena"
	"fmt"
)

func ExampleTree_Arenas() {
	a := arena.NewArena()
	defer a.Free()
	tree := NewComparable[int, int](WithArena(a))
	for _, v := range [...]int{7, 1, 3, 10, 2} {
		valuePtr, inserted := tree.Insert(v, v)
		if !inserted || *valuePtr != v {
			panic("invalid insert result")
		}
	}
	fmt.Println("tree, normal order")
	fwdIt := tree.AscendFromStart()
	for {
		e, ok := fwdIt.Value()
		if !ok {
			break
		}
		fmt.Printf("k: %d, v: %d\n", e.Key, *e.Value)
		fwdIt.Next()
	}
	// Output: tree, normal order
	// k: 1, v: 1
	// k: 2, v: 2
	// k: 3, v: 3
	// k: 7, v: 7
	// k: 10, v: 10
}
