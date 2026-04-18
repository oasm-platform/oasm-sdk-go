package oasm

import (
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/grpc"
)

// Option represents a functional option for configuring the Client.
type Option func(*Client) error

// WithGRPCHost sets the base grpcHost for the Client.
func WithGRPCHost(grpcHost string) Option {
	return func(c *Client) error {
		parsed, err := url.Parse(grpcHost)
		if err != nil {
			return err
		}

		parsed.Path = strings.TrimRight(parsed.Path, "/")
		c.grpcHost = parsed.String()
		return nil
	}
}

// WithApiKey sets the API key for the Client.
func WithApiKey(apiKey string) Option {
	return func(c *Client) error {
		if apiKey == "" {
			return fmt.Errorf("Invalid api key")
		}

		c.apiKey = apiKey
		return nil
	}
}

// WithConn sets a custom grpc.ClientConn for the Client.
func WithConn(conn *grpc.ClientConn) Option {
	return func(c *Client) error {
		if conn == nil {
			return fmt.Errorf("Connection must not nil")
		}

		c.conn = conn
		return nil
	}
}

func WithConfigPath(path string) Option {
	return func(c *Client) error {
		if path == "" {
			return fmt.Errorf("Config path must not empty")
		}

		c.configPath = path
		return nil
	}
}

func WithToolPath(path string) Option {
	return func(c *Client) error {
		if path == "" {
			return fmt.Errorf("Tool path must not empty")
		}

		c.toolPath = path
		return nil
	}
}
