package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"log"
)

func run() error {
	img, err := getBlack()
	if err != nil {
		log.Fatal(err)
	}
	draw2dimg.SaveToPngFile("out/black.png", img)

	img, err = getRed()
	if err != nil {
		log.Fatal(err)
	}
	draw2dimg.SaveToPngFile("out/red.png", img)

	log.Println("done")

	return nil
}
