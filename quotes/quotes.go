package quotes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carlmjohnson/requests"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	CacheLocation string
}

type Client struct {
	client *http.Client
	config Config
}

func New(client *http.Client, config Config) *Client {
	return &Client{client: client, config: config}
}

type Quote struct {
	Status   int    `json:"status"`
	Code     int    `json:"code"`
	Error    string `json:"error"`
	Citation struct {
		Citation string `json:"citation"`
		Infos    struct {
			Auteur     string `json:"auteur"`
			Acteur     string `json:"acteur"`
			Personnage string `json:"personnage"`
			Saison     string `json:"saison"`
			Episode    string `json:"episode"`
		} `json:"infos"`
	} `json:"citation"`
}

var errTooOld = errors.New("quote is too old")

func loadFromDisk(location string) (Quote, error) {
	stat, err := os.Stat(location)
	if err != nil {
		return Quote{}, fmt.Errorf("getting file info: %w", err)
	}

	if stat.ModTime().Add(24 * time.Hour).Before(time.Now()) {
		return Quote{}, errTooOld
	}

	file, err := os.Open(location)
	if err != nil {
		return Quote{}, fmt.Errorf("opening quote: %w", err)
	}

	defer file.Close()

	var res Quote
	if err = json.NewDecoder(file).Decode(&res); err != nil {
		return Quote{}, fmt.Errorf("decoding quote: %w", err)
	}

	return res, nil
}

func (q Quote) dumpToDisk(location string) error {
	file, err := os.Create(location)
	if err != nil {
		return fmt.Errorf("creating quote: %w", err)
	}

	defer file.Close()

	if err = json.NewEncoder(file).Encode(q); err != nil {
		return fmt.Errorf("dumping quote: %w", err)
	}
	return nil
}

func (c *Client) GetQuote(ctx context.Context) (res *Quote, err error) {
	if val, err := loadFromDisk(c.config.CacheLocation); nil == err {
		log.Println("found quote in cache")
		return &val, nil
	}

	log.Println("querying quotes")
	err = requests.URL("https://kaamelott.chaudie.re/api/random").
		Client(c.client).
		ToJSON(&res).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching quotes: %w", err)
	}

	if err := res.dumpToDisk(c.config.CacheLocation); err != nil {
		log.Printf("error dumping files to disk: %v\n", err)
	}

	return res, nil
}
