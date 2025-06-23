package controller

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// BatchMeterError represents a parsed batch meterer error
type BatchMeterError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	AccountID string `json:"account_id,omitempty"`
	QuorumID  uint8  `json:"quorum_id,omitempty"`
}

// BatchMeterErrorCategory represents error categories for handling
type BatchMeterErrorCategory int

const (
	ErrorCategoryUnknown BatchMeterErrorCategory = iota
	ErrorCategoryRateLimit
	ErrorCategoryReservation
	ErrorCategoryValidation
	ErrorCategorySystem
)

// Error categories mapped to error codes
var errorCategoryMap = map[string]BatchMeterErrorCategory{
	"BIN_ALREADY_FULL":           ErrorCategoryRateLimit,
	"USAGE_EXCEEDS_LIMIT":        ErrorCategoryRateLimit,
	"OVERFLOW_PERIOD_LIMIT":      ErrorCategoryRateLimit,
	"OVERFLOW_WINDOW_LIMIT":      ErrorCategoryRateLimit,
	"RESERVATION_NOT_FOUND":      ErrorCategoryReservation,
	"RESERVATION_LOOKUP_FAILED":  ErrorCategoryReservation,
	"RESERVATION_INACTIVE":       ErrorCategoryReservation,
	"RESERVATION_PERIOD_INVALID": ErrorCategoryReservation,
	"BATCH_EMPTY":                ErrorCategoryValidation,
	"BLOB_HEADER_NIL":            ErrorCategoryValidation,
	"PAYMENT_PARAMS_FAILED":      ErrorCategorySystem,
}

// Regex to parse the standardized error format: [CODE] message (account: 0x..., quorum: N)
var batchMeterErrorRegex = regexp.MustCompile(`^\[([A-Z_]+)\]\s+([^(]+)(?:\s+\(account:\s+([0-9a-fA-Fx]+)(?:,\s+quorum:\s+(\d+))?\))?`)

// ParseBatchMeterError attempts to parse a batch meterer error from an error string
func ParseBatchMeterError(errorStr string) (*BatchMeterError, bool) {
	// First try to parse as JSON (if validator sends structured errors)
	var bmErr BatchMeterError
	if err := json.Unmarshal([]byte(errorStr), &bmErr); err == nil && bmErr.Code != "" {
		return &bmErr, true
	}

	// Then try to parse the standardized string format
	matches := batchMeterErrorRegex.FindStringSubmatch(strings.TrimSpace(errorStr))
	if len(matches) < 3 {
		return nil, false
	}

	bmErr = BatchMeterError{
		Code:    strings.TrimSpace(matches[1]),
		Message: strings.TrimSpace(matches[2]),
	}

	// Extract account ID if present
	if len(matches) > 3 && matches[3] != "" {
		bmErr.AccountID = strings.TrimSpace(matches[3])
	}

	// Extract quorum ID if present
	if len(matches) > 4 && matches[4] != "" {
		if quorumID, err := strconv.ParseUint(matches[4], 10, 8); err == nil {
			bmErr.QuorumID = uint8(quorumID)
		}
	}

	return &bmErr, true
}

// GetBatchMeterErrorCategory returns the category of a batch meterer error
func GetBatchMeterErrorCategory(bmErr *BatchMeterError) BatchMeterErrorCategory {
	if category, exists := errorCategoryMap[bmErr.Code]; exists {
		return category
	}
	return ErrorCategoryUnknown
}

// GetBatchMeterErrorSummary returns a human-readable summary of the error
func GetBatchMeterErrorSummary(bmErr *BatchMeterError) string {
	category := GetBatchMeterErrorCategory(bmErr)

	var categoryStr string
	switch category {
	case ErrorCategoryRateLimit:
		categoryStr = "Rate Limit"
	case ErrorCategoryReservation:
		categoryStr = "Reservation"
	case ErrorCategoryValidation:
		categoryStr = "Validation"
	case ErrorCategorySystem:
		categoryStr = "System"
	default:
		categoryStr = "Unknown"
	}

	if bmErr.AccountID != "" && bmErr.QuorumID != 0 {
		return fmt.Sprintf("%s error for account %s on quorum %d: %s",
			categoryStr, bmErr.AccountID, bmErr.QuorumID, bmErr.Message)
	} else if bmErr.AccountID != "" {
		return fmt.Sprintf("%s error for account %s: %s",
			categoryStr, bmErr.AccountID, bmErr.Message)
	}

	return fmt.Sprintf("%s error: %s", categoryStr, bmErr.Message)
}

// LogBatchMeterError logs a parsed batch meterer error with appropriate context
func LogBatchMeterError(logger logging.Logger, bmErr *BatchMeterError, originalErr error) {
	category := GetBatchMeterErrorCategory(bmErr)

	logFields := []any{
		"code", bmErr.Code,
		"category", category,
		"message", bmErr.Message,
	}

	if bmErr.AccountID != "" {
		logFields = append(logFields, "accountID", bmErr.AccountID)
	}

	if bmErr.QuorumID != 0 {
		logFields = append(logFields, "quorumID", bmErr.QuorumID)
	}

	if originalErr != nil {
		logFields = append(logFields, "originalError", originalErr.Error())
	}

	logger.Error("Batch meterer error", logFields...)
}
