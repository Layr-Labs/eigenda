package clients

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Claude generated tests... don't blame the copy paster.
func TestEigenDAClientConfig_CheckAndSetDefaults(t *testing.T) {
	// Helper function to create a valid base config
	newValidConfig := func() *EigenDAClientConfig {
		return &EigenDAClientConfig{
			RPC:            "http://localhost:8080",
			EthRpcUrl:      "http://localhost:8545",
			SvcManagerAddr: "0x1234567890123456789012345678901234567890",
		}
	}

	t.Run("Valid minimal configuration", func(t *testing.T) {
		config := newValidConfig()
		err := config.CheckAndSetDefaults()
		require.NoError(t, err)

		// Check default values are set
		assert.Equal(t, 5*time.Second, config.StatusQueryRetryInterval)
		assert.Equal(t, 25*time.Minute, config.StatusQueryTimeout)
		assert.Equal(t, 30*time.Second, config.ResponseTimeout)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		testCases := []struct {
			name        string
			modifyConf  func(*EigenDAClientConfig)
			expectedErr string
		}{
			{
				name: "Missing RPC",
				modifyConf: func(c *EigenDAClientConfig) {
					c.RPC = ""
				},
				expectedErr: "EigenDAClientConfig.RPC not set",
			},
			{
				name: "Missing EthRpcUrl",
				modifyConf: func(c *EigenDAClientConfig) {
					c.EthRpcUrl = ""
				},
				expectedErr: "EigenDAClientConfig.EthRpcUrl not set. Needed to verify blob confirmed on-chain.",
			},
			{
				name: "Missing SvcManagerAddr",
				modifyConf: func(c *EigenDAClientConfig) {
					c.SvcManagerAddr = ""
				},
				expectedErr: "EigenDAClientConfig.SvcManagerAddr not set. Needed to verify blob confirmed on-chain.",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				config := newValidConfig()
				tc.modifyConf(config)
				err := config.CheckAndSetDefaults()
				assert.EqualError(t, err, tc.expectedErr)
			})
		}
	})

	t.Run("SignerPrivateKeyHex validation", func(t *testing.T) {
		testCases := []struct {
			name        string
			keyHex      string
			shouldError bool
		}{
			{
				name:        "Empty key (valid for read-only)",
				keyHex:      "",
				shouldError: false,
			},
			{
				name:        "Valid length key (64 bytes)",
				keyHex:      "1234567890123456789012345678901234567890123456789012345678901234",
				shouldError: false,
			},
			{
				name:        "Invalid length key",
				keyHex:      "123456",
				shouldError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				config := newValidConfig()
				config.SignerPrivateKeyHex = tc.keyHex
				err := config.CheckAndSetDefaults()
				if tc.shouldError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "SignerPrivateKeyHex")
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Custom timeouts", func(t *testing.T) {
		config := newValidConfig()
		customRetryInterval := 10 * time.Second
		customQueryTimeout := 30 * time.Minute
		customResponseTimeout := 45 * time.Second

		config.StatusQueryRetryInterval = customRetryInterval
		config.StatusQueryTimeout = customQueryTimeout
		config.ResponseTimeout = customResponseTimeout

		err := config.CheckAndSetDefaults()
		require.NoError(t, err)

		assert.Equal(t, customRetryInterval, config.StatusQueryRetryInterval)
		assert.Equal(t, customQueryTimeout, config.StatusQueryTimeout)
		assert.Equal(t, customResponseTimeout, config.ResponseTimeout)
	})

	t.Run("Optional fields", func(t *testing.T) {
		config := newValidConfig()
		config.CustomQuorumIDs = []uint{2, 3, 4}
		config.DisableTLS = true
		config.DisablePointVerificationMode = true

		err := config.CheckAndSetDefaults()
		require.NoError(t, err)

		assert.Equal(t, []uint{2, 3, 4}, config.CustomQuorumIDs)
		assert.True(t, config.DisableTLS)
		assert.True(t, config.DisablePointVerificationMode)
	})
}
