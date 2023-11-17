package eth_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestConfirmerRetry(t *testing.T) {
	tx := coremock.MockTransactor{}
	confirmer, err := eth.NewBatchConfirmer(&tx, 10*time.Second)
	assert.Nil(t, err)
	tx.On("ConfirmBatch").Return(nil, fmt.Errorf("no good")).Twice()
	tx.On("ConfirmBatch").Return(&types.Receipt{
		TxHash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	}, nil).Once()
	_, err = confirmer.ConfirmBatch(context.Background(), &core.BatchHeader{
		ReferenceBlockNumber: 100,
		BatchRoot:            [32]byte{},
	}, map[core.QuorumID]*core.QuorumResult{}, &core.SignatureAggregation{
		NonSigners:       []*core.G1Point{},
		QuorumAggPubKeys: []*core.G1Point{},
		AggPubKey:        nil,
		AggSignature:     nil,
	})
	assert.Nil(t, err)
	tx.AssertNumberOfCalls(t, "ConfirmBatch", 3)
}

func TestConfirmerTimeout(t *testing.T) {
	tx := coremock.MockTransactor{}
	confirmer, err := eth.NewBatchConfirmer(&tx, 100*time.Millisecond)
	assert.Nil(t, err)
	tx.On("ConfirmBatch").Return(nil, fmt.Errorf("EnsureTransactionEvaled: failed to wait for transaction (%s) to mine: %w", "123", context.DeadlineExceeded)).Once()
	_, err = confirmer.ConfirmBatch(context.Background(), &core.BatchHeader{
		ReferenceBlockNumber: 100,
		BatchRoot:            [32]byte{},
	}, map[core.QuorumID]*core.QuorumResult{}, &core.SignatureAggregation{
		NonSigners:       []*core.G1Point{},
		QuorumAggPubKeys: []*core.G1Point{},
		AggPubKey:        nil,
		AggSignature:     nil,
	})
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	tx.AssertNumberOfCalls(t, "ConfirmBatch", 1)
}
