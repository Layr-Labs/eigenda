package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorOnSecondaryInsertFailureValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "Valid: flag OFF, async OFF",
			config: Config{
				AsyncPutWorkers:               0,
				ErrorOnSecondaryInsertFailure: false,
			},
			wantError: false,
		},
		{
			name: "Valid: flag OFF, async ON",
			config: Config{
				AsyncPutWorkers:               5,
				ErrorOnSecondaryInsertFailure: false,
			},
			wantError: false,
		},
		{
			name: "Valid: flag ON, async OFF",
			config: Config{
				AsyncPutWorkers:               0,
				ErrorOnSecondaryInsertFailure: true,
			},
			wantError: false,
		},
		{
			name: "Invalid: flag ON, async ON",
			config: Config{
				AsyncPutWorkers:               5,
				ErrorOnSecondaryInsertFailure: true,
			},
			wantError: true,
			errorMsg:  "requires synchronous writes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Check()
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
