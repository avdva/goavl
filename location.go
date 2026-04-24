package goavl

import (
	"fmt"
)

const (
	dirLeft   direction = -1
	dirCenter direction = 0
	dirRight  direction = 1
)

type direction int8

func (d direction) invert() direction {
	return -d
}

type ptrNode[K, V any] struct {
	node[K, V]
	id                  uint64
	left, right, parent location[K, V]
}

func (n *ptrNode[K, V]) init(k K, v V) {
	n.node.init(k, v)
	n.left = location[K, V]{}
	n.right = location[K, V]{}
	n.parent = location[K, V]{}
}

type location[K, V any] struct {
	*ptrNode[K, V]
}

func (l location[K, V]) isNil() bool {
	return l.ptrNode == nil
}

func (l *location[K, V]) parentAndDir() (parent location[K, V], dir direction) {
	parent = l.parent()
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

func (l *location[K, V]) childDir(child location[K, V]) direction {
	if left := l.left(); child == left {
		return dirLeft
	}
	if right := l.right(); child == right {
		return dirRight
	}
	return dirCenter
}

func (l *location[K, V]) balance() int8 {
	b := int16(0)
	if r := l.right(); !r.isNil() {
		b += int16(r.height()) + 1
	}
	if l := l.left(); !l.isNil() {
		b -= int16(l.height()) + 1
	}
	return int8(b)
}

func (l *location[K, V]) setChild(child location[K, V], dir direction) {
	switch dir {
	case dirLeft:
		l.setLeft(child)
	case dirRight:
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
	l.ptrNode.parent = parent
}

func (l *location[K, V]) id() uint64 {
	return l.ptrNode.id
}

func (l *location[K, V]) setID(id uint64) {
	l.ptrNode.id = id
}

func (l *location[K, V]) setRight(child location[K, V]) {
	l.ptrNode.right = child
	if !child.isNil() {
		child.ptrNode.parent = *l
	}
}

func (l *location[K, V]) setLeft(child location[K, V]) {
	l.ptrNode.left = child
	if !child.isNil() {
		child.ptrNode.parent = *l
	}
}

// addChild panics if there's a child at this direction.
func (l *location[K, V]) addChild(child location[K, V], dir direction) {
	child.ptrNode.parent = *l
	if dir == dirLeft {
		if !l.ptrNode.left.isNil() {
			panic("already has a left child")
		}
		l.ptrNode.left = child
	} else if dir == dirRight {
		if !l.ptrNode.right.isNil() {
			panic("already has a right child")
		}
		l.ptrNode.right = child
	} else {
		panic("wrong dir")
	}
}

func (l *location[K, V]) removeChild(child location[K, V]) {
	if l.left() == child {
		l.ptrNode.left = location[K, V]{}
	} else if l.right() == child {
		l.ptrNode.right = location[K, V]{}
	} else {
		panic("wrong dir")
	}
	child.setParent(location[K, V]{})
}

func (l *location[K, V]) recalcHeight() (heightChanged bool) {
	var height uint8
	if l := l.left(); !l.isNil() {
		height = 1 + l.height()
	}
	if r := l.right(); !r.isNil() {
		height = max2(height, 1+r.height())
	}
	heightChanged = height != l.height()
	l.setHeight(height)
	return heightChanged
}

func (l *location[K, V]) recalcCounts() {
	var nchild uint32
	if left := l.left(); !left.isNil() {
		nchild += 1 + left.childrenCount()
	}
	if right := l.right(); !right.isNil() {
		nchild += 1 + right.childrenCount()
	}
	l.setChildrenCount(nchild)
}

func (l *location[K, V]) parent() location[K, V] {
	return l.ptrNode.parent
}

func (l *location[K, V]) right() location[K, V] {
	return l.ptrNode.right
}

func (l *location[K, V]) left() location[K, V] {
	return l.ptrNode.left
}

func (l *location[K, V]) leftChildrenCount() uint32 {
	if l := l.left(); !l.isNil() {
		return 1 + l.childrenCount()
	}
	return 0
}

func (l *location[K, V]) String() string {
	var parentKey K
	if p := l.parent(); !p.isNil() {
		parentKey = p.key()
	}
	return fmt.Sprintf("{k: %v, v: %v, p: %v b: %d, h: %d, c: %d}",
		l.ptrNode.k, l.ptrNode.v, parentKey, l.balance(), l.height(), l.childrenCount())
}
