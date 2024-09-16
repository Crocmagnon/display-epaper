package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
	"github.com/llgcode/draw2d/draw2dimg"
	"log"
	"time"
)

func run(
	ctx context.Context,
	_ time.Duration,
	transportsClient *transports.Client,
	feteClient *fete.Client,
	weatherClient *weather.Client,
) error {
	img, err := getImg(
		ctx,
		func() time.Time {
			t, err := time.Parse(time.DateOnly, "2024-08-01zzz")
			if err != nil {
				return time.Now()
			}
			return t
		},
		transportsClient,
		feteClient,
		weatherClient,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := draw2dimg.SaveToPngFile("out/black.png", img); err != nil {
		log.Fatalf("error saving image: %v", err)
	}

	return nil
}
