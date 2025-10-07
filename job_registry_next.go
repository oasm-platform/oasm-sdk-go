package oasm

import (
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

type JobsNextParam struct {
	WorkerID string
}

type JobsNextHeader struct {
	WorkerToken string
}

type JobsNextResponse struct {
	JobID    string `json:"jobId,omitempty"`
	Value    string `json:"value,omitempty"`
	Category string `json:"category,omitempty"`
	Command  string `json:"command,omitempty"`
}

func (c *Client) JobsNext(param *JobsNextParam, header *JobsNextHeader) (*JobsNextResponse, error) {
	resp, err := c.GetWithToken(c.getAPIURL("/api/jobs-registry/%s/next", param.WorkerID), header.WorkerToken)
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

	var data JobsNextResponse
	if len(body) > 0 {
		if err = sonic.Unmarshal(body, &data); err != nil {
			return nil, err
		}
	}

	return &data, nil
}
