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

// --- Category-Specific Payload Constructors ---

// NewAssetList creates an AssetList for subdomain results.
func NewAssetList(assets []*pb.Asset) *pb.AssetList {
	return &pb.AssetList{Values: assets}
}

// NewVulnerabilityList creates a VulnerabilityList for vulnerability results.
func NewVulnerabilityList(vulns []*pb.Vulnerability) *pb.VulnerabilityList {
	return &pb.VulnerabilityList{Values: vulns}
}

// NewNumberList creates a NumberList for port scan results.
func NewNumberList(ports []int32) *pb.NumberList {
	return &pb.NumberList{Values: ports}
}

// NewAssetTagList creates an AssetTagList for classifier results.
func NewAssetTagList(tags []*pb.AssetTag) *pb.AssetTagList {
	return &pb.AssetTagList{Values: tags}
}

// --- Legacy Constructors (deprecated) ---

// NewVulnerabilityResult creates a generic DataPayloadResult for vulnerability results.
//
// Deprecated: Use JobsVulnerabilitiesResult with NewVulnerabilityList instead.
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

// NewHttpResult creates a generic DataPayloadResult for HTTP probe results.
//
// Deprecated: Use JobsHttpProbeResult with the HttpResponse directly instead.
func NewHttpResult(httpResp *pb.HttpResponse) *pb.DataPayloadResult {
	return &pb.DataPayloadResult{
		Error: false,
		Payload: &pb.DataPayloadResult_HttpResponse{
			HttpResponse: httpResp,
		},
	}
}

// NewErrorResult creates a generic DataPayloadResult for error results.
//
// Deprecated: Use category-specific methods with error=true and the raw error string instead.
func NewErrorResult(errMsg string) *pb.DataPayloadResult {
	return &pb.DataPayloadResult{
		Error: true,
		Raw:   &errMsg,
	}
}
