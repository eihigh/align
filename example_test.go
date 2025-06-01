package align_test

import "github.com/eihigh/align"

func Example() {
	formSpace := align.WH(100, 100)
	rows := align.Slice[int]{}
	{
		row := align.Map[int]{}
		space := formSpace.Inset(12)
		row["label"] = align.WH(80, 20).Nest(space, 0, .5)
		row["input"] = align.WH(80, 25).Nest(space, 1, .5)
		// row := align.Slice[int]{g.label, g.input}
		row.StackY(rows.Last(), .5, 1).Add(align.XY(0, 4))
		rows = append(rows, row)
		rows.Nest(formSpace, 0, 0)
	}
	// 手持ちの構造体フィールドを突っ込むならMapよりSliceの方が良い
}

func ExampleWrapper() {
	space := align.WH(100, 100)
	widths := []int{10, 20, 30, 40, 50, 60, 70, 80, 90}
	stack := func(a, b *align.Rect[int]) {
		a.StackX(b, 1, .5).Add(align.XY(4, 0))
	}
	wrap := func(a, b *align.Rect[int]) {
		a.StackY(b, 0, 1).Add(align.XY(0, 4))
	}
	w := align.NewWrapper(space, 0, 0, stack, wrap)
	for _, width := range widths {
		r := align.WH(width, 20)
		if !w.Add(r) {
			break // 収まらなかったら終了
		}
	}
	w.Slice().Align(0, 0, space, 0, 0) // 全体をスペースに合わせて配置
}
