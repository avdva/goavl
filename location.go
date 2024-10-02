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
	left, right, parent ptrLocation[K, V]
}

func (n *ptrNode[K, V]) init(k K, v V) {
	n.node.init(k, v)
	n.left = ptrLocation[K, V]{}
	n.right = ptrLocation[K, V]{}
	n.parent = ptrLocation[K, V]{}
}

type ptrLocation[K, V any] struct {
	*ptrNode[K, V]
}

func (l ptrLocation[K, V]) isNil() bool {
	return l.ptrNode == nil
}

func (l *ptrLocation[K, V]) parentAndDir() (parent ptrLocation[K, V], dir direction) {
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

func (l *ptrLocation[K, V]) balance() int8 {
	b := int16(0)
	if r := l.right(); !r.isNil() {
		b += int16(r.height()) + 1
	}
	if l := l.left(); !l.isNil() {
		b -= int16(l.height()) + 1
	}
	return int8(b)
}

func (l *ptrLocation[K, V]) setChild(child ptrLocation[K, V], dir direction) {
	if dir == dirLeft {
		l.setLeft(child)
	} else if dir == dirRight {
		l.setRight(child)
	}
}

func (l *ptrLocation[K, V]) childAt(dir direction) ptrLocation[K, V] {
	if dir == dirCenter {
		panic("invalid direction")
	}
	if dir == dirLeft {
		return l.left()
	}
	return l.right()
}

func (l *ptrLocation[K, V]) setParent(parent ptrLocation[K, V]) {
	l.ptrNode.parent = parent
}

func (l *ptrLocation[K, V]) setRight(child ptrLocation[K, V]) {
	l.ptrNode.right = child
	if !child.isNil() {
		child.ptrNode.parent = *l
	}
}

func (l *ptrLocation[K, V]) setLeft(child ptrLocation[K, V]) {
	l.ptrNode.left = child
	if !child.isNil() {
		child.ptrNode.parent = *l
	}
}

// addChild panics if there's a child at this direction.
func (l *ptrLocation[K, V]) addChild(child ptrLocation[K, V], dir direction) {
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

func (l *ptrLocation[K, V]) removeChild(child ptrLocation[K, V]) {
	if l.left() == child {
		l.ptrNode.left = ptrLocation[K, V]{}
	} else if l.right() == child {
		l.ptrNode.right = ptrLocation[K, V]{}
	} else {
		panic("wrong dir")
	}
	child.setParent(ptrLocation[K, V]{})
}

func (l *ptrLocation[K, V]) recalcHeight() (heightChanged bool) {
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

func (l *ptrLocation[K, V]) recalcCounts() {
	var nchild uint32
	if left := l.left(); !left.isNil() {
		nchild += 1 + left.childrenCount()
	}
	if right := l.right(); !right.isNil() {
		nchild += 1 + right.childrenCount()
	}
	l.setChildrenCount(nchild)
}

func (l *ptrLocation[K, V]) parent() ptrLocation[K, V] {
	return l.ptrNode.parent
}

func (l *ptrLocation[K, V]) right() ptrLocation[K, V] {
	return l.ptrNode.right
}

func (l *ptrLocation[K, V]) left() ptrLocation[K, V] {
	return l.ptrNode.left
}

func (l *ptrLocation[K, V]) leftChildrenCount() uint32 {
	if l := l.left(); !l.isNil() {
		return 1 + l.childrenCount()
	}
	return 0
}

func (l *ptrLocation[K, V]) String() string {
	var parentKey K
	if p := l.parent(); !p.isNil() {
		parentKey = p.key()
	}
	return fmt.Sprintf("{k: %v, v: %v, p: %v b: %d, h: %d, c: %d}",
		l.ptrNode.k, l.ptrNode.v, parentKey, l.balance(), l.height(), l.childrenCount())
}
