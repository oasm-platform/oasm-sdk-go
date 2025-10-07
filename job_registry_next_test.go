package oasm

import (
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestClient_JobsNext(t *testing.T) {
	tests := []struct {
		name   string
		param  *JobsNextParam
		header *JobsNextHeader
		// want    *JobsNextResponse
		wantErr bool
	}{
		{
			name:   "Success",
			param:  &JobsNextParam{WorkerID: "5ef2dbda-4517-41e1-8b1c-4299d7f66412"},
			header: &JobsNextHeader{WorkerToken: "z6PyLY6fHDMJrS0g1xrbkraVMch1vgYY5aXhBooyFkjXGRFy"},
			// want:    &JobsNextResponse{JobID: "job1", Value: "val", Category: "cat", Command: "cmd"},
			wantErr: false,
		},
		{
			name:   "Non200Status",
			param:  &JobsNextParam{WorkerID: "worker1"},
			header: &JobsNextHeader{WorkerToken: "token1"},
			// want:    nil,
			wantErr: true,
		},
		{
			name:   "InvalidJSON",
			param:  &JobsNextParam{WorkerID: "worker1"},
			header: &JobsNextHeader{WorkerToken: "token1"},
			// want:    nil,
			wantErr: true,
		},
		{
			name:   "EmptyWorkerID",
			param:  &JobsNextParam{WorkerID: ""},
			header: &JobsNextHeader{WorkerToken: "token1"},
			// want:    &JobsNextResponse{},
			wantErr: true,
		},
		{
			name:   "EmptyWorkerToken",
			param:  &JobsNextParam{WorkerID: "worker1"},
			header: &JobsNextHeader{WorkerToken: ""},
			// want:    &JobsNextResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				req:    retryablehttp.NewClient(),
				apiURL: "https://3pnb3328-5173.asse.devtunnels.ms",
				apiKey: "LRvageiyjX8boc6OyApx4nigiJSAexXxfzpo",
			}
			got, err := c.JobsNext(tt.param, tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobsNext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("JobsNext() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
