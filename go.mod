module github.com/oasm-platform/oasm-sdk-go

go 1.25.0

require (
	github.com/oasm-platform/open-asm/grpc-client/go v0.0.0-20260716080653-0fa305247d5a
	google.golang.org/grpc v1.80.0
)

replace github.com/oasm-platform/open-asm/grpc-client/go => ../open-asm/grpc-client/go

require (
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
