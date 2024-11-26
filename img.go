package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/Crocmagnon/display-epaper/epd"
	"github.com/Crocmagnon/display-epaper/fete"
	"github.com/Crocmagnon/display-epaper/home_assistant"
	"github.com/Crocmagnon/display-epaper/quotes"
	"github.com/Crocmagnon/display-epaper/transports"
	"github.com/Crocmagnon/display-epaper/weather"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"log/slog"
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

func getImg(ctx context.Context, nowFunc func() time.Time, transportsClient *transports.Client, feteClient *fete.Client, weatherClient *weather.Client, hassClient *home_assistant.Client) (*image.RGBA, error) {
	var (
		bus      *transports.Passages
		tram     *transports.Passages
		velovRoc *transports.Station
		fetes    *fete.Fete
		wthr     *weather.Prevision
		msg      string
	)

	wg := &sync.WaitGroup{}
	wg.Add(6)

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		var err error

		bus, err = transportsClient.GetTCLPassages(ctx, 290)
		if err != nil {
			slog.ErrorContext(ctx, "error getting bus", "err", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		var err error

		tram, err = transportsClient.GetTCLPassages(ctx, 34068)
		if err != nil {
			slog.ErrorContext(ctx, "error getting tram", "err", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		var err error

		velovRoc, err = transportsClient.GetVelovStation(ctx, 10044)
		if err != nil {
			slog.ErrorContext(ctx, "error getting velov", "err", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		fetes, err = feteClient.GetFete(ctx, nowFunc())
		if err != nil {
			slog.ErrorContext(ctx, "error getting fetes", "err", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		wthr, err = weatherClient.GetWeather(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "error getting weather", "err", err)
		}
	}()
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var err error

		msg, err = hassClient.GetState(ctx, "input_text.e_paper_message")
		if err != nil {
			slog.ErrorContext(ctx, "error getting hass message", "err", err)
		}
	}()

	img := newWhite()

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFillRule(draw2d.FillRuleWinding)
	gc.SetFillColor(color.RGBA{255, 255, 255, 255})
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})

	wg.Wait()

	if msg == "" {
		msg = quotes.GetQuote(nowFunc())
	}

	drawTCL(gc, bus, 45)
	drawTCL(gc, tram, 205)
	drawVelov(gc, velovRoc, 365)
	drawDate(gc, nowFunc())
	drawFete(gc, fetes)
	drawWeather(ctx, gc, wthr)
	drawMsg(gc, msg)

	return img, nil
}

func drawMsg(gc *draw2dimg.GraphicContext, quote string) {
	text(gc, quote, 15, leftX, 450, fontItalic)

}

func drawWeather(ctx context.Context, gc *draw2dimg.GraphicContext, wthr *weather.Prevision) {
	if wthr == nil {
		return
	}

	dailyLen := len(wthr.Daily)
	dailyWeatherLen := len(wthr.Daily[0].Weather)
	if dailyLen == 0 || dailyWeatherLen == 0 {
		slog.ErrorContext(ctx, "missing daily or daily weather", "daily_len", dailyLen, "daily_weather_len", dailyWeatherLen)
		return
	}

	daily := wthr.Daily[0]
	dailyWeather := daily.Weather[0]
	err := drawWeatherIcon(gc, wthr.Current.Weather[0])
	if err != nil {
		slog.ErrorContext(ctx, "Failed to draw weather icon", "err", err)
	}

	text(gc, formatTemp(wthr.Current.Temp), 23, leftX, 125, fontRegular)

	const xAlign = 120
	const fontSize = 18

	text(gc, "journée", fontSize, xAlign, 35, fontBold)
	text(gc, "max "+formatTemp(daily.Temp.Max), fontSize, xAlign, 65, fontRegular)
	text(gc, fmt.Sprintf("pluie %v%%", int(math.Round(daily.Pop*100))), fontSize, xAlign, 95, fontRegular)
	text(gc, dailyWeather.Description, fontSize, xAlign, 125, fontRegular)
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

	text(gc, station.Name, 23, rightX, yOffset, fontBold)
	text(gc, fmt.Sprintf("V : %v - P : %v", station.BikesAvailable, station.DocksAvailable), 22, rightX, yOffset+30, fontRegular)
}

func drawDate(gc *draw2dimg.GraphicContext, now time.Time) {
	text(gc, now.Format("15:04"), 110, leftX, 300, fontBold)
	text(gc, getDate(now), 30, leftX, 345, fontRegular)
}

func drawFete(gc *draw2dimg.GraphicContext, fetes *fete.Fete) {
	if fetes == nil {
		return
	}

	text(gc, fmt.Sprintf("On fête les %s", fetes.Name), 18, leftX, 380, fontRegular)
}

func drawTCL(gc *draw2dimg.GraphicContext, passages *transports.Passages, yoffset float64) {
	if passages == nil {
		return
	}

	for i, passage := range passages.Passages {
		x := float64(rightX + i*120)
		text(gc, passage.Ligne, 23, x, yoffset, fontBold)
		for j, delay := range passage.Delays {
			y := yoffset + float64(j+1)*35
			text(gc, delay, 22, x, y, fontRegular)
			if j >= 2 { // limit number of delays displayed
				break
			}
		}
	}
}

func text(gc *draw2dimg.GraphicContext, s string, size, x, y float64, fontName string) {
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
