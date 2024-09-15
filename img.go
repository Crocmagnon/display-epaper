package main

import (
	"context"
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"strconv"
	"time"
)

func getBlack(
	ctx context.Context,
	nowFunc func() time.Time,
	transportsClient *transports.Client,
	feteClient *fete.Client,
) (*image.RGBA, error) {
	passages, err := transportsClient.GetTCLPassages(ctx, 290)
	if err != nil {
		return nil, fmt.Errorf("getting passages: %w", err)
	}

	fetes, err := feteClient.GetFete(ctx, nowFunc())
	if err != nil {
		return nil, fmt.Errorf("getting fetes: %w", err)
	}

	_ = fetes

	img := newWhite()

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFillRule(draw2d.FillRuleWinding)
	gc.SetFillColor(color.RGBA{0, 0, 0, 255})
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})

	drawPassages(gc, passages)

	drawFete(gc, fetes, nowFunc())

	return img, nil
}

func drawFete(gc *draw2dimg.GraphicContext, fetes *fete.Fete, now time.Time) {
	text(gc, getDate(now), 50, 20, 235)
	text(gc, fetes.Name, 18, 20, 270)
}

func drawPassages(gc *draw2dimg.GraphicContext, passages *transports.Passages) {
	for i, passage := range passages.Passages {
		x := float64(600 + i*100)
		text(gc, passage.Ligne, 15, x, 20)
		for j, delay := range passage.Delays {
			y := float64(20 + (j+1)*30)
			text(gc, delay, 15, x, y)
		}
	}
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

func rect(gc *draw2dimg.GraphicContext, x1, y1, x2, y2 float64) {
	gc.BeginPath()
	gc.MoveTo(x1, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x1, y2)
	gc.Close()
	gc.FillStroke()
}

func getDate(now time.Time) string {
	return fmt.Sprintf("%v %v", getDay(now), getMonth(now))
}

func getDay(now time.Time) string {
	if now.Day() == 1 {
		return "1er"
	}

	return strconv.Itoa(now.Day())
}

func getMonth(t time.Time) string {
	switch t.Month() {
	case time.January:
		return "jan."
	case time.February:
		return "fev."
	case time.March:
		return "mars"
	case time.April:
		return "avr."
	case time.May:
		return "mai"
	case time.June:
		return "juin"
	case time.July:
		return "juil."
	case time.August:
		return "août"
	case time.September:
		return "sept."
	case time.October:
		return "oct."
	case time.November:
		return "nov."
	case time.December:
		return "dec."
	}
	return "?"
}
