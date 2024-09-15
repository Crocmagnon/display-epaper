package main

import (
	"context"
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
	"log"
	"periph.io/x/host/v3"
	"time"
)

func run(ctx context.Context, transportsClient *transports.Client, feteClient *fete.Client, weatherClient *weather.Client) error {
	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("initializing host: %w", err)
	}

	display, err := epd.New()
	if err != nil {
		return fmt.Errorf("initializing epd: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("stopping because of context")
			return ctx.Err()
		default:
		}

		err = loop(
			ctx,
			display,
			transportsClient,
			feteClient,
			weatherClient,
		)
		if err != nil {
			log.Printf("error looping: %v\n", err)
		}

		log.Println("time.Sleep(30s)")
		time.Sleep(30 * time.Second)
	}
}

func loop(
	ctx context.Context,
	display *epd.EPD,
	transportsClient *transports.Client,
	feteClient *fete.Client,
	weatherClient *weather.Client,
) error {
	defer func() {
		if err := display.Sleep(); err != nil {
			log.Printf("error sleeping: %v\n", err)
		}
	}()

	err := display.Init()
	if err != nil {
		return fmt.Errorf("initializing display: %w", err)
	}

	display.Clear()

	black, err := getBlack(
		ctx,
		time.Now,
		transportsClient,
		feteClient,
		weatherClient,
	)
	if err != nil {
		return fmt.Errorf("getting black: %w", err)
	}

	display.Send(black, nil)
	display.Refresh()

	return nil
}
