package dispatcher

import (
	"errors"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
)

func TestFailureAggregator_AddTotalStake(t *testing.T) {
	tests := []struct {
		name   string
		stakes []*big.Int
		want   *big.Int
	}{
		{
			name:   "single stake",
			stakes: []*big.Int{big.NewInt(100)},
			want:   big.NewInt(100),
		},
		{
			name:   "multiple stakes",
			stakes: []*big.Int{big.NewInt(100), big.NewInt(200), big.NewInt(300)},
			want:   big.NewInt(600),
		},
		{
			name:   "zero stakes",
			stakes: []*big.Int{big.NewInt(0), big.NewInt(0)},
			want:   big.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)
			for _, stake := range tt.stakes {
				fa.AddTotalStake(stake)
			}
			if fa.TotalStake.Cmp(tt.want) != 0 {
				t.Errorf("TotalStake = %v, want %v", fa.TotalStake, tt.want)
			}
		})
	}
}

func TestFailureAggregator_AddFailure(t *testing.T) {
	tests := []struct {
		name             string
		failures         []OperatorFailure
		wantFailureCount int
		wantFailedStake  *big.Int
	}{
		{
			name: "single failure",
			failures: []OperatorFailure{
				{
					OperatorID: core.OperatorID{1},
					Stake:      big.NewInt(100),
					Error:      errors.New("test error"),
				},
			},
			wantFailureCount: 1,
			wantFailedStake:  big.NewInt(100),
		},
		{
			name: "multiple failures",
			failures: []OperatorFailure{
				{
					OperatorID: core.OperatorID{1},
					Stake:      big.NewInt(100),
					Error:      errors.New("test error 1"),
				},
				{
					OperatorID: core.OperatorID{2},
					Stake:      big.NewInt(200),
					Error:      errors.New("test error 2"),
				},
			},
			wantFailureCount: 2,
			wantFailedStake:  big.NewInt(300),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)
			for _, failure := range tt.failures {
				fa.AddFailure(failure)
			}
			if len(fa.Failures) != tt.wantFailureCount {
				t.Errorf("Failures count = %v, want %v", len(fa.Failures), tt.wantFailureCount)
			}
			if fa.FailedStake.Cmp(tt.wantFailedStake) != 0 {
				t.Errorf("FailedStake = %v, want %v", fa.FailedStake, tt.wantFailedStake)
			}
		})
	}
}

func TestFailureAggregator_GetStakePercentageForAccount(t *testing.T) {
	tests := []struct {
		name        string
		totalStake  *big.Int
		failures    []OperatorFailure
		accountID   string
		wantPercent float64
	}{
		{
			name:       "no failures",
			totalStake: big.NewInt(1000),
			failures:   []OperatorFailure{},
			accountID:  "0x123",
			wantPercent: 0.0,
		},
		{
			name:       "zero total stake",
			totalStake: big.NewInt(0),
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(100)},
			},
			accountID:   "0x123",
			wantPercent: 0.0,
		},
		{
			name:       "single account failure",
			totalStake: big.NewInt(1000),
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(100)},
			},
			accountID:   "0x123",
			wantPercent: 10.0,
		},
		{
			name:       "multiple failures same account",
			totalStake: big.NewInt(1000),
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(100)},
				{AccountID: "0x123", Stake: big.NewInt(200)},
				{AccountID: "0x456", Stake: big.NewInt(150)},
			},
			accountID:   "0x123",
			wantPercent: 30.0,
		},
		{
			name:       "account not found",
			totalStake: big.NewInt(1000),
			failures: []OperatorFailure{
				{AccountID: "0x456", Stake: big.NewInt(100)},
			},
			accountID:   "0x123",
			wantPercent: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)
			fa.TotalStake = new(big.Int).Set(tt.totalStake)
			for _, failure := range tt.failures {
				fa.AddFailure(failure)
			}
			
			got := fa.GetStakePercentageForAccount(tt.accountID)
			if got != tt.wantPercent {
				t.Errorf("GetStakePercentageForAccount() = %v, want %v", got, tt.wantPercent)
			}
		})
	}
}

func TestFailureAggregator_GetStakePercentageForBatchMeterErrors(t *testing.T) {
	tests := []struct {
		name        string
		totalStake  *big.Int
		failures    []OperatorFailure
		accountID   string
		wantPercent float64
	}{
		{
			name:       "no batch meter errors",
			totalStake: big.NewInt(1000),
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(100), IsBatchMeterError: false},
			},
			accountID:   "0x123",
			wantPercent: 0.0,
		},
		{
			name:       "batch meter errors only",
			totalStake: big.NewInt(1000),
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(200), IsBatchMeterError: true},
				{AccountID: "0x123", Stake: big.NewInt(100), IsBatchMeterError: false},
			},
			accountID:   "0x123",
			wantPercent: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)
			fa.TotalStake = new(big.Int).Set(tt.totalStake)
			for _, failure := range tt.failures {
				fa.AddFailure(failure)
			}
			
			got := fa.GetStakePercentageForBatchMeterErrors(tt.accountID)
			if got != tt.wantPercent {
				t.Errorf("GetStakePercentageForBatchMeterErrors() = %v, want %v", got, tt.wantPercent)
			}
		})
	}
}

