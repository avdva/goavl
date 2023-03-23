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
- Provides an efficient way of getting items by index (if `CountChildren` is on).

## API

```go
create:
New[int, int](intCmp, WithCountChildren(true)) // creates a new int --> int tree.
NewComparable[int, int]() // Works for the keys that satisfy constraints.Ordered. 

search:
Find(k K) (v V, found bool) // finds a value for given key.
Min() (k K, v V, found bool) // returns the minimal element of the array.
Max() (k K, v V, found bool) // returns the maximal element of the array.
At(position int) (k K, v V) // returns the ith element.
Len() // returns the number of elements.

modify:
Insert(k K, v V) (inserted bool) // inserts a k,v pair.
Delete(k K) (v V, deleted bool) // deletes a value.
DeleteAt(position int) (k K, v V) // deletes the ith element.
Clear() // deletes all the elements.

iterate:
AscendFromStart() Iterator[K, V] // returns an iterator pointing to the smallest element.
DescendFromEnd() Iterator[K, V] // returns an iterator pointing to the largest element.
Ascend(from K) Iterator[K, V] // returns an iterator pointing to the element that's >= `from`.
Descend(from K) Iterator[K, V] // returns an iterator pointing to the element that's <= `from`.
```

Please see the [examples](/tree_example_test.go) for more details.

## Contact

[Aleksandr Demakin](mailto:alexander.demakin@gmail.com)

## License

Source code is available under the [Apache License Version 2.0](/LICENSE).