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
