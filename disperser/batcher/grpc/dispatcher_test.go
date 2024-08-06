package dispatcher_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	grpcMock "github.com/Layr-Labs/eigenda/api/grpc/mock"
	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	dispatcher "github.com/Layr-Labs/eigenda/disperser/batcher/grpc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/stretchr/testify/assert"
)

func newDispatcher(t *testing.T, config *dispatcher.Config) disperser.Dispatcher {
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	assert.NoError(t, err)
	metrics := batcher.NewMetrics("9091", logger)
	return dispatcher.NewDispatcher(config, logger, metrics.DispatcherMetrics)
}

func TestSendAttestBatchRequest(t *testing.T) {
	dispatcher := newDispatcher(t, &dispatcher.Config{
		Timeout: 5 * time.Second,
	})
	nodeClient := grpcMock.NewMockDispersalClient()
	var X, Y fp.Element
	X = *X.SetBigInt(big.NewInt(1))
	Y = *Y.SetBigInt(big.NewInt(2))
	signature := &core.Signature{
		G1Point: &core.G1Point{
			G1Affine: &bn254.G1Affine{
				X: X,
				Y: Y,
			},
		},
	}
	sigBytes := signature.Bytes()
	nodeClient.On("AttestBatch").Return(&node.AttestBatchReply{
		Signature: sigBytes[:],
	}, nil)
	sigReply, err := dispatcher.SendAttestBatchRequest(context.Background(), nodeClient, [][32]byte{{1}}, &core.BatchHeader{
		ReferenceBlockNumber: 10,
		BatchRoot:            [32]byte{1},
	}, &core.IndexedOperatorInfo{
		PubkeyG1: nil,
		PubkeyG2: nil,
		Socket:   "localhost:8080",
	})
	assert.NoError(t, err)
	assert.Equal(t, signature, sigReply)
}
