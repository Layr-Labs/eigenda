package verification

import (
	"context"
	"sync"
	"testing"
	"time"

	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/test"
	testrandom "github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

var (
	logger = test.GetLogger()
)

func TestWaitForBlockNumber(t *testing.T) {
	ctx := t.Context()

	mockEthClient := &commonmock.MockEthClient{}
	pollRate := time.Millisecond * 50

	blockNumberMonitor, err := NewBlockNumberMonitor(logger, mockEthClient, pollRate)
	require.NoError(t, err)

	// number of goroutines to start, each of which will call WaitForBlockNumber
	callCount := 5

	for i := uint64(0); i < uint64(callCount); i++ {
		// BlockNumber will increment its return value each time it's called, up to callCount-1
		mockEthClient.On("BlockNumber").Return(i).Once()
	}
	// then, all subsequent calls will yield callCount -1
	mockEthClient.On("BlockNumber").Return(uint64(callCount - 1))

	// give plenty of time on the timeout, to get the necessary number of polls in
	timeoutCtx, cancel := context.WithTimeout(ctx, pollRate*time.Duration(callCount*2))
	defer cancel()

	waitGroup := sync.WaitGroup{}

	// start these goroutines in random order, so that it isn't always the same sequence of polling handoffs that gets exercised
	indices := testrandom.NewTestRandom().Perm(callCount)
	for _, index := range indices {
		waitGroup.Add(1)

		go func(i int) {
			defer waitGroup.Done()

			if i == callCount-1 {
				// the last call is set up to fail, by setting the target block to a number that will never be attained
				err := blockNumberMonitor.WaitForBlockNumber(timeoutCtx, uint64(i)+1)
				require.Error(t, err)
			} else {
				// all calls except the final call wait for a block number corresponding to their index
				err := blockNumberMonitor.WaitForBlockNumber(timeoutCtx, uint64(i))
				require.NoError(t, err)
			}
		}(index)
	}

	waitGroup.Wait()
	mockEthClient.AssertExpectations(t)
}