func TestFailureAggregator_GetStakePercentageForAccountAndErrorType(t *testing.T) {
	tests := []struct {
		name        string
		totalStake  *big.Int
		failures    []OperatorFailure
		accountID   string
		errorCode   string
		wantPercent float64
	}{
		{
			name:       "specific error type match",
			totalStake: big.NewInt(1000),
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(100), IsBatchMeterError: true, BatchMeterCode: "RATE_LIMIT"},
				{AccountID: "0x123", Stake: big.NewInt(150), IsBatchMeterError: true, BatchMeterCode: "VALIDATION"},
			},
			accountID:   "0x123",
			errorCode:   "RATE_LIMIT",
			wantPercent: 10.0,
		},
		{
			name:       "no matching error type",
			totalStake: big.NewInt(1000),
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(100), IsBatchMeterError: true, BatchMeterCode: "VALIDATION"},
			},
			accountID:   "0x123",
			errorCode:   "RATE_LIMIT",
			wantPercent: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)
			fa.TotalStake = new(big.Int).Set(tt.totalStake)
			for _, failure := range tt.failures {
				fa.AddFailure(failure)
			}
			
			got := fa.GetStakePercentageForAccountAndErrorType(tt.accountID, tt.errorCode)
			if got != tt.wantPercent {
				t.Errorf("GetStakePercentageForAccountAndErrorType() = %v, want %v", got, tt.wantPercent)
			}
		})
	}
}

func TestFailureAggregator_GetBatchMeterErrorsByAccount(t *testing.T) {
	tests := []struct {
		name      string
		failures  []OperatorFailure
		accountID string
		wantCount int
	}{
		{
			name: "no batch meter errors",
			failures: []OperatorFailure{
				{AccountID: "0x123", IsBatchMeterError: false, Stake: big.NewInt(100)},
			},
			accountID: "0x123",
			wantCount: 0,
		},
		{
			name: "mixed errors",
			failures: []OperatorFailure{
				{AccountID: "0x123", IsBatchMeterError: true, Stake: big.NewInt(100)},
				{AccountID: "0x123", IsBatchMeterError: false, Stake: big.NewInt(200)},
				{AccountID: "0x456", IsBatchMeterError: true, Stake: big.NewInt(150)},
			},
			accountID: "0x123",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)
			for _, failure := range tt.failures {
				fa.AddFailure(failure)
			}
			
			got := fa.GetBatchMeterErrorsByAccount(tt.accountID)
			if len(got) != tt.wantCount {
				t.Errorf("GetBatchMeterErrorsByAccount() count = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

func TestFailureAggregator_createOperatorFailure(t *testing.T) {
	tests := []struct {
		name               string
		errorMsg           string
		wantBatchMeterErr  bool
		wantBatchMeterCode string
		wantAccountID      string
		wantCategory       string
	}{
		{
			name:               "non-batch meter error",
			errorMsg:           "connection timeout",
			wantBatchMeterErr:  false,
			wantBatchMeterCode: "",
			wantAccountID:      "",
			wantCategory:       "",
		},
		{
			name:               "batch meter error",
			errorMsg:           `{"code":"USAGE_EXCEEDS_LIMIT","message":"Usage exceeds limit","account_id":"0x123","quorum_id":1}`,
			wantBatchMeterErr:  true,
			wantBatchMeterCode: "USAGE_EXCEEDS_LIMIT",
			wantAccountID:      "0x123",
			wantCategory:       "rate_limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)
			operatorID := core.OperatorID{1}
			stake := big.NewInt(100)
			err := errors.New(tt.errorMsg)
			
			failure := fa.createOperatorFailure(operatorID, "socket", stake, err)
			
			if failure.IsBatchMeterError != tt.wantBatchMeterErr {
				t.Errorf("IsBatchMeterError = %v, want %v", failure.IsBatchMeterError, tt.wantBatchMeterErr)
			}
			if failure.BatchMeterCode != tt.wantBatchMeterCode {
				t.Errorf("BatchMeterCode = %v, want %v", failure.BatchMeterCode, tt.wantBatchMeterCode)
			}
			if failure.AccountID != tt.wantAccountID {
				t.Errorf("AccountID = %v, want %v", failure.AccountID, tt.wantAccountID)
			}
			if failure.ErrorCategory != tt.wantCategory {
				t.Errorf("ErrorCategory = %v, want %v", failure.ErrorCategory, tt.wantCategory)
			}
		})
	}
}