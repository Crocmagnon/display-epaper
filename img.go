package main

import (
	"context"
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
)

func getBlack(ctx context.Context, transportsClient *transports.Client) (*image.RGBA, error) {
	passages, err := transportsClient.GetTCLPassages(ctx, 290)
	if err != nil {
		return nil, fmt.Errorf("getting passages: %w", err)
	}

	img := newWhite()

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFillRule(draw2d.FillRuleWinding)
	gc.SetFillColor(color.RGBA{0, 0, 0, 255})
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})

	for i, passage := range passages.Passages {
		x := float64(10 + i*100)
		text(gc, passage.Ligne, 15, x, 20)
		for j, delay := range passage.Delays {
			y := float64(20 + (j+1)*30)
			text(gc, delay, 15, x, y)
		}
	}

	return img, nil
}

func text(gc *draw2dimg.GraphicContext, s string, size, x, y float64) {
	gc.SetFontData(draw2d.FontData{Name: fontName})
	gc.SetFontSize(size)
	gc.FillStringAt(s, x, y)
	gc.SetLineWidth(2)
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
