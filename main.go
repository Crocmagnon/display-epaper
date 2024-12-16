package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/fonts"
	"github.com/Crocmagnon/display-epaper/home_assistant"
	"github.com/Crocmagnon/display-epaper/weather"
	"github.com/llgcode/draw2d"
	_ "golang.org/x/image/bmp"
	"log/slog"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	slog.InfoContext(ctx, "starting...")

	fontCache, err := fonts.NewCache()
	if err != nil {
		slog.ErrorContext(ctx, "could not create font cache", "error", err.Error())
		os.Exit(1)
	}

	draw2d.SetFontCache(fontCache)

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

	if err := run(ctx, sleep, initFastThreshold, weatherClient, hassClient); err != nil {
		slog.ErrorContext(ctx, "error", "err", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "done")
}
