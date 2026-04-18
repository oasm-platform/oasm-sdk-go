package oasm

import (
	"context"
	"fmt"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/jobs_registry"
)

func (c *Client) JobsResult(ctx context.Context, jobID string, payload *pb.DataPayloadResult) error {
	req := &pb.JobResultRequest{
		WorkerId: c.workerID,
		Data: &pb.UpdateResultDto{
			JobId: jobID,
			Data:  payload,
		},
	}

	resp, err := c.Jobs().Result(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to submit job result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the result submission")
	}

	return nil
}
