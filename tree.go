package goavl

import (
	"golang.org/x/exp/constraints"
)

// Option is a function type used to configure tree's behavior.
type Option func(o *Options)

// Options defines some parameters of the tree.
type Options struct {
	// countChildren, if set, enables children counts for every node of the tree.
	// the numbers of children in the left and right subtrees allows to locate
	// a node by its position with a guaranteed complexity O(logn).
	countChildren bool

	// syncPoolAllocator, if set makes Tree use sync.Pool to allocate tree nodes,
	// which can improve performance in some use cases.
	syncPoolAllocator bool
}

// WithCountChildren is used to set CountChildren option.
// If set, each node will have a count of children in the left and right sub-trees.
// This enables O(logn) complexity for the functions that operate key positions.
func WithCountChildren(count bool) Option {
	return func(o *Options) {
		o.countChildren = count
	}
}

// WithSyncPoolAllocator makes Tree use sync.Pool to allocate tree nodes.
func WithSyncPoolAllocator(with bool) Option {
	return func(o *Options) {
		o.syncPoolAllocator = with
	}
}

// Tree is a generic avl tree.
// AVL tree (https://en.wikipedia.org/wiki/AVL_tree) is a self-balancing binary search tree.
// For each node of the tree the heights of the left and right sub-trees differ by at most one.
// Find and Delete operations have O(logn) complexity.
type Tree[K, V any, Cmp func(a, b K) int] struct {
	options        Options
	root, min, max ptrLocation[K, V]
	length         int
	cmp            Cmp
	lc             locationCache[K, V]
}

// New returns a new Tree.
// K - key type, V - value type can be any types, one only has to define a comparator func:
// func(a, b K) int that should return
//
//	-1, if a < b
//	0, if a == b
//	1, if a > b
//
// Example:
//
//	func intCmp(a, b int) int {
//		if a < b {
//			return -1
//		}
//		if a > b {
//			return 1
//		}
//		return 0
//	}
//
// tree := New[int, int](intCmp, WithCountChildren(true)).
func New[K, V any, Cmp func(a, b K) int](cmp Cmp, opts ...Option) *Tree[K, V, Cmp] {
	result := &Tree[K, V, Cmp]{
		cmp: cmp,
		options: Options{
			countChildren: false,
		},
	}
	for _, o := range opts {
		o(&result.options)
	}
	if result.options.syncPoolAllocator {
		result.lc = newPooledLocationCache[K, V]()
	} else {
		result.lc = newSimpleLocationCache[K, V]()
	}
	return result
}

// NewComparable returns a new tree.
// This is just a wrapper for New(), where K satisfies constraints.Ordered.
// Example: NewComparable[int, int]().
func NewComparable[K constraints.Ordered, V any](opts ...Option) *Tree[K, V, func(a, b K) int] {
	return New[K, V, func(a, b K) int](func(a, b K) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}, opts...)
}

// Insert inserts a node into the tree.
// Returns a pointer to the value and true, if a new node was added.
// If the key `k` was present in the tree, node's value is updated to `v`.
// Time complexity: O(logn).
func (t *Tree[K, V, Cmp]) Insert(k K, v V) (valuePtr *V, inserted bool) {
	loc, dir := t.locate(k)
	if dir == dirCenter && !loc.isNil() {
		loc.setValue(v)
		return loc.valuePtr(), false
	}
	newNode := t.lc.new(k, v)
	t.length++
	switch dir {
	case dirLeft, dirRight:
		loc.addChild(newNode, dir)
		if dir == dirLeft && loc == t.min {
			t.min = newNode
		} else if dir == dirRight && loc == t.max {
			t.max = newNode
		}
		if loc.recalcHeight() {
			if t.options.countChildren {
				loc.recalcCounts()
			}
			t.checkBalance(loc.parent(), false)
		} else {
			t.updateCounts(loc)
		}
	case dirCenter:
		t.root = newNode
		t.min, t.max = t.root, t.root
	}
	return newNode.valuePtr(), true
}

func (t *Tree[K, V, Cmp]) updateCounts(loc ptrLocation[K, V]) {
	if !t.options.countChildren {
		return
	}
	for !loc.isNil() {
		loc.recalcCounts()
		loc = loc.parent()
	}
}

// Entry is a pair of a key and a pointer to the value.
type Entry[K, V any] struct {
	Key   K
	Value *V
}

