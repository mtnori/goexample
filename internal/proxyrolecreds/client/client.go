package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct {
	Endpoint   string
	HTTPClient HTTPClient
}

type Client struct {
	options Options
}

func New(options Options) *Client {
	client := &Client{
		options: options,
	}

	return client
}

func (c *Client) GetCredentials(_ context.Context) (*GetCredentialsOutput, error) {

	req, _ := http.NewRequest("GET", c.options.Endpoint, nil)
	resp, err := c.options.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // TODO: error handling

	var output GetCredentialsOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}
	return &output, nil
}

type GetCredentialsOutput struct {
	AccessKeyID     string
	SecretAccessKey string
	Token           string
	Expiration      *time.Time
}
