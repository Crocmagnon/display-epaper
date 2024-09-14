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

	if err := epd.resetPin.Out(gpio.Low); err != nil {
		return nil, fmt.Errorf("setting reset pin to low: %w", err)
	}
	if err := epd.dcPin.Out(gpio.Low); err != nil {
		return nil, fmt.Errorf("setting dc pin to low: %w", err)
	}
	if err := epd.csPin.Out(gpio.Low); err != nil {
		return nil, fmt.Errorf("setting cs pin to low: %w", err)
	}
	if err := epd.pwrPin.Out(gpio.High); err != nil {
		return nil, fmt.Errorf("setting pwr pin to low: %w", err)
	}

	var err error

	if epd.spiReg, err = spireg.Open("0"); err != nil {
		return nil, fmt.Errorf("opening SPI: %w", err)
	}

	c, err := epd.spiReg.Connect(4*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, fmt.Errorf("connecting to SPI: %w", err)
	}

	epd.spi = c

	return epd, nil
}

func (e *EPD) TurnOff() error {
	log.Println("turning off...")
	if err := e.spiReg.Close(); err != nil {
		return fmt.Errorf("closing SPI: %w", err)
	}

	e.resetPin.Out(gpio.Low)
	e.dcPin.Out(gpio.Low)
	e.pwrPin.Out(gpio.Low)

	return nil
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
	log.Printf("sending command 0x%02X\n", cmd)
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
	log.Printf("sending data slice %v\n", len(data))
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

func (e *EPD) Init() {
	log.Println("initializing EPD")
	e.reset()

	e.sendCommand(0x01)
	e.sendData(0x07)
	e.sendData(0x07)
	e.sendData(0x3f)
	e.sendData(0x3f)

	e.sendCommand(0x06)
	e.sendData(0x17)
	e.sendData(0x17)
	e.sendData(0x28)
	e.sendData(0x17)

	e.sendCommand(0x04)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()

	e.sendCommand(0x00)
	e.sendData(0x0f)

	e.sendCommand(0x61)
	e.sendData(0x03)
	e.sendData(0x20)
	e.sendData(0x01)
	e.sendData(0xe0)

	e.sendCommand(0x15)
	e.sendData(0x00)

	e.sendCommand(0x50)
	e.sendData(0x11)
	e.sendData(0x07)

	e.sendCommand(0x60)
	e.sendData(0x22)
}

func (e *EPD) Clear() {
	log.Println("clearing epd")
	redBuf := make([]byte, Width*Height/8)
	for i := range redBuf {
		redBuf[i] = 0x00
	}

	blackBuf := make([]byte, Width*Height/8)
	for i := range blackBuf {
		blackBuf[i] = 0xff
	}

	e.sendCommand(0x10)
	e.sendDataSlice(blackBuf)

	e.sendCommand(0x13)
	e.sendDataSlice(redBuf)

	//e.refresh()
}

func (e *EPD) refresh() {
	log.Println("refreshing...")
	e.sendCommand(0x12)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()
}

func (e *EPD) Sleep() error {
	log.Println("sleeping...")
	e.sendCommand(0x02)
	e.readBusy()

	e.sendCommand(0x07)
	e.sendData(0xa5)

	time.Sleep(2 * time.Second)
	if err := e.TurnOff(); err != nil {
		return fmt.Errorf("turning off: %w", err)
	}

	return nil
}

type Color int

const (
	White Color = iota
	Red
	Black
)

func (e *EPD) Fill(c Color) {
	log.Println("filling... (not doing anything yet)")

	//switch c {
	//case White:
	//	e.Draw(nil, nil)
	//case Black:
	//	e.Draw(image.Black, nil)
	//case Red:
	//	e.Draw(nil, image.Black)
	//}
}

func (e *EPD) Draw(black image.Image, red image.Image) {
	log.Println("drawing...")
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
	if black != nil || red != nil {
		e.refresh()
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
