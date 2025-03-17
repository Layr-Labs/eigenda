package fuzz_test

import (
	"log/slog"
	"os"

	"github.com/Layr-Labs/eigenda-proxy/testutils"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
)

// FuzzProxyClientServerV1 will fuzz the proxy client server integration
// and op client keccak256 with malformed inputs. This is never meant to be fuzzed with EigenDA.
func FuzzProxyClientServerV1(f *testing.F) {
	testCfg := testutils.NewTestConfig(true, false)
	fuzzProxyClientServer(f, testCfg)
}

func FuzzProxyClientServerV2(f *testing.F) {
	testCfg := testutils.NewTestConfig(true, true)
	fuzzProxyClientServer(f, testCfg)
}

func fuzzProxyClientServer(f *testing.F, testCfg testutils.TestConfig) {
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	tsSecretConfig := testutils.TestSuiteSecretConfig(testCfg)
	// We want a silent logger for fuzzing because we need to see the output of the fuzzer itself,
	// which tells us each new interesting inputs it finds.
	logger := logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{Level: slog.LevelError})
	ts, kill := testutils.CreateTestSuite(tsConfig, tsSecretConfig, testutils.TestSuiteWithLogger(logger))
	f.Cleanup(kill)

	f.Add([]byte{})
	f.Add([]byte("a"))
	b := make([]byte, 1<<20)
	f.Add(b)

	cfg := &standard_client.Config{
		URL: ts.Address(),
	}

	daClient := standard_client.New(cfg)

	// seed and data are expected. `seed` value is seed: {rune} and data is the one with the random byte(s)
	f.Fuzz(
		func(t *testing.T, data []byte) {
			_, err := daClient.SetData(ts.Ctx, data)
			require.NoError(t, err)
		})
}
