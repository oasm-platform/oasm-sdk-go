package oasm

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
)

func (c *Client) WorkerAlive(ctx context.Context) error {
	logger := NewLogger("Alive")

	req := &pb.AliveRequest{
		WorkerToken: c.token,
	}

	stream, err := c.Workers().Alive(ctx, req)
	if err != nil {
		return err
	}

	logger.Success("Connected to core, start capture alive stream...")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			logger.Warning("Server shut down alive stream.")
			break
		}
		if err != nil {
			return err
		}

		logger.Success("Heartbeat - WorkerID: %s, LastSeen: %s, Alive: %v",
			resp.WorkerId, resp.LastSeenAt, resp.Alive)

		if !resp.Alive {
			logger.Error("Worker dead")
			return nil
		}
	}
	return nil
}
