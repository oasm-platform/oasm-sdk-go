package oasm

import (
	"fmt"
	"sync"

	jobRegistryPb "github.com/oasm-platform/open-asm/grpc-client/go/jobs_registry"
	workerPb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents an grpc wrapper with configurable options.
// with an API URL and API key via functional options.
type Client struct {
	conn       *grpc.ClientConn
	grpcHost   string
	apiKey     string
	configPath string
	token      string
	toolPath   string

	workerOnce sync.Once
	worker     workerPb.WorkersServiceClient

	jobOnce sync.Once
	job     jobRegistryPb.JobsRegistryServiceClient
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
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		grpcHost:   "localhost:16276",
		configPath: "config.json",
		toolPath:   "oasm-tools",
	}

	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, fmt.Errorf("option error: %w", err)
		}
	}

	if c.conn == nil {
		conn, err := grpc.NewClient(c.grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to dial: %w", err)
		}
		c.conn = conn
	}

	return c, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Workers() workerPb.WorkersServiceClient {
	c.workerOnce.Do(func() {
		c.worker = workerPb.NewWorkersServiceClient(c.conn)
	})
	return c.worker
}

func (c *Client) Jobs() jobRegistryPb.JobsRegistryServiceClient {
	c.jobOnce.Do(func() {
		c.job = jobRegistryPb.NewJobsRegistryServiceClient(c.conn)
	})
	return c.job
}
