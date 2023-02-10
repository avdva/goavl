package goavl

const (
	countsMask = 0xFFFFFFF
)

type node[K, V any] struct {
	k      K
	v      V
	counts uint64
}

func (n *node[K, V]) height() uint8 {
	return uint8(n.counts & 0xFF)
}

func (n *node[K, V]) setHeight(height uint8) {
	n.counts = (n.counts & ^uint64(0xFF)) | uint64(height)
}

func (n *node[K, V]) leftCount() uint32 {
	return uint32((n.counts >> 8) & countsMask)
}

func (n *node[K, V]) setLeftNodes(count uint32) {
	n.counts = (n.counts & ^uint64(countsMask<<8)) | (uint64(count) << 8)
}

func (n *node[K, V]) rightCount() uint32 {
	return uint32((n.counts >> 36) & countsMask)
}

func (n *node[K, V]) setRightNodes(count uint32) {
	n.counts = (n.counts & ^uint64(countsMask<<36)) | (uint64(count) << 36)
}

func (n *node[K, V]) childCount() uint32 {
	return n.leftCount() + n.rightCount()
}

func (n *node[K, V]) key() K {
	return n.k
}

func (n *node[K, V]) value() V {
	return n.v
}

func (n *node[K, V]) setValue(v V) {
	n.v = v
}

func makeNode[K, V any](k K, v V) node[K, V] {
	return node[K, V]{
		k:      k,
		v:      v,
		counts: 0,
	}
}
