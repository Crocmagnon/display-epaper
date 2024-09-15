package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"
	"log"
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

	if err := run(); err != nil {
		log.Fatal("error: ", err)
	}

	log.Println("done")
}
