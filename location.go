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

type ptrNode[K, V any] struct {
	node[K, V]
	left, right, parent ptrLocation[K, V]
}

type ptrLocation[K, V any] struct {
	*ptrNode[K, V]
}

func makeLocation[K, V any](k K, v V) ptrLocation[K, V] {
	return ptrLocation[K, V]{
		ptrNode: &ptrNode[K, V]{
			node: makeNode(k, v),
		},
	}
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
	b := int8(0)
	if r := l.right(); !r.isNil() {
		b += int8(r.height()) + 1
	}
	if l := l.left(); !l.isNil() {
		b -= int8(l.height()) + 1
	}
	return b
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
		height = max(height, 1+r.height())
	}
	heightChanged = height != l.height()
	l.setHeight(height)
	return heightChanged
}

func (l *ptrLocation[K, V]) recalcCounts() {
	var leftCount, rightCount uint32
	if left := l.left(); !left.isNil() {
		leftCount = 1 + left.childCount()
	}
	if right := l.right(); !right.isNil() {
		rightCount = 1 + right.childCount()
	}
	l.setLeftNodes(leftCount)
	l.setRightNodes(rightCount)
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

func (l *ptrLocation[K, V]) String() string {
	var parentKey K
	if p := l.parent(); !p.isNil() {
		parentKey = p.key()
	}
	return fmt.Sprintf("{k: %v, v: %v, p: %v b: %d, h: %d, l: %d, r: %d}",
		l.ptrNode.k, l.ptrNode.v, parentKey, l.balance(), l.height(), l.leftCount(), l.rightCount())
}
