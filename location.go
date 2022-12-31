package goavl

import "fmt"

const (
	dirLeft   = -1
	dirCenter = 0
	dirRight  = 1
)

type direction int8

func (d direction) invert() direction {
	return -d
}

type location[K, V any] struct {
	ptr *node[K, V]
}

func makeLocation[K, V any](k K, v V) location[K, V] {
	return location[K, V]{
		ptr: newNode(k, v),
	}
}

func (l location[K, V]) isNil() bool {
	return l.ptr == nil
}

func (l *location[K, V]) key() K {
	return l.ptr.k
}

func (l *location[K, V]) value() V {
	return l.ptr.v
}

func (l *location[K, V]) left() location[K, V] {
	return l.ptr.left
}

func (l *location[K, V]) right() location[K, V] {
	return l.ptr.right
}

func (l *location[K, V]) parent() location[K, V] {
	return l.ptr.parent
}

func (l *location[K, V]) parentAndDir() (parent location[K, V], dir direction) {
	parent = l.ptr.parent
	if parent.isNil() {
		return parent, dirCenter
	}
	if parent.left() == *l {
		dir = dirLeft
	} else if parent.right() == *l {
		dir = dirRight
	} else {
		panic("parents aren't consistent")
	}
	return parent, dir
}

func (l *location[K, V]) height() uint8 {
	return l.ptr.height()
}

func (l *location[K, V]) setValue(v V) {
	l.ptr.v = v
}

func (l *location[K, V]) setKey(k K) {
	l.ptr.k = k
}

func (l *location[K, V]) balance() int8 {
	b := int8(0)
	if r := l.right(); !r.isNil() {
		b += int8(r.height()) + 1
	}
	if l := l.left(); !l.isNil() {
		b -= (int8(l.height()) + 1)
	}
	return b
}

func (l *location[K, V]) setChild(child location[K, V], dir direction) {
	if dir == dirLeft {
		l.setLeft(child)
	} else if dir == dirRight {
		l.setRight(child)
	}
}

func (l *location[K, V]) childAt(dir direction) location[K, V] {
	if dir == dirCenter {
		panic("invalid direction")
	}
	if dir == dirLeft {
		return l.left()
	}
	return l.right()
}

func (l *location[K, V]) setParent(parent location[K, V]) {
	l.ptr.parent = parent
}

func (l *location[K, V]) setRight(child location[K, V]) {
	l.ptr.right = child
	if !child.isNil() {
		child.ptr.parent = *l
	}
}

func (l *location[K, V]) setLeft(child location[K, V]) {
	l.ptr.left = child
	if !child.isNil() {
		child.ptr.parent = *l
	}
}

// addChild panics if there's a child at this direction.
func (l *location[K, V]) addChild(child location[K, V], dir direction) {
	child.ptr.parent = *l
	if dir == dirLeft {
		if !l.ptr.left.isNil() {
			panic("already has a left child")
		}
		l.ptr.left = child
	} else if dir == dirRight {
		if !l.ptr.right.isNil() {
			panic("already has a right child")
		}
		l.ptr.right = child
	} else {
		panic("wrong dir")
	}
}

func (l *location[K, V]) removeChild(child location[K, V]) {
	if l.left() == child {
		l.ptr.left = location[K, V]{}
	} else if l.right() == child {
		l.ptr.right = location[K, V]{}
	} else {
		panic("wrong dir")
	}
	child.setParent(location[K, V]{})
}

func (l *location[K, V]) recalcHeight() (heightChanged bool) {
	var height uint8
	if !l.ptr.left.isNil() {
		height = 1 + l.ptr.left.height()
	}
	if !l.ptr.right.isNil() {
		height = uint8Max(height, 1+l.ptr.right.height())
	}
	heightChanged = height != l.ptr.height()
	l.ptr.setHeight(height)
	return heightChanged
}

func (l *location[K, V]) leftCount() uint32 {
	return l.ptr.leftNodes()
}

func (l *location[K, V]) rightCount() uint32 {
	return l.ptr.rightNodes()
}

func (l *location[K, V]) childCount() uint32 {
	return l.ptr.leftNodes() + l.ptr.rightNodes()
}

func (l *location[K, V]) recalcNodeCounts() {
	var leftCount, rightCount uint32
	if left := l.left(); !left.isNil() {
		leftCount = 1 + left.childCount()
	}
	if right := l.right(); !right.isNil() {
		rightCount = 1 + right.childCount()
	}
	l.ptr.setLeftNodes(leftCount)
	l.ptr.setRightNodes(rightCount)
}

func (l *location[K, V]) String() string {
	var parentKey K
	if p := l.parent(); !p.isNil() {
		parentKey = p.key()
	}
	return fmt.Sprintf("{k: %v, v: %v, p: %v b: %d, h: %d, l: %d, r: %d}",
		l.ptr.k, l.ptr.v, parentKey, l.balance(), l.height(), l.leftCount(), l.rightCount())
}

func uint8Max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
