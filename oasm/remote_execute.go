package oasm

import (
	"context"
	"fmt"
	"io"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
)

type RemoteExecuteHandler struct {
	sessionID string
	id        string
	workerID  string
	stream    pb.WorkersService_RemoteExecuteSubscribeClient
	client    *Client
}

func (c *Client) RemoteExecuteSubscribe(ctx context.Context) (*RemoteExecuteHandler, error) {
	req := &pb.RemoteExecuteSubscribeRequest{}

	stream, err := c.Workers().RemoteExecuteSubscribe(c.WithAuth(ctx), req)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to remote execute: %w", err)
	}

	return &RemoteExecuteHandler{
		stream: stream,
		client: c,
	}, nil
}

func (h *RemoteExecuteHandler) Next(ctx context.Context) (*pb.RemoteExecuteSubscribeResponse, error) {
	resp, err := h.stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to receive remote execute event: %w", err)
	}

	if resp != nil {
		h.sessionID = resp.SessionId
		h.id = resp.Id
		h.workerID = resp.WorkerId
	}

	return resp, nil
}

func (h *RemoteExecuteHandler) SendResult(ctx context.Context, eventType pb.RemoteExecuteResultEventType, data []byte, exitCode int32) error {
	req := &pb.RemoteExecuteResultStream{
		Id:        h.id,
		SessionId: h.sessionID,
		Type:      eventType,
		Data:      data,
		ExitCode:  exitCode,
	}

	resp, err := h.client.Workers().RemoteExecuteResult(h.client.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to send remote execute result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the result: %s", resp.Message)
	}

	return nil
}

func (h *RemoteExecuteHandler) SendStdout(ctx context.Context, data []byte) error {
	return h.SendResult(ctx, pb.RemoteExecuteResultEventType_REMOTE_EXECUTE_RESULT_STDOUT, data, 0)
}

func (h *RemoteExecuteHandler) SendStderr(ctx context.Context, data []byte) error {
	return h.SendResult(ctx, pb.RemoteExecuteResultEventType_REMOTE_EXECUTE_RESULT_STDERR, data, 0)
}

func (h *RemoteExecuteHandler) SendExit(ctx context.Context, exitCode int32) error {
	return h.SendResult(ctx, pb.RemoteExecuteResultEventType_REMOTE_EXECUTE_RESULT_EXIT, nil, exitCode)
}

func (h *RemoteExecuteHandler) SendError(ctx context.Context, errMsg string) error {
	return h.SendResult(ctx, pb.RemoteExecuteResultEventType_REMOTE_EXECUTE_RESULT_ERROR, []byte(errMsg), 0)
}

func (h *RemoteExecuteHandler) SessionID() string {
	return h.sessionID
}

func (h *RemoteExecuteHandler) ID() string {
	return h.id
}

func (h *RemoteExecuteHandler) WorkerID() string {
	return h.workerID
}
