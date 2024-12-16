package main

import (
	"context"
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/Crocmagnon/display-epaper/home_assistant"
	"github.com/Crocmagnon/display-epaper/weather"
	"image"
	"log"
	"log/slog"
	"os"
	"periph.io/x/host/v3"
	"time"
)

func run(
	ctx context.Context,
	sleep time.Duration,
	initFastThreshold time.Duration,
	weatherClient *weather.Client,
	hassClient *home_assistant.Client,
) error {
	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("initializing host: %w", err)
	}

	display, err := epd.New()
	if err != nil {
		return fmt.Errorf("initializing epd: %w", err)
	}

	var currentImg image.Image

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "stopping because of context")
			return ctx.Err()
		default:
		}

		slog.InfoContext(ctx, "running loop")

		img, err := loop(ctx, display, initFastThreshold, currentImg, weatherClient, hassClient)
		if err != nil {
			slog.ErrorContext(ctx, "error looping", "err", err)
		}

		currentImg = img

		slog.InfoContext(ctx, "sleep", "duration", sleep)
		time.Sleep(sleep)
	}
}

func loop(
	ctx context.Context,
	display *epd.EPD,
	initFastThreshold time.Duration,
	currentImg image.Image,
	weatherClient *weather.Client,
	hassClient *home_assistant.Client,
) (image.Image, error) {
	var img image.Image = image.White

	if shouldDisplay(ctx, hassClient) {
		var err error
		img, err = getImg(ctx, time.Now, weatherClient, hassClient)
		if err != nil {
			return nil, fmt.Errorf("getting black: %w", err)
		}
	}

	if imgEqual(currentImg, img, epd.Width, epd.Height) {
		slog.InfoContext(ctx, "Images are equal, doing nothing.")
		return img, nil
	}

	defer func() {
		if err := display.Sleep(); err != nil {
			slog.ErrorContext(ctx, "error sleeping", "err", err)
		}
	}()

	err := initDisplay(ctx, display, initFastThreshold)
	if err != nil {
		return nil, fmt.Errorf("initializing display: %w", err)
	}

	display.Clear()

	display.Send(img)
	display.Refresh()

	return img, nil
}

func shouldDisplay(ctx context.Context, hassClient *home_assistant.Client) bool {
	dayNight, err := hassClient.GetState(ctx, "input_select.house_day_night")
	if err != nil {
		log.Printf("error getting day night: %v ; displaying anyway\n", err)
		return true
	}

	hallLights, err := hassClient.GetState(ctx, "light.couloir")
	if err != nil {
		log.Printf("error getting hall lights: %v ; displaying anyway\n", err)
		return true
	}

	presentAway, err := hassClient.GetState(ctx, "input_select.house_present_away")
	if err != nil {
		log.Printf("error getting day night: %v ; displaying anyway\n", err)
		return true
	}

	slog.InfoContext(ctx, "home assistant states",
		"hall_lights", hallLights,
		"day_night", dayNight,
		"present_away", presentAway)

	res := (hallLights == "on" || dayNight == "day") && presentAway == "present"
	slog.InfoContext(ctx, "computed should display", "should_display", res)
	return res
}

const filename = "/perm/display-epaper-lastFullRefresh"

func initDisplay(ctx context.Context, display *epd.EPD, threshold time.Duration) error {
	if canInitFast(threshold) {
		err := display.InitFast()
		if err != nil {
			return fmt.Errorf("running fast init: %w", err)
		}
		return nil
	}

	err := display.Init()
	if err != nil {
		return fmt.Errorf("running full init: %w", err)
	}

	markInitFull(ctx)

	return nil
}

func canInitFast(threshold time.Duration) bool {
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return stat.ModTime().Add(threshold).After(time.Now())
}

func markInitFull(ctx context.Context) {
	f, err := os.Create(filename)
	if err != nil {
		slog.ErrorContext(ctx, "error marking full refresh", "err", err)
	}

	f.Close()
}

func imgEqual(img1, img2 image.Image, width, height int) bool {
	if img1 == nil || img2 == nil {
		return false
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				return false
			}
		}
	}

	return true
}
