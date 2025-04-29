package healthcheck

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type HeartbeatMessage struct {
	Component string    // e.g., "encodingManager" or "dispatcher"
	Timestamp time.Time // when the heartbeat was sent
}

// HeartbeatMonitor listens for heartbeat messages from different components, updates their last seen timestamps,
// writes a summary to the specified file, and logs warnings if any component stalls.
func HeartbeatMonitor(filePath string, maxStallDuration time.Duration, livenessChan <-chan HeartbeatMessage, logger logging.Logger) error {
	// Map to keep track of last heartbeat per component
	lastHeartbeats := make(map[string]time.Time)
	// Create a timer that periodically checks for stalls
	stallTimer := time.NewTimer(maxStallDuration)

	for {
		select {
		case hb, ok := <-livenessChan:
			if !ok {
				logger.Warn("livenessChan closed, stopping health probe.")
				return nil
			}

			logger.Info("Received heartbeat", "component", hb.Component, "timestamp", hb.Timestamp)

			// Update the last heartbeat for this component
			lastHeartbeats[hb.Component] = hb.Timestamp

			// Write a summary of all components to the health file:
			summary := "Heartbeat summary:\n"
			for comp, ts := range lastHeartbeats {
				summary += fmt.Sprintf("Component %s: Last heartbeat at %v\n", comp, ts.Unix())
			}
			if err := os.WriteFile(filePath, []byte(summary), 0666); err != nil {
				logger.Error("Failed to update heartbeat file", "error", err)
			} else {
				logger.Info("Wrote heartbeat file", "path", filePath)
			}

			stallTimer.Reset(maxStallDuration)

		case <-stallTimer.C:
			// Check for components that haven't sent a heartbeat recently
			now := time.Now()
			var staleComponents []string
			for comp, ts := range lastHeartbeats {
				if now.Sub(ts) > maxStallDuration {
					staleComponents = append(staleComponents, fmt.Sprintf("Component %s: last heartbeat at %v", comp, ts))
				}
			}
			if len(staleComponents) > 0 {
				logger.Warn(
					"Components stalled",
					"components", strings.Join(staleComponents, ","),
					"threshold", maxStallDuration,
				)
			} else {
				logger.Warn("No heartbeat received recently, but no components are stale")
			}
			stallTimer.Reset(maxStallDuration)
		}
	}
}

// SignalHeartbeat sends a non-blocking heartbeat message (with component identifier and timestamp) to the given send-only channel.
func SignalHeartbeat(component string, livenessChan chan<- HeartbeatMessage, logger logging.Logger) {
	hb := HeartbeatMessage{
		Component: component,
		Timestamp: time.Now(),
	}
	select {
	case livenessChan <- hb:
		logger.Info("Heartbeat signal sent", "component", component)
	default:
		logger.Warn("Heartbeat signal skipped, no receiver on the channel", "component", component)
	}
}
