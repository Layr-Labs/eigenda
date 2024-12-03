package main

import (
	"github.com/Layr-Labs/eigenda/common"
	nodegrpc "github.com/Layr-Labs/eigenda/node/grpc"
	"time"
)

// main generates documentation for relay metrics.
func main() {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		panic(err)
	}

	metrics, err := nodegrpc.NewV2Metrics(logger, 0, "", time.Second)
	if err != nil {
		panic(err)
	}

	err = metrics.WriteMetricsDocumentation()
	if err != nil {
		panic(err)
	}
}
