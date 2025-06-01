package align

import (
	"image"
	"iter"
	"math"

	"github.com/eihigh/ng"
)

// --------------------------------------
// Standard functions
// --------------------------------------

// A Rect contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
type Rect[S ng.Scalar] struct {
	Min, Max Point[S]
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r *Rect[S]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r *Rect[S]) Dx() S {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r *Rect[S]) Dy() S {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r *Rect[S]) Size() Point[S] {
	return Point[S]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add translates and returns the rectangle r by p.
func (r *Rect[S]) Add(p Point[S]) *Rect[S] {
	r.Min.X += p.X
	r.Min.Y += p.Y
	r.Max.X += p.X
	r.Max.Y += p.Y
	return r
}

// Sub translates and returns the rectangle r by -p.
func (r *Rect[S]) Sub(p Point[S]) *Rect[S] {
	r.Min.X -= p.X
	r.Min.Y -= p.Y
	r.Max.X -= p.X
	r.Max.Y -= p.Y
	return r
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r *Rect[S]) Inset(n S) *Rect[S] {
	return r.InsetLTRB(n, n, n, n)
}

// Rect.Intersect is removed

// Rect.Union is removed

// Empty reports whether the rectangle contains no points.
func (r *Rect[S]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r *Rect[S]) Eq(s *Rect[S]) bool {
	return (r.Min.Eq(s.Min) && r.Max.Eq(s.Max)) || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r *Rect[S]) Overlaps(s *Rect[S]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r *Rect[S]) In(s Node[S]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	sb := s.Bounds()
	return sb.Min.X <= r.Min.X && r.Max.X <= sb.Max.X &&
		sb.Min.Y <= r.Min.Y && r.Max.Y <= sb.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r *Rect[S]) Canon() *Rect[S] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// --------------------------------------
// Constructors and Converters
// --------------------------------------

// XYXY is shorthand for &[Rect]{XY(x0, y0), XY(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func XYXY[S ng.Scalar](x0, y0, x1, y1 S) *Rect[S] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return &Rect[S]{Point[S]{x0, y0}, Point[S]{x1, y1}}
}

// XYWH creates a rectangle from position (x,y) and dimensions (w,h).
func XYWH[S ng.Scalar](x, y, w, h S) *Rect[S] {
	return &Rect[S]{
		Min: Point[S]{x, y},
		Max: Point[S]{x + w, y + h},
	}
}

// PosSize creates a rectangle from position and size points.
func PosSize[S ng.Scalar](pos, size Point[S]) *Rect[S] {
	return &Rect[S]{
		Min: Point[S]{pos.X, pos.Y},
		Max: Point[S]{pos.X + size.X, pos.Y + size.Y},
	}
}

// WH creates a rectangle from origin (0,0) with dimensions (w,h).
func WH[S ng.Scalar](w, h S) *Rect[S] {
	return &Rect[S]{
		Min: Point[S]{0, 0},
		Max: Point[S]{w, h},
	}
}

// Clone returns a copy of the rectangle.
func (r *Rect[S]) Clone() *Rect[S] {
	return &Rect[S]{
		Min: r.Min,
		Max: r.Max,
	}
}

