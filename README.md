# align
A lightweight, framework-agnostic UI layout library for Go

## Overview

`align` provides a comprehensive set of geometric operations for rectangles and points, making it easy to create complex layouts with minimal code. It uses Go generics to work with any numeric type.

## Quick Reference

### Positioning
- `Align` - Position rectangles relative to each other using anchor points
- `CenterOf` - Center a rectangle within another
- `Nest` - Position a rectangle within another at a relative position
- `StackX/Y` - Stack rectangles horizontally or vertically

### Cutting
- `CutLeft/Top/Right/Bottom` - Split rectangles into strips
- `CutLeftByRate/TopByRate/RightByRate/BottomByRate` - Split by rate

### Sizing
- `Inset/Outset` - Shrink or expand rectangles
- `InsetXY/OutsetXY` - Shrink or expand with different X/Y values
- `InsetLTRB/OutsetLTRB` - Shrink or expand with individual side values

### Division
- `Split` - Divide into a grid with gaps
- `SplitX/Y` - Divide into columns or rows
- `Repeat` - Create a grid of copies
- `RepeatX/Y` - Create horizontal or vertical copies

### Utilities
- `Union` - Get the smallest rectangle containing all inputs
- `Intersect` - Get the largest rectangle contained by all inputs
- `Clamp` - Constrain a rectangle within bounds while keeping its size
- `Anchor` - Get a point at a relative position within a rectangle

### Wrapper
The `Wrapper` type helps arrange multiple rectangles within bounds, automatically handling line wrapping when items don't fit.

## Examples

### Basic Rectangle Creation

```go
// Create rectangles
r1 := align.XYWH(10, 20, 100, 50)  // x=10, y=20, width=100, height=50
r2 := align.WH(200, 100)           // x=0, y=0, width=200, height=100
r3 := align.XYXY(10, 10, 110, 60)  // x=10, y=10, width=100, height=50
```

### Alignment

```go
// Center a button in a container
button := align.WH(100, 40)
container := align.WH(800, 600)
centered := button.CenterOf(container)

// Align to specific positions
// Place at top-right corner
topRight := button.Align(1, 0, container, 1, 0)

// Place at bottom-center
bottomCenter := button.Align(0.5, 1, container, 0.5, 1)
// Equivalent to:
bottomCenter = button.Nest(container, 0.5, 1)
```

### Cutting and Splitting

```go
// Create a header/content layout
screen := align.WH(800, 600)
header, content := screen.CutTop(60)

// Create a sidebar layout
sidebar, main := content.CutLeft(200)

// Split into grid
grid := screen.Split(3, 3, 10, 10)  // 3x3 grid with 10px gaps
```

### Insets and Outsets

```go
// Add padding
padded := rect.Inset(10)  // 10px padding on all sides
customPadded := rect.InsetLTRB(20, 10, 20, 10)  // left, top, right, bottom

// Expand for shadows or borders
withShadow := rect.Outset(5)
```

### Using Wrapper for Flow Layout

```go
bounds := align.WH(800, 600)
wrapper := align.NewWrapper(bounds, 0, 0,
    func(a, b align.Rect[int]) align.Rect[int] {
        // Stack horizontally
        return a.StackX(b, 1, 0.5)
    },
    func(a, b align.Rect[int]) align.Rect[int] {
        // Wrap to next line
        return a.StackY(b, 0, 1)
    },
)

// Add items
for _, size := range sizes {
    item := align.WH(size.X, size.Y)
    if !wrapper.Add(item) {
        break  // No more space
    }
}

// Get all positioned rectangles
rects := wrapper.Rects()
```

## Type Conversions

The library supports easy conversion between numeric types:

```go
// Convert between types
floatRect := intRect.Float64()
intPoint := floatPoint.Int()

// Convert to standard library types
imgRect := rect.Image()  // image.Rectangle
imgPoint := point.Image() // image.Point
```

## License

MIT License - see LICENSE file for details.
