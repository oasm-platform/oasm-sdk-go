package oasm

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestClient_JobsNext(t *testing.T) {
	tests := []struct {
		name    string
		param   *JobsNextParam
		header  *JobsNextHeader
		want    *JobsNextResponse
		wantErr bool
	}{
		{
			name:    "Success",
			param:   &JobsNextParam{WorkerID: "worker1"},
			header:  &JobsNextHeader{WorkerToken: "token1"},
			want:    &JobsNextResponse{JobID: "job1", Value: "val", Category: "cat", Command: "cmd"},
			wantErr: false,
		},
		{
			name:    "Non200Status",
			param:   &JobsNextParam{WorkerID: "worker1"},
			header:  &JobsNextHeader{WorkerToken: "token1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "InvalidJSON",
			param:   &JobsNextParam{WorkerID: "worker1"},
			header:  &JobsNextHeader{WorkerToken: "token1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "EmptyWorkerID",
			param:   &JobsNextParam{WorkerID: ""},
			header:  &JobsNextHeader{WorkerToken: "token1"},
			want:    &JobsNextResponse{},
			wantErr: false,
		},
		{
			name:    "EmptyWorkerToken",
			param:   &JobsNextParam{WorkerID: "worker1"},
			header:  &JobsNextHeader{WorkerToken: ""},
			want:    &JobsNextResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				req:    retryablehttp.NewClient(),
				apiURL: "localhost:6276",
			}
			got, err := c.JobsNext(tt.param, tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobsNext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobsNext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name:  "Non200Status",
			param: &JobsResultParam{WorkerID: "worker1"},
			req:   &JobsResultRequest{JobID: "job1"},

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
