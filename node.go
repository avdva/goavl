package goavl

type node[K, V any] struct {
	k      K
	v      V
	h      uint8
	nchild uint32
}

func (n *node[K, V]) height() uint8 {
	return n.h
}

func (n *node[K, V]) setHeight(height uint8) {
	n.h = height
}

func (n *node[K, V]) childrenCount() uint32 {
	return n.nchild
}

func (n *node[K, V]) setChildrenCount(nchild uint32) {
	n.nchild = nchild
}

func (n *node[K, V]) key() K {
	return n.k
}

func (n *node[K, V]) valuePtr() *V {
	return &n.v
}

func (n *node[K, V]) setValue(v V) {
	n.v = v
}

func (n *node[K, V]) init(k K, v V) {
	n.k = k
	n.v = v
	n.h = 0
	n.nchild = 0
}
