# OASM SDK for Go

![Go](https://img.shields.io/badge/Go-1.26-blue?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![gRPC](https://img.shields.io/badge/gRPC-Integration-orange)

**`oasm-sdk-go`** is the official Go client for interacting with the **OASM (Open Attack Surface Management) Platform** API. It provides convenient wrappers around gRPC-based worker management and job registry endpoints, enabling seamless integration of Go applications with the OASM core infrastructure.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)
- [Configuration](#configuration)
- [License](#license)

## Overview

The **OASM SDK for Go** provides a type-safe, idiomatic Go interface for communicating with the OASM Platform's gRPC backend. It abstracts the complexities of direct gRPC communication, connection management, authentication, and worker lifecycle operations into a clean, easy-to-use client library.

### Value Proposition

- **Simplify Integration**: Eliminate the need to write boilerplate gRPC client code for OASM Platform interactions.
- **Worker Lifecycle Management**: Join workers, maintain heartbeats, download tools, and execute jobs with minimal effort.
- **Resilient Communication**: Built-in exponential backoff reconnection logic ensures robust connectivity in unstable network environments.
- **Job Processing**: Fetch jobs from a central registry, process them, and submit results back seamlessly.

## Installation

### Prerequisites

- **Go** 1.25 or later ([Download](https://go.dev/dl/))
- A valid **OASM Platform API Key**
- Network access to the OASM Platform gRPC endpoint

### Step 1: Install the SDK

```bash
go get -u github.com/oasm-platform/oasm-sdk-go
```

### Step 2: Import in Your Project

```go
import "github.com/oasm-platform/oasm-sdk-go/oasm"
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/oasm-platform/oasm-sdk-go/oasm"
)

func main() {
    // Initialize the client
    client, err := oasm.NewClient(
        oasm.WithGRPCHost("api.oasm.dev:16276"),
        oasm.WithApiKey("your-api-key-here"),
    )
    if err != nil {
        log.Fatalf("failed to create client: %v", err)
    }
    defer client.Close()

    // Connect a worker with automatic reconnection
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    ready := make(chan bool)
    go client.WorkerConnect(ctx, ready)

    isReady := <-ready
    if !isReady {
        log.Fatal("worker failed to connect")
    }

    fmt.Println("Worker connected successfully!")

    // Fetch and process jobs
    for {
        job, err := client.JobsNext(ctx)
        if err != nil {
            log.Printf("error fetching job: %v", err)
            continue
        }
        if job == nil {
            continue // No job available, keep polling
        }

        fmt.Printf("Processing job: %s\n", job.Id)

        // Submit a vulnerability result
        vulnResult := oasm.NewVulnerabilityResult(nil)
        if err := client.JobsResult(ctx, job.Id, vulnResult); err != nil {
            log.Printf("error submitting result for job %s: %v", job.Id, err)
        }
    }
}
```

## Usage Examples

### Client Configuration

#### Basic Configuration

```go
client, err := oasm.NewClient(
    oasm.WithGRPCHost("api.oasm.dev:16276"),
    oasm.WithApiKey("your-api-key"),
)
```

### Worker Lifecycle

#### Join and Send Heartbeat

```go
// Join the platform as a worker
joinResp, err := client.WorkerJoin(ctx)
if err != nil {
    log.Fatalf("worker join failed: %v", err)
}
fmt.Printf("Worker ID: %s\n", client.WorkerID())

// Start a persistent alive stream (heartbeat)
err = client.WorkerAlive(ctx)
if err != nil {
    log.Fatalf("heartbeat failed: %v", err)
}
```

#### Automatic Reconnection

`WorkerConnect` handles automatic reconnection with exponential backoff:

```go
ctx, cancel := context.WithCancel(context.Background())
ready := make(chan bool)

go client.WorkerConnect(ctx, ready)

// Wait until the worker is connected
if isReady := <-ready; !isReady {
    log.Fatal("initial connection failed")
}
```

### Job Registry Operations

#### Fetch Next Job

```go
job, err := client.JobsNext(ctx)
if err != nil {
    log.Fatalf("failed to fetch job: %v", err)
}
if job == nil {
    fmt.Println("No jobs available")
    return
}

// Access job details
fmt.Printf("Job ID: %s\n", job.Id)
fmt.Printf("Asset Type: %s\n", job.Asset.GetTarget().GetType())
```

#### Submit Vulnerability Results

```go
vulns := []*pb.Vulnerability{
    {
        // ... populate vulnerability data
    },
}
result := oasm.NewVulnerabilityResult(vulns)

err = client.JobsResult(ctx, job.Id, result)
if err != nil {
    log.Fatalf("failed to submit results: %v", err)
}
```

#### Submit HTTP Response Results

```go
httpResp := &pb.HttpResponse{
    // ... populate HTTP response data
}
result := oasm.NewHttpResult(httpResp)

err = client.JobsResult(ctx, job.Id, result)
if err != nil {
    log.Fatalf("failed to submit results: %v", err)
}
```

#### Submit Error Results

```go
result := oasm.NewErrorResult("scan timed out")

err = client.JobsResult(ctx, job.Id, result)
if err != nil {
    log.Fatalf("failed to submit error: %v", err)
}
```

### Helper Functions

```go
// Extract DNS records from a job asset
dnsRecords := client.GetDNSRecordsMap(job)

// Get asset creation timestamp
createdAt := client.GetCreatedAtTime(job)
fmt.Printf("Asset created at: %v\n", createdAt)
```

## API Reference

### `NewClient(opts ...Option) (*Client, error)`

Creates a new `Client` instance with the provided configuration options. Default gRPC host is `localhost:16276` if not specified.

### `Client` Methods

| Method                      | Description                                                  |
| --------------------------- | ------------------------------------------------------------ |
| `WorkerJoin(ctx)`           | Sends a join request to the platform; returns `JoinResponse` |
| `WorkerAlive(ctx)`          | Establishes a persistent heartbeat stream                    |
| `WorkerConnect(ctx, ready)` | Manages reconnection with exponential backoff                |
| `WorkerDownloadTools(ctx)`  | Downloads and extracts platform tools                        |
| `JobsNext(ctx)`             | Fetches the next available job from the registry             |
| `JobsResult(ctx, id, data)` | Submits job results back to the registry                     |
| `GetDNSRecordsMap(job)`     | Extracts DNS records from a job's asset                      |
| `GetCreatedAtTime(job)`     | Parses the asset creation timestamp                          |
| `Workers()`                 | Returns the raw gRPC `WorkersServiceClient`                  |
| `Jobs()`                    | Returns the raw gRPC `JobsRegistryServiceClient`             |
| `WorkerID()`                | Returns the current worker's unique identifier               |
| `Token()`                   | Returns the current worker authentication token              |
| `Close()`                   | Closes the underlying gRPC connection                        |

### `WithAuth(ctx) context.Context`

Attaches the worker's authentication token to the gRPC metadata of the given context. Used internally for authenticated job registry calls.

## Configuration Options

| Option               | Type               | Default           | Description                       |
| -------------------- | ------------------ | ----------------- | --------------------------------- |
| `WithGRPCHost(host)` | `string`           | `localhost:16276` | gRPC server address               |
| `WithApiKey(key)`    | `string`           | —                 | API key for worker authentication |
| `WithConn(conn)`     | `*grpc.ClientConn` | Auto-created      | Custom gRPC connection            |

## License

This project is licensed under the **[MIT License](LICENSE)**.

```
MIT License

Copyright (c) 2025 Open Attack Surface Management

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
