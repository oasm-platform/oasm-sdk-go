package oasm

import (
	"fmt"
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data *WorkerAliveResponse
	if err = sonic.Unmarshal(body, data); err != nil {
		return nil, err
	}

	return data, nil
}

type WorkerJoinRequest struct {
	Token string `json:"token"`
}

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

func (c *Client) WorkerJoin(req *WorkerJoinRequest) (*WorkerJoinResponse, error) {
	resp, err := c.req.Post(c.apiURL+"/api/workers/join", "application/json", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data *WorkerJoinResponse
	if err = sonic.Unmarshal(body, data); err != nil {
		return nil, err
	}

	return data, nil
}
