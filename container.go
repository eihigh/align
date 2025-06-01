package align

import "github.com/eihigh/ng"

// Node represents a node in a tree structure.
// The leaf rectangle [Rect] and the container [Slice]
// and [Map] implement this interface.
type Node[S ng.Scalar] interface {
	Bounds() *Rect[S]
	Shift(Point[S])
}

// Slice is a slice of nodes that can be aligned and moved together.
type Slice[S ng.Scalar] []Node[S]

// Bounds returns the bounding rectangle that contains all nodes in the slice.
func (s Slice[S]) Bounds() *Rect[S] {
	if len(s) == 0 {
		return &Rect[S]{}
	}
	r := s[0].Bounds().Clone()
	for _, n := range s[1:] {
		r.Min.X = min(r.Min.X, n.Bounds().Min.X)
		r.Min.Y = min(r.Min.Y, n.Bounds().Min.Y)
		r.Max.X = max(r.Max.X, n.Bounds().Max.X)
		r.Max.Y = max(r.Max.Y, n.Bounds().Max.Y)
	}
	return r
}

// Shift moves all nodes in the slice by the given offset.
func (s Slice[S]) Shift(p Point[S]) {
	for _, n := range s {
		n.Shift(p)
	}
}

// Add translates all nodes in the slice by p and returns the slice.
func (s Slice[S]) Add(p Point[S]) Slice[S] {
	s.Shift(p)
	return s
}

