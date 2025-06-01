package align

import (
	"fmt"
	"image"

	"github.com/eihigh/ng"
)

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point[S ng.Scalar] struct {
	X, Y S
}

// String returns a string representation of p like "(3,4)".
func (p Point[S]) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

// Add returns the vector p+q.
func (p Point[S]) Add(q Point[S]) Point[S] {
	return Point[S]{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point[S]) Sub(q Point[S]) Point[S] {
	return Point[S]{p.X - q.X, p.Y - q.Y}
}

// Scale returns the vector p*k.
func (p Point[S]) Scale(k S) Point[S] {
	return Point[S]{p.X * k, p.Y * k}
}

// Mul returns the vector p*q.
func (p Point[S]) Mul(q Point[S]) Point[S] {
	return Point[S]{p.X * q.X, p.Y * q.Y}
}

// Div returns the vector p/q.
func (p Point[S]) Div(q Point[S]) Point[S] {
	return Point[S]{p.X / q.X, p.Y / q.Y}
}

// In reports whether p is in r.
func (p Point[S]) In(r Rect[S]) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// TODO: Point.Mod?

// Eq reports whether p and q are equal.
func (p Point[S]) Eq(q Point[S]) bool {
	return p == q
}

// XY is shorthand for [Point[S]]{X, Y}.
func XY[S ng.Scalar](X, Y S) Point[S] {
	return Point[S]{X, Y}
}

// --------------------------------------
// Additional functions
// --------------------------------------

// XY unpacks the X and Y coordinates of p.
func (p Point[S]) XY() (S, S) {
	return p.X, p.Y
}

// Int converts p to Point[int].
func (p Point[S]) Int() Point[int] {
	return Point[int]{int(p.X), int(p.Y)}
}

// Float64 converts p to Point[float64].
func (p Point[S]) Float64() Point[float64] {
	return Point[float64]{float64(p.X), float64(p.Y)}
}

// Float32 converts p to Point[float32].
func (p Point[S]) Float32() Point[float32] {
	return Point[float32]{float32(p.X), float32(p.Y)}
}

// Image converts p to an [image.Point].
func (p Point[S]) Image() image.Point {
	return image.Pt(int(p.X), int(p.Y))
}
