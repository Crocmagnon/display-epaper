package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/quotes"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
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

	transportsClient := transports.New(nil, transports.Config{})

	feteClient := fete.New(nil, fete.Config{
		APIKey:        os.Getenv("FETE_API_KEY"),
		CacheLocation: os.Getenv("FETE_CACHE_LOCATION"),
	})

	weatherClient := weather.New(nil, weather.Config{
		APIKey: os.Getenv("WEATHER_API_KEY"),
	})

	quotesClient := quotes.New(nil, quotes.Config{
		CacheLocation: os.Getenv("QUOTES_CACHE_LOCATION"),
	})

	if err := run(ctx, transportsClient, feteClient, weatherClient, quotesClient); err != nil {
		log.Fatal("error: ", err)
	}

	log.Println("done")
}
