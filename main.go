package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"
	"log"
	"os"
)

const fontName = "default"

func main() {
	log.Println("starting...")

	ctx := context.Background()

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("loading font: %v\n", err)
	}
	fontCache := MyFontCache{}
	fontCache.Store(draw2d.FontData{Name: fontName}, font)
	draw2d.SetFontCache(fontCache)

	transportsClient := transports.New(nil, transports.Config{
		Authorization: os.Getenv("TRANSPORTS_AUTHORIZATION"),
	})

	if err := run(ctx, transportsClient); err != nil {
		log.Fatal("error: ", err)
	}

	log.Println("done")
}
