package epd

import (
	"fmt"
	"image"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3/bcm283x"
	"time"
)

const (
	Width  = 800
	Height = 480
)

type EPD struct {
	width  int
	height int

	resetPin gpio.PinOut
	dcPin    gpio.PinOut
	busyPin  gpio.PinIn
	csPin    gpio.PinOut
	pwrPin   gpio.PinOut

	partFlag byte

	spi    spi.Conn
	spiReg spi.PortCloser
}

func New() (*EPD, error) {
	epd := &EPD{
		width:    Width,
		height:   Height,
		resetPin: bcm283x.GPIO17,
		dcPin:    bcm283x.GPIO25,
		busyPin:  bcm283x.GPIO24,
		csPin:    bcm283x.GPIO8,
		pwrPin:   bcm283x.GPIO18,
		partFlag: 1,
	}

	return epd, nil
}

func (e *EPD) reset() {
	e.resetPin.Out(gpio.High)
	time.Sleep(200 * time.Millisecond)
	e.resetPin.Out(gpio.Low)
	time.Sleep(4 * time.Millisecond)
	e.resetPin.Out(gpio.High)
	time.Sleep(200 * time.Millisecond)
}

func (e *EPD) sendCommand(cmd byte) {
	e.dcPin.Out(gpio.Low)
	e.csPin.Out(gpio.Low)
	if _, err := e.spiWrite([]byte{cmd}); err != nil {
		log.Fatalf("writing to spi: %v", err)
	}
	e.csPin.Out(gpio.High)
}

func (e *EPD) sendData(data byte) {
	e.dcPin.Out(gpio.High)
	e.csPin.Out(gpio.Low)
	if _, err := e.spiWrite([]byte{data}); err != nil {
		log.Fatalf("writing to spi: %v", err)
	}
	e.csPin.Out(gpio.High)
}

func (e *EPD) sendDataSlice(data []byte) {
	e.dcPin.Out(gpio.High)
	toSend := len(data)
	const maxSize = 4096
	if toSend <= maxSize {
		e.csPin.Out(gpio.Low)
		if _, err := e.spiWrite(data); err != nil {
			log.Fatalf("writing to spi: %v", err)
		}
		e.csPin.Out(gpio.High)
		return
	}

	cursor := 0
	for cursor < toSend {
		chunk := data[cursor:min(cursor+maxSize, toSend)]
		e.csPin.Out(gpio.Low)
		if _, err := e.spiWrite(chunk); err != nil {
			log.Fatalf("writing to spi: %v", err)
		}
		e.csPin.Out(gpio.High)
		log.Printf("sent chunk %v\n", cursor)
		cursor = min(cursor+maxSize, toSend)
	}
	log.Printf("sent chunk %v\n", cursor)
}

func (e *EPD) spiWrite(write []byte) ([]byte, error) {
	read := make([]byte, len(write))

	if err := e.spi.Tx(write, read); err != nil {
		return nil, fmt.Errorf("tx: %w", err)
	}

	return read, nil
}

func (e *EPD) readBusy() {
	e.sendCommand(0x71)
	busy := e.busyPin.Read()
	for busy == gpio.Low {
		e.sendCommand(0x71)
		busy = e.busyPin.Read()
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)
}

func (e *EPD) turnOn() error {
	log.Println("turning on")
	if err := e.resetPin.Out(gpio.Low); err != nil {
		return fmt.Errorf("setting reset pin to low: %w", err)
	}
	if err := e.dcPin.Out(gpio.Low); err != nil {
		return fmt.Errorf("setting dc pin to low: %w", err)
	}
	if err := e.csPin.Out(gpio.Low); err != nil {
		return fmt.Errorf("setting cs pin to low: %w", err)
	}
	if err := e.pwrPin.Out(gpio.High); err != nil {
		return fmt.Errorf("setting pwr pin to low: %w", err)
	}

	var err error

	if e.spiReg, err = spireg.Open("0"); err != nil {
		return fmt.Errorf("opening SPI: %w", err)
	}

	c, err := e.spiReg.Connect(4*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return fmt.Errorf("connecting to SPI: %w", err)
	}

	e.spi = c

	return nil
}

