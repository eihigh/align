package align

import (
	"image"
	"iter"
	"math"

	"github.com/eihigh/ng"
)

// A Rect contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
type Rect[S ng.Scalar] struct {
	Min, Max Point[S]
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rect[S]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rect[S]) Dx() S {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rect[S]) Dy() S {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rect[S]) Size() Point[S] {
	return Point[S]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rect[S]) Add(p Point[S]) Rect[S] {
	return Rect[S]{
		Point[S]{r.Min.X + p.X, r.Min.Y + p.Y},
		Point[S]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rect[S]) Sub(p Point[S]) Rect[S] {
	return Rect[S]{
		Point[S]{r.Min.X - p.X, r.Min.Y - p.Y},
		Point[S]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rect[S]) Inset(n S) Rect[S] {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

// Rect.Intersect is removed

// Rect.Union is removed

// Empty reports whether the rectangle contains no points.
func (r Rect[S]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rect[S]) Eq(s Rect[S]) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rect[S]) Overlaps(s Rect[S]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rect[S]) In(s Rect[S]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rect[S]) Canon() Rect[S] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// XYXY is shorthand for [Rect]{XY(x0, y0), XY(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func XYXY[S ng.Scalar](x0, y0, x1, y1 S) Rect[S] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rect[S]{Point[S]{x0, y0}, Point[S]{x1, y1}}
}

// --------------------------------------
// Additional functions
// --------------------------------------

// XYWH creates a rectangle from position (x,y) and dimensions (w,h).
func XYWH[S ng.Scalar](x, y, w, h S) Rect[S] {
	return Rect[S]{
		Min: Point[S]{x, y},
		Max: Point[S]{x + w, y + h},
	}
}

// PosSize creates a rectangle from position and size points.
func PosSize[S ng.Scalar](pos, size Point[S]) Rect[S] {
	return Rect[S]{
		Min: Point[S]{pos.X, pos.Y},
		Max: Point[S]{pos.X + size.X, pos.Y + size.Y},
	}
}

// PosWH creates a rectangle from position point and dimensions.
func PosWH[S ng.Scalar](pos Point[S], w, h S) Rect[S] {
	return Rect[S]{
		Min: Point[S]{pos.X, pos.Y},
		Max: Point[S]{pos.X + w, pos.Y + h},
	}
}

// XYSize creates a rectangle from coordinates and size point.
func XYSize[S ng.Scalar](x, y S, size Point[S]) Rect[S] {
	return Rect[S]{
		Min: Point[S]{x, y},
		Max: Point[S]{x + size.X, y + size.Y},
	}
}

// WH creates a rectangle from origin (0,0) with dimensions (w,h).
func WH[S ng.Scalar](w, h S) Rect[S] {
	return Rect[S]{
		Min: Point[S]{0, 0},
		Max: Point[S]{w, h},
	}
}

// Int converts r to Rect[int].
func (r Rect[S]) Int() Rect[int] {
	return XYXY(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

// Float64 converts r to Rect[float64].
func (r Rect[S]) Float64() Rect[float64] {
	return XYXY(float64(r.Min.X), float64(r.Min.Y), float64(r.Max.X), float64(r.Max.Y))
}

// Float32 converts r to Rect[float32].
func (r Rect[S]) Float32() Rect[float32] {
	return XYXY(float32(r.Min.X), float32(r.Min.Y), float32(r.Max.X), float32(r.Max.Y))
}

// Image converts r to an [image.Rectangle].
func (r Rect[S]) Image() image.Rectangle {
	return image.Rect(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

// Points returns an iterator that yields the points in r in row-major order.
func (r Rect[S]) Points() iter.Seq[Point[int]] {
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
