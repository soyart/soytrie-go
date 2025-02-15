package soytrie_test

import (
	"fmt"
	"testing"

	"github.com/soyart/soytrie-go"
)

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

func TestInsertStrict(t *testing.T) {
	root := soytrie.New[int, string]()
	var err error
	_, err = root.InsertStrict("1", 1)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertStrict("2", 2)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertStrict("1", 1)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	_, err = root.InsertStrict("1,2", 1, 2)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertStrict("1,2,3", 1, 2, 3)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertStrict("new_value_1,2,3", 1, 2, 3)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	_, err = root.InsertStrict("1,2,3,4", 1, 2, 3, 4)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
}

func TestInsertNoOverwrite(t *testing.T) {
	root := soytrie.New[int, string]()
	var err error
	_, err = root.InsertNoOverwrite("1", 1)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertNoOverwrite("2", 2)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertNoOverwrite("1", 1)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	_, err = root.InsertNoOverwrite("1,2", 1, 2)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertNoOverwrite("1,2,3,4,5", 1, 2, 3, 4, 5)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertNoOverwrite("1,2,3", 1, 2, 3)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertNoOverwrite("new_value_1,2,3", 1, 2, 3)
	if err == nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertNoOverwrite("1,2,3,4", 1, 2, 3, 4)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	_, err = root.InsertNoOverwrite("10,20,30", 10, 20, 30)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
}

func TestRemove(t *testing.T) {
	t.Run("remove", func(t *testing.T) {
		root := soytrie.New[int, string]()
		_ = root.Insert("1,2,3", 1, 2, 3)
		_ = root.Insert("1,2,2", 1, 2, 2)
		_ = root.Insert("1,2,7", 1, 2, 7)
		_ = root.Insert("1,3,7", 1, 3, 7)
		_ = root.Insert("1,3,8", 1, 3, 8)
		_ = root.Insert("1,10,20", 1, 10, 20)
		_ = root.Insert("0,2,3", 0, 2, 3)

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
	})

	t.Run("remove and predict", func(t *testing.T) {
		root := soytrie.New[int, string]()
		_ = root.Insert("1,2,3", 1, 2, 3)
		_ = root.Insert("1,2,2", 1, 2, 2)
		_ = root.Insert("1,2,7", 1, 2, 7)
		_ = root.Insert("1,3,7", 1, 3, 7)
		_ = root.Insert("1,3,8", 1, 3, 8)
		_ = root.Insert("1,10,20", 1, 10, 20)
		_ = root.Insert("0,2,10,20,30", 0, 2, 10, 20, 30)
		_ = root.Insert("0,2,100,200,300", 0, 2, 100, 200, 300)
		_ = root.Insert("0,2,10,15", 0, 2, 10, 15)
		_ = root.Insert("0,2,10,15,16", 0, 2, 10, 15, 16)
		_ = root.Insert("0,2,10,15,17", 0, 2, 10, 15, 17)

		collector, ok := root.Predict(soytrie.ModePrefix, 0) // [0] [0,2] [0,2,10] [0,2,100] [0,2,10,20] [0,2,10,20,30] [...200] [...300] [... 15] [...16] [...17]
		if !ok {
			panic("unexpected Predict behavior")
		}
		if len(collector) != 11 {
			panic(fmt.Sprintf("unexpected Predict behavior, got len=%d", len(collector)))
		}

		root.Remove(0, 2, 100, 200)
		collector, ok = root.Predict(soytrie.ModePrefix, 0) // [0,2] [0,2,10] [0,2,100] [0,2,10,20] [0,2,10,20,30] [...200] [...300] [... 15] [...16] [...17]
		if !ok {
			panic("unexpected Predict behavior")
		}
		if len(collector) != 9 {
			panic(fmt.Sprintf("unexpected Predict behavior, got len=%d", len(collector)))
		}

		root.Remove(0, 2, 10, 15)
		collector, ok = root.Predict(soytrie.ModePrefix, 0) // [0,2] [0,2,10] [0,2,100] [0,2,10,20] [0,2,10,20,30] [...200] [...300]
		if !ok {
			panic("unexpected Predict behavior")
		}
		if len(collector) != 6 {
			panic(fmt.Sprintf("unexpected Predict behavior, got len=%d", len(collector)))
		}

		collector, ok = root.Predict(soytrie.ModePrefix) // prev 6 + 10 (from p0=1) + root = 17
		if !ok {
			panic("unexpected Predict behavior")
		}
		if len(collector) != 17 {
			panic(fmt.Sprintf("unexpected Predict behavior, got len=%d", len(collector)))
		}
	})
}

