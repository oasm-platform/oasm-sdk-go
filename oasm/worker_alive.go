package oasm

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
)

func (c *Client) WorkerAlive(ctx context.Context) error {
	req := &pb.AliveRequest{
		WorkerToken: c.token,
	}

	stream, err := c.Workers().Alive(ctx, req)
	if err != nil {
		return err
	}

	log.Println("Connected to core, start capture alive stream...")

	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			log.Println("Server shut down alive stream.")
			break
		}

		if err != nil {
			return err
		}

		Logger("Alive").Success(fmt.Sprintf("Heartbeat - WorkerID: %s, LastSeen: %s, Alive: %v",
			resp.WorkerId, resp.LastSeenAt, resp.Alive))

		if !resp.Alive {
			Logger("Alive").Error("Worker dead")
			return nil
		}
	}

	return nil
}