// Find returns a value for key k.
// Time complexity: O(logn).
func (t *Tree[K, V, Cmp]) Find(k K) (v *V, found bool) {
	loc, dir := t.locate(k)
	if dir != dirCenter || loc.isNil() {
		return v, false
	}
	return loc.valuePtr(), true
}

// Min returns the minimum of the tree.
// If the tree is empty, `found` value will be false.
// Time complexity: O(1).
func (t *Tree[K, V, Cmp]) Min() (entry Entry[K, V], found bool) {
	if found = !t.min.isNil(); found {
		entry.Key = t.min.key()
		entry.Value = t.min.valuePtr()
	}
	return entry, found
}

// Max returns the maximum of the tree.
// If the tree is empty, `found` value will be false.
// Time complexity: O(1).
func (t *Tree[K, V, Cmp]) Max() (entry Entry[K, V], found bool) {
	if found = !t.max.isNil(); found {
		entry.Key = t.max.key()
		entry.Value = t.max.valuePtr()
	}
	return entry, found
}

// At returns a (key, value) pair at the ith position of the sorted array.
// Panics if position >= tree.Len().
// Time complexity:
//
//	O(logn) - if children node counts are enabled.
//	O(n) - otherwise.
func (t *Tree[K, V, Cmp]) At(position int) Entry[K, V] {
	node := t.locateAt(position)
	return Entry[K, V]{Key: node.key(), Value: node.valuePtr()}
}

func (t *Tree[K, V, Cmp]) shouldLocateAtLinearly(position int) bool {
	position = min2(position, t.length-position-1)
	return position <= 8
}

func (t *Tree[K, V, Cmp]) locateAt(position int) ptrLocation[K, V] {
	if position < 0 || position >= t.Len() {
		panic("index out of range")
	}
	if !t.options.countChildren || t.shouldLocateAtLinearly(position) {
		if position < t.length/2 {
			return advance(t.min, position)
		}
		return advanceBack(t.max, t.length-position-1)
	}
	node := t.root
	for {
		leftCount := int(node.leftChildrenCount())
		switch {
		case position == leftCount:
			return node
		case position < leftCount:
			node = node.left()
		default:
			position -= (leftCount + 1)
			node = node.right()
		}
	}
}

// AscendAt returns an iterator pointing to the i'th element.
// Panics if position >= tree.Len().
// Time complexity:
//
//	O(logn) - if children node counts are enabled.
//	O(n) - otherwise.
func (t *Tree[K, V, Cmp]) AscendAt(position int) Iterator[K, V] {
	loc := t.locateAt(position)
	return Iterator[K, V]{
		head: t.min,
		tail: t.max,
		loc:  loc,
	}
}

// Delete deletes a node from the tree.
// Returns node's value and true, if the node was present in the tree.
// Time complexity: O(logn).
func (t *Tree[K, V, Cmp]) Delete(k K) (v V, deleted bool) {
	loc, dir := t.locate(k)
	if dir != dirCenter || loc.isNil() {
		return v, false
	}
	v = *loc.valuePtr()
	t.deleteAndReplace(loc)
	return v, true
}

// DeleteIterator deletes the element referenced by the iterator.
// Returns iterator to the next element.
// Time complexity: O(logn).
func (t *Tree[K, V, Cmp]) DeleteIterator(it Iterator[K, V]) Iterator[K, V] {
	if !t.isValidloc(it.loc) {
		return Iterator[K, V]{}
	}
	next := nextLocation(it.loc)
	t.deleteAndReplace(it.loc)
	return Iterator[K, V]{
		head: t.min,
		tail: t.max,
		loc:  next,
	}
}

func (t *Tree[K, V, Cmp]) isValidloc(loc ptrLocation[K, V]) bool {
	return !loc.isNil()
}

// DeleteAt deletes a node at the given position.
// Returns node's value. Panics if position >= tree.Len().
// Time complexity:
//
//	O(logn) - if children node counts are enabled.
//	O(n) - otherwise.
func (t *Tree[K, V, Cmp]) DeleteAt(position int) (k K, v V) {
	loc := t.locateAt(position)
	k = loc.key()
	v = *loc.valuePtr()
	t.deleteAndReplace(loc)
	return k, v
}

