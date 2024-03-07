package weth_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ethereumcm "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/Layr-Labs/eigenda/indexer/eth"
	"github.com/Layr-Labs/eigenda/indexer/test/mock"

	"github.com/Layr-Labs/eigenda/indexer/inmem"
)

var logger = logging.NewNoopLogger()

func newTestFilterer(sc *mock.ContractSimulator, isFastMode bool) *Filterer {
	return &Filterer{
		Filterer: sc.Client,
		Address:  sc.WethAddr,
		Accounts: []ethereumcm.Address{sc.DeployerAddr},
		FastMode: isFastMode,
	}
}

func newTestAccumlatorHandlers(filterer *Filterer, acc *Accumulator, status indexer.Status) []indexer.AccumulatorHandler {
	return []indexer.AccumulatorHandler{
		{
			Acc:      acc,
			Filterer: filterer,
			Status:   status,
		},
	}
}

func TestIndex(t *testing.T) {
	t.Skip("Skipping this test after the simulated backend upgrade broke this test. Enable it after fixing the issue.")
	sc := mock.MustNewContractSimulator()

	upgrader := &Upgrader{}
	acc := &Accumulator{}

	filterer := newTestFilterer(sc, true)
	handlers := newTestAccumlatorHandlers(filterer, acc, indexer.Good)
	headerSrvc := eth.NewHeaderService(logger, sc.Client)
	headerStore := inmem.NewHeaderStore()
	config := indexer.Config{
		PullInterval: 100 * time.Millisecond,
	}
	indexer := indexer.New(
		&config,
		handlers,
		headerSrvc,
		headerStore,
		upgrader,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())

	// Start Blockchain Events
	sc.Start(time.Millisecond, cancel)

	err := indexer.Index(ctx)
	assert.NoError(t, err)

	select {
	case <-ctx.Done():
		assert.Equal(t, 4, len(headerStore.Chain), "header store chain should have 4 headers")
		assert.Equal(t, uint64(1), headerStore.Chain[0].Header.Number, "header number should have number 1")
		assert.Equal(t, uint64(2), headerStore.Chain[1].Header.Number, "header number should have number 2")
		assert.Equal(t, uint64(3), headerStore.Chain[2].Header.Number, "header number should have number 3")
		assert.Equal(t, uint64(4), headerStore.Chain[3].Header.Number, "header number should have number 4")

		ao, h, err := headerStore.GetLatestObject(acc, false)
		assert.NoError(t, err)
		assert.Equal(t, uint64(8), ao.(AccountBalanceV1).Balance, "balance for the latest Object should have value 8")

		ao, _, err = headerStore.GetObject(h, acc)
		assert.NoError(t, err)
		assert.Equal(t, uint64(8), ao.(AccountBalanceV1).Balance, "balance should have value 8")

		ao, _, err = headerStore.GetObject(headerStore.Chain[0].Header, acc)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), ao.(AccountBalanceV1).Balance, "balance at Header number 1 should have value 0")

		ao, _, err = headerStore.GetObject(headerStore.Chain[1].Header, acc)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), ao.(AccountBalanceV1).Balance, "balance at Header number 2 should have value 1")

		ao, _, err = headerStore.GetObject(headerStore.Chain[2].Header, acc)
		assert.NoError(t, err)
		assert.Equal(t, uint64(4), ao.(AccountBalanceV1).Balance, "balance at Header number 3 should have value 4")

		ao, _, err = headerStore.GetObject(headerStore.Chain[3].Header, acc)
		assert.NoError(t, err)
		assert.Equal(t, uint64(8), ao.(AccountBalanceV1).Balance, "balance at Header number 4 should have value 8")

	case <-time.After(time.Second * 40):
		t.Fatalf("expected call to Index method")
	}
}
