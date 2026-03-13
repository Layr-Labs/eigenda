package ejector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEjectorConfig_HasSufficientOnChainMirror(t *testing.T) {
	tests := []struct {
		name                      string
		ejectionRetryDelay        time.Duration
		ejectionFinalizationDelay time.Duration
		onChainCooldown           uint64 // seconds
		onChainFinalizationDelay  uint64 // seconds
		expectError               bool
		errorContains             string
	}{
		{
			name:                      "valid configuration",
			ejectionRetryDelay:        48 * time.Hour,
			ejectionFinalizationDelay: 2 * time.Hour,
			onChainCooldown:           uint64((24 * time.Hour).Seconds()),
			onChainFinalizationDelay:  uint64((1 * time.Hour).Seconds()),
			expectError:               false,
		},
		{
			name:                      "valid configuration - delays exactly equal to on-chain values",
			ejectionRetryDelay:        24 * time.Hour,
			ejectionFinalizationDelay: 1 * time.Hour,
			onChainCooldown:           uint64((24 * time.Hour).Seconds()),
			onChainFinalizationDelay:  uint64((1 * time.Hour).Seconds()),
			expectError:               false,
		},
		{
			name:                      "invalid - ejection retry delay too small",
			ejectionRetryDelay:        12 * time.Hour,
			ejectionFinalizationDelay: 2 * time.Hour,
			onChainCooldown:           uint64((24 * time.Hour).Seconds()),
			onChainFinalizationDelay:  uint64((1 * time.Hour).Seconds()),
			expectError:               true,
			errorContains:             "EjectionRetryDelay must be >= the on-chain cooldown period",
		},
		{
			name:                      "invalid - ejection finalization delay too small",
			ejectionRetryDelay:        48 * time.Hour,
			ejectionFinalizationDelay: 30 * time.Minute,
			onChainCooldown:           uint64((24 * time.Hour).Seconds()),
			onChainFinalizationDelay:  uint64((1 * time.Hour).Seconds()),
			expectError:               true,
			errorContains:             "EjectionFinalizationDelay must be >= the on-chain finalization delay period",
		},
		{
			name:                      "invalid - both delays too small",
			ejectionRetryDelay:        12 * time.Hour,
			ejectionFinalizationDelay: 30 * time.Minute,
			onChainCooldown:           uint64((24 * time.Hour).Seconds()),
			onChainFinalizationDelay:  uint64((1 * time.Hour).Seconds()),
			expectError:               true,
			errorContains:             "EjectionRetryDelay must be >= the on-chain cooldown period",
		},
		{
			name:                      "edge case - zero on-chain values",
			ejectionRetryDelay:        1 * time.Second,
			ejectionFinalizationDelay: 1 * time.Second,
			onChainCooldown:           0,
			onChainFinalizationDelay:  0,
			expectError:               false,
		},
		{
			name:                      "edge case - off by one second on retry delay",
			ejectionRetryDelay:        24*time.Hour - 1*time.Second,
			ejectionFinalizationDelay: 2 * time.Hour,
			onChainCooldown:           uint64((24 * time.Hour).Seconds()),
			onChainFinalizationDelay:  uint64((1 * time.Hour).Seconds()),
			expectError:               true,
			errorContains:             "EjectionRetryDelay must be >= the on-chain cooldown period",
		},
		{
			name:                      "edge case - off by one second on finalization delay",
			ejectionRetryDelay:        48 * time.Hour,
			ejectionFinalizationDelay: 1*time.Hour - 1*time.Second,
			onChainCooldown:           uint64((24 * time.Hour).Seconds()),
			onChainFinalizationDelay:  uint64((1 * time.Hour).Seconds()),
			expectError:               true,
			errorContains:             "EjectionFinalizationDelay must be >= the on-chain finalization delay period",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &EjectorConfig{
				EjectionRetryDelay:        tt.ejectionRetryDelay,
				EjectionFinalizationDelay: tt.ejectionFinalizationDelay,
			}

			err := config.HasSufficientOnChainMirror(tt.onChainCooldown, tt.onChainFinalizationDelay)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
