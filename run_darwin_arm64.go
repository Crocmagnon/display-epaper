package main

import (
	"context"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/llgcode/draw2d/draw2dimg"
	"log"
)

func run(ctx context.Context, transportsClient *transports.Client) error {
	img, err := getBlack(ctx, transportsClient)
	if err != nil {
		log.Fatal(err)
	}
	draw2dimg.SaveToPngFile("out/black.png", img)

	log.Println("done")

	return nil
}
