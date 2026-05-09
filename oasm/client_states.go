package oasm

// import (
// 	"encoding/json" // Hoặc dùng yaml
// 	"os"
//
// 	pb "github.com/oasm-platform/open-asm/grpc-client/go/workers"
// )
//
// // helper to save worker state
// func (c *Client) saveWorkerState(resp *pb.JoinResponse) error {
// 	if c.configPath == "" {
// 		return nil
// 	}
//
// 	data, err := json.MarshalIndent(resp, "", "  ")
// 	if err != nil {
// 		return err
// 	}
//
// 	return os.WriteFile(c.configPath, data, 0o644)
// }
//
// // helper to read state
// func (c *Client) loadWorkerState() (*pb.JoinResponse, error) {
// 	if _, err := os.Stat(c.configPath); os.IsNotExist(err) {
// 		return nil, nil
// 	}
//
// 	data, err := os.ReadFile(c.configPath)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var state pb.JoinResponse
// 	if err := json.Unmarshal(data, &state); err != nil {
// 		return nil, err
// 	}
// 	return &state, nil
// }
