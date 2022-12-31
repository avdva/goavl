package goavl

import (
	"io"
)

type tree[K, V any, A allocator[K, V]] struct {
	root     location[K, V]
	min, max location[K, V]
	length   int
	cmp      func(a, b K) int
	alloc    A
}

type Option func(o *Options)

type Options struct {
	CountChildren bool
}

func WithCountChildren(count bool) Option {
	return func(o *Options) {
		o.CountChildren = count
	}
}

type Tree[K, V any] struct {
	options Options
	*tree[K, V, *preAllocator[K, V]]
}

func New[K, V any, Cmp func(a, b K) int](cmp Cmp, opts ...Option) *Tree[K, V] {
	result := &Tree[K, V]{
		tree: &tree[K, V, *preAllocator[K, V]]{
			cmp:   cmp,
			alloc: &preAllocator[K, V]{},
		},
	}
	for _, o := range opts {
		o(&result.options)
	}
	return result
}

func (t *Tree[K, V]) Insert(k K, v V) (inserted bool) {
	loc, dir := t.locate(k)
	if dir == dirCenter && !loc.isNil() {
		loc.setValue(v)
		return
	}
	newNode := t.alloc.new(k, v)
	t.length++
	inserted = true
	switch dir {
	case dirLeft, dirRight:
		loc.addChild(newNode, dir)
		if dir == dirLeft && loc == t.min {
			t.min = newNode
		} else if dir == dirRight && loc == t.max {
			t.max = newNode
		}
		if loc.recalcHeight() {
			if t.options.CountChildren {
				loc.recalcNodeCounts()
			}
			t.checkBalance(loc.parent(), false)
		} else {
			t.updateCounts(loc)
		}
	case dirCenter:
		t.root = newNode
		t.min, t.max = t.root, t.root
	}
	return inserted
}

func (t *Tree[K, V]) updateCounts(loc location[K, V]) {
	if !t.options.CountChildren {
		return
	}
	for !loc.isNil() {
		loc.recalcNodeCounts()
		loc = loc.parent()
	}
}

func (t *Tree[K, V]) Find(k K) (v V, found bool) {
	loc, dir := t.locate(k)
	if dir != dirCenter {
		return v, false
	}
	return loc.value(), true
}

func (t *Tree[K, V]) Min() (k K, v V, found bool) {
	if found = !t.min.isNil(); found {
		k = t.min.key()
		v = t.min.value()
	}
	return k, v, found
}

func (t *Tree[K, V]) Max() (k K, v V, found bool) {
	if found = !t.max.isNil(); found {
		k = t.max.key()
		v = t.max.value()
	}
	return k, v, found
}

// At returns a (key, value) pair at the ith position of the sorted array.
// It panics if position >= tree.Len() or children node counts are not enabled for this tree.
// Complexity: O(log(n))
func (t *Tree[K, V]) At(position int) (k K, v V) {
	if position >= t.Len() {
		panic("index out of range")
	}
	if !t.options.CountChildren {
		panic("unsupported operation")
	}
	node := t.root
	for {
		leftCount := int(node.leftCount())
		switch {
		case position == leftCount:
			return node.key(), node.value()
		case position < leftCount:
			node = node.left()
		default:
			position -= (leftCount + 1)
			node = node.right()
		}
	}
}

func (t *Tree[K, V]) Delete(k K) (v V, deleted bool) {
	loc, dir := t.locate(k)
	if dir != dirCenter || loc.isNil() {
		return v, false
	}
	v = loc.value()
	t.deleteAndReplace(loc)
	t.length--
	return v, true
}

func (t *Tree[K, V]) findReplacement(loc location[K, V]) location[K, V] {
	left, right := loc.left(), loc.right()
	var replacement location[K, V]
	if left.isNil() {
		if !right.isNil() {
			replacement = goLeft(right)
		}
	} else if right.isNil() {
		replacement = goRight(left)
	} else {
		if left.height() <= right.height() { // TODO(avd) - find a better estimation
			replacement = goRight(left)
		} else {
			replacement = goLeft(right)
		}
	}
	return replacement
}

func (t *Tree[K, V]) deleteAndReplace(loc location[K, V]) {
	replacement := t.findReplacement(loc)
	parent, dir := loc.parentAndDir()
	if replacement.isNil() {
		if parent.isNil() {
			// the last element. the tree is now empty.
			t.setRoot(parent)
		} else {
			// no children. just remove the node from parent and check balance.
			parent.removeChild(loc)
			t.checkBalance(parent, false)
		}
	} else {
		replacementParent, replacementDir := replacement.parentAndDir()
		if replacementParent == loc {
			// replacement is one of the node's children.
			if parent.isNil() { // no parent, replacement becomes the root.
				t.setRoot(replacement)
			} else {
				// replacement takes place of the deleted node.
				// it takes the other node's child as its own child.
				parent.setChild(replacement, dir)
			}
			inverted := replacementDir.invert()
			replacement.setChild(loc.childAt(inverted), inverted)
			t.checkBalance(replacement, true)
		} else {
			replacementChild := replacement.childAt(replacementDir.invert())
			replacementParent.setChild(replacementChild, replacementDir)
			if parent.isNil() {
				t.setRoot(replacement)
			} else {
				parent.setChild(replacement, dir)
			}
			replacement.setLeft(loc.left())
			replacement.setRight(loc.right())
			t.checkBalance(replacementParent, true)
		}
	}
	t.alloc.free(loc)
}

