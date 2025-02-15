package soytrie

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

func (n *Node[K, V]) GetDirect(k K) (*Node[K, V], bool) {
	n, ok := n.Children[k]
	return n, ok
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
	curr.Value, curr.Valued = v, true
	return curr
}
