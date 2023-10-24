package goavl

import (
	"fmt"
)

func ExampleTree() {
	// Define a new tree explicitly.
	// Note that the comparator is a third generic argument.
	// It allows Go compiler (but not forces it) to inline the comparator.
	var _ *Tree[string, string, func(a, b string) int]
	// if you use New(), the third generic argument can be omitted.
	// in the options we specify `WithCountChildren` allowing `At` operation.
	tree := New[string, string](func(a, b string) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}, WithCountChildren(true))
	// insert some values
	tree.Insert("a", "a")
	tree.Insert("z", "z")
	tree.Insert("m", "m")
	tree.Insert("l", "l")
	tree.Insert("b", "b")
	// print tree, ascending
	fmt.Println("tree, normal order")
	fwdIt := tree.AscendFromStart()
	for {
		e, ok := fwdIt.Next()
		if !ok {
			break
		}
		fmt.Printf("k: %s, v: %s\n", e.Key, *e.Value)
	}
	// print tree, descending
	fmt.Println("tree, reverse order")
	revIt := tree.DescendFromEnd()
	for {
		e, ok := revIt.Prev()
		if !ok {
			break
		}
		fmt.Printf("k: %s, v: %s\n", e.Key, *e.Value)
	}
	v, found := tree.Find("b")
	if found {
		fmt.Printf("the value for 'b' is '%s'\n", *v)
	}
	e := tree.At(2)
	fmt.Printf("the kv at position 2 is %s: %s", e.Key, *e.Value)
	// Output: tree, normal order
	//k: a, v: a
	//k: b, v: b
	//k: l, v: l
	//k: m, v: m
	//k: z, v: z
	//tree, reverse order
	//k: z, v: z
	//k: m, v: m
	//k: l, v: l
	//k: b, v: b
	//k: a, v: a
	//the value for 'b' is 'b'
	//the kv at position 2 is l: l
}

func ExampleNewComparable() {
	// no need to specify a comparator for NewComparable().
	tree := NewComparable[int, int]()
	for _, v := range [...]int{7, 1, 3, 10, 2} {
		tree.Insert(v, v)
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
	//k: 1, v: 1
	//k: 2, v: 2
	//k: 3, v: 3
	//k: 7, v: 7
	//k: 10, v: 10
}
