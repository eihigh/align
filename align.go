package align

import "github.com/eihigh/ng"

// Union returns the smallest rectangle containing all input rectangles.
func Union[S ng.Scalar](rs ...Rect[S]) Rect[S] {
	r := Rect[S]{}
	for _, rr := range rs {
		r.Min.X = min(r.Min.X, rr.Min.X)
		r.Min.Y = min(r.Min.Y, rr.Min.Y)
		r.Max.X = max(r.Max.X, rr.Max.X)
		r.Max.Y = max(r.Max.Y, rr.Max.Y)
	}
	return r
}

// Intersect returns the largest rectangle contained by all input rectangles.
func Intersect[S ng.Scalar](rs ...Rect[S]) Rect[S] {
	r := Rect[S]{}
	for _, rr := range rs {
		r.Min.X = max(r.Min.X, rr.Min.X)
		r.Min.Y = max(r.Min.Y, rr.Min.Y)
		r.Max.X = min(r.Max.X, rr.Max.X)
		r.Max.Y = min(r.Max.Y, rr.Max.Y)
	}
	return r
}

// CutLeft splits r vertically, returning a left strip of width w and the remainder.
func (r Rect[S]) CutLeft(w S) (left, rest Rect[S]) {
	w = max(0, min(w, r.Dx()))
	left = Rect[S]{
		Min: Point[S]{r.Min.X, r.Min.Y},
		Max: Point[S]{r.Min.X + w, r.Max.Y},
	}
	rest = Rect[S]{
		Min: Point[S]{r.Min.X + w, r.Min.Y},
		Max: Point[S]{r.Max.X, r.Max.Y},
	}
	return left, rest
}

// CutTop splits r horizontally, returning a top strip of height h and the remainder.
func (r Rect[S]) CutTop(h S) (top, rest Rect[S]) {
	h = max(0, min(h, r.Dy()))
	top = Rect[S]{
		Min: Point[S]{r.Min.X, r.Min.Y},
		Max: Point[S]{r.Max.X, r.Min.Y + h},
	}
	rest = Rect[S]{
		Min: Point[S]{r.Min.X, r.Min.Y + h},
		Max: Point[S]{r.Max.X, r.Max.Y},
	}
	return top, rest
}

// CutRight splits r vertically, returning a right strip of width w and the remainder.
func (r Rect[S]) CutRight(w S) (right, rest Rect[S]) {
	w = max(0, min(w, r.Dx()))
	right = Rect[S]{
		Min: Point[S]{r.Max.X - w, r.Min.Y},
		Max: Point[S]{r.Max.X, r.Max.Y},
	}
	rest = Rect[S]{
		Min: Point[S]{r.Min.X, r.Min.Y},
		Max: Point[S]{r.Max.X - w, r.Max.Y},
	}
	return right, rest
}

// CutBottom splits r horizontally, returning a bottom strip of height h and the remainder.
func (r Rect[S]) CutBottom(h S) (bottom, rest Rect[S]) {
	h = max(0, min(h, r.Dy()))
	bottom = Rect[S]{
		Min: Point[S]{r.Min.X, r.Max.Y - h},
		Max: Point[S]{r.Max.X, r.Max.Y},
	}
	rest = Rect[S]{
		Min: Point[S]{r.Min.X, r.Min.Y},
		Max: Point[S]{r.Max.X, r.Max.Y - h},
	}
	return bottom, rest
}

// CutLeftByRate splits r vertically by ratio f (0.0-1.0).
func (r Rect[S]) CutLeftByRate(f float64) (left, rest Rect[S]) {
	w := S(float64(r.Dx()) * f)
	return r.CutLeft(w)
}

// CutTopByRate splits r horizontally by ratio f (0.0-1.0).
func (r Rect[S]) CutTopByRate(f float64) (top, rest Rect[S]) {
	h := S(float64(r.Dy()) * f)
	return r.CutTop(h)
}

// CutRightByRate splits r vertically by ratio f (0.0-1.0).
func (r Rect[S]) CutRightByRate(f float64) (right, rest Rect[S]) {
	w := S(float64(r.Dx()) * f)
	return r.CutRight(w)
}

