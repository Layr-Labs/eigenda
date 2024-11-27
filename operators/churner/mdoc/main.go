package main

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/operators/churner"
)

// main generates documentation for churner metrics.
func main() {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		panic(err)
	}

	metrics, err := churner.NewMetrics(0, logger)
	if err != nil {
		panic(err)
	}

	err = metrics.WriteMetricsDocumentation()
	if err != nil {
		panic(err)
	}
}
