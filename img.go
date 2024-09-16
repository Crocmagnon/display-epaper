package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/quotes"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed icons
var icons embed.FS

const (
	leftX  = 20
	rightX = 530
)

func getBlack(
	ctx context.Context,
	nowFunc func() time.Time,
	transportsClient *transports.Client,
	feteClient *fete.Client,
	weatherClient *weather.Client,
) (*image.RGBA, error) {
	var (
		bus      *transports.Passages
		tram     *transports.Passages
		velovRoc *transports.Station
		fetes    *fete.Fete
		wthr     *weather.Prevision
	)

	wg := &sync.WaitGroup{}
	wg.Add(5)

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		bus, err = transportsClient.GetTCLPassages(ctx, 290)
		if err != nil {
			log.Println("error getting bus:", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		tram, err = transportsClient.GetTCLPassages(ctx, 34068)
		if err != nil {
			log.Println("error getting tram:", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		velovRoc, err = transportsClient.GetVelovStation(ctx, 10044)
		if err != nil {
			log.Println("error getting velov:", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		fetes, err = feteClient.GetFete(ctx, nowFunc())
		if err != nil {
			log.Println("error getting fetes:", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		wthr, err = weatherClient.GetWeather(ctx)
		if err != nil {
			log.Println("error getting weather:", err)
		}
	}()

	quote := quotes.GetQuote(nowFunc())

	img := newWhite()

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFillRule(draw2d.FillRuleWinding)
	gc.SetFillColor(color.RGBA{255, 255, 255, 255})
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})

	wg.Wait()

	drawTCL(gc, bus, 55)
	drawTCL(gc, tram, 190)
	drawVelov(gc, velovRoc, 350)
	drawDate(gc, nowFunc())
	drawFete(gc, fetes)
	drawWeather(gc, wthr)
	drawQuote(gc, quote)

	return img, nil
}

func drawQuote(gc *draw2dimg.GraphicContext, quote string) {
	text(gc, quote, 15, leftX, 450)

}

func drawWeather(gc *draw2dimg.GraphicContext, wthr *weather.Prevision) {
	if wthr == nil {
		return
	}

	if len(wthr.Daily) == 0 || len(wthr.Daily[0].Weather) == 0 {
		log.Println("missing daily or daily weather")
		return
	}

	daily := wthr.Daily[0]
	dailyWeather := daily.Weather[0]
	err := drawWeatherIcon(gc, dailyWeather)
	if err != nil {
		log.Println("Failed to draw weather icon:", err)
	}

	text(gc, formatTemp(wthr.Current.Temp), 23, leftX, 120)

	text(gc, "max "+formatTemp(daily.Temp.Max), 18, 120, 45)
	text(gc, fmt.Sprintf("pluie %v%%", int(math.Round(daily.Pop*100))), 18, 120, 80)
	text(gc, dailyWeather.Description, 18, 120, 115)
}

func drawWeatherIcon(gc *draw2dimg.GraphicContext, dailyWeather weather.Weather) error {
	icon := strings.TrimSuffix(dailyWeather.Icon, "d")
	icon = strings.TrimSuffix(icon, "n")
	f, err := icons.Open(fmt.Sprintf("icons/%sd.png", icon))
	if err != nil {
		return fmt.Errorf("opening icon: %w", err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decoding icon: %w", err)
	}

	gc.DrawImage(img)
	return nil
}

func formatTemp(temp float64) string {
	return fmt.Sprintf("%v°C", int(math.Round(temp)))
}

func drawVelov(gc *draw2dimg.GraphicContext, station *transports.Station, yOffset float64) {
	if station == nil {
		return
	}

	text(gc, station.Name, 23, rightX, yOffset)
	text(gc, fmt.Sprintf("V : %v - P : %v", station.BikesAvailable, station.DocksAvailable), 23, rightX, yOffset+30)
}

func drawDate(gc *draw2dimg.GraphicContext, now time.Time) {
	text(gc, now.Format("15:04"), 110, leftX, 300)
	text(gc, getDate(now), 30, leftX, 345)
}

func drawFete(gc *draw2dimg.GraphicContext, fetes *fete.Fete) {
	if fetes == nil {
		return
	}

	text(gc, fmt.Sprintf("On fête les %s", fetes.Name), 18, leftX, 380)
}

func drawTCL(gc *draw2dimg.GraphicContext, passages *transports.Passages, yoffset float64) {
	if passages == nil {
		return
	}

	for i, passage := range passages.Passages {
		x := float64(rightX + i*120)
		text(gc, passage.Ligne, 23, x, yoffset)
		for j, delay := range passage.Delays {
			y := yoffset + float64(j+1)*35
			text(gc, delay, 23, x, y)
			if j >= 2 { // limit number of delays displayed
				break
			}
		}
	}
}

func text(gc *draw2dimg.GraphicContext, s string, size, x, y float64) {
	gc.SetFillColor(color.RGBA{0, 0, 0, 255})
	gc.SetFontData(draw2d.FontData{Name: fontName})
	gc.SetFontSize(size)
	gc.FillStringAt(s, x, y)
}

func newWhite() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, epd.Width, epd.Height))
	for y := 0; y < epd.Height; y++ {
		for x := 0; x < epd.Width; x++ {
			img.Set(x, y, color.White)
		}
	}
	return img
}

func rect(gc *draw2dimg.GraphicContext, x1, y1, x2, y2 float64) {
	gc.BeginPath()
	gc.MoveTo(x1, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x1, y2)
	gc.Close()
	gc.FillStroke()
}

func getDate(now time.Time) string {
	return fmt.Sprintf("%v %v %v", getDow(now), getDay(now), getMonth(now))
}

func getDow(now time.Time) string {
	switch now.Weekday() {
	case time.Monday:
		return "lun"
	case time.Tuesday:
		return "mar"
	case time.Wednesday:
		return "mer"
	case time.Thursday:
		return "jeu"
	case time.Friday:
		return "ven"
	case time.Saturday:
		return "sam"
	case time.Sunday:
		return "dim"
	}

	return "?"
}

func getDay(now time.Time) string {
	if now.Day() == 1 {
		return "1er"
	}

	return strconv.Itoa(now.Day())
}

func getMonth(t time.Time) string {
	switch t.Month() {
	case time.January:
		return "jan."
	case time.February:
		return "fev."
	case time.March:
		return "mars"
	case time.April:
		return "avr."
	case time.May:
		return "mai"
	case time.June:
		return "juin"
	case time.July:
		return "juil."
	case time.August:
		return "août"
	case time.September:
		return "sept."
	case time.October:
		return "oct."
	case time.November:
		return "nov."
	case time.December:
		return "dec."
	}
	return "?"
}
