package main

import (
	"context"
	"fmt"
	"github.com/Crocmagnon/display-epaper/home_assistant"
	"github.com/Crocmagnon/display-epaper/weather"
	"github.com/llgcode/draw2d/draw2dimg"
	"time"
)

func run(
	ctx context.Context,
	_ time.Duration,
	_ time.Duration,
	weatherClient *weather.Client,
	hassClient *home_assistant.Client,
) error {
	img, err := getImg(ctx, func() time.Time {
		t, err := time.Parse(time.DateOnly, "2024-08-01zzz")
		if err != nil {
			return time.Now()
		}
		return t
	}, weatherClient, hassClient)
	if err != nil {
		return err
	}

	if err := draw2dimg.SaveToPngFile("out/black.png", img); err != nil {
		return fmt.Errorf("saving img: %w", err)
	}

	return nil
}
