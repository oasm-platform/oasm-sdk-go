package oasm

import (
	"context"
	"fmt"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/jobs_registry"
)

func (c *Client) JobsNext(ctx context.Context) (*pb.Job, error) {
	job, err := c.Jobs().Next(c.WithAuth(ctx), &pb.Worker{
		Id: c.workerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch next job: %w", err)
	}

	if job == nil || job.Id == "" {
		return nil, nil
	}

	return job, nil
}
