package align

import (
	"slices"
	"testing"
)

func TestAlign(t *testing.T) {
	s := XYXY(-5, -5, 5, 5)
	tests := []struct {
		name      string
		got, want *Rect[int]
	}{
		{
			name: "CenterOf",
			got:  WH(5, 5).CenterOf(s),
			want: XYWH(-2, -2, 5, 5),
		},
		{
			name: "Nest",
			got:  WH(5, 5).Nest(s, 1, 1),
			want: XYWH(0, 0, 5, 5),
		},
		{
			name: "StackX",
			got:  WH(5, 5).StackX(s, 1, .5),
			want: XYWH(5, -2, 5, 5),
		},
		{
			name: "StackY",
			got:  WH(5, 5).StackY(s, .5, 1),
			want: XYWH(-2, 5, 5, 5),
		},
		{
			name: "Overlap",
			got:  WH(5, 5).Align(.5, .5, s, 0, 0),
			want: XYWH(-7, -7, 5, 5),
		},
		{
			name: "Clamp min",
			got:  XYWH(-10, -10, 5, 5).Clamp(s),
			want: XYWH(-5, -5, 5, 5),
		},
		{
			name: "Clamp max",
			got:  XYWH(20, 20, 5, 5).Clamp(s),
			want: XYWH(0, 0, 5, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Eq(tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestSplitRepeat(t *testing.T) {
	s := XYWH(-30, -30, 60, 60)
	tests := []struct {
		name      string
		got, want Slice[int]
	}{
		{
			name: "SplitX",
			got:  s.SplitX(5, 2), // (10+2)*5 = 60
			want: Slice[int]{
				XYWH(-30, -30, 10, 60),
				XYWH(-18, -30, 10, 60),
				XYWH(-6, -30, 10, 60),
				XYWH(6, -30, 10, 60),
				XYWH(18, -30, 10, 60),
			},
		},
		{
			name: "Split",
			got:  s.Split(3, 2, 10, 10), // (10+10)*3 = 60, (20+10)*2 = 60
			want: Slice[int]{
				XYWH(-30, -30, 20, 10),
				XYWH(-10, -30, 20, 10),
				XYWH(10, -30, 20, 10),
				XYWH(-30, 0, 20, 10),
				XYWH(-10, 0, 20, 10),
				XYWH(10, 0, 20, 10),
			},
		},
		{
			name: "RepeatX",
			got:  WH(5, 5).RepeatX(3, 3), // (5+3)*3 = 24
			want: Slice[int]{
				XYWH(0, 0, 5, 5),
				XYWH(8, 0, 5, 5),
				XYWH(16, 0, 5, 5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slices.EqualFunc(tt.got, tt.want, func(a, b Node[int]) bool {
				return a.Bounds().Eq(b.Bounds())
			})
		})
	}
}
