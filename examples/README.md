## examples

<table>
<tr>
<th>Code</th>
<th>Screenshot</th>
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
