# goavl
An [AVL tree](https://en.wikipedia.org/wiki/AVL_tree) implementation in Go.

## Badges

![Build Status](https://github.com/avdva/goavl/workflows/golangci-lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/avdva/goavl)](https://goreportcard.com/report/github.com/avdva/goavl)

## Installation

To start using this package, run:

```sh
$ go get github.com/avdva/goavl
```

## Features

- Support for Go generics (Go 1.18+).
- Forward and reverse iterators.
- Go 1.23 style iterators support.
- Provides an efficient way of getting items by index (if `CountChildren` is on).

## API

```go
// Create a tree:
// New creates a tree with a user-defined comparator:  
// intCmp := func(a, b int) int {
//    if a < b {
//      return -1
//    }
//    if a > b {
//      return 1
//    }
//    return 0
//  }
New[int, int](intCmp, WithCountChildren(true)) {}
//  NewComparable works for the keys that satisfy constraints.Ordered.
NewComparable[int, int]() {}

// Search for elements:
// Find finds a value for given key.
Find(k K) (v *V, found bool) {}
// Min returns the minimum element of the tree.
Min() (entry Entry[K, V], found bool) {}
// Max returns the maximum element of the tree.
Max() (entry Entry[K, V], found bool) {}
// At returns the i'th element of the tree.
At(position int) Entry[K, V] {}
// Len returns the number of elements.
Len() int {}

// Tree modifications:
// Insert inserts a kv pair.
Insert(k K, v V) (v *V, inserted bool) {}
// Delete deletes a key.
Delete(k K) (v V, deleted bool) {}
// DeleteAt deletes i'th element.
DeleteAt(position int) (k K, v V) {}
// DeleteIterator deletes the element pointed at by it.
DeleteIterator(it Iterator[K, V]) {}
// Clear deletes all the elements.
Clear() {}

// Iterators:
// AscendFromStart returns an iterator pointing to the smallest element.
AscendFromStart() Iterator[K, V] {}
// DescendFromEnd returns an iterator pointing to the largest element.
DescendFromEnd() Iterator[K, V] {}
// Ascend returns an iterator pointing to the element that's >= `from`.
Ascend(from K) Iterator[K, V] {}
// Descend returns an iterator pointing to the element that's <= `from`.
Descend(from K) Iterator[K, V] {}
// AscendAt returns an iterator pointing to the i'th element.
AscendAt(position int) Iterator[K, V]
/*
Go 1.23 iterators are also supported:
for k, v := range tree.All() {
	fmt.Printf("k: %d, v: %d\n", k, v)
}
*/
```

Please see the [examples](/tree_example_test.go) and new Go 1.23 [examples](/tree_example_go123_test.go) for more details.

## Contact

[Aleksandr Demakin](mailto:alexander.demakin@gmail.com)

## License

Source code is available under the [Apache License Version 2.0](/LICENSE).