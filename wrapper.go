package align

import "github.com/eihigh/ng"

// Wrapper arranges rectangles within bounds using stack and wrap functions.
type Wrapper[S ng.Scalar] struct {
	bounds      *Rect[S]
	x, y        float64
	stack, wrap func(a, b *Rect[S])
	s           Slice[S]
	lineFirst   *Rect[S]
}

// NewWrapper creates a Wrapper that arranges rectangles within bounds.
// (x,y) specifies initial position, stack arranges adjacent items, wrap handles line breaks.
func NewWrapper[S ng.Scalar](bounds *Rect[S], x, y float64, stack, wrap func(a, b *Rect[S])) *Wrapper[S] {
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
func (w *Wrapper[S]) Add(r *Rect[S]) (ok bool) {
	if len(w.s) == 0 {
		r = r.Nest(w.bounds, w.x, w.y)
		if !r.In(w.bounds) {
			return false
		}
		w.lineFirst = r
	} else {
		w.stack(r, w.s.Last().Bounds())
		if !r.In(w.bounds) {
			w.wrap(r, w.lineFirst)
			if !r.In(w.bounds) {
				return false
			}
			w.lineFirst = r
		}
	}
	w.s = append(w.s, r)
	return true
}

// Slice returns all rectangles added to the wrapper.
func (w *Wrapper[S]) Slice() Slice[S] { return w.s }
