package common

import (
	"fmt"
	"sync"
	"time"
)

const (
	// For an original 64-bit nanosecond timestamp, we use its lower bits for instance ID
	// and the higher bits for timestamp.
	//
	// When instance ID is 10 bits:
	// - The time resolution is 2^10 = 1,024 nanoseconds (1.024 microseconds)
	// - This means an instance can handle at most ~976,562 requests per second
	//   (1 second / 1.024 microseconds)
	// - The max number of instances will be 2^10 = 1,024
	DefaultInstanceIDBits = 10 // Default number of bits for instance ID

	// When instance ID is 16 bits (the max allowed):
	//   - The time resolution is 2^16 = 65,536 nanoseconds (65.536 microseconds)
	//   - This means an instance can handle at most ~15,259 requests per second
	//     (1 second / 65.536 microseconds)
	//   - The max number of instances will be 2^16 = 65,536
	MaxInstanceIDBits = 16 // Maximum allowed bits for instance ID
)

// TimestampOracle generates unique timestamps across multiple service replicas
type TimestampOracle struct {
	mu             sync.Mutex
	lastTimestamp  uint64
	instanceID     uint64
	instanceIDBits uint
	instanceIDMask uint64
}

func NewTimestampOracle(instanceID uint64, instanceIDBits uint) (*TimestampOracle, error) {
	if instanceIDBits > MaxInstanceIDBits {
		return nil, fmt.Errorf("instance ID bits must be between 1 and %d, got %d", MaxInstanceIDBits, instanceIDBits)
	}

	// Calculate the maximum allowed value for instance ID with the given bits
	maxInstanceID := uint64(1<<instanceIDBits) - 1

	if instanceID > maxInstanceID {
		return nil, fmt.Errorf("instance ID must be between 0 and %d for %d bits, got %d",
			maxInstanceID, instanceIDBits, instanceID)
	}

	return &TimestampOracle{
		lastTimestamp:  0,
		instanceID:     instanceID,
		instanceIDBits: instanceIDBits,
		instanceIDMask: maxInstanceID,
	}, nil
}

// GetUniqueTimestamp returns a unique timestamp
// - Uses top (64-instanceIDBits) bits for timestamp (nanoseconds since epoch shifted)
// - Uses bottom instanceIDBits bits for inID
func (g *TimestampOracle) GetUniqueTimestamp() uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Get current time in nanoseconds
	current := uint64(time.Now().UnixNano())

	// Extract the timestamp part (top bits)
	timestampPart := current >> g.instanceIDBits << g.instanceIDBits

	// If this timestamp part is less than or equal to the last timestamp,
	// increment the timestamp part to ensure uniqueness
	if timestampPart <= g.lastTimestamp {
		// Increment by 1 in the timestamp part
		timestampPart = g.lastTimestamp + (1 << g.instanceIDBits)
	}

	g.lastTimestamp = timestampPart

	return timestampPart | g.instanceID
}
