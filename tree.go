package goavl

import (
	"sync"

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

	// at is the allocator type used to allocate nodes.
	at int8

	s *sync.Pool

	ao arenaOptions
}

const (
	allocBasic = iota
	allocSyncPool
	allocArenas
)

// WithCountChildren is used to set CountChildren option.
// If set, each node will have a count of children in the left and right sub-trees.
// This enables O(logn) complexity for the functions that operate key positions.
func WithCountChildren(count bool) Option {
	return func(o *Options) {
		o.countChildren = count
	}
}

// WithSyncPoolAllocator makes Tree use sync.Pool to allocate tree nodes.
// Deprecated: use WithSyncPool instead.
func WithSyncPoolAllocator(bool) Option {
	return func(o *Options) {
		o.at = allocSyncPool
	}
}

// WithSyncPool makes Tree use sync.Pool to allocate tree nodes.
// 's' can be nil. In this case new pool is created for each tree.
// if 's' is not nil, one pool can be shared between several instances of a tree, however
//   - all instances should be of the same generic type.
//   - s.New must be nil.
func WithSyncPool(s *sync.Pool) Option {
	return func(o *Options) {
		o.at = allocSyncPool
		o.s = s
	}
}

// Tree is a generic avl tree.
// AVL tree (https://en.wikipedia.org/wiki/AVL_tree) is a self-balancing binary search tree.
// For each node of the tree the heights of the left and right sub-trees differ by at most one.
// Find and Delete operations have O(logn) complexity.
type Tree[K, V any, Cmp func(a, b K) int] struct {
	options        Options
	root, min, max location[K, V]
	length         int
	nextID         uint64
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
	switch result.options.at {
	case allocBasic:
		result.lc = newBasicLocationCache[K, V]()
	case allocSyncPool:
		result.lc = newPooledLocationCache[K, V](result.options.s)
	case allocArenas:
		result.lc = newArenaLocationCache[K, V](result.options.ao)
	}
	return result
}

