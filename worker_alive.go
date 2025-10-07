package oasm

import (
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

// WorkerAliveRequest represents the request payload for a worker keep-alive signal.
type WorkerAliveRequest struct {
	Token string `json:"token"`
}

// WorkerAliveResponse represents the response returned by the API
// after receiving a keep-alive signal from a worker.
type WorkerAliveResponse struct {
	Alive string `json:"alive"`
}

// WorkerAlive sends a keep-alive request to the API to indicate that the worker is still active.
// It returns a WorkerAliveResponse on success, or an error if the request fails.
func (c *Client) WorkerAlive(req *WorkerAliveRequest) (*WorkerAliveResponse, error) {
	reqBody, err := sonic.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(c.getAPIURL("/api/workers/alive"), "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, ErrorResponse(body)
	}

	var data WorkerAliveResponse
	if err = sonic.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
