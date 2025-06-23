package dispatcher

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var retryableBatchMeterErrors = []string{
	"RESERVATION_NOT_FOUND",
	"RESERVATION_INACTIVE",
	"RESERVATION_PERIOD_INVALID",
	"BIN_ALREADY_FULL",
}

// OperatorFailure represents a single operator's failure with stake and error context
type OperatorFailure struct {
	OperatorID        core.OperatorID
	Socket            string
	Stake             *big.Int
	Error             error
	IsBatchMeterError bool
	BatchMeterCode    string
	AccountID         string
	ErrorCategory     string
}

// FailureAggregator collects and analyzes operator failures during batch processing
type FailureAggregator struct {
	Failures    []OperatorFailure
	TotalStake  *big.Int // Total stake of all operators in the batch
	FailedStake *big.Int // Total stake of operators that failed
	logger      logging.Logger
}

// NewFailureAggregator creates a new failure aggregator with logger
func NewFailureAggregator(logger logging.Logger) *FailureAggregator {
	return &FailureAggregator{
		Failures:    make([]OperatorFailure, 0),
		TotalStake:  big.NewInt(0),
		FailedStake: big.NewInt(0),
		logger:      logger,
	}
}

// AddFailure records an operator failure and updates failed stake
func (fa *FailureAggregator) AddFailure(failure OperatorFailure) {
	fa.Failures = append(fa.Failures, failure)
	fa.FailedStake.Add(fa.FailedStake, failure.Stake)
}

// AddTotalStake adds operator stake to total stake pool for percentage calculations
func (fa *FailureAggregator) AddTotalStake(stake *big.Int) {
	fa.TotalStake.Add(fa.TotalStake, stake)
}

// GetBatchMeterErrorsByAccount returns all batch meterer errors for a specific account
func (fa *FailureAggregator) GetBatchMeterErrorsByAccount(accountID string) []OperatorFailure {
	var result []OperatorFailure
	for _, failure := range fa.Failures {
		if failure.IsBatchMeterError && failure.AccountID == accountID {
			result = append(result, failure)
		}
	}
	return result
}

// GetBatchMeterErrorsByAccountAndType returns batch meterer errors filtered by account and error type
func (fa *FailureAggregator) GetBatchMeterErrorsByAccountAndType(accountID, errorCode string) []OperatorFailure {
	var result []OperatorFailure
	for _, failure := range fa.Failures {
		if failure.IsBatchMeterError && failure.AccountID == accountID && failure.BatchMeterCode == errorCode {
			result = append(result, failure)
		}
	}
	return result
}

// GetStakePercentageForBatchMeterErrors calculates percentage of stake affected by batch meterer errors for an account
func (fa *FailureAggregator) GetStakePercentageForBatchMeterErrors(accountID string) float64 {
	if fa.TotalStake.Cmp(big.NewInt(0)) == 0 {
		return 0.0
	}

	failedStake := big.NewInt(0)
	for _, failure := range fa.Failures {
		if failure.IsBatchMeterError && failure.AccountID == accountID {
			failedStake.Add(failedStake, failure.Stake)
		}
	}

	percentage := new(big.Float).Quo(
		new(big.Float).SetInt(failedStake),
		new(big.Float).SetInt(fa.TotalStake),
	)
	result, _ := percentage.Float64()
	return result * 100.0
}

// GetStakePercentageForAccountAndErrorType calculates stake percentage for specific account and error type combination
func (fa *FailureAggregator) GetStakePercentageForAccountAndErrorType(accountID, errorCode string) float64 {
	if fa.TotalStake.Cmp(big.NewInt(0)) == 0 {
		return 0.0
	}

	failedStake := big.NewInt(0)
	for _, failure := range fa.Failures {
		if failure.IsBatchMeterError && failure.AccountID == accountID && failure.BatchMeterCode == errorCode {
			failedStake.Add(failedStake, failure.Stake)
		}
	}

	percentage := new(big.Float).Quo(
		new(big.Float).SetInt(failedStake),
		new(big.Float).SetInt(fa.TotalStake),
	)
	result, _ := percentage.Float64()
	return result * 100.0
}

