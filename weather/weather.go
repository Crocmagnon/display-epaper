package weather

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
)

type Config struct {
	APIKey string
}

type Client struct {
	config Config
	client *http.Client
}

func New(httpClient *http.Client, config Config) *Client {
	return &Client{
		config: config,
		client: httpClient,
	}
}

type Prevision struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		Dt         int     `json:"dt"`
		Sunrise    int     `json:"sunrise"`
		Sunset     int     `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        float64 `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		WindGust   float64 `json:"wind_gust"`
		Weather    []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Daily  []Daily `json:"daily"`
	Alerts []struct {
		SenderName  string   `json:"sender_name"`
		Event       string   `json:"event"`
		Start       int      `json:"start"`
		End         int      `json:"end"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	} `json:"alerts"`
}

type Daily struct {
	Dt        int     `json:"dt"`
	Sunrise   int     `json:"sunrise"`
	Sunset    int     `json:"sunset"`
	Moonrise  int     `json:"moonrise"`
	Moonset   int     `json:"moonset"`
	MoonPhase float64 `json:"moon_phase"`
	Summary   string  `json:"summary"`
	Temp      struct {
		Day   float64 `json:"day"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	} `json:"temp"`
	FeelsLike struct {
		Day   float64 `json:"day"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	} `json:"feels_like"`
	Pressure  int       `json:"pressure"`
	Humidity  int       `json:"humidity"`
	DewPoint  float64   `json:"dew_point"`
	WindSpeed float64   `json:"wind_speed"`
	WindDeg   int       `json:"wind_deg"`
	WindGust  float64   `json:"wind_gust"`
	Weather   []Weather `json:"weather"`
	Clouds    int       `json:"clouds"`
	Pop       float64   `json:"pop"`
	Rain      float64   `json:"rain"`
	Uvi       float64   `json:"uvi"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (c *Client) GetWeather(ctx context.Context) (res *Prevision, err error) {
	err = requests.URL("https://api.openweathermap.org/data/3.0/onecall").
		Client(c.client).
		Param("lat", "45.78").
		Param("lon", "4.89").
		Param("appid", c.config.APIKey).
		Param("units", "metric").
		Param("lang", "fr").
		Param("exclude", "minutely,hourly").
		ToJSON(&res).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("calling openweathermap: %w", err)
	}

	return res, nil
}
