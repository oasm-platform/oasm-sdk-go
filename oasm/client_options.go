package oasm

import (
	"net/url"
	"strings"
)

// Option represents a functional option for configuring the Client.
type Option func(*Client) error

// WithApiURL sets the base API URL for the Client.
func WithApiURL(apiUrl string) Option {
	return func(c *Client) error {
		parsed, err := url.Parse(apiUrl)
		if err != nil {
			return err
		}

		parsed.Path = strings.TrimRight(parsed.Path, "/")
		c.apiURL = parsed.String()
		return nil
	}
}

// WithApiKey sets the API key for the Client.
func WithApiKey(apiKey string) Option {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}
