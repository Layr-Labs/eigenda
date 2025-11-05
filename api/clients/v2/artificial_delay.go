package clients

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

var latencyMap = map[uint64]time.Duration{
	0: 0 * time.Millisecond,   // simulated latency for something in the same region
	1: 20 * time.Millisecond,  // simulated latency for Canada
	2: 100 * time.Millisecond, // simulated latency for Germany
	3: 150 * time.Millisecond, // simulated latency for Brazil
	4: 170 * time.Millisecond, // simulated latency for Nigeria
	5: 250 * time.Millisecond, // simulated latency for China
}

// SimulateLatency introduces an artificial delay based on the validator ID's simulated region. Each validator
// ID is mapped to a region, with each region having equal probability.
func SimulateLatency(
	ctx context.Context,
	validatorId core.OperatorID,
) {

	region := binary.BigEndian.Uint64(validatorId[:]) % uint64(len(latencyMap))
	delay := latencyMap[region]
	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
		}
	}
}
