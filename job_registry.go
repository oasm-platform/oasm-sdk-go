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

	resp, err := c.req.Post(c.apiURL+"/api/jobs-registry/"+param.WorkerID+"/result", "application/json", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrorResponse(body)
	}

	return nil
}
