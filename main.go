package main

import (
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"log"
	"os"
	"periph.io/x/host/v3"
)

const fontName = "default"

func main() {
	log.Println("starting...")

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("loading font: %v\n", err)
	}
	fontCache := MyFontCache{}
	fontCache.Store(draw2d.FontData{Name: fontName}, font)
	draw2d.SetFontCache(fontCache)

	if len(os.Args) < 2 {
		if err := run(); err != nil {
			log.Fatal("error: ", err)
		}
	} else {
		img, err := getBlack()
		if err != nil {
			log.Fatal(err)
		}
		draw2dimg.SaveToPngFile("black.png", img)

		img, err = getRed()
		if err != nil {
			log.Fatal(err)
		}
		draw2dimg.SaveToPngFile("red.png", img)

	}

	log.Println("done")
}

func run() error {
	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("initializing host: %w", err)
	}

	display, err := epd.New()
	if err != nil {
		return fmt.Errorf("initializing epd: %w", err)
	}

	display.Init()
	display.Clear()

	black, err := getBlack()
	if err != nil {
		return fmt.Errorf("getting black: %w", err)
	}

	red, err := getRed()
	if err != nil {
		return fmt.Errorf("getting red: %w", err)
	}

	display.Draw(black, red)

	log.Println("sleeping...")

	if err := display.Sleep(); err != nil {
		return fmt.Errorf("sleeping: %w", err)
	}

	log.Println("done")

	return nil
}

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

type MyFontCache map[string]*truetype.Font

func (fc MyFontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc MyFontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("font %s is not stored in font cache", fd.Name)
	}
	return font, nil
}
