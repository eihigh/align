package align

import (
	"fmt"
	"image"
	"iter"

	"github.com/eihigh/ng"
)

type Point[S ng.Scalar] struct {
	X, Y S
}

// XY is shorthand for Point{X, Y}.
func XY[S ng.Scalar](x, y S) Point[S] {
	return Point[S]{X: x, Y: y}
}

func (p Point[S]) XY() (S, S) {
	return p.X, p.Y
}

// String returns a string representation of p like "(3,4)".
func (p Point[S]) String() string {
	return "(" + fmt.Sprint(p.X) + "," + fmt.Sprint(p.Y) + ")"
}

// Add returns the vector p+q.
func (p Point[S]) Add(q Point[S]) Point[S] {
	return Point[S]{X: p.X + q.X, Y: p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point[S]) Sub(q Point[S]) Point[S] {
	return Point[S]{X: p.X - q.X, Y: p.Y - q.Y}
}

func (p Point[S]) Mul(q Point[S]) Point[S] {
	return Point[S]{X: p.X * q.X, Y: p.Y * q.Y}
}

func (p Point[S]) Div(q Point[S]) Point[S] {
	return Point[S]{X: p.X / q.X, Y: p.Y / q.Y}
}

func (p Point[S]) Scale(k S) Point[S] {
	return Point[S]{X: p.X * k, Y: p.Y * k}
}

// In reports whether p is in r.
// func (p Point[S]) In(r Rect[S]) bool {
// 	return r.Min.X <= p.X && p.X < r.Max.X &&
// 		r.Min.Y <= p.Y && p.Y < r.Max.Y
// }

// TODO: Point.Mod?

// Eq reports whether p and q are equal.
func (p Point[S]) Eq(q Point[S]) bool {
	return p.X == q.X && p.Y == q.Y
}

// Image returns the point as an image.Point.
func (p Point[S]) Image() image.Point {
	return image.Point{X: int(p.X), Y: int(p.Y)}
}

// Int returns the point as an int point.
func (p Point[S]) Int() Point[int] {
	return Point[int]{X: int(p.X), Y: int(p.Y)}
}

// Float64 returns the point as a float64 point.
func (p Point[S]) Float64() Point[float64] {
	return Point[float64]{X: float64(p.X), Y: float64(p.Y)}
}

// Float32 returns the point as a float32 point.
func (p Point[S]) Float32() Point[float32] {
	return Point[float32]{X: float32(p.X), Y: float32(p.Y)}
}

type Node[S ng.Scalar] struct {
	Pos, Size  Point[S]
	prev, next *Node[S] // linked list
}

func XYWH[S ng.Scalar](x, y, w, h S) *Node[S] {
	return &Node[S]{
		Pos:  XY(x, y),
		Size: XY(w, h),
	}
}

func XYXY[S ng.Scalar](x1, y1, x2, y2 S) *Node[S] {
	return &Node[S]{
		Pos:  XY(x1, y1),
		Size: XY(x2-x1, y2-y1),
	}
}

func PosSize[S ng.Scalar](pos, size Point[S]) *Node[S] {
	return &Node[S]{Pos: pos, Size: size}
}

func PosWH[S ng.Scalar](pos Point[S], w, h S) *Node[S] {
	return &Node[S]{Pos: pos, Size: XY(w, h)}
}

func XYSize[S ng.Scalar](x, y S, size Point[S]) *Node[S] {
	return &Node[S]{Pos: XY(x, y), Size: size}
}

func WH[S ng.Scalar](w, h S) *Node[S] {
	return &Node[S]{Size: XY(w, h)}
}

func Union[S ng.Scalar](nodes ...*Node[S]) *Node[S] {
	if len(nodes) == 0 {
		return nil
	}
	n := &Node[S]{}
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if node.Empty() {
			continue
		}
		n.Pos.X = min(n.Pos.X, node.Pos.X)
		n.Pos.Y = min(n.Pos.Y, node.Pos.Y)
		n.Size.X = max(n.Size.X, node.Pos.X+node.Size.X-n.Pos.X)
		n.Size.Y = max(n.Size.Y, node.Pos.Y+node.Size.Y-n.Pos.Y)
		n.Link(node)
	}
	return n
}

