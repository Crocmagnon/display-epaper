package fete

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	APIKey        string
	CacheLocation string
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

type Fete struct {
	Day   int    `json:"day"`
	Month int    `json:"month"`
	Name  string `json:"name"`
}

func loadFromDisk(location string) (Fete, error) {
	file, err := os.Open(location)
	if err != nil {
		return Fete{}, fmt.Errorf("opening fetes: %w", err)
	}

	defer file.Close()

	var res Fete
	if err = json.NewDecoder(file).Decode(&res); err != nil {
		return Fete{}, fmt.Errorf("decoding fetes: %w", err)
	}

	return res, nil
}

func (f Fete) dumpToDisk(location string) error {
	file, err := os.Create(location)
	if err != nil {
		return fmt.Errorf("creating fetes: %w", err)
	}

	defer file.Close()

	if err = json.NewEncoder(file).Encode(f); err != nil {
		return fmt.Errorf("dumping fetes: %w", err)
	}
	return nil
}

func (c *Client) GetFete(ctx context.Context, date time.Time) (res *Fete, err error) {
	if val, err := loadFromDisk(c.config.CacheLocation); nil == err {
		log.Println("found fete in cache")
		if val.Day == date.Day() && val.Month == int(date.Month()) {
			log.Println("fete cache is up to date")
			return &val, nil
		}
		log.Println("fete cache is old, fetching...")
	}

	log.Println("querying fete")
	err = requests.URL("https://fetedujour.fr").
		Pathf("/api/v2/%v/json-normal-%d-%d", c.config.APIKey, date.Day(), date.Month()).
		UserAgent("e-paper-display").
		AddValidator(func(resp *http.Response) error {
			if resp.StatusCode >= 400 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				return fmt.Errorf("Fete API error: %s, %s", resp.Header, string(body))
			}
			return nil
		}).
		Client(c.client).
		ToJSON(&res).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("calling API: %w", err)
	}

	if err := res.dumpToDisk(c.config.CacheLocation); err != nil {
		log.Printf("error dumping files to disk: %v\n", err)
	}

	return res, nil
}
