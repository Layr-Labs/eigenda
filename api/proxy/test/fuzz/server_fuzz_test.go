package fuzz_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

// FuzzProxyClientServerV1 will fuzz the proxy client server integration
// and op client keccak256 with malformed inputs. This is never meant to be fuzzed with EigenDA.
func FuzzProxyClientServerV1(f *testing.F) {
	fuzzProxyClientServer(f, common.V1EigenDABackend)
}

func FuzzProxyClientServerV2(f *testing.F) {
	fuzzProxyClientServer(f, common.V2EigenDABackend)
}

func fuzzProxyClientServer(f *testing.F, dispersalBackend common.EigenDABackend) {
	testCfg := testutils.NewTestConfig(testutils.MemstoreBackend, dispersalBackend, nil)
	testCfg.MaxBlobLength = "16mib"
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	// We want a silent logger for fuzzing because we need to see the output of the fuzzer itself,
	// which tells us each new interesting inputs it finds.
	logger := logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{Level: slog.LevelError})
	ts, kill := testutils.CreateTestSuite(tsConfig, testutils.TestSuiteWithLogger(logger))
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
