package align

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

const (
	mid    = 0.5
	top    = 0.0
	bottom = 1.0
	left   = 0.0
	right  = 1.0
)

var (
	update = flag.Bool("update", false, "update test images")

	red     = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green   = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue    = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	cyan    = color.RGBA{R: 0, G: 255, B: 255, A: 255}
	yellow  = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	magenta = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	gray    = color.RGBA{R: 128, G: 128, B: 128, A: 255}
)

type testImg struct {
	name string
	img  image.RGBA
}

func newTestImg(t *testing.T, name string, w, h int) *testImg {
	t.Helper()
	return &testImg{
		name: name,
		img:  *image.NewRGBA(image.Rect(0, 0, w, h)),
	}
}

func (i *testImg) fill(n *Node[int], color color.Color) {
	for y := range n.Size.Y {
		for x := range n.Size.X {
			i.img.Set(n.Pos.X+x, n.Pos.Y+y, color)
		}
	}
}

func (i *testImg) save(t *testing.T) {
	t.Helper()
	f, err := os.Create(filepath.Join("testdata", i.name+".png"))
	if err != nil {
		t.Fatalf("failed to create image file: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, &i.img); err != nil {
		t.Fatalf("failed to encode image: %v", err)
	}
	t.Logf("saved image: %s", filepath.Join("testdata", i.name+".png"))
}

func (i *testImg) test(t *testing.T) {
	if *update {
		i.save(t)
		return
	}

	f, err := os.Open(filepath.Join("testdata", i.name+".png"))
	if err != nil {
		t.Fatalf("failed to open image file: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode image: %v", err)
	}
	i.compare(t, img)
}

func (i *testImg) compare(t *testing.T, other image.Image) {
	t.Helper()
	if i.img.Bounds() != other.Bounds() {
		t.Fatalf("image bounds mismatch: %v != %v", i.img.Bounds(), other.Bounds())
	}
	err := false
	for y := 0; y < i.img.Bounds().Dy(); y++ {
		for x := 0; x < i.img.Bounds().Dx(); x++ {
			if !compareColor(i.img.At(x, y), other.At(x, y)) {
				err = true
			}
		}
	}
	t.Errorf("image comparison failed: %v", err)
}

func compareColor(a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func TestSimple(t *testing.T) {
	screen := WH(100, 100)
	n := screen.Inset(20)

	img := newTestImg(t, "simple", 100, 100)
	img.fill(screen, blue)
	img.fill(n, red)
	img.test(t)
}

func TestComplex(t *testing.T) {
	screen := WH(100, 100)
	portrait := WH(20, 30).Nest(screen.Inset(10), left, top)
	btns := screen.InsetXY(20, 10)
	okBtn := WH(20, 10).Nest(btns, left, bottom)
	cancelBtn := WH(20, 10).Nest(btns, right, bottom)

	img := newTestImg(t, "complex", 100, 100)
	img.fill(screen, blue)
	img.fill(portrait, red)
	img.fill(okBtn, green)
	img.fill(cancelBtn, green)
	img.test(t)
}

// func testNodes(t *testing.T, n int) []*Node[int] {
// 	t.Helper()
// 	nodes := make([]*Node[int], n)
// 	for i := range n {
// 		nodes[i] = XYWH(i*10, i*10, 10, 10)
// 	}
// 	return nodes
// }

// func TestLink(t *testing.T) {
// 	t.Run("basic linking", func(t *testing.T) {
// 		nodes := testNodes(t, 3)
// 		n1, n2, n3 := nodes[0], nodes[1], nodes[2]

// 		n1.Link(n2, n3)

// 		// Check forward links
// 		if n1.next != n2 {
// 			t.Errorf("n1.next should be n2")
// 		}
// 		if n2.next != n3 {
// 			t.Errorf("n2.next should be n3")
// 		}
// 		if n3.next != nil {
// 			t.Errorf("n3.next should be nil")
// 		}

// 		// Check backward links
// 		if n1.prev != nil {
// 			t.Errorf("n1.prev should be nil")
// 		}
// 		if n2.prev != n1 {
// 			t.Errorf("n2.prev should be n1")
// 		}
// 		if n3.prev != n2 {
// 			t.Errorf("n3.prev should be n2")
// 		}
// 	})

// 	t.Run("linking to existing list", func(t *testing.T) {
// 		nodes := testNodes(t, 4)
// 		n1, n2, n3, n4 := nodes[0], nodes[1], nodes[2], nodes[3]

// 		n1.Link(n2)
// 		n1.Link(n3, n4)

// 		// Check that n3 and n4 are added to the end
// 		if n2.next != n3 {
// 			t.Errorf("n2.next should be n3")
// 		}
// 		if n3.next != n4 {
// 			t.Errorf("n3.next should be n4")
// 		}
// 		if n4.next != nil {
// 			t.Errorf("n4.next should be nil")
// 		}
// 	})

// 	t.Run("skip nodes already in a list", func(t *testing.T) {
// 		nodes := testNodes(t, 4)
// 		n1, n2, n3, n4 := nodes[0], nodes[1], nodes[2], nodes[3]

// 		// Create a separate list
// 		n3.Link(n2, n4)

// 		// Try to link n3 and n4 to n1's list
// 		n1.Link(n2, n4, n3)

// 		// want: n1 <-> n3 <-> n2 <-> n4 (n3のリストがn1に連結される)
// 		testNodeLinking(t, []int{0, 20, 10, 30}, n1)
// 	})

// 	t.Run("skip duplicate nodes", func(t *testing.T) {
// 		nodes := testNodes(t, 3)
// 		n1, n2, n3 := nodes[0], nodes[1], nodes[2]

// 		n1.Link(n2)
// 		n1.Link(n2, n3) // n2 is already in the list

// 		// n2 should not be duplicated
// 		testNodeLinkingFromNodes(t, nodes[:3], n1)
// 	})

// 	t.Run("handle nil nodes", func(t *testing.T) {
// 		nodes := testNodes(t, 2)
// 		n1, n2 := nodes[0], nodes[1]

// 		n1.Link(nil, n2, nil)

// 		// Only n2 should be linked
// 		testNodeLinkingFromNodes(t, nodes[:2], n1)
// 	})

// 	t.Run("empty link call", func(t *testing.T) {
// 		nodes := testNodes(t, 1)
// 		n1 := nodes[0]
// 		n1.Link() // No nodes to link

// 		if n1.prev != nil || n1.next != nil {
// 			t.Errorf("n1 should have no links")
// 		}
// 	})

// 	t.Run("single node remains isolated", func(t *testing.T) {
// 		nodes := testNodes(t, 1)
// 		n1 := nodes[0]

// 		// A single node should have both prev and next as nil
// 		if n1.prev != nil {
// 			t.Errorf("single node should have prev = nil")
// 		}
// 		if n1.next != nil {
// 			t.Errorf("single node should have next = nil")
// 		}
// 	})
// }

// func testNodeLinking(t *testing.T, xs []int, n *Node[int]) {
// 	t.Helper()
// 	i := 0
// 	for current := range n.Linked() {
// 		if i >= len(xs) {
// 			t.Errorf("more nodes than expected")
// 			break
// 		}
// 		if current.Pos.X != xs[i] {
// 			t.Errorf("want %d, got %d", xs[i], current.Pos.X)
// 		}
// 		i++
// 	}
// }

// // testNodeLinkingFromNodes は testNodes で生成したノードを使って連結をテストする
// // ノードのPos.Xがインデックス*10になっていることを利用
// func testNodeLinkingFromNodes(t *testing.T, nodes []*Node[int], start *Node[int]) {
// 	t.Helper()
// 	expected := make([]int, len(nodes))
// 	for i, node := range nodes {
// 		expected[i] = node.Pos.X
// 	}
// 	testNodeLinking(t, expected, start)
// }

// // Helper function to traverse and collect nodes in a list
// func collectNodes[S ng.Scalar](start *Node[S]) []*Node[S] {
// 	return slices.Collect(start.Linked())
// }

// func TestCollectNodes(t *testing.T) {
// 	nodes := testNodes(t, 3)
// 	n1, n2, n3 := nodes[0], nodes[1], nodes[2]

// 	n1.Link(n2, n3)

// 	collected := collectNodes(n2) // Start from middle
// 	if len(collected) != 3 {
// 		t.Errorf("expected 3 nodes, got %d", len(collected))
// 	}
// 	if collected[0] != n1 || collected[1] != n2 || collected[2] != n3 {
// 		t.Errorf("incorrect node order")
// 	}
// }