// GetStakePercentageForAccount calculates percentage of stake that failed due to any error for a specific account
func (fa *FailureAggregator) GetStakePercentageForAccount(accountID string) float64 {
	if fa.TotalStake.Cmp(big.NewInt(0)) == 0 {
		return 0.0
	}

	failedStake := big.NewInt(0)
	for _, failure := range fa.Failures {
		if failure.AccountID == accountID {
			failedStake.Add(failedStake, failure.Stake)
		}
	}

	percentage := new(big.Float).Quo(
		new(big.Float).SetInt(failedStake),
		new(big.Float).SetInt(fa.TotalStake),
	)
	result, _ := percentage.Float64()
	return result * 100.0
}

// createOperatorFailure parses error details and creates an OperatorFailure with batch meterer context
func (fa *FailureAggregator) createOperatorFailure(operatorID core.OperatorID, socket string, stake *big.Int, err error) OperatorFailure {
	failure := OperatorFailure{
		OperatorID: operatorID,
		Socket:     socket,
		Stake:      new(big.Int).Set(stake),
		Error:      err,
	}

	// Parse batch meterer errors for enhanced tracking
	if bmErr, parsed := controller.ParseBatchMeterError(err.Error()); parsed {
		failure.IsBatchMeterError = true
		failure.BatchMeterCode = bmErr.Code
		failure.AccountID = bmErr.AccountID

		// Categorize error for retry decision logic
		category := controller.GetBatchMeterErrorCategory(bmErr)
		switch category {
		case controller.ErrorCategoryRateLimit:
			failure.ErrorCategory = "rate_limit"
		case controller.ErrorCategoryReservation:
			failure.ErrorCategory = "reservation"
		case controller.ErrorCategoryValidation:
			failure.ErrorCategory = "validation"
		case controller.ErrorCategorySystem:
			failure.ErrorCategory = "system"
		default:
			failure.ErrorCategory = "unknown"
		}
	}

	return failure
}

// RetryDecision represents the decision on whether to retry a batch based on failure analysis
type RetryDecision struct {
	ShouldRetry        bool
	Reason             string
	FailedStakePercent float64
	MaxAccountFailures map[string]float64 // accountID -> failure percentage
	TriggeringAccounts []string           // accountIDs that triggered retry and should be filtered out
}

// ShouldRetryBatch analyzes failures and determines if batch should be retried based on config
func (fa *FailureAggregator) ShouldRetryBatch(
	maxAccountFailurePercentage float64,
) *RetryDecision {
	decision := &RetryDecision{
		ShouldRetry:        false,
		MaxAccountFailures: make(map[string]float64),
		TriggeringAccounts: make([]string, 0),
	}

	// No failures means no retry needed
	if len(fa.Failures) == 0 {
		decision.Reason = "no failures detected"
		return decision
	}

	// Calculate overall failed stake percentage
	totalFailedStake := new(big.Float).SetInt(fa.FailedStake)
	totalStake := new(big.Float).SetInt(fa.TotalStake)
	if fa.TotalStake.Cmp(big.NewInt(0)) > 0 {
		percentage := new(big.Float).Quo(totalFailedStake, totalStake)
		decision.FailedStakePercent, _ = percentage.Float64()
		decision.FailedStakePercent *= 100.0
	}

	// Count retryable errors and check account-specific thresholds
	for _, failure := range fa.Failures {
		// Track account-specific failure percentages
		if failure.AccountID != "" && fa.isRetryableBatchMeterError(failure.BatchMeterCode, retryableBatchMeterErrors) {
			accountFailurePercent := fa.GetStakePercentageForAccount(failure.AccountID)
			decision.MaxAccountFailures[failure.AccountID] = accountFailurePercent

			if accountFailurePercent > maxAccountFailurePercentage {
				decision.ShouldRetry = true
				decision.TriggeringAccounts = append(decision.TriggeringAccounts, failure.AccountID)
			}
		}
	}

	// Set appropriate reason based on results
	if decision.ShouldRetry {
		if len(decision.TriggeringAccounts) == 1 {
			decision.Reason = fmt.Sprintf("account %s failure %.2f%% exceeds threshold %.2f%%",
				decision.TriggeringAccounts[0], decision.MaxAccountFailures[decision.TriggeringAccounts[0]], maxAccountFailurePercentage)
		} else {
			decision.Reason = fmt.Sprintf("%d accounts exceeded failure threshold %.2f%%",
				len(decision.TriggeringAccounts), maxAccountFailurePercentage)
		}
	} else {
		decision.Reason = "no retry conditions met"
	}

	return decision
}

