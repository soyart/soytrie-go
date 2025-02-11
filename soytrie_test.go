package soytrie_test

import (
	"fmt"
	"testing"

	"github.com/soyart/soytrie-go"
)

func pront[K comparable, V any](n *soytrie.Node[K, V]) {
	curr := n
	for p, c := range curr.Children {
		fmt.Println("p", p, "child", c.Value)
		pront(c)
	}
}

func TestInsertAndGet(t *testing.T) {
	root := soytrie.New[int, string]()

	v012 := "0,1,2"
	v12 := "1,2"
	v13 := "1,3"
	v12345 := "1,2,3,4,5"
	newPrefix := "_new_"

	_ = root.Insert(v012, 0, 1, 2)
	_ = root.Insert(v12, 1, 2)
	_ = root.Insert(v13, 1, 3)
	_ = root.Insert(v12345, 1, 2, 3, 4, 5)

	if l := len(root.Children); l != 2 {
		t.Fatalf("unexpected number of root children %d, expecting %d", l, 2)
	}

	node1, ok := root.Get(1)
	if !ok {
		t.Fatalf("expecting ok")
	}
	if node1.Valued {
		t.Fatalf("unexpected valued=true, expecting false")
	}
	if node1.Value != "" {
		t.Fatalf("unexpected value '%s', expecting empty string", node1.Value)
	}
	if l := len(node1.Children); l != 2 {
		t.Fatalf("unexpected number of root children %d, expecting %d", l, 2)
	}

	node12, ok := node1.Get(2)
	if !ok {
		t.Fatalf("expecting ok")
	}
	if !node12.Valued {
		t.Fatalf("unexpected valued=false, expecting true")
	}
	if node12.Value != v12 {
		t.Fatalf("unexpected value '%s', expecting '%s'", node12.Value, v12)
	}

	node12345, ok := node12.Get(3, 4, 5)
	if !ok {
		t.Fatalf("expecting ok")
	}
	node12345Also, ok := root.Get(1, 2, 3, 4, 5)
	if !ok {
		t.Fatalf("expecting ok")
	}
	if node12345 != node12345Also {
		t.Fatalf("unexpected pointer value")
	}
	if !node12345.Valued {
		t.Fatalf("unexpected valued=false, expecting true")
	}
	if node12345.Value != v12345 {
		t.Fatalf("unexpected value '%s', expecting '%s'", node12345.Value, v12345)
	}

	// TODO: is this the best approach to overwrite value? What about dangling children?
	//
	// Assert that it overwrites the value, **but not the node**
	newV12345 := newPrefix + v12345
	nodeOverwritten := node12.Insert(newV12345, 3, 4, 5)
	if nodeOverwritten != node12345 {
		t.Fatalf("unexpected pointer value")
	}
	if node12345 != node12345Also {
		t.Fatalf("unexpected pointer value")
	}
	if node12345.Value != newV12345 {
		t.Fatalf("unexpected value '%s', expecting '%s'", node12345.Value, newV12345)
	}
}

func TestRemove(t *testing.T) {
	root := soytrie.New[int, string]()

	root.Insert("1,2,3", 1, 2, 3)
	root.Insert("1,2,2", 1, 2, 2)
	root.Insert("1,2,7", 1, 2, 7)
	root.Insert("1,3,7", 1, 3, 7)
	root.Insert("1,3,8", 1, 3, 8)
	root.Insert("1,10,20", 1, 10, 20)
	root.Insert("0,2,3", 0, 2, 3)

	root.Remove(1, 2)
	if expected, actual := 2, len(root.Children); expected != actual {
		t.Fatalf("unexpected length of children: expecting=%d, got %d", expected, actual)
	}

	root.Remove(0)
	if expected, actual := 1, len(root.Children); expected != actual {
		t.Fatalf("unexpected length of children: expecting=%d, got %d", expected, actual)
	}

	node13, ok := root.Get(1, 3)
	if !ok {
		t.Fatal("unexpected false")
	}
	if expected, actual := 2, len(node13.Children); expected != actual {
		t.Fatalf("unexpected length of children: expecting=%d, got %d", expected, actual)
	}
}
