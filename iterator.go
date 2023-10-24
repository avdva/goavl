package goavl

const (
	itStateBeforeHead = iota + 1
	itStateAfterEnd
)

// Iterator allows to iterate over a tree in ascending or descending order.
type Iterator[K, V any] struct {
	loc, head, tail ptrLocation[K, V]
	state           uint8
}

// Value returns current value and true, if the value is valid.
func (it *Iterator[K, V]) Value() (entry Entry[K, V], found bool) {
	if !it.loc.isNil() {
		found = true
		entry.Key, entry.Value = it.loc.key(), it.loc.valuePtr()
	}
	return entry, found
}

// Next returns current entry and advances the iterator.
func (it *Iterator[K, V]) Next() (entry Entry[K, V], found bool) {
	if it.loc.isNil() {
		if it.state == itStateBeforeHead && !it.head.isNil() {
			it.loc = it.head
		} else {
			return entry, false
		}
	}
	entry.Key, entry.Value = it.loc.key(), it.loc.valuePtr()
	it.loc = nextLocation(it.loc)
	if it.loc.isNil() {
		it.state = itStateAfterEnd
	}
	return entry, true
}

// Prev returns current entry and moves to the previous one.
func (it *Iterator[K, V]) Prev() (entry Entry[K, V], found bool) {
	if it.loc.isNil() {
		if it.state == itStateAfterEnd && !it.tail.isNil() {
			it.loc = it.tail
		} else {
			return entry, false
		}
	}
	entry.Key, entry.Value = it.loc.key(), it.loc.valuePtr()
	it.loc = prevLocation(it.loc)
	if it.loc.isNil() {
		it.state = itStateBeforeHead
	}
	return entry, true
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

func advance[K, V any](loc ptrLocation[K, V], count int) ptrLocation[K, V] {
	for count > 0 {
		loc = nextLocation(loc)
		count--
	}
	return loc
}

func advanceBack[K, V any](loc ptrLocation[K, V], count int) ptrLocation[K, V] {
	for count > 0 {
		loc = prevLocation(loc)
		count--
	}
	return loc
}
