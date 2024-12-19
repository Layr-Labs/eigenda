package e2e_test

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"unicode"

	"github.com/Layr-Labs/eigenda-proxy/client"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
)

// FuzzProxyClientServerIntegrationAndOpClientKeccak256MalformedInputs will fuzz the proxy client server integration
// and op client keccak256 with malformed inputs. This is never meant to be fuzzed with EigenDA.
func FuzzProxyClientServerIntegration(f *testing.F) {
	if !runFuzzTests {
		f.Skip("Skipping test as FUZZ env var not set")
	}

	tsConfig := e2e.TestSuiteConfig(e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(tsConfig)

	for r := rune(0); r <= unicode.MaxRune; r++ {
		if unicode.IsPrint(r) {
			f.Add([]byte(string(r))) // Add each printable Unicode character as a seed
		}
	}

	cfg := &client.Config{
		URL: ts.Address(),
	}

	daClient := client.New(cfg)

	// seed and data are expected. `seed` value is seed: {rune} and data is the one with the random byte(s)
	f.Fuzz(func(t *testing.T, data []byte) {
		_, err := daClient.SetData(ts.Ctx, data)
		assert.NoError(t, err)
		if err != nil {
			t.Errorf("Failed to set data: %v", err)
		}
	})

	f.Cleanup(kill)

}
