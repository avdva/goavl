package goavl

const (
	itStateBeforeHead = iota + 1
	itStateAfterEnd
)

// Iterator allows to iterate over a tree in ascending or descending order.
type Iterator[K, V any, Cmp func(a, b K) int] struct {
	loc   location[K, V]
	t     *Tree[K, V, Cmp]
	id    uint64
	state uint8
}

// Value returns current value and true, if the value is valid.
func (it *Iterator[K, V, Cmp]) Value() (entry Entry[K, V], found bool) {
	if !it.loc.isNil() {
		found = true
		entry.Key, entry.Value = it.loc.key(), it.loc.valuePtr()
	}
	return entry, found
}

// Next returns current entry and advances the iterator.
func (it *Iterator[K, V, Cmp]) Next() (entry Entry[K, V], found bool) {
	if it.loc.isNil() {
		if it.state == itStateBeforeHead && it.t != nil && !it.t.min.isNil() {
			it.loc = it.t.min
			it.id = it.loc.id()
		} else {
			return entry, false
		}
	}
	entry.Key, entry.Value = it.loc.key(), it.loc.valuePtr()
	it.loc = nextLocation(it.loc)
	if it.loc.isNil() {
		it.state = itStateAfterEnd
		it.id = 0
	} else {
		it.id = it.loc.id()
	}
	return entry, true
}

// Prev returns current entry and moves to the previous one.
func (it *Iterator[K, V, Cmp]) Prev() (entry Entry[K, V], found bool) {
	if it.loc.isNil() {
		if it.state == itStateAfterEnd && it.t != nil && !it.t.max.isNil() {
			it.loc = it.t.max
			it.id = it.loc.id()
		} else {
			return entry, false
		}
	}
	entry.Key, entry.Value = it.loc.key(), it.loc.valuePtr()
	it.loc = prevLocation(it.loc)
	if it.loc.isNil() {
		it.state = itStateBeforeHead
		it.id = 0
	} else {
		it.id = it.loc.id()
	}
	return entry, true
}

func nextLocation[K, V any](loc location[K, V]) location[K, V] {
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

func prevLocation[K, V any](loc location[K, V]) location[K, V] {
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

func advance[K, V any](loc location[K, V], count int) location[K, V] {
	for count > 0 {
		loc = nextLocation(loc)
		count--
	}
	return loc
}

func advanceBack[K, V any](loc location[K, V], count int) location[K, V] {
	for count > 0 {
		loc = prevLocation(loc)
		count--
	}
	return loc
}
