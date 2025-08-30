package oasm

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestClient_WorkerJoin(t *testing.T) {
	type fields struct {
		req    *retryablehttp.Client
		apiURL string
		apiKey string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *WorkerJoinResponse
		wantErr bool
	}{
		{
			name: "TestError",
			fields: fields{
				req:    retryablehttp.NewClient(),
				apiURL: "http://localhost:6276",
				apiKey: "test",
			},
			wantErr: true,
		},
		{
			name: "TestSuccess",
			fields: fields{
				req:    retryablehttp.NewClient(),
				apiURL: "http://localhost:6276",
				apiKey: "aaCzNTmDi6J9A6OzXURHkpgQ5dDJTK4j",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				req:    tt.fields.req,
				apiURL: tt.fields.apiURL,
				apiKey: tt.fields.apiKey,
			}
			got, err := c.WorkerJoin()
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkerJoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want.Scope != got.Scope || tt.want.Type != got.Type {
				t.Errorf("WorkerJoin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_WorkerAlive(t *testing.T) {
	type fields struct {
		req    *retryablehttp.Client
		apiURL string
		apiKey string
	}
	type args struct {
		req *WorkerAliveRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *WorkerAliveResponse
		wantErr bool
	}{
		{
			name: "TestError",
			fields: fields{
				req:    retryablehttp.NewClient(),
				apiURL: "http://localhost:6276",
				apiKey: "test",
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "TestSuccess",
			fields: fields{
				req:    retryablehttp.NewClient(),
				apiURL: "http://localhost:6276",
				apiKey: "aaCzNTmDi6J9A6OzXURHkpgQ5dDJTK4j",
			},
			args: args{
				req: &WorkerAliveRequest{
					Token: "hNk47sYDCf9HfjPkviN9bRZh3fU72qRKg45kFhUxkSmPLS3k",
				},
			},
			want: &WorkerAliveResponse{
				Alive: "OK",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				req:    tt.fields.req,
				apiURL: tt.fields.apiURL,
				apiKey: tt.fields.apiKey,
			}
			got, err := c.WorkerAlive(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkerAlive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WorkerAlive() got = %v, want %v", got, tt.want)
			}
		})
	}
}
