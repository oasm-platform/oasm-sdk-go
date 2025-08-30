# OASM SDK for Go

`oasm-sdk-go` is the official Go client for interacting with the **OASM Platform** API.  
It provides convenient wrappers for worker management endpoints such as **join** and **keep-alive**.

---

## Installation

Use `go get` to install:

```bash
go get -u github.com/oasm-platform/oasm-sdk-go
````

Then import it in your project:

```go
import oasm "github.com/oasm-platform/oasm-sdk-go"
```

---

## Usage

### Initialize Client

```go
package main

import (
	"fmt"
	"log"

	oasm "github.com/oasm-platform/oasm-sdk-go"
)

func main() {
	// Create a new client with API URL and API key
	client := oasm.NewClient(
		oasm.WithApiURL("https://api.oasm.dev"),
		oasm.WithApiKey("your-api-key"),
	)

	// Join worker
	joinResp, err := client.WorkerJoin()
	if err != nil {
		log.Fatalf("failed to join worker: %v", err)
	}
	fmt.Println("Worker joined:", joinResp.Id)

	// Send keep-alive
	aliveResp, err := client.WorkerAlive(&oasm.WorkerAliveRequest{
		Token: joinResp.Token,
	})
	if err != nil {
		log.Fatalf("failed to send keep-alive: %v", err)
	}
	fmt.Println("Worker alive:", aliveResp.Alive)
}
```

---

## License

[MIT](./LICENSE)
