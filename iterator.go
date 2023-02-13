package goavl

const (
	itStateBeforeHead = iota + 1
	itStateAfterEnd
)

type Iterator[K, V any] struct {
	loc, head, tail ptrLocation[K, V]
	state           uint8
}

func (it *Iterator[K, V]) Next() (k K, v V, ok bool) {
	if it.loc.isNil() {
		if it.state == itStateBeforeHead && !it.head.isNil() {
			it.loc = it.head
		} else {
			return k, v, ok
		}
	}
	k, v, ok = it.loc.key(), it.loc.value(), true
	it.loc = nextLocation(it.loc)
	if it.loc.isNil() {
		it.state = itStateAfterEnd
	}
	return k, v, ok
}

func (it *Iterator[K, V]) Prev() (k K, v V, ok bool) {
	if it.loc.isNil() {
		if it.state == itStateAfterEnd && !it.tail.isNil() {
			it.loc = it.tail
		} else {
			return k, v, ok
		}
	}
	k, v, ok = it.loc.key(), it.loc.value(), true
	it.loc = prevLocation(it.loc)
	if it.loc.isNil() {
		it.state = itStateBeforeHead
	}
	return k, v, ok
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