// CutBottomByRate splits r horizontally by ratio f (0.0-1.0).
func (r Rect[S]) CutBottomByRate(f float64) (bottom, rest Rect[S]) {
	h := S(float64(r.Dy()) * f)
	return r.CutBottom(h)
}

// Split divides r into a grid of xs×ys rectangles with specified gaps.
func (r Rect[S]) Split(xs, ys int, xGap, yGap S) []Rect[S] {
	if xs <= 0 || ys <= 0 {
		return nil
	}
	if xs == 1 && ys == 1 {
		return []Rect[S]{r}
	}

	w := (r.Dx() - S(xGap)*(S(xs)-1)) / S(xs)
	h := (r.Dy() - S(yGap)*(S(ys)-1)) / S(ys)

	rs := make([]Rect[S], 0, xs*ys)
	for y := range ys {
		for x := range xs {
			minX := r.Min.X + S(x)*(w+xGap)
			minY := r.Min.Y + S(y)*(h+yGap)
			maxX := minX + w
			maxY := minY + h
			rs = append(rs, Rect[S]{
				Min: Point[S]{minX, minY},
				Max: Point[S]{maxX, maxY},
			})
		}
	}
	return rs
}

// SplitX divides r into xs columns with gap xGap.
func (r Rect[S]) SplitX(xs int, xGap S) []Rect[S] {
	return r.Split(xs, 1, xGap, 0)
}

// SplitY divides r into ys rows with gap yGap.
func (r Rect[S]) SplitY(ys int, yGap S) []Rect[S] {
	return r.Split(1, ys, 0, yGap)
}

// Repeat creates a grid of xs×ys copies of r with specified gaps.
func (r Rect[S]) Repeat(xs, ys int, xGap, yGap S) []Rect[S] {
	if xs <= 0 || ys <= 0 {
		return nil
	}
	if xs == 1 && ys == 1 {
		return []Rect[S]{r}
	}

	w := r.Dx() + xGap
	h := r.Dy() + yGap

	rs := make([]Rect[S], 0, xs*ys)
	for y := range ys {
		for x := range xs {
			minX := r.Min.X + S(x)*w
			minY := r.Min.Y + S(y)*h
			maxX := minX + r.Dx()
			maxY := minY + r.Dy()
			rs = append(rs, Rect[S]{
				Min: Point[S]{minX, minY},
				Max: Point[S]{maxX, maxY},
			})
		}
	}
	return rs
}

// RepeatX creates xs horizontal copies of r with gap xGap.
func (r Rect[S]) RepeatX(xs int, xGap S) []Rect[S] {
	return r.Repeat(xs, 1, xGap, 0)
}

// RepeatY creates ys vertical copies of r with gap yGap.
func (r Rect[S]) RepeatY(ys int, yGap S) []Rect[S] {
	return r.Repeat(1, ys, 0, yGap)
}

// InsetXY shrinks r by x horizontally and y vertically.
func (r Rect[S]) InsetXY(x, y S) Rect[S] {
	return r.InsetLTRB(x, y, x, y)
}

// InsetLTRB shrinks r by specified amounts on each side.
func (r Rect[S]) InsetLTRB(left, top, right, bottom S) Rect[S] {
	if r.Dx() < left+right {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += left
		r.Max.X -= right
	}
	if r.Dy() < top+bottom {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += top
		r.Max.Y -= bottom
	}
	return r
}

// Outset expands r by s on all sides.
func (r Rect[S]) Outset(s S) Rect[S] {
	return r.Inset(-s)
}

// OutsetXY expands r by x horizontally and y vertically.
func (r Rect[S]) OutsetXY(x, y S) Rect[S] {
	return r.InsetXY(-x, -y)
}

// OutsetLTRB expands r by specified amounts on each side.
func (r Rect[S]) OutsetLTRB(left, top, right, bottom S) Rect[S] {
	return r.InsetLTRB(-left, -top, -right, -bottom)
}

