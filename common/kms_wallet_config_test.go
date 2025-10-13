package common

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestKMSWalletCLIFlags(t *testing.T) {
	envVarPrefix := "TEST"
	flagPrefix := "test"

	flags := KMSWalletCLIFlags(envVarPrefix, flagPrefix)

	// Check that we have the expected number of flags
	assert.Len(t, flags, 4) // provider, key-id, key-region, disable

	flagNames := make(map[string]bool)
	for _, flag := range flags {
		switch f := flag.(type) {
		case cli.StringFlag:
			flagNames[f.Name] = true
		case cli.BoolFlag:
			flagNames[f.Name] = true
		}
	}

	// Debug: print actual flag names
	t.Logf("Actual flag names: %v", flagNames)

	// Check that all expected flags are present
	expectedFlags := []string{
		"test.kms-provider",
		"test.kms-key-id",
		"test.kms-key-region",
		"test.kms-key-disable",
	}

	for _, expectedFlag := range expectedFlags {
		assert.True(t, flagNames[expectedFlag], "Missing flag: %s", expectedFlag)
	}
}

func TestReadKMSKeyConfig(t *testing.T) {
	// Create a test CLI context
	app := cli.NewApp()
	flagPrefix := "test"
	envVarPrefix := "TEST"
	
	app.Flags = KMSWalletCLIFlags(envVarPrefix, flagPrefix)

	tests := []struct {
		name     string
		args     []string
		expected KMSKeyConfig
	}{
		{
			name: "AWS provider config",
			args: []string{"app", 
				"--test.kms-provider", "aws",
				"--test.kms-key-id", "test-key-123",
				"--test.kms-key-region", "us-east-1",
			},
			expected: KMSKeyConfig{
				Provider: "aws",
				KeyID:    "test-key-123",
				Region:   "us-east-1",
				Disable:  false,
			},
		},
		{
			name: "OCI provider config",
			args: []string{"app",
				"--test.kms-provider", "oci",
				"--test.kms-key-id", "ocid1.key.oc1.test",
			},
			expected: KMSKeyConfig{
				Provider: "oci",
				KeyID:    "ocid1.key.oc1.test",
				Disable:  false,
			},
		},
		{
			name: "disabled KMS",
			args: []string{"app", "--test.kms-key-disable"},
			expected: KMSKeyConfig{
				Provider: "aws", // Default value
				Disable:  true,
			},
		},
		{
			name: "empty config (defaults)",
			args: []string{"app"},
			expected: KMSKeyConfig{
				Provider: "aws", // Default value
				Disable:  false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Parse the flags
			set := flag.NewFlagSet("test", 0)
			for _, f := range app.Flags {
				f.Apply(set)
			}
			err := set.Parse(test.args[1:])
			require.NoError(t, err)

			// Create CLI context
			ctx := cli.NewContext(app, set, nil)

			// Test ReadKMSKeyConfig
			config := ReadKMSKeyConfig(ctx, flagPrefix)

			assert.Equal(t, test.expected.Provider, config.Provider)
			assert.Equal(t, test.expected.KeyID, config.KeyID)
			assert.Equal(t, test.expected.Region, config.Region)
			assert.Equal(t, test.expected.Disable, config.Disable)
		})
	}
}

func TestReadKMSKeyConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"TEST_KMS_PROVIDER":   "oci",
		"TEST_KMS_KEY_ID":     "ocid1.key.oc1.env",
		"TEST_KMS_KEY_REGION": "us-phoenix-1",
	}

	// Set environment variables
	for key, value := range envVars {
		_ = os.Setenv(key, value)
		defer func(k string) { _ = os.Unsetenv(k) }(key)
	}

	// Create a test CLI context
	app := cli.NewApp()
	flagPrefix := "test"
	envVarPrefix := "TEST"
	
	app.Flags = KMSWalletCLIFlags(envVarPrefix, flagPrefix)

	// Parse with no command line arguments (should use env vars)
	set := flag.NewFlagSet("test", 0)
	for _, f := range app.Flags {
		f.Apply(set)
	}
	err := set.Parse([]string{})
	require.NoError(t, err)

	ctx := cli.NewContext(app, set, nil)

	// Test ReadKMSKeyConfig
	config := ReadKMSKeyConfig(ctx, flagPrefix)

	assert.Equal(t, "oci", config.Provider)
	assert.Equal(t, "ocid1.key.oc1.env", config.KeyID)
	assert.Equal(t, "us-phoenix-1", config.Region)
	assert.False(t, config.Disable)
}

func TestKMSKeyConfig_Struct(t *testing.T) {
	config := KMSKeyConfig{
		Provider: "oci",
		KeyID:    "test-key",
		Region:   "us-west-1",
		Disable:  false,
	}

	assert.Equal(t, "oci", config.Provider)
	assert.Equal(t, "test-key", config.KeyID)
	assert.Equal(t, "us-west-1", config.Region)
	assert.False(t, config.Disable)
}

func TestKMSWalletCLIFlags_FlagProperties(t *testing.T) {
	envVarPrefix := "TEST"
	flagPrefix := "test"

	flags := KMSWalletCLIFlags(envVarPrefix, flagPrefix)

	// Test individual flag properties
	for _, flag := range flags {
		switch f := flag.(type) {
		case cli.StringFlag:
			assert.NotEmpty(t, f.Name, "StringFlag should have a name")
			assert.NotEmpty(t, f.Usage, "StringFlag should have usage text")
			assert.NotEmpty(t, f.EnvVar, "StringFlag should have env var")
			assert.Contains(t, f.EnvVar, envVarPrefix, "EnvVar should contain prefix")
		case cli.BoolFlag:
			assert.NotEmpty(t, f.Name, "BoolFlag should have a name")
			assert.NotEmpty(t, f.Usage, "BoolFlag should have usage text")
			assert.NotEmpty(t, f.EnvVar, "BoolFlag should have env var")
			assert.Contains(t, f.EnvVar, envVarPrefix, "EnvVar should contain prefix")
		}
	}
}

func TestReadKMSKeyConfig_CommandLineOverridesEnvironment(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("TEST_KMS_PROVIDER", "aws")
	_ = os.Setenv("TEST_KMS_KEY_ID", "env-key")
	defer func() { _ = os.Unsetenv("TEST_KMS_PROVIDER") }()
	defer func() { _ = os.Unsetenv("TEST_KMS_KEY_ID") }()

	// Create CLI app with flags that should override env vars
	app := cli.NewApp()
	flagPrefix := "test"
	envVarPrefix := "TEST"
	
	app.Flags = KMSWalletCLIFlags(envVarPrefix, flagPrefix)

	// Parse with command line arguments that should override env vars
	args := []string{"app", 
		"--test.kms-provider", "oci",
		"--test.kms-key-id", "cli-key",
	}

	set := flag.NewFlagSet("test", 0)
	for _, f := range app.Flags {
		f.Apply(set)
	}
	err := set.Parse(args[1:])
	require.NoError(t, err)

	ctx := cli.NewContext(app, set, nil)

	// Test that CLI args override environment variables
	config := ReadKMSKeyConfig(ctx, flagPrefix)

	assert.Equal(t, "oci", config.Provider) // Should be CLI value, not env
	assert.Equal(t, "cli-key", config.KeyID) // Should be CLI value, not env
}