# goavl
An [AVL tree](https://en.wikipedia.org/wiki/AVL_tree) implementation in Go.

## Badges

![Go build](https://github.com/avdva/goavl/actions/workflows/go.yml/badge.svg)
![Golangci-lint](https://github.com/avdva/goavl/workflows/golangci-lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/avdva/goavl)](https://goreportcard.com/report/github.com/avdva/goavl)

## Installation

Package `goavl` requires Go 1.20+. To start, run:

```sh
$ go get github.com/avdva/goavl
```

## Features

- Support for Go generics.
- Forward and reverse iterators.
- Go 1.23 style iterators support.
- Optional O(logn) access by sorted position with `WithCountChildren(true)`.
- Optional `sync.Pool` and experimental arena allocators.

## API

```go
// Constructors:
// New creates a tree with a user-defined comparator:
//
// Comparators follow the usual convention: return a negative value if a < b,
// zero if a == b, and a positive value if a > b.
//
// func intCmp(a, b int) int {
// 	if a < b {
// 		return -1
// 	}
// 	if a > b {
// 		return 1
// 	}
// 	return 0
// }
// tree := New[int, int](intCmp, WithCountChildren(true))
//
// Options:
// - WithCountChildren(bool) enables O(logn) complexity for the functions that operate
// on element positions.
// - WithSyncPool(*sync.Pool) makes Tree use sync.Pool to allocate tree nodes.
//   If the pool is shared, all trees using it must have the same K and V types.
// - WithArena(*arena.Arena) makes Tree use arenas (currently experimental) to allocate
// tree nodes. This requires GOEXPERIMENT=arenas to be set.
New[K, V any, Cmp func(a, b K) int](cmp Cmp, opts ...Option) *Tree[K, V, Cmp] {}
//  NewComparable works for the keys that satisfy constraints.Ordered.
NewComparable[K constraints.Ordered, V any](opts ...Option) *Tree[K, V, func(a, b K) int] {}

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
// DeleteIterator deletes the element pointed at by it and returns an iterator
// pointing to the next element.
DeleteIterator(it Iterator[K, V, Cmp]) Iterator[K, V, Cmp] {}
// Clear deletes all the elements in O(1) time without returning nodes to the allocator.
Clear() {}

// Iterators:
// AscendFromStart returns an iterator pointing to the smallest element.
AscendFromStart() Iterator[K, V, Cmp] {}
// DescendFromEnd returns an iterator pointing to the largest element.
DescendFromEnd() Iterator[K, V, Cmp] {}
// Ascend returns an iterator pointing to the element that's >= `from`.
Ascend(from K) Iterator[K, V, Cmp] {}
// Descend returns an iterator pointing to the element that's <= `from`.
Descend(from K) Iterator[K, V, Cmp] {}
// AscendAt returns an iterator pointing to the i'th element.
AscendAt(position int) Iterator[K, V, Cmp]
/*
Go 1.23 iterators are also supported:
for k, v := range tree.All() {
	fmt.Printf("k: %d, v: %d\n", k, v)
}
*/
```

## Notes

- `At`, `AscendAt`, and `DeleteAt` are O(logn) only when `WithCountChildren(true)` is enabled. Without it they may scan from the nearest end.
- Tree mutations can invalidate existing iterators. Use the iterator returned by `DeleteIterator` to continue after deleting through an iterator.
- `Clear` is O(1): it drops tree references but does not walk nodes or return them to allocator-specific storage. Delete elements explicitly if you need `sync.Pool` reuse before clearing.
- Arena allocation requires the experimental Go arenas feature. Free the arena only after all trees and values allocated from it are no longer used.

Please see the [examples](/tree_example_test.go), new Go 1.23 [examples](/tree_example_go123_test.go) and arena [examples](/tree_arena_example_test.go) for more details.

## Contact

[Aleksandr Demakin](mailto:alexander.demakin@gmail.com)

## License

Source code is available under the [Apache License Version 2.0](/LICENSE).
