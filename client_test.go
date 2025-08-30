package oasm

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{

		{
			name: "TestSuccess",
			args: args{
				opts: []Option{
					WithApiURL("http://localhost:6276"),
					WithApiKey("test"),
				},
			},
			want: &Client{
				apiURL: "http://localhost:6276",
				apiKey: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.opts...); got.apiURL != tt.want.apiURL || got.apiKey != tt.want.apiKey {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
