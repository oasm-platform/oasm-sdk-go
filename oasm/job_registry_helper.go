package oasm

import (
	"time"

	pb "github.com/oasm-platform/open-asm/grpc-client/go/jobs_registry"
)

func (c *Client) GetDNSRecordsMap(job *pb.Job) map[string]interface{} {
	if job.Asset == nil || job.Asset.DnsRecords == nil {
		return nil
	}
	return job.Asset.DnsRecords.AsMap()
}

func (c *Client) GetCreatedAtTime(job *pb.Job) time.Time {
	if job.Asset == nil || job.Asset.CreatedAt == nil {
		return time.Time{}
	}
	return job.Asset.CreatedAt.AsTime()
}

func NewVulnerabilityResult(vulns []*pb.Vulnerability) *pb.DataPayloadResult {
	return &pb.DataPayloadResult{
		Error: false,
		Payload: &pb.DataPayloadResult_Vulnerabilities{
			Vulnerabilities: &pb.VulnerabilityList{
				Values: vulns,
			},
		},
	}
}

func NewHttpResult(httpResp *pb.HttpResponse) *pb.DataPayloadResult {
	return &pb.DataPayloadResult{
		Error: false,
		Payload: &pb.DataPayloadResult_HttpResponse{
			HttpResponse: httpResp,
		},
	}
}

func NewErrorResult(errMsg string) *pb.DataPayloadResult {
	return &pb.DataPayloadResult{
		Error: true,
		Raw:   &errMsg,
	}
}
