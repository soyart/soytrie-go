package soytrie

import (
	"errors"
	"fmt"
)

type Mode uint8

const (
	ModeExact = iota
	ModePrefix
)

type Node[K comparable, V any] struct {
	Value    V
	Valued   bool
	Children map[K]*Node[K, V]
}

func New[K comparable, V any]() *Node[K, V] {
	return &Node[K, V]{}
}

func NewWithValue[K comparable, V any](v V) *Node[K, V] {
	return &Node[K, V]{Value: v, Valued: true}
}

func Zero[T any]() T {
	var zero T
	return zero
}

func (n *Node[K, V]) HasDirect(k K) bool {
	_, ok := n.Children[k]
	return ok
}

func (n *Node[K, V]) GetDirect(k K) (*Node[K, V], bool) {
	child, ok := n.Children[k]
	return child, ok
}

func (n *Node[K, V]) RemoveDirect(k K) (*Node[K, V], bool) {
	target, ok := n.Children[k]
	if !ok {
		return nil, false
	}
	delete(n.Children, k)
	return target, true
}

func (n *Node[K, V]) Get(path ...K) (*Node[K, V], bool) {
	curr := n
	for i := range path {
		p := path[i]
		next, ok := curr.GetDirect(p)
		if !ok {
			return nil, false
		}
		curr = next
	}
	return curr, true
}

func (n *Node[K, V]) Search(mode Mode, path ...K) bool {
	target, ok := n.Get(path...)
	if !ok {
		return false
	}
	if mode == ModePrefix {
		return true
	}
	return target.Valued
}

func (n *Node[K, V]) Predict(mode Mode, path ...K) ([]*Node[K, V], bool) {
	target, ok := n.Get(path...)
	if !ok {
		return nil, false
	}

	collector := []*Node[K, V]{}
	switch mode {
	case ModeExact:
		CollectChildrenValued(target, &collector)

	case ModePrefix:
		CollectChildren(target, &collector)
	}

	return collector, true
}

// Collect collects all values under node to collector
func Collect[K comparable, V any](
	testFn func(*Node[K, V]) bool,
	node *Node[K, V],
	collector *[]*Node[K, V],
) {
	if testFn == nil || testFn(node) {
		*collector = append(*collector, node)
	}
	for _, child := range node.Children {
		Collect(testFn, child, collector)
	}
}

func CollectChildren[K comparable, V any](
	node *Node[K, V],
	collector *[]*Node[K, V],
) {
	Collect(nil, node, collector)
}

func CollectChildrenValued[K comparable, V any](
	node *Node[K, V],
	collector *[]*Node[K, V],
) {
	Collect(func(n *Node[K, V]) bool {
		return n.Valued
	}, node, collector)
}

func CollectChildrenLeaf[K comparable, V any](
	node *Node[K, V],
	collector *[]*Node[K, V],
) {
	Collect(func(n *Node[K, V]) bool {
		return len(n.Children) == 0
	}, node, collector)
}

// Unique returns whether the path is a unique path
// or a prefix to a valued node.
func (n *Node[K, V]) Unique(path ...K) bool {
	collected, ok := n.Predict(ModeExact, path...)
	if !ok {
		return false
	}
	return len(collected) == 1
}

func (n *Node[K, V]) Remove(path ...K) (*Node[K, V], bool) {
	l := len(path)
	if l == 0 {
		return nil, false
	}
	last, ok := n.Get(path[:l-1]...)
	if !ok {
		return nil, false
	}
	return last.RemoveDirect(path[l-1])
}

func (n *Node[K, V]) GetOrInsertDirectValue(k K, v V) (*Node[K, V], bool) {
	old, ok := n.GetDirect(k)
	if ok {
		return old, true
	}

	child := &Node[K, V]{Value: v}
	if n.Children == nil {
		n.Children = make(map[K]*Node[K, V])
	}

	n.Children[k] = child
	return child, false
}

func (n *Node[K, V]) GetOrInsertDirect(k K, node *Node[K, V]) (*Node[K, V], bool) {
	old, ok := n.GetDirect(k)
	if ok {
		return old, ok
	}
	if n.Children == nil {
		n.Children = make(map[K]*Node[K, V])
	}
	n.Children[k] = node
	return node, false
}

func (n *Node[K, V]) Insert(v V, p0 K, pRest ...K) *Node[K, V] {
	path := append([]K{p0}, pRest...)
	curr := n
	for i := range path {
		p := path[i]
		next, _ := curr.GetOrInsertDirect(p, New[K, V]())
		curr = next
	}
	curr.Valued, curr.Value = true, v
	return curr
}

// InsertStrict inserts v to p0+pRest if and only if
// the insertion would create a new node
func (n *Node[K, V]) InsertStrict(v V, p0 K, pRest ...K) (*Node[K, V], error) {
	path := append([]K{p0}, pRest...)
	last := len(path) - 1

	parent, ok := n.Get(path[:last]...)
	if !ok {
		return nil, errors.New("missing some path")
	}
	if parent.HasDirect(path[last]) {
		return nil, fmt.Errorf("node already had path %v", path[last])
	}
	return parent.Insert(v, path[last]), nil
}

// InsertNoOverwrite inserts v to p0+pRest only if
// the insertion to p0+pRest does not overwrite existing value
func (n *Node[K, V]) InsertNoOverwrite(v V, p0 K, pRest ...K) (*Node[K, V], error) {
	path := append([]K{p0}, pRest...)
	curr := n
	for i := range path {
		next, _ := curr.GetOrInsertDirect(path[i], New[K, V]())
		curr = next
	}
	if curr.Valued {
		return nil, fmt.Errorf("valued node exists: %v", curr.Value)
	}
	curr.Valued, curr.Value = true, v
	return curr, nil
}