func goLeft[K, V any](loc location[K, V]) location[K, V] {
	if loc.isNil() {
		return loc
	}
	for !loc.left().isNil() {
		loc = loc.left()
	}
	return loc
}

func goRight[K, V any](loc location[K, V]) location[K, V] {
	if loc.isNil() {
		return loc
	}
	for !loc.right().isNil() {
		loc = loc.right()
	}
	return loc
}

func (t *Tree[K, V]) setRoot(root location[K, V]) {
	t.root = root
	if !t.root.isNil() {
		t.root.setParent(location[K, V]{})
	}
}

func (t *Tree[K, V]) Clear() {
	t.root = location[K, V]{}
	t.min = t.root
	t.max = t.root
	t.length = 0
}

func (t *Tree[K, V]) Len() int {
	return t.length
}

func (t *Tree[K, V]) locate(k K) (loc location[K, V], dir direction) {
	loc = t.root
	if loc.isNil() {
		return loc, dirCenter
	}
	for {
		var next location[K, V]
		switch t.cmp(k, loc.key()) {
		case -1:
			next = loc.left()
			dir = dirLeft
		case 0:
			return loc, dirCenter
		case 1:
			next = loc.right()
			dir = dirRight
		}
		if next.isNil() {
			break
		}
		loc = next
	}
	return loc, dir
}

func (t *Tree[K, V]) checkBalance(loc location[K, V], fullWayUp bool) {
	for {
		if loc.isNil() {
			return
		}
		heightChanged := loc.recalcHeight()
		parent := loc.parent()
		switch loc.balance() {
		case -2:
			left := loc.left()
			switch left.balance() {
			case -1, 0:
				t.rr(loc)
			case 1:
				t.lr(loc)
			default:
				panic("wrong balance" + loc.String())
			}
		case 2:
			right := loc.right()
			switch right.balance() {
			case -1:
				t.rl(loc)
			case 1, 0:
				t.ll(loc)
			default:
				panic("wrong balance" + loc.String())
			}
		default:
			if !heightChanged && !fullWayUp {
				t.updateCounts(loc)
				return
			}
			if t.options.CountChildren {
				loc.recalcNodeCounts()
			}
		}
		loc = parent
	}
}

func (t *Tree[K, V]) rr(loc location[K, V]) {
	left := loc.left()
	leftRight := left.right()
	parent, dir := loc.parentAndDir()
	if dir != dirCenter {
		parent.setChild(left, dir)
	} else {
		t.setRoot(left)
	}

	loc.setLeft(leftRight)
	left.setRight(loc)

	loc.recalcHeight()
	left.recalcHeight()

	if t.options.CountChildren {
		loc.recalcNodeCounts()
		left.recalcNodeCounts()
	}
}

func (t *Tree[K, V]) lr(loc location[K, V]) {
	left := loc.left()
	leftRight := left.right()

	parent, dir := loc.parentAndDir()
	if dir != dirCenter {
		parent.setChild(leftRight, dir)
	} else {
		t.setRoot(leftRight)
	}

	leftRightRight := leftRight.right()
	leftRightLeft := leftRight.left()

	leftRight.setRight(loc)
	leftRight.setLeft(left)

	loc.setLeft(leftRightRight)
	left.setRight(leftRightLeft)

	loc.recalcHeight()
	left.recalcHeight()
	leftRight.recalcHeight()

	if t.options.CountChildren {
		loc.recalcNodeCounts()
		left.recalcNodeCounts()
		leftRight.recalcNodeCounts()
	}
}

func (t *Tree[K, V]) rl(loc location[K, V]) {
	right := loc.right()
	rightLeft := right.left()

	parent, dir := loc.parentAndDir()
	if dir != dirCenter {
		parent.setChild(rightLeft, dir)
	} else {
		t.setRoot(rightLeft)
	}

	rightLeftLeft := rightLeft.left()
	rightLeftRight := rightLeft.right()

	rightLeft.setLeft(loc)
	rightLeft.setRight(right)

	loc.setRight(rightLeftLeft)
	right.setLeft(rightLeftRight)

	loc.recalcHeight()
	right.recalcHeight()
	rightLeft.recalcHeight()

	if t.options.CountChildren {
		loc.recalcNodeCounts()
		right.recalcNodeCounts()
		rightLeft.recalcNodeCounts()
	}
}

func (t *Tree[K, V]) ll(loc location[K, V]) {
	right := loc.right()
	rightLeft := right.left()
	parent, dir := loc.parentAndDir()
	if dir != dirCenter {
		parent.setChild(right, dir)
	} else {
		t.setRoot(right)
	}
	loc.setRight(rightLeft)
	right.setLeft(loc)

	loc.recalcHeight()
	right.recalcHeight()

	if t.options.CountChildren {
		loc.recalcNodeCounts()
		right.recalcNodeCounts()
	}
}

func (t *Tree[K, V]) traverse(f func(loc location[K, V]) bool) {
	if t.root.isNil() {
		return
	}
	traverseLocation(t.root, f)
}

func traverseLocation[K, V any](loc location[K, V], f func(loc location[K, V]) bool) {
	if !loc.left().isNil() {
		traverseLocation(loc.left(), f)
	}
	f(loc)
	if !loc.right().isNil() {
		traverseLocation(loc.right(), f)
	}
}

func printTree[K, V any](t *Tree[K, V], w io.Writer) {
	t.traverse(func(loc location[K, V]) bool {
		w.Write([]byte(loc.String()))
		w.Write([]byte{'\n'})
		return true
	})
}
