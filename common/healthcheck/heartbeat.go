package healthcheck

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// HeartbeatMonitorConfig configures the heartbeat monitoring system that tracks component health.
type HeartbeatMonitorConfig struct {
	// FilePath is the path to the file where heartbeat status will be written. Required.
	FilePath string
	// MaxStallDuration is the maximum time allowed between heartbeats before a component is considered stalled. Required.
	MaxStallDuration time.Duration
}

var _ config.VerifiableConfig = &HeartbeatMonitorConfig{}

// GetDefaultHeartbeatMonitorConfig returns a HeartbeatMonitorConfig with sensible default values.
func GetDefaultHeartbeatMonitorConfig() *HeartbeatMonitorConfig {
	return &HeartbeatMonitorConfig{
		FilePath:         "/tmp/controller-health",
		MaxStallDuration: 4 * time.Minute,
	}
}

// Validate checks that the configuration is valid, returning an error if it is not.
func (c *HeartbeatMonitorConfig) Verify() error {
	if c.FilePath == "" {
		return fmt.Errorf("FilePath is required")
	}
	if c.MaxStallDuration <= 0 {
		return fmt.Errorf("MaxStallDuration must be positive, got %v", c.MaxStallDuration)
	}
	return nil
}

type HeartbeatMessage struct {
	Component string    // e.g., "encodingManager" or "dispatcher"
	Timestamp time.Time // when the heartbeat was sent
}

// HeartbeatMonitor listens for heartbeat messages from different components, updates their last seen timestamps,
// writes a summary to the specified file, and logs warnings if any component stalls.
func NewHeartbeatMonitor(
	logger logging.Logger,
	livenessChan <-chan HeartbeatMessage,
	config HeartbeatMonitorConfig,
) error {
	if err := config.Verify(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Map to keep track of last heartbeat per component
	lastHeartbeats := make(map[string]time.Time)
	// Create a timer that periodically checks for stalls
	stallTicker := time.NewTicker(config.MaxStallDuration)
	defer stallTicker.Stop()

	for {
		select {
		case hb, ok := <-livenessChan:
			if !ok {
				logger.Warn("livenessChan closed, stopping health probe.")
				return nil
			}

			// Update the last heartbeat for this component
			lastHeartbeats[hb.Component] = hb.Timestamp

			// Write a summary of all components to the health file:
			summary := "Heartbeat summary:\n"
			for comp, ts := range lastHeartbeats {
				summary += fmt.Sprintf("Component %s: Last heartbeat at %v\n", comp, ts.Unix())
			}
			if err := os.WriteFile(config.FilePath, []byte(summary), 0666); err != nil {
				logger.Error("Failed to update heartbeat file", "error", err)
			}

			stallTicker.Reset(config.MaxStallDuration)

		case <-stallTicker.C:
			// Check for components that haven't sent a heartbeat recently
			now := time.Now()
			var staleComponents []string
			for comp, ts := range lastHeartbeats {
				if now.Sub(ts) > config.MaxStallDuration {
					staleComponents = append(staleComponents, fmt.Sprintf("Component %s: last heartbeat at %v", comp, ts))
				}
			}
			if len(staleComponents) > 0 {
				logger.Warn(
					"Components stalled",
					"components", strings.Join(staleComponents, ","),
					"threshold", config.MaxStallDuration,
				)
			} else {
				logger.Warn("No heartbeat received recently, but no components are stale")
			}
		}
	}
}

// SignalHeartbeat sends a non-blocking heartbeat message (with component identifier and timestamp) to the given send-only channel.
func SignalHeartbeat(logger logging.Logger, component string, livenessChan chan<- HeartbeatMessage) {
	hb := HeartbeatMessage{
		Component: component,
		Timestamp: time.Now(),
	}
	select {
	case livenessChan <- hb:
	default:
		logger.Warn("Heartbeat signal skipped, no receiver on the channel", "component", component)
	}
}
