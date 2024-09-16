package main

import (
	"context"
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
	"log"
	"os"
	"periph.io/x/host/v3"
	"time"
)

func run(
	ctx context.Context,
	sleep time.Duration,
	transportsClient *transports.Client,
	feteClient *fete.Client,
	weatherClient *weather.Client,
) error {
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

		log.Println("running loop")

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

		log.Printf("time.Sleep(%v)\n", sleep)
		time.Sleep(sleep)
	}
}

func loop(
	ctx context.Context,
	display *epd.EPD,
	transportsClient *transports.Client,
	feteClient *fete.Client,
	weatherClient *weather.Client,
) error {
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

	defer func() {
		if err := display.Sleep(); err != nil {
			log.Printf("error sleeping: %v\n", err)
		}
	}()

	err = initDisplay(display)
	if err != nil {
		return fmt.Errorf("initializing display: %w", err)
	}

	display.Clear()

	display.Send(black)
	display.Refresh()

	return nil
}

const filename = "/perm/display-epaper-lastFullRefresh"

func initDisplay(display *epd.EPD) error {
	if canInitFast() {
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

	markInitFull()

	return nil
}

func canInitFast() bool {
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return stat.ModTime().Add(12 * time.Hour).After(time.Now())
}

func markInitFull() {
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("error marking full refresh: %v\n", err)
	}

	f.Close()
}