func (n *Node[S]) String() string {
	return fmt.Sprintf("Pos: %s, Size: %s", n.Pos.String(), n.Size.String())
}

func (n *Node[S]) Empty() bool {
	return n == nil || (n.Size.X <= 0 && n.Size.Y <= 0)
}

func (n *Node[S]) In(other *Node[S]) bool {
	p0 := n.Pos
	p1 := n.Pos.Add(n.Size)
	q0 := other.Pos
	q1 := other.Pos.Add(other.Size)
	return q0.X <= p0.X && p1.X <= q1.X &&
		q0.Y <= p0.Y && p1.Y <= q1.Y
}

func (n *Node[S]) Link(other *Node[S]) (merged bool) {
	// Find the tail of n's list
	tail := n
	for tail.next != nil {
		tail = tail.next
	}

	// Link each node or merge their lists
	if other == nil {
		return false
	}

	// Check if the node is already in n's list
	for current := n; current != nil; current = current.next {
		if current == other {
			return false
		}
	}

	// Find the head of other's list
	head := other
	for head.prev != nil {
		head = head.prev
	}

	// Disconnect the other's list from any existing connections
	// and connect it to our tail
	tail.next = head
	head.prev = tail

	// Find the new tail
	for tail.next != nil {
		tail = tail.next
	}

	return true
}

func (n *Node[S]) Unlink() {
	// Unlink this node from its list
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	// Clear the links
	n.prev = nil
	n.next = nil
}

func (n *Node[S]) Linked() iter.Seq[*Node[S]] {
	return func(yield func(*Node[S]) bool) {
		head := n
		for head.prev != nil {
			head = head.prev
		}
		for current := head; current != nil; current = current.next {
			if !yield(current) {
				return
			}
		}
	}
}

func (n *Node[S]) Anchor(anchorX, anchorY float64) Point[S] {
	offset := XY(S(float64(n.Size.X)*anchorX), S(float64(n.Size.Y)*anchorY))
	return n.Pos.Add(offset)
}

func (n *Node[S]) Max() Point[S] {
	return n.Pos.Add(n.Size)
}

// 相対位置を保ったままリンクされたノード全てを移動する
func (n *Node[S]) MoveTo(p Point[S]) *Node[S] {
	delta := p.Sub(n.Pos)
	n.Pos = p
	for l := range n.Linked() {
		l.Pos = l.Pos.Add(delta)
	}
	return n
}

// リンクされたノード全てを移動する
func (n *Node[S]) add(p Point[S]) *Node[S] {
	for l := range n.Linked() {
		l.Pos = l.Pos.Add(p)
	}
	return n
}

func (n *Node[S]) Offset(p Point[S]) *Node[S] {
	// 自分自身の位置を移動する
	n.Pos = n.Pos.Add(p)
	return n
}

func (n *Node[S]) Sub(p Point[S]) *Node[S] {
	for l := range n.Linked() {
		l.Pos = l.Pos.Sub(p)
	}
	return n
}

func (n *Node[S]) Outset(s S) *Node[S] {
	return n.OutsetLTRB(s, s, s, s)
}

func (n *Node[S]) OutsetXY(x, y S) *Node[S] {
	return n.OutsetLTRB(x, y, x, y)
}

func (n *Node[S]) OutsetLTRB(l, t, r, b S) *Node[S] {
	if n == nil {
		return nil
	}
	o := XYWH(n.Pos.X-l, n.Pos.Y-t, n.Size.X+r+l, n.Size.Y+b+t)
	o.Link(n)
	return o
}

func (n *Node[S]) Inset(s S) *Node[S] {
	return n.InsetLTRB(s, s, s, s)
}

func (n *Node[S]) InsetXY(x, y S) *Node[S] {
	return n.InsetLTRB(x, y, x, y)
}

