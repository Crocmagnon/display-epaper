package transports

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
)

type Stop struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Passage struct {
	Ligne       string   `json:"ligne"`
	Delays      []string `json:"delais"`
	Destination Stop     `json:"destination"`
}

type Passages struct {
	Passages []Passage `json:"passages"`
	Stop     Stop      `json:"stop"`
}

type Config struct {
}

type Client struct {
	client *http.Client
	config Config
}

func New(httpClient *http.Client, config Config) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		client: httpClient,
		config: config,
	}
}

func (c *Client) GetTCLPassages(ctx context.Context, stop int) (res *Passages, err error) {
	err = requests.URL("https://tcl.augendre.info").
		Pathf("/tcl/stop/%v", stop).
		Client(c.client).
		ToJSON(&res).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("calling api: %w", err)
	}
	return res, nil
}

type Station struct {
	Name             string `json:"name"`
	BikesAvailable   int    `json:"bikes_available"`
	DocksAvailable   int    `json:"docks_available"`
	AvailabilityCode int    `json:"availability_code"`
}

func (c *Client) GetVelovStation(ctx context.Context, station int) (res *Station, err error) {
	err = requests.URL("https://tcl.augendre.info").
		Pathf("/velov/station/%v", station).
		Client(c.client).
		ToJSON(&res).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("calling api: %w", err)
	}
	return res, nil
}
