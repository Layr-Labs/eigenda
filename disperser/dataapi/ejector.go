package dataapi

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	walletsdk "github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	maxSendTransactionRetry = 3
	queryTickerDuration     = 3 * time.Second
)

// stakeShareToSLA returns the SLA for a given stake share in a quorum.
// The caller should ensure "stakeShare" is in range (0, 1].
func stakeShareToSLA(stakeShare float64) float64 {
	switch {
	case stakeShare > 0.15:
		return 0.995
	case stakeShare > 0.1:
		return 0.98
	case stakeShare > 0.05:
		return 0.95
	default:
		return 0.9
	}
}

// operatorPerfScore scores an operator based on its stake share and nonsigning rate. The
// performance score will be in range [0, 1], with higher score indicating better performance.
func operatorPerfScore(stakeShare float64, nonsigningRate float64) float64 {
	if nonsigningRate == 0 {
		return 1.0
	}
	sla := stakeShareToSLA(stakeShare / 100.0)
	perf := (1 - sla) / nonsigningRate
	return perf / (1.0 + perf)
}

func computePerfScore(metric *OperatorNonsigningPercentageMetrics) float64 {
	return operatorPerfScore(metric.StakePercentage, metric.Percentage)
}

type Ejector struct {
	wallet     walletsdk.Wallet
	ethClient  common.EthClient
	logger     logging.Logger
	transactor core.Transactor
	metrics    *Metrics
	txnTimeout time.Duration

	// For serializing the ejection requests.
	mu sync.Mutex
}

func NewEjector(wallet walletsdk.Wallet, ethClient common.EthClient, logger logging.Logger, tx core.Transactor, metrics *Metrics, txnTimeout time.Duration) *Ejector {
	return &Ejector{
		wallet:     wallet,
		ethClient:  ethClient,
		logger:     logger.With("component", "Ejector"),
		transactor: tx,
		metrics:    metrics,
		txnTimeout: txnTimeout,
	}
}

func (e *Ejector) Eject(ctx context.Context, nonsigningRate *OperatorsNonsigningPercentage) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	nonsigners := make([]*OperatorNonsigningPercentageMetrics, 0)
	for _, metric := range nonsigningRate.Data {
		// Collect only the nonsigners who violate the SLA.
		if metric.Percentage/100.0 > 1-stakeShareToSLA(metric.StakePercentage/100.0) {
			nonsigners = append(nonsigners, metric)
		}
	}

	// Rank the operators for each quorum by the operator performance score.
	// The operators with lower perf score will get ejected with priority in case of
	// rate limiting.
	sort.Slice(nonsigners, func(i, j int) bool {
		if nonsigners[i].QuorumId == nonsigners[j].QuorumId {
			if computePerfScore(nonsigners[i]) == computePerfScore(nonsigners[j]) {
				return float64(nonsigners[i].TotalUnsignedBatches)*nonsigners[i].StakePercentage > float64(nonsigners[j].TotalUnsignedBatches)*nonsigners[j].StakePercentage
			}
			return computePerfScore(nonsigners[i]) < computePerfScore(nonsigners[j])
		}
		return nonsigners[i].QuorumId < nonsigners[j].QuorumId
	})

	operatorsByQuorum, err := e.convertOperators(nonsigners)
	if err != nil {
		return err
	}

	txn, err := e.transactor.BuildEjectOperatorsTxn(ctx, operatorsByQuorum)
	if err != nil {
		e.logger.Error("Failed to build ejection transaction", "err", err)
		return err
	}

	var txID walletsdk.TxID
	retryFromFailure := 0
	for retryFromFailure < maxSendTransactionRetry {
		gasTipCap, gasFeeCap, err := e.ethClient.GetLatestGasCaps(ctx)
		if err != nil {
			return fmt.Errorf("failed to get latest gas caps: %w", err)
		}

		txn, err = e.ethClient.UpdateGas(ctx, txn, big.NewInt(0), gasTipCap, gasFeeCap)
		if err != nil {
			return fmt.Errorf("failed to update gas price: %w", err)
		}
		txID, err = e.wallet.SendTransaction(ctx, txn)
		var urlErr *url.Error
		didTimeout := false
		if errors.As(err, &urlErr) {
			didTimeout = urlErr.Timeout()
		}
		if didTimeout || errors.Is(err, context.DeadlineExceeded) {
			e.logger.Warn("failed to send txn due to timeout", "hash", txn.Hash().Hex(), "numRetries", retryFromFailure, "maxRetry", maxSendTransactionRetry, "err", err)
			retryFromFailure++
			continue
		} else if err != nil {
			return fmt.Errorf("failed to send txn %s: %w", txn.Hash().Hex(), err)
		} else {
			e.logger.Debug("successfully sent txn", "txID", txID, "txHash", txn.Hash().Hex())
			break
		}
	}

	queryTicker := time.NewTicker(queryTickerDuration)
	defer queryTicker.Stop()
	ctxWithTimeout, cancelCtx := context.WithTimeout(ctx, e.txnTimeout)
	defer cancelCtx()
	var receipt *types.Receipt
	for {
		receipt, err = e.wallet.GetTransactionReceipt(ctxWithTimeout, txID)
		if err == nil {
			break
		}

		if errors.Is(err, ethereum.NotFound) || errors.Is(err, walletsdk.ErrReceiptNotYetAvailable) {
			e.logger.Debug("Transaction not yet mined", "txID", txID, "txHash", txn.Hash().Hex(), "err", err)
		} else if errors.Is(err, walletsdk.ErrNotYetBroadcasted) {
			e.logger.Warn("Transaction has not been broadcasted to network but attempted to retrieve receipt", "err", err)
		} else if errors.Is(err, walletsdk.ErrTransactionFailed) {
			e.logger.Error("Transaction failed", "txID", txID, "txHash", txn.Hash().Hex(), "err", err)
			return err
		} else {
			e.logger.Error("Transaction receipt retrieval failed", "err", err)
			return err
		}

		// Wait for the next round.
		select {
		case <-ctxWithTimeout.Done():
			return ctxWithTimeout.Err()
		case <-queryTicker.C:
		}
	}

	e.logger.Info("Ejection transaction succeeded", "receipt", receipt)

	e.metrics.UpdateEjectionGasUsed(receipt.GasUsed)

	// TODO: get the txn response and update the metrics.

	return nil
}

func (e *Ejector) convertOperators(nonsigners []*OperatorNonsigningPercentageMetrics) ([][]core.OperatorID, error) {
	var maxQuorumId uint8
	for _, metric := range nonsigners {
		if metric.QuorumId > maxQuorumId {
			maxQuorumId = metric.QuorumId
		}
	}

	numOperatorByQuorum := make(map[uint8]int)
	stakeShareByQuorum := make(map[uint8]float64)

	result := make([][]core.OperatorID, maxQuorumId+1)
	for _, metric := range nonsigners {
		id, err := core.OperatorIDFromHex(metric.OperatorId)
		if err != nil {
			return nil, err
		}
		result[metric.QuorumId] = append(result[metric.QuorumId], id)
		numOperatorByQuorum[metric.QuorumId]++
		stakeShareByQuorum[metric.QuorumId] += metric.StakePercentage
	}
	e.metrics.UpdateRequestedOperatorMetric(numOperatorByQuorum, stakeShareByQuorum)

	return result, nil
}
