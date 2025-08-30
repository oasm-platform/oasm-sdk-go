package oasm

import (
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
)

type WorkerAliveRequest struct {
	Token string `json:"token"`
}

type WorkerAliveResponse struct {
	Message string `json:"message"`
}

func (c *Client) WorkerAlive(req *WorkerAliveRequest) (*WorkerAliveResponse, error) {
	resp, err := c.req.Post(c.apiURL+"/api/workers/alive", "application/json", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorResponse(body)
	}

	var data WorkerAliveResponse
	if err = sonic.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// WorkerJoinRequest represents the request payload for joining a worker.
type WorkerJoinRequest struct {
	ApiKey string `json:"apiKey"`
}

// WorkerJoinResponse represents the response returned after a worker successfully joins.
type WorkerJoinResponse struct {
	Id               string    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	LastSeenAt       time.Time `json:"lastSeenAt"`
	Token            string    `json:"token"`
	CurrentJobsCount int       `json:"currentJobsCount"`
	Type             string    `json:"type"`
	Scope            string    `json:"scope"`
}

// WorkerJoin sends a join request to the API using the client's ApiKey.
// It returns a WorkerJoinResponse on success, or an error if the request fails.
func (c *Client) WorkerJoin() (*WorkerJoinResponse, error) {
	reqBody, err := sonic.Marshal(&WorkerJoinRequest{
		ApiKey: c.apiKey,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.req.Post(c.apiURL+"/api/workers/join", "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusCreated {
		return nil, ErrorResponse(body)
	}

	var data WorkerJoinResponse
	if err = sonic.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