// NewComparable returns a new tree.
// This is just a wrapper for New(), where K satisfies constraints.Ordered.
// Example: NewComparable[int, int]().
func NewComparable[K constraints.Ordered, V any](opts ...Option) *Tree[K, V, func(a, b K) int] {
	return New[K, V](func(a, b K) int {
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
	newNode.setID(t.newLocationID())
	t.insertLocation(loc, dir, newNode)
	return newNode.valuePtr(), true
}

func (t *Tree[K, V, Cmp]) insertLocation(loc location[K, V], dir direction, newNode location[K, V]) {
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
}

func (t *Tree[K, V, Cmp]) updateCounts(loc location[K, V]) {
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

// Rank returns the position of k in the sorted sequence.
// Returns false if k is not present.
// Time complexity:
//
//	O(logn) - if children node counts are enabled.
//	O(n) - otherwise.
func (t *Tree[K, V, Cmp]) Rank(k K) (rank int, found bool) {
	if !t.options.countChildren {
		return t.rankLinearly(k)
	}
	return t.rankWithCountChildren(k)
}

func (t *Tree[K, V, Cmp]) rankWithCountChildren(k K) (rank int, found bool) {
	loc := t.root
	for !loc.isNil() {
		switch cmp := t.cmp(k, loc.key()); {
		case cmp < 0:
			loc = loc.left()
		case cmp == 0:
			return rank + int(loc.leftChildrenCount()), true
		case cmp > 0:
			rank += int(loc.leftChildrenCount()) + 1
			loc = loc.right()
		}
	}
	return 0, false
}

func (t *Tree[K, V, Cmp]) rankLinearly(k K) (rank int, found bool) {
	it := t.IteratorAtFirst()
	for entry, ok := it.Value(); ok; entry, ok = it.Value() {
		switch cmp := t.cmp(k, entry.Key); {
		case cmp < 0:
			return 0, false
		case cmp == 0:
			return rank, true
		case cmp > 0:
			rank++
			it.Next()
		}
	}
	return 0, false
}

// RankDistance returns the absolute distance between sorted positions of k1 and k2.
// Returns false if k1 or k2 is not present.
// Time complexity:
//
//	O(logn) - if children node counts are enabled.
//	O(n) - otherwise.
func (t *Tree[K, V, Cmp]) RankDistance(k1 K, k2 K) (distance int, found bool) {
	r1, found := t.Rank(k1)
	if !found {
		return 0, false
	}
	r2, found := t.Rank(k2)
	if !found {
		return 0, false
	}
	if r2 >= r1 {
		return r2 - r1, true
	}
	return r1 - r2, true
}

// CountInRange returns the number of elements on the inclusive interval [k1, k2].
// k1 and k2 themselves may not be present in the tree.
// Time complexity:
//
//	O(logn) - if children node counts are enabled.
//	O(n) - otherwise.
func (t *Tree[K, V, Cmp]) CountInRange(k1 K, k2 K) int {
	r1, found := t.lowerBoundRank(k1)
	if !found {
		return 0
	}
	r2, found := t.floorRank(k2)
	if !found {
		return 0
	}
	if r2 >= r1 {
		return r2 - r1 + 1
	}
	return 0
}

func (t *Tree[K, V, Cmp]) lowerBoundRank(k K) (rank int, found bool) {
	if t.options.countChildren {
		return t.lowerBoundRankWithCountChildren(k)
	}
	return t.lowerBoundRankLinearly(k)
}

func (t *Tree[K, V, Cmp]) floorRank(k K) (rank int, found bool) {
	if t.options.countChildren {
		return t.floorRankWithCountChildren(k)
	}
	return t.floorRankLinearly(k)
}

func (t *Tree[K, V, Cmp]) lowerBoundRankLinearly(k K) (rank int, found bool) {
	it := t.IteratorAtFirst()
	for entry, ok := it.Value(); ok; entry, ok = it.Value() {
		if t.cmp(k, entry.Key) <= 0 {
			return rank, true
		}
		rank++
		it.Next()
	}
	return 0, false
}

func (t *Tree[K, V, Cmp]) floorRankLinearly(k K) (rank int, found bool) {
	it := t.IteratorAtFirst()
	for entry, ok := it.Value(); ok; entry, ok = it.Value() {
		if t.cmp(k, entry.Key) < 0 {
			if rank == 0 {
				return 0, false
			}
			return rank - 1, true
		}
		rank++
		it.Next()
	}
	if rank == 0 {
		return 0, false
	}
	return rank - 1, true
}

func (t *Tree[K, V, Cmp]) lowerBoundRankWithCountChildren(k K) (rank int, found bool) {
	loc := t.root
	candidate := -1
	for !loc.isNil() {
		switch cmp := t.cmp(k, loc.key()); {
		case cmp < 0:
			candidate = rank + int(loc.leftChildrenCount())
			loc = loc.left()
		case cmp == 0:
			return rank + int(loc.leftChildrenCount()), true
		case cmp > 0:
			rank += int(loc.leftChildrenCount()) + 1
			loc = loc.right()
		}
	}
	if candidate < 0 {
		return 0, false
	}
	return candidate, true
}

func (t *Tree[K, V, Cmp]) floorRankWithCountChildren(k K) (rank int, found bool) {
	loc := t.root
	candidate := -1
	for !loc.isNil() {
		switch cmp := t.cmp(k, loc.key()); {
		case cmp < 0:
			loc = loc.left()
		case cmp == 0:
			return rank + int(loc.leftChildrenCount()), true
		case cmp > 0:
			candidate = rank + int(loc.leftChildrenCount())
			rank += int(loc.leftChildrenCount()) + 1
			loc = loc.right()
		}
	}
	if candidate < 0 {
		return 0, false
	}
	return candidate, true
}

func (t *Tree[K, V, Cmp]) shouldLocateAtLinearly(position int) bool {
	position = min2(position, t.length-position-1)
	return position <= 8
}

func (t *Tree[K, V, Cmp]) locateAt(position int) location[K, V] {
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

func (t *Tree[K, V, Cmp]) newLocationID() uint64 {
	t.nextID++
	return t.nextID
}

func (t *Tree[K, V, Cmp]) iteratorAt(loc location[K, V]) Iterator[K, V, Cmp] {
	it := Iterator[K, V, Cmp]{
		loc: loc,
		t:   t,
	}
	if !loc.isNil() {
		it.id = loc.id()
	}
	return it
}

// IteratorAt returns an iterator pointing to the i'th element.
// Panics if position >= tree.Len().
// Time complexity:
//
//	O(logn) - if children node counts are enabled.
//	O(n) - otherwise.
func (t *Tree[K, V, Cmp]) IteratorAt(position int) Iterator[K, V, Cmp] {
	loc := t.locateAt(position)
	return t.iteratorAt(loc)
}

// AscendAt returns an iterator pointing to the i'th element.
//
// Deprecated: use IteratorAt instead.
func (t *Tree[K, V, Cmp]) AscendAt(position int) Iterator[K, V, Cmp] {
	return t.IteratorAt(position)
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

func (t *Tree[K, V, Cmp]) canUpdateKeyInPlace(loc location[K, V], newKey K) bool {
	if prev := prevLocation(loc); !prev.isNil() && t.cmp(prev.key(), newKey) >= 0 {
		return false
	}
	if next := nextLocation(loc); !next.isNil() && t.cmp(newKey, next.key()) >= 0 {
		return false
	}
	return true
}

func (t *Tree[K, V, Cmp]) resetDetachedLocation(loc location[K, V], k K, v V) {
	loc.init(k, v)
}

// UpdateKey changes a node key while preserving its value.
// If newKey already exists, the old value replaces the existing value and oldKey is removed.
// Returns a pointer to the final value and true if oldKey was present.
// Time complexity: O(logn).
func (t *Tree[K, V, Cmp]) UpdateKey(oldKey K, newKey K) (valuePtr *V, updated bool) {
	oldLoc, oldDir := t.locate(oldKey)
	if oldDir != dirCenter || oldLoc.isNil() {
		return nil, false
	}
	if t.cmp(oldLoc.key(), newKey) == 0 {
		oldLoc.k = newKey
		return oldLoc.valuePtr(), true
	}

	newLoc, newDir := t.locate(newKey)
	if newDir == dirCenter && !newLoc.isNil() {
		oldValue := *oldLoc.valuePtr()
		newLoc.setValue(oldValue)
		t.deleteAndReplace(oldLoc)
		return newLoc.valuePtr(), true
	}

	if t.canUpdateKeyInPlace(oldLoc, newKey) {
		oldLoc.k = newKey
		return oldLoc.valuePtr(), true
	}

	oldValue := *oldLoc.valuePtr()
	t.detachAndReplace(oldLoc)
	t.resetDetachedLocation(oldLoc, newKey, oldValue)
	t.insertLocation(newLoc, newDir, oldLoc)
	return oldLoc.valuePtr(), true
}

// DeleteIterator deletes the element referenced by the iterator.
// Returns iterator to the next element.
// Time complexity: O(logn).
func (t *Tree[K, V, Cmp]) DeleteIterator(it Iterator[K, V, Cmp]) Iterator[K, V, Cmp] {
	if it.t != t || !t.isValidloc(it.loc, it.id) {
		return Iterator[K, V, Cmp]{}
	}
	next := nextLocation(it.loc)
	t.deleteAndReplace(it.loc)
	return t.iteratorAt(next)
}

func (t *Tree[K, V, Cmp]) isValidloc(loc location[K, V], id uint64) bool {
	if loc.isNil() || loc.id() != id {
		return false
	}
	for {
		parent := loc.parent()
		if parent.isNil() {
			return loc == t.root
		}
		loc = parent
	}
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

func (t *Tree[K, V, Cmp]) findReplacement(loc location[K, V]) location[K, V] {
	var replacement location[K, V]
	left, right := loc.left(), loc.right()
	if !left.isNil() {
		if !right.isNil() {
			// Russell A. Brown, Optimized Deletion From an AVL Tree.
			// https://arxiv.org/pdf/2406.05162v5
			if loc.balance() <= 0 {
				replacement = goRight(left)
			} else {
				replacement = goLeft(right)
			}
		} else {
			replacement = left
		}
	} else if !right.isNil() {
		replacement = right
	}
	return replacement
}

func (t *Tree[K, V, Cmp]) detachAndReplace(loc location[K, V]) {
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
	t.length--
}

func (t *Tree[K, V, Cmp]) deleteAndReplace(loc location[K, V]) {
	t.detachAndReplace(loc)
	loc.ptrNode.left = location[K, V]{}
	loc.ptrNode.right = location[K, V]{}
	loc.ptrNode.parent = location[K, V]{}
	t.lc.release(loc)
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

func (t *Tree[K, V, Cmp]) setRoot(root location[K, V]) {
	t.root = root
	if !t.root.isNil() {
		t.root.setParent(location[K, V]{})
	}
}

// Clear clears the tree in O(1) time.
// Allocated nodes are not returned to the allocator. Delete elements explicitly
// if you want allocator-specific release behavior, such as sync.Pool reuse.
func (t *Tree[K, V, Cmp]) Clear() {
	t.root = location[K, V]{}
	t.min = t.root
	t.max = t.root
	t.length = 0
}

// Len returns the number of elements.
func (t *Tree[K, V, Cmp]) Len() int {
	return t.length
}

// IteratorAtFirst returns an iterator pointing to the minimum element.
func (t *Tree[K, V, Cmp]) IteratorAtFirst() Iterator[K, V, Cmp] {
	return t.iteratorAt(t.min)
}

// IteratorAtLast returns an iterator pointing to the maximum element.
func (t *Tree[K, V, Cmp]) IteratorAtLast() Iterator[K, V, Cmp] {
	return t.iteratorAt(t.max)
}

// LowerBound returns an iterator pointing to the first element whose key is not less than k.
func (t *Tree[K, V, Cmp]) LowerBound(k K) Iterator[K, V, Cmp] {
	loc := t.root
	var candidate location[K, V]
	for !loc.isNil() {
		switch cmp := t.cmp(k, loc.key()); {
		case cmp < 0:
			candidate = loc
			loc = loc.left()
		case cmp == 0:
			return t.iteratorAt(loc)
		case cmp > 0:
			loc = loc.right()
		}
	}
	return t.iteratorAt(candidate)
}

// UpperBound returns an iterator pointing to the first element whose key is greater than k.
func (t *Tree[K, V, Cmp]) UpperBound(k K) Iterator[K, V, Cmp] {
	loc := t.root
	var candidate location[K, V]
	for !loc.isNil() {
		switch cmp := t.cmp(k, loc.key()); {
		case cmp < 0:
			candidate = loc
			loc = loc.left()
		default:
			loc = loc.right()
		}
	}
	return t.iteratorAt(candidate)
}

// Floor returns an iterator pointing to the last element whose key is not greater than k.
func (t *Tree[K, V, Cmp]) Floor(k K) Iterator[K, V, Cmp] {
	loc := t.root
	var candidate location[K, V]
	for !loc.isNil() {
		switch cmp := t.cmp(k, loc.key()); {
		case cmp < 0:
			loc = loc.left()
		case cmp == 0:
			return t.iteratorAt(loc)
		case cmp > 0:
			candidate = loc
			loc = loc.right()
		}
	}
	return t.iteratorAt(candidate)
}

// AscendFromStart returns an iterator pointing to the minimum element.
//
// Deprecated: use IteratorAtFirst instead.
func (t *Tree[K, V, Cmp]) AscendFromStart() Iterator[K, V, Cmp] {
	return t.IteratorAtFirst()
}

// DescendFromEnd returns an iterator pointing to the maximum element.
//
// Deprecated: use IteratorAtLast instead.
func (t *Tree[K, V, Cmp]) DescendFromEnd() Iterator[K, V, Cmp] {
	return t.IteratorAtLast()
}

// Ascend returns an iterator pointing to the element that's >= from.
//
// Deprecated: use LowerBound instead.
func (t *Tree[K, V, Cmp]) Ascend(from K) Iterator[K, V, Cmp] {
	return t.LowerBound(from)
}

// Descend returns an iterator pointing to the element that's <= from.
//
// Deprecated: use Floor instead.
func (t *Tree[K, V, Cmp]) Descend(from K) Iterator[K, V, Cmp] {
	return t.Floor(from)
}

func (t *Tree[K, V, Cmp]) locate(k K) (loc location[K, V], dir direction) {
	loc = t.root
	dir = dirCenter
	if loc.isNil() {
		return loc, dir
	}
	for {
		var next location[K, V]
		switch cmp := t.cmp(k, loc.key()); {
		case cmp < 0:
			next = loc.left()
			dir = dirLeft
		case cmp == 0:
			return loc, dirCenter
		case cmp > 0:
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

func (t *Tree[K, V, Cmp]) treeRotated(parent, oldRoot, newRoot location[K, V]) {
	if !parent.isNil() {
		parent.setChild(newRoot, parent.childDir(oldRoot))
	} else {
		t.setRoot(newRoot)
	}
}

func (t *Tree[K, V, Cmp]) checkBalance(loc location[K, V], fullWayUp bool) {
	for !loc.isNil() {
		parent := loc.parent()
		switch loc.balance() {
		case -2:
			left := loc.left()
			switch left.balance() {
			case -1, 0:
				t.treeRotated(parent, loc, rr(loc, t.options.countChildren))
			case 1:
				t.treeRotated(parent, loc, lr(loc, t.options.countChildren))
			default:
				panic("wrong balance" + loc.String())
			}
		case 2:
			right := loc.right()
			switch right.balance() {
			case -1:
				t.treeRotated(parent, loc, rl(loc, t.options.countChildren))
			case 1, 0:
				t.treeRotated(parent, loc, ll(loc, t.options.countChildren))
			default:
				panic("wrong balance" + loc.String())
			}
		default:
			if !loc.recalcHeight() && !fullWayUp {
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

func rr[K, V any](loc location[K, V], recalcCounts bool) location[K, V] {
	left := loc.left()
	leftRight := left.right()

	loc.setLeft(leftRight)
	left.setRight(loc)

	loc.recalcHeight()
	left.recalcHeight()

	if recalcCounts {
		loc.recalcCounts()
		left.recalcCounts()
	}

	return left
}

func lr[K, V any](loc location[K, V], recalcCounts bool) location[K, V] {
	left := loc.left()
	leftRight := left.right()

	leftRightRight := leftRight.right()
	leftRightLeft := leftRight.left()

	leftRight.setRight(loc)
	leftRight.setLeft(left)

	loc.setLeft(leftRightRight)
	left.setRight(leftRightLeft)

	loc.recalcHeight()
	left.recalcHeight()
	leftRight.recalcHeight()

	if recalcCounts {
		loc.recalcCounts()
		left.recalcCounts()
		leftRight.recalcCounts()
	}

	return leftRight
}

func rl[K, V any](loc location[K, V], recalcCounts bool) location[K, V] {
	right := loc.right()
	rightLeft := right.left()

	rightLeftLeft := rightLeft.left()
	rightLeftRight := rightLeft.right()

	rightLeft.setLeft(loc)
	rightLeft.setRight(right)

	loc.setRight(rightLeftLeft)
	right.setLeft(rightLeftRight)

	loc.recalcHeight()
	right.recalcHeight()
	rightLeft.recalcHeight()

	if recalcCounts {
		loc.recalcCounts()
		right.recalcCounts()
		rightLeft.recalcCounts()
	}

	return rightLeft
}

func ll[K, V any](loc location[K, V], recalcCounts bool) location[K, V] {
	right := loc.right()
	rightLeft := right.left()

	loc.setRight(rightLeft)
	right.setLeft(loc)

	loc.recalcHeight()
	right.recalcHeight()

	if recalcCounts {
		loc.recalcCounts()
		right.recalcCounts()
	}

	return right
}
