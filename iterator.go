package goavl

type iterator[K, V any, Cmp func(a, b K) int] struct {
	tree *Tree[K, V, Cmp]
	loc  ptrLocation[K, V]
}

func (it *iterator[K, V, Cmp]) next() (k K, v V, ok bool) {
	if it.loc.isNil() {
		return k, v, ok
	}
	k, v, ok = it.loc.key(), it.loc.value(), true
	it.loc = nextLocation(it.loc)
	return k, v, ok
}

func (it *iterator[K, V, Cmp]) prev() (k K, v V, ok bool) {
	if it.loc.isNil() {
		return k, v, ok
	}
	k, v, ok = it.loc.key(), it.loc.value(), true
	it.loc = prevLocation(it.loc)
	return k, v, ok
}

type Iterator[K, V any, Cmp func(a, b K) int] struct {
	iterator[K, V, Cmp]
}

func (it *Iterator[K, V, Cmp]) Next() (k K, v V, ok bool) {
	return it.next()
}

type ReverseIterator[K, V any, Cmp func(a, b K) int] struct {
	iterator[K, V, Cmp]
}

func (it *ReverseIterator[K, V, Cmp]) Next() (k K, v V, ok bool) {
	return it.prev()
}

func nextLocation[K, V any](loc ptrLocation[K, V]) ptrLocation[K, V] {
	if r := loc.right(); !r.isNil() {
		return goLeft(r)
	}
	var dir direction
	for {
		loc, dir = loc.parentAndDir()
		if dir == dirLeft || dir == dirCenter {
			return loc
		}
	}
}

func prevLocation[K, V any](loc ptrLocation[K, V]) ptrLocation[K, V] {
	if l := loc.left(); !l.isNil() {
		return goRight(l)
	}
	var dir direction
	for {
		loc, dir = loc.parentAndDir()
		if dir == dirRight || dir == dirCenter {
			return loc
		}
	}
}
