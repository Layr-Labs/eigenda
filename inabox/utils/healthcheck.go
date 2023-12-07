package utils

import (
	"context"
	"fmt"
	"net"
	"time"
)

func WaitForServer(socketAddr string, timeout time.Duration) error {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a ticker for periodic attempts
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // Check if the context deadline is exceeded
			return fmt.Errorf("timeout after %v waiting for server", timeout)
		case <-ticker.C: // Attempt to connect at each tick
			conn, err := net.DialTimeout("tcp", socketAddr, 1*time.Second)
			if err == nil {
				conn.Close() // Close the connection if successful
				return nil
			}
		}
	}
}
