package traffic_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/tools/traffic"
	traffic_mock "github.com/Layr-Labs/eigenda/tools/traffic/mock"

	"github.com/stretchr/testify/mock"
)

func TestTrafficGenerator(t *testing.T) {
	disperserClient := traffic_mock.NewMockDisperserClient()
	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	if err != nil {
		panic("failed to create new logger")
	}
	trafficGenerator := &traffic.TrafficGenerator{
		Logger: logger,
		Config: &traffic.Config{
			DataSize:        1000_000,
			RequestInterval: 2 * time.Second,
			Timeout:         1 * time.Second,
		},
		DisperserClient: disperserClient,
	}

	processing := disperser.Processing
	disperserClient.On("DisperseBlob", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&processing, []byte{1}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = trafficGenerator.StartTraffic(ctx)
	}()
	time.Sleep(5 * time.Second)
	cancel()
	disperserClient.AssertNumberOfCalls(t, "DisperseBlob", 2)
}
