package home_assistant

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
	"strconv"
	"time"
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
	var resp struct {
		State string `json:"state"`
	}

	err := c.getState(ctx, entityID, &resp)
	if err != nil {
		return "", err
	}

	return resp.State, nil
}

func (c *Client) GetTimeState(ctx context.Context, entityID string) (time.Time, error) {
	var resp struct {
		State time.Time `json:"state"`
	}

	err := c.getState(ctx, entityID, &resp)
	if err != nil {
		return time.Time{}, err
	}

	return resp.State, nil
}

func (c *Client) GetFloatState(ctx context.Context, entityID string) (float64, error) {
	var resp struct {
		State string `json:"state"`
	}

	err := c.getState(ctx, entityID, &resp)
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseFloat(resp.State, 64)
	if err != nil {
		return 0, fmt.Errorf("converting to float: %w", err)
	}

	return val, nil
}

func (c *Client) getState(ctx context.Context, entityID string, resp any) error {
	err := requests.URL(c.config.BaseURL).
		Header("Authorization", "Bearer "+c.config.Token).
		Pathf("/api/states/%s", entityID).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return err
	}

	return nil
}
