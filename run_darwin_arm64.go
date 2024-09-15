package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/llgcode/draw2d/draw2dimg"
	"log"
	"time"
)

func run(ctx context.Context, transportsClient *transports.Client, feteClient *fete.Client) error {
	img, err := getBlack(ctx, func() time.Time {
		t, err := time.Parse(time.DateOnly, "2024-08-01zzz")
		if err != nil {
			return time.Now()
		}
		return t
	}, transportsClient, feteClient)
	if err != nil {
		log.Fatal(err)
	}
	if err := draw2dimg.SaveToPngFile("out/black.png", img); err != nil {
		log.Fatalf("error saving image: %v", err)
	}

	log.Println("done")

	return nil
}
