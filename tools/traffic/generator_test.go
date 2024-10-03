package traffic_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/tools/traffic"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/stretchr/testify/mock"
)

func TestTrafficGenerator(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	logger := logging.NewNoopLogger()
	trafficGenerator := &traffic.TrafficGenerator{
		Logger: logger,
		Config: &traffic.Config{
			Config: clients.Config{
				Timeout: 1 * time.Second,
			},
			DataSize:        1000_000,
			RequestInterval: 2 * time.Second,
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

func TestTrafficGeneratorAuthenticated(t *testing.T) {
	disperserClient := clientsmock.NewMockDisperserClient()
	logger := logging.NewNoopLogger()

	trafficGenerator := &traffic.TrafficGenerator{
		Logger: logger,
		Config: &traffic.Config{
			Config: clients.Config{
				Timeout: 1 * time.Second,
			},
			DataSize:         1000_000,
			RequestInterval:  2 * time.Second,
			SignerPrivateKey: "Hi",
		},
		DisperserClient: disperserClient,
	}

	processing := disperser.Processing
	disperserClient.On("DisperseBlobAuthenticated", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&processing, []byte{1}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = trafficGenerator.StartTraffic(ctx)
	}()
	time.Sleep(5 * time.Second)
	cancel()
	disperserClient.AssertNumberOfCalls(t, "DisperseBlobAuthenticated", 2)
}
