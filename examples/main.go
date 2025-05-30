package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/eihigh/align"
)

const (
	top    = 0.0
	mid    = 0.5
	bottom = 1.0
	left   = 0.0
	right  = 1.0
)

var (
	red     = color.RGBA{255, 0, 0, 255}
	green   = color.RGBA{0, 255, 0, 255}
	blue    = color.RGBA{0, 0, 255, 255}
	yellow  = color.RGBA{255, 255, 0, 255}
	cyan    = color.RGBA{0, 255, 255, 255}
	magenta = color.RGBA{255, 0, 255, 255}
)

type img struct {
	img  *image.RGBA
	name string
}

func newImg(name string, w, h int) img {
	return img{
		img:  image.NewRGBA(image.Rect(0, 0, w, h)),
		name: name,
	}
}

func (i *img) fill(r align.Rect[int], c color.Color) {
	for p := range r.Points() {
		i.img.Set(p.X, p.Y, c)
	}
}

func (i *img) save() error {
	f, err := os.Create(i.name + ".png")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, i.img); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := inset(); err != nil {
		panic(err)
	}
	if err := status(); err != nil {
		panic(err)
	}
	if err := wrap(); err != nil {
		panic(err)
	}
	if err := title(); err != nil {
		panic(err)
	}
}

func inset() error {
	screen := align.WH(100, 100)
	r := screen.Inset(20)

	img := newImg("inset", 100, 100)
	img.fill(screen, blue)
	img.fill(r, red)
	return img.save()
}

func status() error {
	screen := align.WH(100, 100)
	topBar, screen := screen.CutTop(10)
	title := align.WH(30, 10).CenterOf(topBar)
	portrait := align.WH(20, 30).Nest(screen.Inset(10), left, top)
	btn := align.WH(20, 10).Nest(screen.Inset(10), right, bottom)

	img := newImg("status", 100, 100)
	img.fill(screen, blue)
	img.fill(topBar, cyan)
	img.fill(title, red)
	img.fill(portrait, green)
	img.fill(btn, magenta)
	return img.save()
}

func wrap() error {
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

	img := newImg("wrap", 100, 100)
	img.fill(screen, blue)
	for _, r := range w.Rects() {
		img.fill(r, magenta)
	}
	return img.save()
}

func title() error {
	screen := align.WH(100, 100) // blue
	logo, menu := screen.CutTopByRate(0.5)
	logo = align.WH(60, 30).CenterOf(logo) // yellow
	menu = menu.Inset(5)
	off := align.XY(0, 4)
	newGame := align.WH(30, 8).Nest(menu, mid, top)                        // red
	continueGame := align.WH(30, 8).StackY(newGame, mid, bottom).Add(off)  // red
	exitGame := align.WH(30, 8).StackY(continueGame, mid, bottom).Add(off) // red
	menu = align.Union(newGame, continueGame, exitGame).Outset(5)          // cyan

	img := newImg("title", 100, 100)
	img.fill(screen, blue)
	img.fill(logo, yellow)
	img.fill(menu, cyan)
	img.fill(newGame, red)
	img.fill(continueGame, red)
	img.fill(exitGame, red)
	return img.save()
}
