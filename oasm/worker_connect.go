package oasm

import (
	"context"
	"time"
)

func (c *Client) WorkerConnect(ctx context.Context, ready chan<- bool) {
	l := NewLogger("Worker.Connect")
	const (
		baseDelay = 2 * time.Second
		maxDelay  = 30 * time.Second
	)
	currentDelay := baseDelay

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		l.Verbose("Attempting to connect to Open ASM Core...")

		_, err := c.WorkerJoin(ctx)
		if err != nil {
			select {
			case ready <- false:
			default:
			}

			l.ErrorE("Join failed, retrying in %v", err, currentDelay)

			if !c.waitWithContext(ctx, currentDelay, l) {
				return
			}

			currentDelay *= 2
			if currentDelay > maxDelay {
				currentDelay = maxDelay
			}
			continue
		}

		currentDelay = baseDelay
		l.Success("Join successful. Worker ID: %s", c.workerID)

		select {
		case ready <- true:
		default:
		}

		err = c.WorkerAlive(ctx)

		select {
		case ready <- false:
			l.Warning("Worker connection state shifted to offline")
		default:
		}

		if err != nil {
			l.Warning("Alive stream interrupted: %v. Reconnecting...", err)
		} else {
			l.Info("Alive stream closed by server. Re-joining...")
		}

		if !c.waitWithContext(ctx, 1*time.Second, l) {
			return
		}
	}
}

func (c *Client) waitWithContext(ctx context.Context, delay time.Duration, l *LoggerType) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		l.Warning("WorkerConnect stopping: context cancelled.")
		return false
	case <-timer.C:
		return true
	}
}