func (t *Tree[K, V, Cmp]) findReplacement(loc ptrLocation[K, V]) ptrLocation[K, V] {
	var replacement ptrLocation[K, V]
	left, right := loc.left(), loc.right()
	if !left.isNil() && (right.isNil() || left.height() <= right.height()) {
		replacement = goRight(left)
	} else if !right.isNil() {
		replacement = goLeft(right)
	}
	return replacement
}

func (t *Tree[K, V, Cmp]) deleteAndReplace(loc ptrLocation[K, V]) {
	replacement := t.findReplacement(loc)
	parent, dir := loc.parentAndDir()
	if loc == t.min {
		t.min = nextLocation(loc)
	}
	if loc == t.max {
		t.max = prevLocation(loc)
	}
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
	t.lc.release(loc)
	t.length--
}

func goLeft[K, V any](loc ptrLocation[K, V]) ptrLocation[K, V] {
	if loc.isNil() {
		return loc
	}
	for !loc.left().isNil() {
		loc = loc.left()
	}
	return loc
}

func goRight[K, V any](loc ptrLocation[K, V]) ptrLocation[K, V] {
	if loc.isNil() {
		return loc
	}
	for !loc.right().isNil() {
		loc = loc.right()
	}
	return loc
}

func (t *Tree[K, V, Cmp]) setRoot(root ptrLocation[K, V]) {
	t.root = root
	if !t.root.isNil() {
		t.root.setParent(ptrLocation[K, V]{})
	}
}

// Clear clears the tree.
func (t *Tree[K, V, Cmp]) Clear() {
	t.root = ptrLocation[K, V]{}
	t.min = t.root
	t.max = t.root
	t.length = 0
}

// Len returns the number of elements.
func (t *Tree[K, V, Cmp]) Len() int {
	return t.length
}

// AscendFromStart returns an iterator pointing to the minimum element.
func (t *Tree[K, V, Cmp]) AscendFromStart() Iterator[K, V] {
	return Iterator[K, V]{
		head: t.min,
		tail: t.max,
		loc:  t.min,
	}
}

// DescendFromEnd returns an iterator pointing to the maximum element.
func (t *Tree[K, V, Cmp]) DescendFromEnd() Iterator[K, V] {
	return Iterator[K, V]{
		head: t.min,
		tail: t.max,
		loc:  t.max,
	}
}

// Ascend returns an iterator pointing to the element that's >= `from`.
func (t *Tree[K, V, Cmp]) Ascend(from K) Iterator[K, V] {
	loc, dir := t.locate(from)
	if dir == dirRight {
		for !loc.isNil() && dir == dirRight {
			loc, dir = loc.parentAndDir()
		}
	}
	return Iterator[K, V]{
		head: t.min,
		tail: t.max,
		loc:  loc,
	}
}

// Descend returns an iterator pointing to the element that's <= `from`.
func (t *Tree[K, V, Cmp]) Descend(from K) Iterator[K, V] {
	loc, dir := t.locate(from)
	if dir == dirLeft {
		for !loc.isNil() && dir == dirLeft {
			loc, dir = loc.parentAndDir()
		}
	}
	return Iterator[K, V]{
		head: t.min,
		tail: t.max,
		loc:  loc,
	}
}

func (t *Tree[K, V, Cmp]) locate(k K) (loc ptrLocation[K, V], dir direction) {
	loc = t.root
	dir = dirCenter
	if loc.isNil() {
		return loc, dir
	}
	for {
		var next ptrLocation[K, V]
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

func (t *Tree[K, V, Cmp]) checkBalance(loc ptrLocation[K, V], fullWayUp bool) {
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
			if t.options.countChildren {
				loc.recalcCounts()
			}
		}
		loc = parent
	}
}

func (t *Tree[K, V, Cmp]) rr(loc ptrLocation[K, V]) {
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

	if t.options.countChildren {
		loc.recalcCounts()
		left.recalcCounts()
	}
}

func (t *Tree[K, V, Cmp]) lr(loc ptrLocation[K, V]) {
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

	if t.options.countChildren {
		loc.recalcCounts()
		left.recalcCounts()
		leftRight.recalcCounts()
	}
}

func (t *Tree[K, V, Cmp]) rl(loc ptrLocation[K, V]) {
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

	if t.options.countChildren {
		loc.recalcCounts()
		right.recalcCounts()
		rightLeft.recalcCounts()
	}
}

func (t *Tree[K, V, Cmp]) ll(loc ptrLocation[K, V]) {
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

	if t.options.countChildren {
		loc.recalcCounts()
		right.recalcCounts()
	}
}
