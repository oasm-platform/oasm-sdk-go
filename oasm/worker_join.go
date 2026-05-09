package oasm

import (
	"context"
	"os"
	"runtime"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
)

func getMetadata() *pb.WorkerMetadata {
	osName := runtime.GOOS
	hostname, err := os.Hostname()
	if err != nil {
		return nil
	}

	return &pb.WorkerMetadata{
		Name: &hostname,
		Os:   &osName,
	}
}

// WorkerJoin sends a join request to the API using the client's ApiKey.
// It returns a WorkerJoinResponse on success, or an error if the request fails.
func (c *Client) WorkerJoin(ctx context.Context) (*pb.JoinResponse, error) {
	req := &pb.JoinRequest{
		ApiKey:   c.apiKey,
		Metadata: getMetadata(),
	}

	resp, err := c.Workers().Join(ctx, req)
	if err != nil {
		return nil, err
	}

	c.workerID = resp.WorkerId
	c.token = resp.WorkerToken

	return resp, nil
}