func TestSearch(t *testing.T) {
	dirRoot := soytrie.NewWithValue[string]("/")
	_ = dirRoot.Insert(
		"source code",
		"/src",
	)
	_ = dirRoot.Insert(
		"test data",
		"/src", "/testdata",
	)
	_ = dirRoot.Insert(
		"test case for race condition #7",
		"/src", "/testdata", "/race", "/7",
	)
	_ = dirRoot.Insert(
		"main go program",
		"/src", "/cmd", "/main.go",
	)
	_ = dirRoot.Insert(
		"binary releases",
		"/release",
	)
	_ = dirRoot.Insert(
		"binary release foo for amd64",
		"/release", "/amd64", "/bin", "/foo",
	)
	_ = dirRoot.Insert(
		"binary release foo for aarch64",
		"/release", "/aarch64", "/bin", "/foo",
	)

	dirSrc, ok := dirRoot.Get("/src")
	if !ok || dirSrc == nil {
		panic("nil /src node")
	}
	dirAmd64, ok := dirRoot.Get("/release", "/amd64")
	if !ok || dirAmd64 == nil {
		panic("nil /src node")
	}

	type testCase struct {
		path     []string
		mode     soytrie.Mode
		expected bool
	}

	tests := map[*soytrie.Node[string, string]][]testCase{
		dirRoot: {
			{
				path:     []string{"/src"},
				mode:     soytrie.ModeExact,
				expected: true,
			},
			{
				path:     []string{"/src"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/src", "/testdata", "/race"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/src", "/testdata", "/race"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/release"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/release", "/badpath"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/release", "/badpath"},
				mode:     soytrie.ModePrefix,
				expected: false,
			},
			{
				path:     []string{"/release", "/amd64", "/bin"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/release", "/amd64", "/bin"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/release", "/amd64", "/bin", "/foo"},
				mode:     soytrie.ModeExact,
				expected: true,
			},
			{
				path:     []string{"/release", "/amd64", "/bin", "/foo"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/release", "/badarch", "/bin", "/foo"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/release", "/badarch", "/bin", "/foo"},
				mode:     soytrie.ModePrefix,
				expected: false,
			},
			{
				path:     []string{"/release", "/badarch", "/foo"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/release", "/badarch", "/foo"},
				mode:     soytrie.ModePrefix,
				expected: false,
			},
		},
		dirAmd64: {
			{
				path:     []string{"/src"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/foo"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/foo"},
				mode:     soytrie.ModePrefix,
				expected: false,
			},
			{
				path:     []string{"/bin"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/bin"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/bin", "/foo"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/bin", "/foo"},
				mode:     soytrie.ModeExact,
				expected: true,
			},
			{
				path:     []string{"/bin", "/bar"},
				mode:     soytrie.ModePrefix,
				expected: false,
			},
			{
				path:     []string{"/bin", "/bar"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
		},
		dirSrc: {
			{
				path:     []string{"/testdata"},
				mode:     soytrie.ModeExact,
				expected: true,
			},
			{
				path:     []string{"/testdata"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/testdata", "/race"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/testdata", "/race"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/testdata", "/race", "/7"},
				mode:     soytrie.ModeExact,
				expected: true,
			},
			{
				path:     []string{"/testdata", "/race", "/7"},
				mode:     soytrie.ModePrefix,
				expected: true,
			},
			{
				path:     []string{"/testdata", "/race", "/no_such"},
				mode:     soytrie.ModeExact,
				expected: false,
			},
			{
				path:     []string{"/testdata", "/race", "/no_such"},
				mode:     soytrie.ModePrefix,
				expected: false,
			},
		},
	}

	for root := range tests {
		rootTests := tests[root]
		for i := range rootTests {
			tc := &rootTests[i]
			actual := root.Search(tc.mode, tc.path...)
			if actual != tc.expected {
				t.Logf("root=%+v", root)
				t.Fatalf("unexpected value %v for tests[%d] with path=%v,mode=%d,expected=%v", actual, i, tc.path, tc.mode, tc.expected)
			}
		}
	}
}

func TestCollectChildren(t *testing.T) {
	root := soytrie.New[int, string]()
	_ = root.Insert("1", 1)
	_ = root.Insert("1,2", 1, 2)
	_ = root.Insert("1,2,3", 1, 2, 3)
	_ = root.Insert("2", 2)
	_ = root.Insert("2,3", 2, 3)
	_ = root.Insert("10,20,30,40,50", 10, 20, 30, 40, 50)

	t.Run("CollectChildren", func(t *testing.T) {
		collector := []*soytrie.Node[int, string]{}
		soytrie.CollectChildren(root, &collector)
		if l := len(collector); l != 11 { // 10 nodes plus root
			t.Fatalf("unexpected length of collected result %d", l)
		}
	})
	t.Run("CollectChildrenValued", func(t *testing.T) {
		collector := []*soytrie.Node[int, string]{}
		soytrie.CollectChildrenValued(root, &collector)
		if l := len(collector); l != 6 { // 6 nodes (root is not valued)
			t.Fatalf("unexpected length of collected result %d", l)
		}

	})
	t.Run("CollectChildrenLeaf", func(t *testing.T) {
		collector := []*soytrie.Node[int, string]{}
		soytrie.CollectChildrenLeaf(root, &collector)
		if l := len(collector); l != 3 {
			t.Fatalf("unexpected length of collected result %d", l)
		}
	})
}

func TestPredict(t *testing.T) {
	root := soytrie.New[int, string]()
	_, _ = root.InsertNoOverwrite("1", 1)
	_, _ = root.InsertNoOverwrite("1,2", 1, 2)
	_, _ = root.InsertNoOverwrite("1,2,3", 1, 2, 3)
	_, _ = root.InsertNoOverwrite("2", 2)
	_, _ = root.InsertNoOverwrite("2,3", 2, 3)
	_, _ = root.InsertNoOverwrite("10,20,30,40,50", 10, 20, 30, 40, 50)

	type testCase struct {
		path        []int
		mode        soytrie.Mode
		expectedOk  bool
		expectedLen int // len
	}

	tests := []testCase{
		{
			path:        []int{1},
			mode:        soytrie.ModePrefix,
			expectedOk:  true,
			expectedLen: 3,
		},
		{
			path:        []int{1},
			mode:        soytrie.ModeExact,
			expectedOk:  true,
			expectedLen: 3,
		},
		{
			path:        []int{1, 2},
			mode:        soytrie.ModePrefix,
			expectedOk:  true,
			expectedLen: 2,
		},
		{
			path:        []int{1, 2},
			mode:        soytrie.ModeExact,
			expectedOk:  true,
			expectedLen: 2,
		},
		{
			path:        []int{2},
			mode:        soytrie.ModeExact,
			expectedOk:  true,
			expectedLen: 2,
		},
		{
			path:        []int{2},
			mode:        soytrie.ModePrefix,
			expectedOk:  true,
			expectedLen: 2,
		},
		{
			path:        []int{2, 3},
			mode:        soytrie.ModeExact,
			expectedOk:  true,
			expectedLen: 1,
		},
		{
			path:        []int{2, 3},
			mode:        soytrie.ModePrefix,
			expectedOk:  true,
			expectedLen: 1,
		},
		{
			path:        []int{2, 3, 4},
			mode:        soytrie.ModeExact,
			expectedOk:  false,
			expectedLen: 0,
		},
		{
			path:        []int{2, 3, 4},
			mode:        soytrie.ModePrefix,
			expectedOk:  false,
			expectedLen: 0,
		},
		{
			path:        []int{2, 5},
			mode:        soytrie.ModeExact,
			expectedOk:  false,
			expectedLen: 0,
		},
		{
			path:        []int{2, 5},
			mode:        soytrie.ModePrefix,
			expectedOk:  false,
			expectedLen: 0,
		},
		{
			path:        []int{10, 20},
			mode:        soytrie.ModeExact,
			expectedOk:  true,
			expectedLen: 1,
		},
		{
			path:        []int{10, 20},
			mode:        soytrie.ModePrefix,
			expectedOk:  true,
			expectedLen: 4,
		},
	}

	for i := range tests {
		tc := &tests[i]
		actual, ok := root.Predict(tc.mode, tc.path...)
		if ok != tc.expectedOk {
			t.Fatalf("[case %d] unexpected ok, expecting=%v, actual=%v", i, tc.expectedOk, actual)
		}
		if len(actual) != tc.expectedLen {
			t.Fatalf("[case %d] unexpected value %d, expecting %d", i, len(actual), tc.expectedLen)
		}
	}
}

func TestUnique(t *testing.T) {
	root := soytrie.New[int, string]()
	_ = root.Insert("1,2,3", 1, 2, 3)
	_ = root.Insert("1,2,2", 1, 2, 2)
	_ = root.Insert("1,2,7", 1, 2, 7)
	_ = root.Insert("1,3,7", 1, 3, 7)
	_ = root.Insert("1,3,8", 1, 3, 8)
	_ = root.Insert("1,10,20", 1, 10, 20)
	_ = root.Insert("0,2,3", 0, 2, 3)

	type testCase struct {
		path     []int
		expected bool
	}

	tests := []testCase{
		{
			path:     []int{1},
			expected: false,
		},
		{
			path:     []int{1, 2},
			expected: false,
		},
		{
			path:     []int{1, 2, 7},
			expected: true,
		},
		{
			path:     []int{1, 2, 7, 8},
			expected: false,
		},
		{
			path:     []int{1, 10},
			expected: true,
		},
		{
			path:     []int{1, 10, 20},
			expected: true,
		},
		{
			path:     []int{0},
			expected: true,
		},
		{
			path:     []int{0, 2},
			expected: true,
		},
		{
			path:     []int{0, 2, 3},
			expected: true,
		},
		{
			path:     []int{0, 2, 1, 3},
			expected: false,
		},
	}

	for i := range tests {
		tc := &tests[i]
		actual := root.Unique(tc.path...)
		if actual != tc.expected {
			t.Fatalf("[case %d] unexpected value %v, expecting %v with path=%v", i, actual, tc.expected, tc.path)
		}
	}
}

func pront[K comparable, V any](n *soytrie.Node[K, V]) {
	curr := n
	for p, c := range curr.Children {
		fmt.Println("p", p, "child", c.Value)
		pront(c)
	}
}
