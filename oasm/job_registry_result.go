package oasm

import (
	"context"
	"fmt"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/jobs_registry"
)

// JobsResult submits a job result using the legacy generic endpoint.
//
// Deprecated: Use category-specific methods instead (e.g., JobsSubdomainsResult,
// JobsHttpProbeResult, JobPortsResult, etc.). The category-specific endpoints
// provide type-safe payloads and are the recommended path forward.
// This method is kept for backward compatibility during worker migration.
func (c *Client) JobsResult(ctx context.Context, jobID string, payload *pb.DataPayloadResult) error {
	req := &pb.JobResultRequest{
		WorkerId: c.workerID,
		Data: &pb.UpdateResultDto{
			JobId: jobID,
			Data:  payload,
		},
	}

	resp, err := c.Jobs().Result(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit job result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the result submission")
	}

	return nil
}

// JobsSubdomainsResult submits subdomain discovery results for a job.
func (c *Client) JobsSubdomainsResult(ctx context.Context, jobID string, error bool, raw string, assets []*pb.Asset) error {
	req := &pb.SubdomainResultRequest{
		WorkerId: c.workerID,
		JobId:    jobID,
		Error:    error,
		Assets:   &pb.AssetList{Values: assets},
	}
	if raw != "" {
		req.Raw = &raw
	}

	resp, err := c.Jobs().ResultSubdomains(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit subdomains result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the subdomains result submission")
	}

	return nil
}

// JobsHttpProbeResult submits HTTP probe scan results for a job.
func (c *Client) JobsHttpProbeResult(ctx context.Context, jobID string, error bool, raw string, httpResponse *pb.HttpResponse) error {
	req := &pb.HttpProbeResultRequest{
		WorkerId:     c.workerID,
		JobId:        jobID,
		Error:        error,
		HttpResponse: httpResponse,
	}
	if raw != "" {
		req.Raw = &raw
	}

	resp, err := c.Jobs().ResultHttpProbe(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit http-probe result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the http-probe result submission")
	}

	return nil
}

// JobsPortsResult submits port scanner results for a job.
func (c *Client) JobsPortsResult(ctx context.Context, jobID string, error bool, raw string, ports []int32) error {
	req := &pb.PortsResultRequest{
		WorkerId: c.workerID,
		JobId:    jobID,
		Error:    error,
		Numbers:  &pb.NumberList{Values: ports},
	}
	if raw != "" {
		req.Raw = &raw
	}

	resp, err := c.Jobs().ResultPorts(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit ports result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the ports result submission")
	}

	return nil
}

// JobsVulnerabilitiesResult submits vulnerability scan results for a job.
func (c *Client) JobsVulnerabilitiesResult(ctx context.Context, jobID string, error bool, raw string, vulns []*pb.Vulnerability) error {
	req := &pb.VulnerabilitiesResultRequest{
		WorkerId:       c.workerID,
		JobId:          jobID,
		Error:          error,
		Vulnerabilities: &pb.VulnerabilityList{Values: vulns},
	}
	if raw != "" {
		req.Raw = &raw
	}

	resp, err := c.Jobs().ResultVulnerabilities(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit vulnerabilities result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the vulnerabilities result submission")
	}

	return nil
}

// JobsScreenshotResult submits screenshot capture results for a job.
func (c *Client) JobsScreenshotResult(ctx context.Context, jobID string, error bool, raw string) error {
	req := &pb.ScreenshotResultRequest{
		WorkerId: c.workerID,
		JobId:    jobID,
		Error:    error,
	}
	if raw != "" {
		req.Raw = &raw
	}

	resp, err := c.Jobs().ResultScreenshot(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit screenshot result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the screenshot result submission")
	}

	return nil
}

// JobsClassifierResult submits asset classification results for a job.
func (c *Client) JobsClassifierResult(ctx context.Context, jobID string, error bool, raw string, assetTags []*pb.AssetTag) error {
	req := &pb.ClassifierResultRequest{
		WorkerId:  c.workerID,
		JobId:     jobID,
		Error:     error,
		AssetTags: &pb.AssetTagList{Values: assetTags},
	}
	if raw != "" {
		req.Raw = &raw
	}

	resp, err := c.Jobs().ResultClassifier(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit classifier result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the classifier result submission")
	}

	return nil
}

// JobsAssistantResult submits AI assistant results for a job.
func (c *Client) JobsAssistantResult(ctx context.Context, jobID string, error bool, raw string) error {
	req := &pb.AssistantResultRequest{
		WorkerId: c.workerID,
		JobId:    jobID,
		Error:    error,
	}
	if raw != "" {
		req.Raw = &raw
	}

	resp, err := c.Jobs().ResultAssistant(c.WithAuth(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to submit assistant result: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server rejected the assistant result submission")
	}

	return nil
}
