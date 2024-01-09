package eth

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	maxRetries = 3
	baseDelay  = 1 * time.Second
)

type BatchConfirmer struct {
	Transactor core.Transactor
	timeout    time.Duration
}

// NewBatchConfirmer returns a new BatchConfirmer
func NewBatchConfirmer(tx core.Transactor, timeout time.Duration) (disperser.BatchConfirmer, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("failed to create new Confirmer because timeout is not greater than 0")
	}
	return &BatchConfirmer{
		Transactor: tx,
		timeout:    timeout,
	}, nil
}

var _ disperser.BatchConfirmer = (*BatchConfirmer)(nil)

func (c *BatchConfirmer) ConfirmBatch(ctx context.Context, header *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, sigAgg *core.SignatureAggregation) (*types.Receipt, error) {
	var (
		txReceipt *types.Receipt
		err       error
	)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	for i := 0; i < maxRetries; i++ {
		txReceipt, err = c.Transactor.ConfirmBatch(ctxWithTimeout, header, quorums, sigAgg)
		if err == nil {
			break
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		if strings.Contains(err.Error(), "execution reverted") {
			return nil, err
		}

		retrySec := math.Pow(2, float64(i))
		time.Sleep(time.Duration(retrySec) * baseDelay)
	}

	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}
