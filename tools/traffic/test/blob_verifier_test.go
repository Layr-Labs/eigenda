package test

import (
	"context"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
	"time"
)

func TestBlobVerifier(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)
	startTime := time.Unix(rand.Int63()%2_000_000_000, 0)
	ticker := newMockTicker(startTime)

	config := &workers.Config{}

	blobTable := table.NewBlobTable()

	verifierMetrics := metrics.NewMockMetrics()

	verifier := workers.NewBlobVerifier(
		&ctx,
		&waitGroup,
		logger,
		ticker,
		config,
		&blobTable,
		nil,
		verifierMetrics)

	verifier.Start()

	cancel()
	tu.ExecuteWithTimeout(func() {
		waitGroup.Wait()
	}, time.Second)
}
