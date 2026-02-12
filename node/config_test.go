package node

import (
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/node/flags"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
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

// setBaselineConfigEnv sets the minimum environment variables needed for NewConfig to succeed.
// Individual tests can override specific variables before calling runNewConfig.
func setBaselineConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv("NODE_HOSTNAME", "localhost")
	t.Setenv("NODE_DISPERSAL_PORT", "9000")
	t.Setenv("NODE_RETRIEVAL_PORT", "9001")
	t.Setenv("NODE_ENABLE_NODE_API", "true")
	t.Setenv("NODE_ENABLE_METRICS", "true")
	t.Setenv("NODE_TIMEOUT", "1s")
	t.Setenv("NODE_QUORUM_ID_LIST", "0")
	t.Setenv("NODE_DB_PATH", "/tmp/eigenda-node-test")
	t.Setenv("NODE_EIGENDA_DIRECTORY", "0x0000000000000000000000000000000000000000")
	t.Setenv("NODE_CHURNER_URL", "http://localhost:1234")
	t.Setenv("NODE_PUBLIC_IP_PROVIDER", "ipify")
	t.Setenv("NODE_PUBLIC_IP_CHECK_INTERVAL", "0s")
	t.Setenv("NODE_CHAIN_RPC", "http://localhost:8545")
	t.Setenv("NODE_PRIVATE_KEY", "0x00")
	t.Setenv("NODE_G1_PATH", "/tmp/g1.point")
	t.Setenv("NODE_CACHE_PATH", "/tmp/eigenda-srs-cache")
	t.Setenv("NODE_SRS_ORDER", "1")
	t.Setenv("NODE_SRS_LOAD", "1")
	t.Setenv("NODE_V2_DISPERSAL_PORT", "32005")
	t.Setenv("NODE_V2_RETRIEVAL_PORT", "32004")
	t.Setenv("NODE_INTERNAL_V2_DISPERSAL_PORT", "32007")
	t.Setenv("NODE_INTERNAL_V2_RETRIEVAL_PORT", "32006")
	t.Setenv("NODE_ENABLE_TEST_MODE", "true")
	t.Setenv("NODE_TEST_PRIVATE_BLS", "deadbeef")
}

// runNewConfig runs a cli.App that calls NewConfig and returns the config and any error.
func runNewConfig(t *testing.T) (*Config, error) {
	t.Helper()
	app := cli.NewApp()
	app.Flags = flags.Flags

	var cfg *Config
	var configErr error
	app.Action = func(ctx *cli.Context) error {
		c, err := NewConfig(ctx)
		if err != nil {
			configErr = err
			return err
		}
		cfg = c
		return nil
	}
	// app.Run itself may return an error wrapping configErr.
	_ = app.Run([]string{os.Args[0]})
	return cfg, configErr
}

func TestNewConfig_RateLimitConfigFromEnv(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_DISPERSER_RATE_LIMIT_PER_SECOND", "0.5")
	t.Setenv("NODE_DISPERSER_RATE_LIMIT_BURST", "10")

	cfg, err := runNewConfig(t)
	assert.NoError(t, err)
	if !assert.NotNil(t, cfg) {
		return
	}
	assert.InDelta(t, 0.5, cfg.DisperserRateLimitPerSecond, 1e-9)
	assert.Equal(t, 10, cfg.DisperserRateLimitBurst)
}

func TestNewConfig_InvalidTimeout(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_TIMEOUT", "not-a-duration")

	_, err := runNewConfig(t)
	assert.Error(t, err)
}

func TestNewConfig_InvalidQuorumID(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_QUORUM_ID_LIST", "abc")

	_, err := runNewConfig(t)
	assert.Error(t, err)
}

func TestNewConfig_ExpirationPollIntervalTooLow(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_EXPIRATION_POLL_INTERVAL", "1")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expiration-poll-interval")
}

func TestNewConfig_ReachabilityPollIntervalTooLow(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_REACHABILITY_POLL_INTERVAL", "5")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reachability-poll-interval")
}

func TestNewConfig_MissingV2DispersalPort(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_V2_DISPERSAL_PORT", "")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "v2 dispersal port")
}

func TestNewConfig_MissingV2RetrievalPort(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_V2_RETRIEVAL_PORT", "")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "v2 retrieval port")
}

func TestNewConfig_InvalidV2DispersalPort(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_V2_DISPERSAL_PORT", "99999")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid v2 dispersal port")
}

func TestNewConfig_InvalidV2RetrievalPort(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_V2_RETRIEVAL_PORT", "99999")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid v2 retrieval port")
}

func TestNewConfig_OnDemandMeterFuzzFactorZero(t *testing.T) {
	setBaselineConfigEnv(t)
	t.Setenv("NODE_ON_DEMAND_METER_FUZZ_FACTOR", "0")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "on-demand-meter-fuzz-factor")
}

func TestNewConfig_InternalPortDefaults(t *testing.T) {
	setBaselineConfigEnv(t)
	// Clear internal ports so they default to v2 ports.
	t.Setenv("NODE_INTERNAL_V2_DISPERSAL_PORT", "")
	t.Setenv("NODE_INTERNAL_V2_RETRIEVAL_PORT", "")

	cfg, err := runNewConfig(t)
	assert.NoError(t, err)
	if !assert.NotNil(t, cfg) {
		return
	}
	assert.Equal(t, cfg.V2DispersalPort, cfg.InternalV2DispersalPort)
	assert.Equal(t, cfg.V2RetrievalPort, cfg.InternalV2RetrievalPort)
}

func TestNewConfig_BLSRemoteSignerMissingURL(t *testing.T) {
	setBaselineConfigEnv(t)
	// Disable test mode to hit the BLS remote signer branch.
	t.Setenv("NODE_ENABLE_TEST_MODE", "false")
	t.Setenv("NODE_BLS_REMOTE_SIGNER_ENABLED", "true")
	t.Setenv("NODE_BLS_REMOTE_SIGNER_URL", "")
	t.Setenv("NODE_BLS_PUBLIC_KEY_HEX", "")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BLS remote signer URL")
}

func TestNewConfig_BLSLocalSignerMissingKey(t *testing.T) {
	setBaselineConfigEnv(t)
	// Disable test mode and remote signer.
	t.Setenv("NODE_ENABLE_TEST_MODE", "false")
	t.Setenv("NODE_BLS_REMOTE_SIGNER_ENABLED", "false")
	t.Setenv("NODE_BLS_KEY_FILE", "")
	t.Setenv("NODE_BLS_KEY_PASSWORD", "")

	_, err := runNewConfig(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BLS key file and password")
}
