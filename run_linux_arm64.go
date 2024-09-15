package main

import (
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"log"
	"periph.io/x/host/v3"
)

func run() error {
	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("initializing host: %w", err)
	}

	display, err := epd.New()
	if err != nil {
		return fmt.Errorf("initializing epd: %w", err)
	}

	display.Init()
	//display.InitFast()
	display.Clear()
	display.Refresh()

	black, err := getBlack()
	if err != nil {
		return fmt.Errorf("getting black: %w", err)
	}

	red, err := getRed()
	if err != nil {
		return fmt.Errorf("getting red: %w", err)
	}

	display.Send(black, red)

	log.Println("sleeping...")

	if err := display.Sleep(); err != nil {
		return fmt.Errorf("sleeping: %w", err)
	}

	log.Println("done")

	return nil
}