// Int converts r to Rect[int].
func (r *Rect[S]) Int() *Rect[int] {
	return XYXY(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

// Float64 converts r to Rect[float64].
func (r *Rect[S]) Float64() *Rect[float64] {
	return XYXY(float64(r.Min.X), float64(r.Min.Y), float64(r.Max.X), float64(r.Max.Y))
}

// Float32 converts r to *Rect[float32].
func (r *Rect[S]) Float32() *Rect[float32] {
	return XYXY(float32(r.Min.X), float32(r.Min.Y), float32(r.Max.X), float32(r.Max.Y))
}

// Image converts r to an [image.Rectangle].
func (r *Rect[S]) Image() image.Rectangle {
	return image.Rect(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

// --------------------------------------
// Additional functions
// --------------------------------------

// Points returns an iterator that yields the points in r in row-major order.
func (r *Rect[S]) Points() iter.Seq[Point[int]] {
	return func(yield func(Point[int]) bool) {
		x0 := int(math.Ceil(float64(r.Min.X)))
		y0 := int(math.Ceil(float64(r.Min.Y)))
		// Note that the bounds are inclusive for Min and exclusive for Max.
		for y := y0; S(y) < r.Max.Y; y++ {
			for x := x0; S(x) < r.Max.X; x++ {
				if !yield(Point[int]{x, y}) {
					return
				}
			}
		}
	}
}

// Bounds implements [Node] interface.
func (r *Rect[S]) Bounds() *Rect[S] {
	return r
}

// Shift implements [Node] interface.
func (r *Rect[S]) Shift(p Point[S]) {
	r.Min.X += p.X
	r.Min.Y += p.Y
	r.Max.X += p.X
	r.Max.Y += p.Y
}

// Anchor returns a point at the relative position (ax, ay) within r.
// Values 0 and 1 represent Min and Max bounds respectively.
func (r *Rect[S]) Anchor(ax, ay float64) Point[S] {
	return Point[S]{
		X: r.Min.X + S(float64(r.Max.X-r.Min.X)*ax),
		Y: r.Min.Y + S(float64(r.Max.Y-r.Min.Y)*ay),
	}
}

// Align aligns r's anchor point (ax, ay) to target's anchor point (tax, tay).
func (r *Rect[S]) Align(ax, ay float64, target Node[S], tax, tay float64) *Rect[S] {
	tb := target.Bounds()
	d := tb.Anchor(tax, tay).Sub(r.Anchor(ax, ay))
	return r.Add(d)
}

// CenterOf centers r within target.
func (r *Rect[S]) CenterOf(target Node[S]) *Rect[S] {
	return r.Align(0.5, 0.5, target, 0.5, 0.5)
}

// Nest positions r within target at the relative position (ax, ay).
func (r *Rect[S]) Nest(target Node[S], ax, ay float64) *Rect[S] {
	return r.Align(ax, ay, target, ax, ay)
}

// StackX stacks r horizontally against target at the relative position (tax, tay).
func (r *Rect[S]) StackX(target Node[S], tax, tay float64) *Rect[S] {
	return r.Align(1-tax, tay, target, tax, tay)
}

// StackY stacks r vertically against target at the relative position (tax, tay).
func (r *Rect[S]) StackY(target Node[S], tax, tay float64) *Rect[S] {
	return r.Align(tax, 1-tay, target, tax, tay)
}

// Clamp constrains r to fit within s while preserving r's size.
func (r *Rect[S]) Clamp(s *Rect[S]) *Rect[S] {
	w, h := r.Dx(), r.Dy()
	r.Min.X = max(s.Min.X, min(s.Max.X-w, r.Min.X))
	r.Min.Y = max(s.Min.Y, min(s.Max.Y-h, r.Min.Y))
	r.Max.X = r.Min.X + w
	r.Max.Y = r.Min.Y + h
	return r
}

// InsetXY returns r inset by x horizontally and y vertically.
func (r *Rect[S]) InsetXY(x, y S) *Rect[S] {
	return r.InsetLTRB(x, y, x, y)
}

// InsetLTRB returns r inset by the specified amounts on each side.
// If the insets exceed the rectangle's dimensions, returns an empty rectangle.
func (r *Rect[S]) InsetLTRB(left, top, right, bottom S) *Rect[S] {
	s := &Rect[S]{Min: r.Min, Max: r.Max}
	if s.Dx() < left+right {
		t := float64(left) / float64(left+right)
		s.Min.X = s.Min.X + S(t*float64(s.Max.X-s.Min.X))
		s.Max.X = s.Min.X
	} else {
		s.Min.X += left
		s.Max.X -= right
	}
	if s.Dy() < top+bottom {
		t := float64(top) / float64(top+bottom)
		s.Min.Y = s.Min.Y + S(t*float64(s.Max.Y-s.Min.Y))
		s.Max.Y = s.Min.Y
	} else {
		s.Min.Y += top
		s.Max.Y -= bottom
	}
	return s
}

// Outset returns r expanded by n on all sides.
func (r *Rect[S]) Outset(n S) *Rect[S] {
	return r.InsetLTRB(-n, -n, -n, -n)
}

// OutsetXY returns r expanded by x horizontally and y vertically.
func (r *Rect[S]) OutsetXY(x, y S) *Rect[S] {
	return r.InsetLTRB(-x, -y, -x, -y)
}

// OutsetLTRB returns r expanded by the specified amounts on each side.
func (r *Rect[S]) OutsetLTRB(left, top, right, bottom S) *Rect[S] {
	return r.InsetLTRB(-left, -top, -right, -bottom)
}

// CutX divides the rectangle r into the left and right rectangles
// at the given width w.
func (r *Rect[S]) CutX(w S) (left, right *Rect[S]) {
	w = max(0, min(r.Dx(), w))
	left = XYXY(r.Min.X, r.Min.Y, r.Min.X+w, r.Max.Y)
	right = XYXY(r.Min.X+w, r.Min.Y, r.Max.X, r.Max.Y)
	return
}

// CutXByRate divides r horizontally at the given rate (0.0 to 1.0).
func (r *Rect[S]) CutXByRate(rate float64) (left, right *Rect[S]) {
	w := S(float64(r.Dx()) * rate)
	return r.CutX(w)
}

// CutY divides the rectangle r into the top and bottom rectangles
// at the given height h.
func (r *Rect[S]) CutY(h S) (top, bottom *Rect[S]) {
	h = max(0, min(r.Dy(), h))
	top = XYXY(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+h)
	bottom = XYXY(r.Min.X, r.Min.Y+h, r.Max.X, r.Max.Y)
	return
}

// CutYByRate divides r vertically at the given rate (0.0 to 1.0).
func (r *Rect[S]) CutYByRate(rate float64) (top, bottom *Rect[S]) {
	h := S(float64(r.Dy()) * rate)
	return r.CutY(h)
}

// Split divides the rectangle r into a grid of rectangles with xs columns and
// ys rows, with gaps of xGap and yGap between them.
func (r *Rect[S]) Split(xs, ys int, xGap, yGap S) Slice[S] {
	if xs <= 0 || ys <= 0 {
		return nil
	}
	if xs == 1 && ys == 1 {
		return Slice[S]{r}
	}

	w := (r.Dx() - S(xGap)*(S(xs)-1)) / S(xs)
	h := (r.Dy() - S(yGap)*(S(ys)-1)) / S(ys)

	rs := make(Slice[S], 0, xs*ys)
	for y := range ys {
		for x := range xs {
			minX := r.Min.X + S(x)*(w+xGap)
			minY := r.Min.Y + S(y)*(h+yGap)
			maxX := minX + w
			maxY := minY + h
			rs = append(rs, XYXY(minX, minY, maxX, maxY))
		}
	}
	return rs
}

// SplitX divides r into xs columns with xGap spacing between them.
func (r *Rect[S]) SplitX(xs int, xGap S) Slice[S] {
	return r.Split(xs, 1, xGap, 0)
}

// SplitY divides r into ys rows with yGap spacing between them.
func (r *Rect[S]) SplitY(ys int, yGap S) Slice[S] {
	return r.Split(1, ys, 0, yGap)
}

// Repeat creates a grid of copies of r, with xs columns and ys rows,
// with gaps of xGap and yGap between them.
func (r *Rect[S]) Repeat(xs, ys int, xGap, yGap S) Slice[S] {
	if xs <= 0 || ys <= 0 {
		return nil
	}
	if xs == 1 && ys == 1 {
		return Slice[S]{r}
	}

	w := r.Dx() + xGap
	h := r.Dy() + yGap

	rs := make(Slice[S], 0, xs*ys)
	for y := range ys {
		for x := range xs {
			minX := r.Min.X + S(x)*w
			minY := r.Min.Y + S(y)*h
			maxX := minX + r.Dx()
			maxY := minY + r.Dy()
			rs = append(rs, XYXY(minX, minY, maxX, maxY))
		}
	}
	return rs
}

// RepeatX creates xs copies of r arranged horizontally with xGap spacing.
func (r *Rect[S]) RepeatX(xs int, xGap S) Slice[S] {
	return r.Repeat(xs, 1, xGap, 0)
}

// RepeatY creates ys copies of r arranged vertically with yGap spacing.
func (r *Rect[S]) RepeatY(ys int, yGap S) Slice[S] {
	return r.Repeat(1, ys, 0, yGap)
}
