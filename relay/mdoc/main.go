package main

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/relay/metrics"
)

// main generates documentation for relay metrics.
func main() {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		panic(err)
	}

	metrics, err := metrics.NewRelayMetrics(logger, 0)
	if err != nil {
		panic(err)
	}

	err = metrics.WriteMetricsDocumentation()
	if err != nil {
		panic(err)
	}
}
