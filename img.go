package main

import (
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
)

func getBlack() (*image.RGBA, error) {
	img := newWhite()

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFillColor(color.RGBA{0, 0, 0, 255})
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})
	gc.SetLineWidth(2)

	text(gc, "Hello, world", 18, 110, 50)

	return img, nil
}

func getRed() (*image.RGBA, error) {
	img := newBlack()
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.RGBA{0, 0, 0, 255})
	gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})
	gc.SetLineWidth(2)

	rect(gc, 10, 10, 50, 50)
	rect(gc, 60, 10, 100, 50)

	return img, nil
}

func text(gc *draw2dimg.GraphicContext, s string, size, x, y float64) {
	gc.SetFontData(draw2d.FontData{Name: fontName})
	gc.SetFontSize(size)
	gc.StrokeStringAt(s, x, y)
}

func newWhite() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, epd.Width, epd.Height))
	for y := 0; y < epd.Height; y++ {
		for x := 0; x < epd.Width; x++ {
			img.Set(x, y, color.White)
		}
	}
	return img
}

func newBlack() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, epd.Width, epd.Height))
	for y := 0; y < epd.Height; y++ {
		for x := 0; x < epd.Width; x++ {
			img.Set(x, y, color.Black)
		}
	}
	return img
}

func rect(gc *draw2dimg.GraphicContext, x1, y1, x2, y2 float64) {
	gc.BeginPath()
	gc.MoveTo(x1, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x1, y2)
	gc.Close()
	gc.FillStroke()
}
