## examples

<table>
<tr>
<th>Code</th>
<th>Result</th>
</tr>
<tr>
<td>

```go
screen := align.WH(100, 100) // blue
r := screen.Inset(20) // red
```

</td>
<td>

![](/examples/inset.png)

</td>
</tr>
<tr>
<td>

```go
screen := align.WH(100, 100) // blue
topBar, screen := screen.CutTop(10) // cyan
title := align.WH(30, 10).CenterOf(topBar) // red
portrait := align.WH(20, 30).Nest(screen.Inset(10), left, top) // green
btn := align.WH(20, 10).Nest(screen.Inset(10), right, bottom) // magenta
```

</td>
<td>

![](/examples/status.png)

</td>
</tr>
<tr>
<td>

```go
screen := align.WH(100, 100) // blue
w := align.NewWrapper(screen.Inset(10), 0, 0,
	func(a, b align.Rect[int]) align.Rect[int] {
		return a.StackX(b, 1, 0).Add(align.XY(5, 0))
	},
	func(a, b align.Rect[int]) align.Rect[int] {
		return a.StackY(b, 0, 1).Add(align.XY(0, 5))
	},
)
for range 8 {
	if !w.Add(align.WH(20, 20)) { // magenta
		break
	}
}
```

</td>
<td>

![](/examples/wrap.png)

</td>
</tr>