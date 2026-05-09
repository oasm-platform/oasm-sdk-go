package oasm

import (
	"context"
	"fmt"
	"log"
	"time"
)

// WorkerConnect manages the full lifecycle of the worker's connection.
// It handles initial Join, maintains the Alive stream, and automatically reconnects on failure.
func (c *Client) WorkerConnect(ctx context.Context, ready chan<- bool) {
	logger := NewLogger("Connect")
	const (
		baseDelay = 2 * time.Second
		maxDelay  = 30 * time.Second
	)
	currentDelay := baseDelay

	for {
		logger.Verbose("Attempting to connect to Open ASM Core...")

		_, err := c.WorkerJoin(ctx)
		if err != nil {
			select {
			case ready <- false:
			default:
			}

			logger.Error("Join failed: %v. Retrying in %v...", err, currentDelay)

			if !c.waitWithContext(ctx, currentDelay, logger) {
				return
			}

			currentDelay *= 2
			if currentDelay > maxDelay {
				currentDelay = maxDelay
			}
			continue
		}

		currentDelay = baseDelay
		logger.Success("Join successful. Worker ID: %s", c.workerID)

		select {
		case ready <- true:
		default:
		}

		err = c.WorkerAlive(ctx)

		select {
		case ready <- false:
		default:
		}

		if err != nil {
			logger.Warning("Alive stream interrupted: %v. Reconnecting...", err)
		} else {
			logger.Info("Alive stream closed by server. Re-joining...")
		}

		if !c.waitWithContext(ctx, 1*time.Second, logger) {
			return
		}
	}
}

func (c *Client) waitWithContext(ctx context.Context, delay time.Duration, logger *LoggerType) bool {
	select {
	case <-ctx.Done():
		logger.Warning("WorkerConnect stopping: context cancelled.")
		return false
	case <-time.After(delay):
		return true
	}
}