// Last returns the last node in the slice, or nil if the slice is empty.
func (s Slice[S]) Last() Node[S] {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

// Align positions the slice relative to the target using anchor points.
// (ax, ay) is the anchor point on the slice (0-1), and (tax, tay) is the anchor point on the target.
func (s Slice[S]) Align(ax, ay float64, target Node[S], tax, tay float64) Slice[S] {
	tb := target.Bounds()
	sb := s.Bounds()
	d := tb.Anchor(tax, tay).Sub(sb.Anchor(ax, ay))
	s.Shift(d)
	return s
}

// Nest positions the slice within the target at the given relative position.
func (s Slice[S]) Nest(target Node[S], ax, ay float64) Slice[S] {
	return s.Align(ax, ay, target, ax, ay)
}

// CenterOf centers the slice within the target.
func (s Slice[S]) CenterOf(target Node[S]) Slice[S] {
	return s.Align(0.5, 0.5, target, 0.5, 0.5)
}

// StackX stacks the slice horizontally relative to the target.
func (s Slice[S]) StackX(target Node[S], tax, tay float64) Slice[S] {
	return s.Align(1-tax, tay, target, tax, tay)
}

// StackY stacks the slice vertically relative to the target.
func (s Slice[S]) StackY(target Node[S], tax, tay float64) Slice[S] {
	return s.Align(tax, 1-tay, target, tax, tay)
}

// Clamp constrains the slice within the target bounds while keeping its size.
func (s Slice[S]) Clamp(target Node[S]) Slice[S] {
	oldPos := s.Bounds().Min
	b := s.Bounds().Clamp(target.Bounds())
	s.Shift(b.Min.Sub(oldPos))
	return s
}

// Inset returns the bounding rectangle inset by n.
func (s Slice[S]) Inset(n S) *Rect[S] {
	return s.Bounds().Inset(n)
}

// InsetXY returns the bounding rectangle inset by x horizontally and y vertically.
func (s Slice[S]) InsetXY(x, y S) *Rect[S] {
	return s.Bounds().InsetXY(x, y)
}

// InsetLTRB returns the bounding rectangle inset by the given amounts on each side.
func (s Slice[S]) InsetLTRB(left, top, right, bottom S) *Rect[S] {
	return s.Bounds().InsetLTRB(left, top, right, bottom)
}

// Outset returns the bounding rectangle expanded by n.
func (s Slice[S]) Outset(n S) *Rect[S] {
	return s.Bounds().Outset(n)
}

// OutsetXY returns the bounding rectangle expanded by x horizontally and y vertically.
func (s Slice[S]) OutsetXY(x, y S) *Rect[S] {
	return s.Bounds().OutsetXY(x, y)
}

// OutsetLTRB returns the bounding rectangle expanded by the given amounts on each side.
func (s Slice[S]) OutsetLTRB(left, top, right, bottom S) *Rect[S] {
	return s.Bounds().OutsetLTRB(left, top, right, bottom)
}

// Map is a map of named nodes that can be aligned and moved together.
type Map[S ng.Scalar] map[string]Node[S]

// Bounds returns the bounding rectangle that contains all nodes in the map.
func (m Map[S]) Bounds() *Rect[S] {
	if len(m) == 0 {
		return &Rect[S]{}
	}
	var r *Rect[S]
	for _, n := range m {
		if r == nil {
			r = n.Bounds().Clone()
		} else {
			r.Min.X = min(r.Min.X, n.Bounds().Min.X)
			r.Min.Y = min(r.Min.Y, n.Bounds().Min.Y)
			r.Max.X = max(r.Max.X, n.Bounds().Max.X)
			r.Max.Y = max(r.Max.Y, n.Bounds().Max.Y)
		}
	}
	return r
}

// Shift moves all nodes in the map by the given offset.
func (m Map[S]) Shift(p Point[S]) {
	for _, n := range m {
		n.Shift(p)
	}
}

// Add translates all nodes in the map by p and returns the map.
func (m Map[S]) Add(p Point[S]) Map[S] {
	m.Shift(p)
	return m
}

// Align positions the map relative to the target using anchor points.
// (ax, ay) is the anchor point on the map (0-1), and (tax, tay) is the anchor point on the target.
func (m Map[S]) Align(ax, ay float64, target Node[S], tax, tay float64) Map[S] {
	tb := target.Bounds()
	mb := m.Bounds()
	d := tb.Anchor(tax, tay).Sub(mb.Anchor(ax, ay))
	m.Shift(d)
	return m
}

// Nest positions the map within the target at the given relative position.
func (m Map[S]) Nest(target Node[S], ax, ay float64) Map[S] {
	return m.Align(ax, ay, target, ax, ay)
}

// StackX stacks the map horizontally relative to the target.
func (m Map[S]) StackX(target Node[S], tax, tay float64) Map[S] {
	return m.Align(1-tax, tay, target, tax, tay)
}

// StackY stacks the map vertically relative to the target.
func (m Map[S]) StackY(target Node[S], tax, tay float64) Map[S] {
	return m.Align(tax, 1-tay, target, tax, tay)
}

// Inset returns the bounding rectangle inset by n.
func (m Map[S]) Inset(n S) *Rect[S] {
	return m.Bounds().Inset(n)
}

// InsetXY returns the bounding rectangle inset by x horizontally and y vertically.
func (m Map[S]) InsetXY(x, y S) *Rect[S] {
	return m.Bounds().InsetXY(x, y)
}

// InsetLTRB returns the bounding rectangle inset by the given amounts on each side.
func (m Map[S]) InsetLTRB(left, top, right, bottom S) *Rect[S] {
	return m.Bounds().InsetLTRB(left, top, right, bottom)
}

// Outset returns the bounding rectangle expanded by n.
func (m Map[S]) Outset(n S) *Rect[S] {
	return m.Bounds().Outset(n)
}

// OutsetXY returns the bounding rectangle expanded by x horizontally and y vertically.
func (m Map[S]) OutsetXY(x, y S) *Rect[S] {
	return m.Bounds().OutsetXY(x, y)
}

// OutsetLTRB returns the bounding rectangle expanded by the given amounts on each side.
func (m Map[S]) OutsetLTRB(left, top, right, bottom S) *Rect[S] {
	return m.Bounds().OutsetLTRB(left, top, right, bottom)
}
