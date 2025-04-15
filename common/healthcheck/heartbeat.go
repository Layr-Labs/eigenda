package healthcheck

import (
	"fmt"
	"log"
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
func HeartbeatMonitor(filePath string, maxStallDuration time.Duration, livenessChan <-chan HeartbeatMessage) {
	// Map to keep track of last heartbeat per component
	lastHeartbeats := make(map[string]time.Time)
	// Create a timer that periodically checks for stalls
	stallTimer := time.NewTimer(maxStallDuration)

	for {
		select {
		case hb, ok := <-livenessChan:
			if !ok {
				log.Println("livenessChan closed, stopping health probe.")
				return
			}
			log.Printf("Received heartbeat from component '%s' at %v", hb.Component, hb.Timestamp)
			// Update the last heartbeat for this component
			lastHeartbeats[hb.Component] = hb.Timestamp

			// Write a summary of all components to the health file:
			summary := "Heartbeat summary:\n"
			for comp, ts := range lastHeartbeats {
				summary += fmt.Sprintf("Component %s: Last heartbeat at %v\n", comp, ts)
			}
			if err := os.WriteFile(filePath, []byte(summary), 0666); err != nil {
				log.Printf("Failed to update heartbeat file: %v", err)
			} else {
				log.Printf("Updated heartbeat file: %v with summary:\n%s", filePath, summary)
			}

			stallTimer.Reset(maxStallDuration)

		case <-stallTimer.C:
			// Check for components that haven't sent a heartbeat recently
			now := time.Now()
			var staleComponents []string
			for comp, ts := range lastHeartbeats {
				if now.Sub(ts) > maxStallDuration {
					staleComponents = append(staleComponents, fmt.Sprintf("Component %s last heartbeat at %v", comp, ts))
				}
			}
			if len(staleComponents) > 0 {
				log.Printf("Warning: No heartbeat received within max stall duration for: %s", strings.Join(staleComponents, ", "))
			} else {
				log.Println("Warning: No heartbeat received within max stall duration, but no stale components identified")
			}
			stallTimer.Reset(maxStallDuration)
		}
	}
}

// SignalHeartbeat sends a non-blocking heartbeat message (with component identifier and timestamp) on the given channel.
func SignalHeartbeat(component string, livenessChan chan HeartbeatMessage, logger logging.Logger) {
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
