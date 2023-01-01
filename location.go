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

type ptrLocation[K, V any] struct {
	ptr *node[K, V]
}

func makeLocation[K, V any](k K, v V) ptrLocation[K, V] {
	return ptrLocation[K, V]{
		ptr: newNode(k, v),
	}
}

func (l ptrLocation[K, V]) isNil() bool {
	return l.ptr == nil
}

func (l *ptrLocation[K, V]) key() K {
	return l.ptr.k
}

func (l *ptrLocation[K, V]) value() V {
	return l.ptr.v
}

func (l *ptrLocation[K, V]) left() ptrLocation[K, V] {
	return l.ptr.left
}

func (l *ptrLocation[K, V]) right() ptrLocation[K, V] {
	return l.ptr.right
}

func (l *ptrLocation[K, V]) parent() ptrLocation[K, V] {
	return l.ptr.parent
}

func (l *ptrLocation[K, V]) parentAndDir() (parent ptrLocation[K, V], dir direction) {
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

func (l *ptrLocation[K, V]) height() uint8 {
	return l.ptr.height()
}

func (l *ptrLocation[K, V]) setValue(v V) {
	l.ptr.v = v
}

func (l *ptrLocation[K, V]) setKey(k K) {
	l.ptr.k = k
}

func (l *ptrLocation[K, V]) balance() int8 {
	b := int8(0)
	if r := l.right(); !r.isNil() {
		b += int8(r.height()) + 1
	}
	if l := l.left(); !l.isNil() {
		b -= (int8(l.height()) + 1)
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
	l.ptr.parent = parent
}

func (l *ptrLocation[K, V]) setRight(child ptrLocation[K, V]) {
	l.ptr.right = child
	if !child.isNil() {
		child.ptr.parent = *l
	}
}

func (l *ptrLocation[K, V]) setLeft(child ptrLocation[K, V]) {
	l.ptr.left = child
	if !child.isNil() {
		child.ptr.parent = *l
	}
}

// addChild panics if there's a child at this direction.
func (l *ptrLocation[K, V]) addChild(child ptrLocation[K, V], dir direction) {
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

func (l *ptrLocation[K, V]) removeChild(child ptrLocation[K, V]) {
	if l.left() == child {
		l.ptr.left = ptrLocation[K, V]{}
	} else if l.right() == child {
		l.ptr.right = ptrLocation[K, V]{}
	} else {
		panic("wrong dir")
	}
	child.setParent(ptrLocation[K, V]{})
}

func (l *ptrLocation[K, V]) recalcHeight() (heightChanged bool) {
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

func (l *ptrLocation[K, V]) leftCount() uint32 {
	return l.ptr.leftNodes()
}

func (l *ptrLocation[K, V]) rightCount() uint32 {
	return l.ptr.rightNodes()
}

func (l *ptrLocation[K, V]) childCount() uint32 {
	return l.ptr.leftNodes() + l.ptr.rightNodes()
}

func (l *ptrLocation[K, V]) recalcNodeCounts() {
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

func (l *ptrLocation[K, V]) String() string {
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
