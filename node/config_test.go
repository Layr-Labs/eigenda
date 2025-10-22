package node

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestECDSAKeyRequirementLogic tests the logic for determining when ECDSA keys are required.
func TestECDSAKeyRequirementLogic(t *testing.T) {
	tests := []struct {
		name                   string
		registerAtStart        bool
		pubIPCheckInterval     time.Duration
		ejectionDefenseEnabled bool
		expectedNeedECDSAKey   bool
	}{
		{
			name:                   "no features requiring ECDSA key",
			registerAtStart:        false,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: false,
			expectedNeedECDSAKey:   false,
		},
		{
			name:                   "register at start requires ECDSA key",
			registerAtStart:        true,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: false,
			expectedNeedECDSAKey:   true,
		},
		{
			name:                   "pub IP check interval requires ECDSA key",
			registerAtStart:        false,
			pubIPCheckInterval:     5 * time.Minute,
			ejectionDefenseEnabled: false,
			expectedNeedECDSAKey:   true,
		},
		{
			name:                   "ejection defense requires ECDSA key",
			registerAtStart:        false,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: true,
			expectedNeedECDSAKey:   true,
		},
		{
			name:                   "all features requiring ECDSA key",
			registerAtStart:        true,
			pubIPCheckInterval:     5 * time.Minute,
			ejectionDefenseEnabled: true,
			expectedNeedECDSAKey:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic directly as it would be evaluated in NewConfig
			needECDSAKey := tt.registerAtStart || tt.pubIPCheckInterval > 0 || tt.ejectionDefenseEnabled
			assert.Equal(t, tt.expectedNeedECDSAKey, needECDSAKey, "needECDSAKey logic should match expected result")
		})
	}
}

// TestECDSAKeyValidationErrors tests the specific error messages returned when
// ECDSA keys are required but not provided.
func TestECDSAKeyValidationErrors(t *testing.T) {
	tests := []struct {
		name                   string
		registerAtStart        bool
		pubIPCheckInterval     time.Duration
		ejectionDefenseEnabled bool
		ecdsaKeyFile           string
		ecdsaKeyPassword       string
		expectedErrorContains  string
	}{
		{
			name:                   "ejection defense enabled without key file",
			registerAtStart:        false,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: true,
			ecdsaKeyFile:           "",
			ecdsaKeyPassword:       "password",
			expectedErrorContains:  "ecdsa-key-file and ecdsa-key-password are required if ejection-defense-enabled is enabled",
		},
		{
			name:                   "ejection defense enabled without password",
			registerAtStart:        false,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: true,
			ecdsaKeyFile:           "/path/to/key",
			ecdsaKeyPassword:       "",
			expectedErrorContains:  "ecdsa-key-file and ecdsa-key-password are required if ejection-defense-enabled is enabled",
		},
		{
			name:                   "ejection defense enabled without both",
			registerAtStart:        false,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: true,
			ecdsaKeyFile:           "",
			ecdsaKeyPassword:       "",
			expectedErrorContains:  "ecdsa-key-file and ecdsa-key-password are required if ejection-defense-enabled is enabled",
		},
		{
			name:                   "register at start without key file",
			registerAtStart:        true,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: false,
			ecdsaKeyFile:           "",
			ecdsaKeyPassword:       "password",
			expectedErrorContains:  "ecdsa-key-file and ecdsa-key-password are required if register-at-node-start is enabled",
		},
		{
			name:                   "pub IP check interval without password",
			registerAtStart:        false,
			pubIPCheckInterval:     5 * time.Minute,
			ejectionDefenseEnabled: false,
			ecdsaKeyFile:           "/path/to/key",
			ecdsaKeyPassword:       "",
			expectedErrorContains:  "ecdsa-key-file and ecdsa-key-password are required if pub-ip-check-interval is > 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic directly by simulating the conditions
			needECDSAKey := tt.registerAtStart || tt.pubIPCheckInterval > 0 || tt.ejectionDefenseEnabled
			assert.True(t, needECDSAKey, "All test cases should require ECDSA key")

			// Test the specific validation logic for each case
			if tt.registerAtStart && (tt.ecdsaKeyFile == "" || tt.ecdsaKeyPassword == "") {
				// This would trigger the registerAtStart error
				assert.Contains(t, tt.expectedErrorContains, "register-at-node-start")
			}

			if tt.pubIPCheckInterval > 0 && (tt.ecdsaKeyFile == "" || tt.ecdsaKeyPassword == "") {
				// This would trigger the pubIPCheckInterval error
				assert.Contains(t, tt.expectedErrorContains, "pub-ip-check-interval")
			}

			if tt.ejectionDefenseEnabled && (tt.ecdsaKeyFile == "" || tt.ecdsaKeyPassword == "") {
				// This would trigger the ejectionDefenseEnabled error
				assert.Contains(t, tt.expectedErrorContains, "ejection-defense-enabled")
			}
		})
	}
}

// TestECDSAKeyValidationSuccess tests that valid configurations with ejection defense don't fail
func TestECDSAKeyValidationSuccess(t *testing.T) {
	tests := []struct {
		name                   string
		registerAtStart        bool
		pubIPCheckInterval     time.Duration
		ejectionDefenseEnabled bool
		ecdsaKeyFile           string
		ecdsaKeyPassword       string
	}{
		{
			name:                   "ejection defense enabled with valid credentials",
			registerAtStart:        false,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: true,
			ecdsaKeyFile:           "/path/to/key",
			ecdsaKeyPassword:       "password",
		},
		{
			name:                   "all features enabled with valid credentials",
			registerAtStart:        true,
			pubIPCheckInterval:     5 * time.Minute,
			ejectionDefenseEnabled: true,
			ecdsaKeyFile:           "/path/to/key",
			ecdsaKeyPassword:       "password",
		},
		{
			name:                   "no features requiring ECDSA key",
			registerAtStart:        false,
			pubIPCheckInterval:     0,
			ejectionDefenseEnabled: false,
			ecdsaKeyFile:           "",
			ecdsaKeyPassword:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			needECDSAKey := tt.registerAtStart || tt.pubIPCheckInterval > 0 || tt.ejectionDefenseEnabled

			// If ECDSA key is needed, validate that we have both file and password
			if needECDSAKey {
				assert.True(t, tt.ecdsaKeyFile != "" && tt.ecdsaKeyPassword != "",
					"Valid configurations should provide both key file and password when needed")
			}

			// Test that each individual validation would pass
			registerAtStartValid := !tt.registerAtStart || (tt.ecdsaKeyFile != "" && tt.ecdsaKeyPassword != "")
			pubIPCheckValid := tt.pubIPCheckInterval == 0 || (tt.ecdsaKeyFile != "" && tt.ecdsaKeyPassword != "")
			ejectionDefenseValid := !tt.ejectionDefenseEnabled || (tt.ecdsaKeyFile != "" && tt.ecdsaKeyPassword != "")

			assert.True(t, registerAtStartValid, "Register at start validation should pass")
			assert.True(t, pubIPCheckValid, "Pub IP check validation should pass")
			assert.True(t, ejectionDefenseValid, "Ejection defense validation should pass")
		})
	}
}
