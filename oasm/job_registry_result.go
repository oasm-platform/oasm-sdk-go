package oasm

import (
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

type JobsResultParam struct {
	WorkerID string
}

type JobsResultRequest struct {
	JobID string `json:"jobId,omitempty"`
	Data  struct {
		Error   bool   `json:"error,omitempty"`
		Raw     string `json:"raw,omitempty"`
		Payload any    `json:"payload,omitempty"`
	} `json:"data,omitempty"`
}

func (c *Client) JobsResult(param *JobsResultParam, req *JobsResultRequest) error {
	reqBody, err := sonic.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.Post(c.getAPIURL("/api/jobs-registry/%s/result", param.WorkerID), "application/json", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return ErrorResponse(body)
	}

	return nil
}