func (e *EPD) Init() error {
	log.Println("initializing EPD")

	if err := e.turnOn(); err != nil {
		return fmt.Errorf("turning on: %w", err)
	}

	e.reset()

	e.sendCommand(0x01)
	e.sendDataSlice([]byte{0x07, 0x07, 0x3f, 0x3f})

	e.sendCommand(0x06)
	e.sendDataSlice([]byte{0x17, 0x17, 0x28, 0x17})

	e.sendCommand(0x04)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()

	e.sendCommand(0x00)
	e.sendData(0x0f)

	e.sendCommand(0x61)
	e.sendDataSlice([]byte{0x03, 0x20, 0x01, 0xe0})

	e.sendCommand(0x15)
	e.sendData(0x00)

	e.sendCommand(0x50)
	e.sendDataSlice([]byte{0x11, 0x07})

	e.sendCommand(0x60)
	e.sendData(0x22)

	return nil
}

func (e *EPD) InitFast() error {
	log.Println("initializing Fast EPD")

	if err := e.turnOn(); err != nil {
		return fmt.Errorf("turning on: %w", err)
	}

	e.reset()

	e.sendCommand(0x00)
	e.sendData(0x0f)

	e.sendCommand(0x04)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()

	e.sendCommand(0x06)
	e.sendDataSlice([]byte{0x27, 0x27, 0x18, 0x17})

	e.sendCommand(0xe0)
	e.sendData(0x02)

	e.sendCommand(0xe5)
	e.sendData(0x5a)

	e.sendCommand(0x50)
	e.sendDataSlice([]byte{0x11, 0x07})

	return nil
}

func (e *EPD) Clear() {
	log.Println("clearing epd")
	e.Fill(White)
}

func (e *EPD) Refresh() {
	log.Println("refreshing...")
	e.sendCommand(0x12)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()
}

func (e *EPD) Sleep() error {
	log.Println("sleeping display...")
	e.sendCommand(0x02)
	e.readBusy()

	e.sendCommand(0x07)
	e.sendData(0xa5)

	time.Sleep(2 * time.Second)
	if err := e.turnOff(); err != nil {
		return fmt.Errorf("turning off: %w", err)
	}

	return nil
}

func (e *EPD) turnOff() error {
	log.Println("turning off...")
	if err := e.spiReg.Close(); err != nil {
		return fmt.Errorf("closing SPI: %w", err)
	}

	e.resetPin.Out(gpio.Low)
	e.dcPin.Out(gpio.Low)
	e.pwrPin.Out(gpio.Low)

	return nil
}

type Color int

const (
	White Color = iota
	Red
	Black
)

func (e *EPD) Fill(c Color) {
	log.Println("filling...")

	switch c {
	case White:
		e.Send(image.White, image.Black)
	case Black:
		e.Send(image.Black, image.Black)
	case Red:
		e.Send(image.White, image.White)
	}
}

func (e *EPD) Send(black image.Image, red image.Image) {
	if black != nil {
		log.Println("sending black")
		e.sendCommand(0x10) // write bw data
		e.sendImg(black)
	}
	if red != nil {
		log.Println("sending red")
		e.sendCommand(0x13) // write red data
		e.sendImg(red)
	}
}

func (e *EPD) sendImg(img image.Image) {
	log.Println("sending img...")
	// TODO check img size
	for row := 0; row < e.height; row++ {
		for col := 0; col < e.width; col += 8 {
			// this loop converts individual pixels into a single byte
			// 8-pixels at a time and then sends that byte to render
			var b byte = 0xFF
			for px := 0; px < 8; px++ {
				var pixel = img.At(col+px, row)
				if isdark(pixel.RGBA()) {
					b &= ^(0x80 >> (px % 8))
				}
			}
			e.sendData(b)
		}
	}
}

func isdark(r, g, b, _ uint32) bool {
	return r < 255 || g < 255 || b < 255
}
