package bitbank

import (
	"net/http"
	"time"
)

const DefaultBaseURL = "https://bitbank.nz"

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}
