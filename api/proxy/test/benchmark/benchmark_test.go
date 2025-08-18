package benchmark

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
)

// BenchmarkPutsWithSecondaryV1  ... Takes in an async worker count and profiles blob insertions using
// constant blob sizes in parallel. Exercises V1 code pathways
func BenchmarkPutsWithSecondaryV1(b *testing.B) {
	testCfg := testutils.NewTestConfig(testutils.MemstoreBackend, common.V1EigenDABackend, nil)
	putsWithSecondary(b, testCfg)
}

// BenchmarkPutsWithSecondaryV2  ... Takes in an async worker count and profiles blob insertions using
// constant blob sizes in parallel. Exercises V2 code pathways
func BenchmarkPutsWithSecondaryV2(b *testing.B) {
	testCfg := testutils.NewTestConfig(testutils.MemstoreBackend, common.V2EigenDABackend, nil)
	putsWithSecondary(b, testCfg)
}

func putsWithSecondary(b *testing.B, testCfg testutils.TestConfig) {
	testCfg.UseS3Caching = true
	writeThreadCount := os.Getenv("WRITE_THREAD_COUNT")
	threadInt, err := strconv.Atoi(writeThreadCount)
	if err != nil {
		panic(fmt.Errorf("Could not parse WRITE_THREAD_COUNT field %w", err))
	}
	testCfg.WriteThreadCount = threadInt

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &standard_client.Config{
		URL: ts.Address(),
	}
	daClient := standard_client.New(cfg)

	for i := 0; i < b.N; i++ {
		_, err := daClient.SetData(
			b.Context(),
			[]byte("I am a blob and I only live for 14 days on EigenDA"))
		if err != nil {
			panic(err)
		}
	}
}
