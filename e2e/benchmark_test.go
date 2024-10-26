package e2e

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/client"
)

// BenchmarkPutsWithSecondary  ... Takes in an async worker count and profiles blob insertions using
// constant blob sizes in parallel
func BenchmarkPutsWithSecondary(b *testing.B) {
	testCfg := TestConfig(true)
	testCfg.UseS3Caching = true
	writeThreadCount := os.Getenv("WRITE_THREAD_COUNT")
	threadInt, err := strconv.Atoi(writeThreadCount)
	if err != nil {
		panic(fmt.Errorf("Could not parse WRITE_THREAD_COUNT field %w", err))
	}
	testCfg.WriteThreadCount = threadInt

	tsConfig := TestSuiteConfig(testCfg)
	ts, kill := CreateTestSuite(tsConfig)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)

	for i := 0; i < b.N; i++ {
		_, err := daClient.SetData(context.Background(), []byte("I am a blob and I only live for 14 days on EigenDA"))
		if err != nil {
			panic(err)
		}
	}
}
