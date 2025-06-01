package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/eihigh/align"
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

func (i *img) fill(r align.Node[int], c color.Color) {
	for p := range r.Bounds().Points() {
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
	screen := align.WH(100, 100)      // blue
	topBar, screen := screen.CutY(10) // cyan
	uiSpace := screen.Inset(10)
	title := align.WH(30, 10).CenterOf(topBar)       // red
	portrait := align.WH(20, 30).Nest(uiSpace, 0, 0) // green
	btn := align.WH(20, 10).Nest(uiSpace, 1, 1)      // magenta

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
	w := align.NewWrapper(screen, 0, 0,
		func(a, b *align.Rect[int]) { // stack
			a.StackX(b, 1, .5).Add(align.XY(4, 0))
		},
		func(a, b *align.Rect[int]) { // wrap
			a.StackY(b, 0, 1).Add(align.XY(0, 4))
		},
	)
	for range 8 {
		if !w.Add(align.WH(25, 25)) { // magenta
			break
		}
	}
	w.Slice().CenterOf(screen)

	img := newImg("wrap", 100, 100)
	img.fill(screen, blue)
	for _, r := range w.Slice() {
		img.fill(r, magenta)
	}
	return img.save()
}

func title() error {
	screen := align.WH(100, 100) // blue
	logo, menuSpace := screen.CutYByRate(0.5)
	logo = align.WH(60, 30).CenterOf(logo) // yellow

	off := align.XY(0, 4)
	newGame := align.WH(30, 8)                                       // red
	continueGame := align.WH(40, 8).StackY(newGame, .5, 1).Add(off)  // red
	exitGame := align.WH(25, 8).StackY(continueGame, .5, 1).Add(off) // red
	menuItems := align.Slice[int]{newGame, continueGame, exitGame}
	menuItems.CenterOf(menuSpace)
	menuWindow := menuItems.Outset(5) // cyan

	img := newImg("title", 100, 100)
	img.fill(screen, blue)
	img.fill(logo, yellow)
	img.fill(menuWindow, cyan)
	img.fill(newGame, red)
	img.fill(continueGame, red)
	img.fill(exitGame, red)
	return img.save()
}