// InsetLTRB returns the rectangle r inset by left, top, right, and bottom.
// If either of r's dimensions is less than left+right or top+bottom then an
// empty rectangle near the center of r will be returned.
func (n *Node[S]) InsetLTRB(l, t, r, b S) *Node[S] {
	if n == nil {
		return nil
	}
	o := XYWH(n.Pos.X+l, n.Pos.Y+t, n.Size.X-r-l, n.Size.Y-b-t)
	o.Link(n)
	return o
}

func (n *Node[S]) Align(ax, ay float64, target *Node[S], targetAx, targetAy float64) *Node[S] {
	d := target.Anchor(targetAx, targetAy).Sub(n.Anchor(ax, ay))
	sep := true
	for l := range n.Linked() {
		if l == target {
			sep = false // 自分自身はリンクしない
			break
		}
	}
	if sep {
		n.add(d)
		n.Link(target)
	} else {
		n.Pos = n.Pos.Add(d)
	}
	return n
}

func (n *Node[S]) AlignSelf(ax, ay float64, target *Node[S], targetAx, targetAy float64) *Node[S] {
	if n == nil || target == nil {
		return n
	}

	// 位置を計算する
	nPos := n.Anchor(ax, ay)
	tPos := target.Anchor(targetAx, targetAy)

	// 移動する（自分だけ）
	delta := tPos.Sub(nPos)
	n.Pos = n.Pos.Add(delta)

	// リンクしない
	// n.Link(target)

	return n
}

func (n *Node[S]) Nest(target *Node[S], ax, ay float64) *Node[S] {
	return n.Align(ax, ay, target, ax, ay)
}

func (n *Node[S]) CenterOf(target *Node[S]) *Node[S] {
	return n.Align(.5, .5, target, 0.5, 0.5)
}

func (n *Node[S]) StackX(target *Node[S], targetAx, targetAy float64) *Node[S] {
	return n.Align(1-targetAx, targetAy, target, targetAx, targetAy)
}

func (n *Node[S]) StackY(target *Node[S], targetAx, targetAy float64) *Node[S] {
	return n.Align(targetAx, 1-targetAy, target, targetAx, targetAy)
}

func (n *Node[S]) Clamp(target *Node[S]) *Node[S] {
	// 位置を計算する
	p := n.Pos
	targetMax := target.Max()
	p.X = max(p.X, target.Pos.X)
	p.Y = max(p.Y, target.Pos.Y)
	p.X = min(p.X, targetMax.X-n.Size.X)
	p.Y = min(p.Y, targetMax.Y-n.Size.Y)
	n.Pos = p

	return n
}

type Wrapper[S ng.Scalar] struct {
	bounds      *Node[S]            // bounds of the wrapper
	ax, ay      float64             // alignment factors
	stack, wrap func(a, b *Node[S]) // functions to stack or wrap nodes
	nodes       []*Node[S]          // linked list of nodes
	lineFirst   *Node[S]            // first node in the current line
}

func NewWrapper[S ng.Scalar](bounds *Node[S], ax, ay float64,
	stack, wrap func(a, b *Node[S])) *Wrapper[S] {
	return &Wrapper[S]{
		bounds: bounds,
		ax:     ax,
		ay:     ay,
		stack:  stack,
		wrap:   wrap,
	}
}

func (w *Wrapper[S]) Add(n *Node[S]) bool {
	if len(w.nodes) == 0 {
		n.Nest(w.bounds, w.ax, w.ay)
		if !n.In(w.bounds) {
			n.Unlink() // Unlink if it doesn't fit
			return false
		}
		w.lineFirst = n
	} else {
		w.stack(n, w.nodes[len(w.nodes)-1])
		if !n.In(w.bounds) {
			w.wrap(n, w.lineFirst)
			if !n.In(w.bounds) {
				n.Unlink() // Unlink if it doesn't fit
				return false
			}
			w.lineFirst = n
		}
	}
	w.nodes = append(w.nodes, n)
	return true
}

func (w *Wrapper[S]) Nodes() []*Node[S] { return w.nodes }
