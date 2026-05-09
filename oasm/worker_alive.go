package oasm

import (
	"context"
	"io"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
)

func (c *Client) WorkerAlive(ctx context.Context) error {
	l := NewLogger("Worker.Alive")

	req := &pb.AliveRequest{
		WorkerToken: c.token,
	}

	stream, err := c.Workers().Alive(ctx, req)
	if err != nil {
		return err
	}

	l.Success("Connected to core, start capture alive stream...")

	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			l.Warning("Server shut down alive stream.")
			break
		}

		if err != nil {
			return err
		}

		l.Success("Heartbeat - WorkerID: %s, LastSeen: %s, Alive: %v", resp.WorkerId, resp.LastSeenAt, resp.Alive)

		if !resp.Alive {
			l.Error("Worker dead")
			return nil
		}
	}

	return nil
}
