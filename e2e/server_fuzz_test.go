package e2e_test

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"unicode"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
)

// FuzzProxyClientServerIntegrationAndOpClientKeccak256MalformedInputs will fuzz the proxy client server integration
// and op client keccak256 with malformed inputs. This is never meant to be fuzzed with EigenDA.
func FuzzProxyClientServerIntegration(f *testing.F) {
	if !runFuzzTests {
		f.Skip("Skipping test as FUZZ env var not set")
	}

	testCfg := e2e.TestConfig(useMemory(), runIntegrationTestsV2)

	tsConfig := e2e.TestSuiteConfig(testCfg)
	tsSecretConfig := e2e.TestSuiteSecretConfig(testCfg)
	ts, kill := e2e.CreateTestSuite(tsConfig, tsSecretConfig)

	for r := rune(0); r <= unicode.MaxRune; r++ {
		if unicode.IsPrint(r) {
			f.Add([]byte(string(r))) // Add each printable Unicode character as a seed
		}
	}

	cfg := &standard_client.Config{
		URL: ts.Address(),
	}

	daClient := standard_client.New(cfg)

	// seed and data are expected. `seed` value is seed: {rune} and data is the one with the random byte(s)
	f.Fuzz(
		func(t *testing.T, data []byte) {
			_, err := daClient.SetData(ts.Ctx, data)
			assert.NoError(t, err)
			if err != nil {
				t.Errorf("Failed to set data: %v", err)
			}
		})

	f.Cleanup(kill)

}
