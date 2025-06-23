package dispatcher

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
)

func TestFailureAggregator_ShouldRetryBatch(t *testing.T) {
	tests := []struct {
		name                        string
		failures                    []OperatorFailure
		totalStake                  *big.Int
		maxAccountFailurePercentage float64
		wantShouldRetry             bool
		wantReason                  string
		wantTriggeringAccounts      []string
	}{
		{
			name:            "no failures",
			failures:        []OperatorFailure{},
			wantShouldRetry: false,
			wantReason:      "no failures detected",
		},
		{
			name: "account failure exceeds threshold",
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(300)}, // 30% for account
			},
			totalStake:                  big.NewInt(1000),
			maxAccountFailurePercentage: 20.0,
			wantShouldRetry:             true,
			wantReason:                  "account 0x123 failure 30.00% exceeds threshold 20.00%",
			wantTriggeringAccounts:      []string{"0x123"},
		},
		{
			name: "multiple accounts exceed threshold",
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(300)}, // 30% for account
				{AccountID: "0x456", Stake: big.NewInt(250)}, // 25% for account
			},
			totalStake:                  big.NewInt(1000),
			maxAccountFailurePercentage: 20.0,
			wantShouldRetry:             true,
			wantReason:                  "2 accounts exceeded failure threshold 20.00%",
			wantTriggeringAccounts:      []string{"0x123", "0x456"},
		},
		{
			name: "no retry conditions met",
			failures: []OperatorFailure{
				{AccountID: "0x123", Stake: big.NewInt(100)}, // 10% - under threshold
			},
			totalStake:                  big.NewInt(1000),
			maxAccountFailurePercentage: 20.0,
			wantShouldRetry:             false,
			wantReason:                  "no retry conditions met",
			wantTriggeringAccounts:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFailureAggregator(nil)

			if tt.totalStake != nil {
				fa.TotalStake = new(big.Int).Set(tt.totalStake)
			}

			for _, failure := range tt.failures {
				fa.AddFailure(failure)
			}

			decision := fa.ShouldRetryBatch(tt.maxAccountFailurePercentage)

			if decision.ShouldRetry != tt.wantShouldRetry {
				t.Errorf("ShouldRetry = %v, want %v", decision.ShouldRetry, tt.wantShouldRetry)
			}

			if decision.Reason != tt.wantReason {
				t.Errorf("Reason = %v, want %v", decision.Reason, tt.wantReason)
			}

			if len(decision.TriggeringAccounts) != len(tt.wantTriggeringAccounts) {
				t.Errorf("TriggeringAccounts count = %v, want %v", len(decision.TriggeringAccounts), len(tt.wantTriggeringAccounts))
			}

			for _, expectedAccount := range tt.wantTriggeringAccounts {
				found := false
				for _, actualAccount := range decision.TriggeringAccounts {
					if actualAccount == expectedAccount {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected triggering account %v not found in %v", expectedAccount, decision.TriggeringAccounts)
				}
			}
		})
	}
}

func TestDispatcher_FilterBlobsByAccountIDs(t *testing.T) {
	tests := []struct {
		name               string
		blobs              []core.EncodedBlob
		excludeAccountIDs  []string
		expectedBlobCount  int
		expectedExcluded   int
	}{
		{
			name:               "no exclusions",
			blobs:              createTestBlobs([]string{"0x123", "0x456"}),
			excludeAccountIDs:  []string{},
			expectedBlobCount:  2,
			expectedExcluded:   0,
		},
		{
			name:               "exclude one account",
			blobs:              createTestBlobs([]string{"0x123", "0x456", "0x789"}),
			excludeAccountIDs:  []string{"0x456"},
			expectedBlobCount:  2,
			expectedExcluded:   1,
		},
		{
			name:               "exclude multiple accounts",
			blobs:              createTestBlobs([]string{"0x123", "0x456", "0x789"}),
			excludeAccountIDs:  []string{"0x123", "0x789"},
			expectedBlobCount:  1,
			expectedExcluded:   2,
		},
		{
			name:               "exclude all accounts",
			blobs:              createTestBlobs([]string{"0x123", "0x456"}),
			excludeAccountIDs:  []string{"0x123", "0x456"},
			expectedBlobCount:  0,
			expectedExcluded:   2,
		},
		{
			name:               "exclude non-existent account",
			blobs:              createTestBlobs([]string{"0x123", "0x456"}),
			excludeAccountIDs:  []string{"0x999"},
			expectedBlobCount:  2,
			expectedExcluded:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dispatcher{
				logger: nil, // Will be handled by the function
			}
			
			filteredBlobs := d.filterBlobsByAccountIDs(tt.blobs, tt.excludeAccountIDs)
			
			if len(filteredBlobs) != tt.expectedBlobCount {
				t.Errorf("Filtered blob count = %v, want %v", len(filteredBlobs), tt.expectedBlobCount)
			}
			
			excludedCount := len(tt.blobs) - len(filteredBlobs)
			if excludedCount != tt.expectedExcluded {
				t.Errorf("Excluded blob count = %v, want %v", excludedCount, tt.expectedExcluded)
			}
			
			// Verify that excluded accounts are not in filtered blobs
			for _, blob := range filteredBlobs {
				if blob.BlobHeader != nil {
					for _, excludeID := range tt.excludeAccountIDs {
						if blob.BlobHeader.AccountID == excludeID {
							t.Errorf("Found excluded account %v in filtered blobs", excludeID)
						}
					}
				}
			}
		})
	}
}

// Helper function to create test blobs with specific account IDs
func createTestBlobs(accountIDs []string) []core.EncodedBlob {
	blobs := make([]core.EncodedBlob, len(accountIDs))
	for i, accountID := range accountIDs {
		blobs[i] = core.EncodedBlob{
			BlobHeader: &core.BlobHeader{
				AccountID: accountID,
			},
			EncodedBundlesByOperator: make(map[core.OperatorID]core.EncodedBundles),
		}
	}
	return blobs
}
