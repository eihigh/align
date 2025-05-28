package align

import (
	"iter"
	"slices"
	"testing"
)

func TestRect_Points(t *testing.T) {
	tests := []struct {
		r    interface{ Points() iter.Seq[Point[int]] }
		want []Point[int]
	}{
		{
			r: XYXY(0, 0, 3, 2),
			want: []Point[int]{
				{0, 0}, {1, 0}, {2, 0},
				{0, 1}, {1, 1}, {2, 1},
			},
		},
		{
			r: XYXY(-1, -2, 1, 0),
			want: []Point[int]{
				{-1, -2}, {0, -2},
				{-1, -1}, {0, -1},
			},
		},
		{
			r:    XYXY(-2, -2, -2, -2),
			want: []Point[int]{},
		},
		{
			r: XYXY(0.1, 0.1, 3, 3),
			want: []Point[int]{
				{1, 1}, {2, 1},
				{1, 2}, {2, 2},
			},
		},
	}
	for _, tt := range tests {
		got := slices.Collect(tt.r.Points())
		if !slices.Equal(got, tt.want) {
			t.Errorf("%v: got %v, want %v", tt.r, got, tt.want)
		}
	}
}
