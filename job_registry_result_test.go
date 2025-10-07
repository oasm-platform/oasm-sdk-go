package oasm

import (
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestClient_JobsResult(t *testing.T) {
	tests := []struct {
		name    string
		param   *JobsResultParam
		req     *JobsResultRequest
		wantErr bool
	}{
		{
			name:  "Success",
			param: &JobsResultParam{WorkerID: "worker1"},
			req: &JobsResultRequest{
				JobID: "job1",
				Data: struct {
					Error   bool   `json:"error,omitempty"`
					Raw     string `json:"raw,omitempty"`
					Payload any    `json:"payload,omitempty"`
				}{
					Error:   false,
					Raw:     "raw data",
					Payload: "payload",
				},
			},

			wantErr: false,
		},
		{
			name:    "Non200Status",
			param:   &JobsResultParam{WorkerID: "worker1"},
			req:     &JobsResultRequest{JobID: "job1"},
			wantErr: true,
		},
		{
			name:    "EmptyWorkerID",
			param:   &JobsResultParam{WorkerID: ""},
			req:     &JobsResultRequest{JobID: "job1"},
			wantErr: false,
		},
		{
			name:    "NilRequest",
			param:   &JobsResultParam{WorkerID: "worker1"},
			req:     nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				req:    retryablehttp.NewClient(),
				apiURL: "localhost:6276",
			}
			err := c.JobsResult(tt.param, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobsResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