// Anchor returns a point at relative position (x,y) within r, where (0,0) is top-left and (1,1) is bottom-right.
func (r Rect[S]) Anchor(x, y float64) Point[S] {
	return Point[S]{
		X: r.Min.X + S(float64(r.Dx())*x),
		Y: r.Min.Y + S(float64(r.Dy())*y),
	}
}

// Align positions r relative to dst using anchor points.
// (x,y) specifies r's anchor point, (dx,dy) specifies dst's anchor point.
// Both use normalized coordinates where (0,0) is top-left and (1,1) is bottom-right.
func (r Rect[S]) Align(x, y float64, dst Rect[S], dx, dy float64) Rect[S] {
	return Rect[S]{
		Min: Point[S]{
			X: dst.Min.X + S(float64(dst.Dx())*dx) - S(float64(r.Dx())*x),
			Y: dst.Min.Y + S(float64(dst.Dy())*dy) - S(float64(r.Dy())*y),
		},
		Max: Point[S]{
			X: dst.Min.X + S(float64(dst.Dx())*dx) + S(float64(r.Dx())*(1-x)),
			Y: dst.Min.Y + S(float64(dst.Dy())*dy) + S(float64(r.Dy())*(1-y)),
		},
	}
}

// CenterOf centers r within dst.
func (r Rect[S]) CenterOf(dst Rect[S]) Rect[S] {
	return r.Align(0.5, 0.5, dst, 0.5, 0.5)
}

// Nest positions r within dst at relative position (dx,dy).
func (r Rect[S]) Nest(dst Rect[S], dx, dy float64) Rect[S] {
	return r.Align(dx, dy, dst, dx, dy)
}

// StackX positions r horizontally adjacent to dst.
func (r Rect[S]) StackX(dst Rect[S], dx, dy float64) Rect[S] {
	return r.Align(1-dx, dy, dst, dx, dy)
}

// StackY positions r vertically adjacent to dst.
func (r Rect[S]) StackY(dst Rect[S], dx, dy float64) Rect[S] {
	return r.Align(dx, 1-dy, dst, dx, dy)
}

// Clamp adjusts the rectangle r to fit within dst while keeping its size.
func (r Rect[S]) Clamp(dst Rect[S]) Rect[S] {
	size := r.Size()
	r.Min.X = max(dst.Min.X, min(dst.Max.X-size.X, r.Min.X))
	r.Min.Y = max(dst.Min.Y, min(dst.Max.Y-size.Y, r.Min.Y))
	r.Max.X = r.Min.X + size.X
	r.Max.Y = r.Min.Y + size.Y
	return r
}

// Wrapper arranges rectangles within bounds using stack and wrap functions.
type Wrapper[S ng.Scalar] struct {
	bounds      Rect[S]
	x, y        float64
	stack, wrap func(a, b Rect[S]) Rect[S]
	rects       []Rect[S]
	lineFirst   Rect[S]
}

// NewWrapper creates a Wrapper that arranges rectangles within bounds.
// (x,y) specifies initial position, stack arranges adjacent items, wrap handles line breaks.
func NewWrapper[S ng.Scalar](bounds Rect[S], x, y float64, stack, wrap func(a, b Rect[S]) Rect[S]) *Wrapper[S] {
	return &Wrapper[S]{
		bounds: bounds,
		x:      x,
		y:      y,
		stack:  stack,
		wrap:   wrap,
	}
}

// Add places r using stack function, wrapping to next line if needed.
// Returns false if r cannot fit within bounds.
func (w *Wrapper[S]) Add(r Rect[S]) (ok bool) {
	if len(w.rects) == 0 {
		r = r.Nest(w.bounds, w.x, w.y)
		if !r.In(w.bounds) {
			return false
		}
		w.lineFirst = r
	} else {
		r = w.stack(r, w.rects[len(w.rects)-1])
		if !r.In(w.bounds) {
			r = w.wrap(r, w.lineFirst)
			if !r.In(w.bounds) {
				return false
			}
			w.lineFirst = r
		}
	}
	w.rects = append(w.rects, r)
	return true
}

// Rects returns all rectangles added to the wrapper.
func (w *Wrapper[S]) Rects() []Rect[S] { return w.rects }
