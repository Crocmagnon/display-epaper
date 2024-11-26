package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/fonts"
	"github.com/Crocmagnon/display-epaper/home_assistant"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	_ "golang.org/x/image/bmp"
	"log/slog"
	"os"
	"time"
)

const (
	fontRegular = "regular"
	fontBold    = "bold"
	fontItalic  = "italic"
)

func main() {
	ctx := context.Background()

	slog.InfoContext(ctx, "starting...")
	fontCache := MyFontCache{}

	loadFont(ctx, fontCache, fonts.Regular, fontRegular)
	loadFont(ctx, fontCache, fonts.Bold, fontBold)
	loadFont(ctx, fontCache, fonts.Italic, fontItalic)

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

	slog.InfoContext(ctx, "config",
		"sleep_duration", sleep,
		"init_fast_threshold", initFastThreshold)

	hassClient := home_assistant.New(nil, home_assistant.Config{
		Token:   os.Getenv("HOME_ASSISTANT_TOKEN"),
		BaseURL: os.Getenv("HOME_ASSISTANT_BASE_URL"),
	})

	if err := run(
		ctx,
		sleep,
		initFastThreshold,
		transportsClient,
		feteClient,
		weatherClient,
		hassClient,
	); err != nil {
		slog.ErrorContext(ctx, "error", "err", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "done")
}

func loadFont(ctx context.Context, fontCache MyFontCache, ttf []byte, name string) {
	font, err := truetype.Parse(ttf)
	if err != nil {
		slog.ErrorContext(ctx, "error loading font", "err", err)
		os.Exit(1)
	}
	fontCache.Store(draw2d.FontData{Name: name}, font)
}