// isRetryableBatchMeterError checks if a batch meterer error code is configured as retryable
func (fa *FailureAggregator) isRetryableBatchMeterError(errorCode string, retryableCodes []string) bool {
	for _, retryableCode := range retryableCodes {
		if errorCode == retryableCode {
			return true
		}
	}
	return false
}

// LogRetryDecision logs the retry decision analysis for monitoring
func (fa *FailureAggregator) LogRetryDecision(decision *RetryDecision) {
	if fa.logger == nil {
		return
	}

	if decision.ShouldRetry {
		fa.logger.Warn("Batch retry recommended",
			"reason", decision.Reason,
			"failed_stake_percent", fmt.Sprintf("%.2f%%", decision.FailedStakePercent),
			"account_failures", decision.MaxAccountFailures,
		)
	} else {
		fa.logger.Info("No batch retry needed",
			"reason", decision.Reason,
			"failed_stake_percent", fmt.Sprintf("%.2f%%", decision.FailedStakePercent),
		)
	}
}

// LogFailureStatistics provides comprehensive batch meterer error reporting by account and stake
func (fa *FailureAggregator) LogFailureStatistics() {
	if fa.logger == nil || len(fa.Failures) == 0 {
		return
	}

	batchMeterFailures := 0
	accountErrorMap := make(map[string]map[string]int)

	// Aggregate batch meterer errors by account and error type
	for _, failure := range fa.Failures {
		if failure.IsBatchMeterError {
			batchMeterFailures++

			if accountErrorMap[failure.AccountID] == nil {
				accountErrorMap[failure.AccountID] = make(map[string]int)
			}
			accountErrorMap[failure.AccountID][failure.BatchMeterCode]++
		}
	}

	// Log comprehensive error statistics for monitoring
	if batchMeterFailures > 0 {
		fa.logger.Info("Batch meterer error summary",
			"total_failures", len(fa.Failures),
			"batch_meterer_failures", batchMeterFailures,
			"total_stake", fa.TotalStake.String(),
			"failed_stake", fa.FailedStake.String(),
		)

		// Provide account-level breakdown with stake impact analysis
		for accountID, errorCounts := range accountErrorMap {
			stakePercentage := fa.GetStakePercentageForBatchMeterErrors(accountID)
			fa.logger.Info("Account-specific batch meterer errors",
				"account_id", accountID,
				"stake_percentage", fmt.Sprintf("%.2f%%", stakePercentage),
				"error_counts", errorCounts,
			)
		}
	}
}

// LogBatchMeterError logs individual batch meterer errors for account-level monitoring
func (fa *FailureAggregator) LogBatchMeterError(operatorFailure OperatorFailure, operatorID core.OperatorID, socket string) {
	if fa.logger == nil || !operatorFailure.IsBatchMeterError {
		return
	}

	fa.logger.Info("Batch meterer error detected",
		"operator", operatorID.Hex(),
		"code", operatorFailure.BatchMeterCode,
		"category", operatorFailure.ErrorCategory,
		"account", operatorFailure.AccountID,
		"socket", socket,
	)
}
