package oasm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

// Client represents an HTTP client wrapper with configurable options.
// It uses retryablehttp.Client for automatic retries and can be customized
// with an API URL and API key via functional options.
type Client struct {
	req    *retryablehttp.Client
	apiURL string
	apiKey string
}

// NewClient creates a new Client instance with the given options.
//
// Example:
//
//	client, err := NewClient(
//	    WithApiURL("https://api.example.com"),
//	    WithApiKey("my-secret-key"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The functional options (WithApiURL, WithApiKey, etc.) allow you to
// configure the client in a flexible and extensible way without
// changing the constructor signature.
func NewClient(opts ...Option) *Client {
	c := &Client{
		req:    retryablehttp.NewClient(),
		apiURL: "http://localhost:6277",
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) getAPIURL(path string, a ...any) string {
	base := strings.TrimRight(c.apiURL, "/")
	if len(a) > 0 {
		path = fmt.Sprintf(path, a...)
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return base + path
}
