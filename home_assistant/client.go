package home_assistant

import (
	"context"
	"github.com/carlmjohnson/requests"
	"net/http"
)

type Config struct {
	Token   string
	BaseURL string
}

type Client struct {
	config Config
	client *http.Client
}

func New(client *http.Client, config Config) *Client {
	return &Client{config: config, client: client}
}

func (c *Client) GetState(ctx context.Context, entityID string) (string, error) {
	type stateResponse struct {
		State string `json:"state"`
	}

	var resp stateResponse

	err := requests.URL(c.config.BaseURL).
		Header("Authorization", "Bearer "+c.config.Token).
		Pathf("/api/states/%s", entityID).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return "", err
	}

	return resp.State, nil
}
