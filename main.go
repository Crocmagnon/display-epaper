package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"
	"log"
	"os"
	"time"
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
		APIKey:        os.Getenv("WEATHER_API_KEY"),
		CacheLocation: os.Getenv("WEATHER_CACHE_LOCATION"),
	})

	const minSleep = 1 * time.Second

	sleep, err := time.ParseDuration(os.Getenv("SLEEP_DURATION"))
	if err != nil || sleep < minSleep {
		sleep = minSleep
	}

	const minInitFastThreshold = 1 * time.Second

	initFastThreshold, err := time.ParseDuration(os.Getenv("INIT_FAST_THRESHOLD"))
	if err != nil || initFastThreshold < minInitFastThreshold {
		initFastThreshold = minInitFastThreshold
	}

	log.Printf("sleep duration: %v\n", sleep)

	if err := run(
		ctx,
		sleep,
		initFastThreshold,
		transportsClient,
		feteClient,
		weatherClient,
	); err != nil {
		log.Fatal("error: ", err)
	}

	log.Println("done")
}
