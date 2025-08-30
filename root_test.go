package oasm

import (
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestClient_Health(t *testing.T) {
	type fields struct {
		req    *retryablehttp.Client
		apiURL string
		apiKey string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "TestError",
			fields: fields{
				req:    retryablehttp.NewClient(),
				apiURL: "http://localhost:6277",
			},

			wantErr: true,
		},
		{
			name: "TestSuccess",
			fields: fields{
				req:    retryablehttp.NewClient(),
				apiURL: "http://localhost:6276",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				req:    tt.fields.req,
				apiURL: tt.fields.apiURL,
				apiKey: tt.fields.apiKey,
			}
			got, err := c.Health()
			if (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Health() got = %v, want %v", got, tt.want)
			}
		})
	}
}
