package e2e_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
	"github.com/Layr-Labs/eigenda-proxy/store"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"

	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/stretchr/testify/require"
)

// Integration tests are run against memstore whereas
// testnet integration tests are run against eigenda backend talking to testnet disperser.
// Some of the assertions in the tests are different based on the backend as well.
// e.g, in TestProxyServerCaching we only assert to read metrics with EigenDA
// when referencing memstore since we don't profile the eigenDAClient interactions
var (
	runTestnetIntegrationTests bool // holesky tests
	runIntegrationTests        bool // memstore V1 tests
	runIntegrationTestsV2      bool // memstore V2 tests
	runFuzzTests               bool // fuzz tests
)

func flagActivated(envVar string) bool {
	return os.Getenv(envVar) == "true" || os.Getenv(envVar) == "1"
}

// ParseEnv ... reads testing cfg fields. Go test flags don't work for this library due to the dependency on Optimism's E2E framework
// which initializes test flags per init function which is called before an init in this package.
func ParseEnv() {
	runFuzzTests = flagActivated("FUZZ")
	runIntegrationTestsV2 = flagActivated("INTEGRATION_V2")
	runIntegrationTests = flagActivated("INTEGRATION")
	runTestnetIntegrationTests = flagActivated("TESTNET")

	if runIntegrationTests && runTestnetIntegrationTests {
		panic("only one of INTEGRATION=true or TESTNET=true env var can be set")
	}

	if runIntegrationTests && runIntegrationTestsV2 {
		panic("only one of INTEGRATION=true or INTEGRATION_V2=true env var can be set")
	}

	fmt.Print("fuzz_tests", runFuzzTests, "integration_tests_v1", runIntegrationTests,
		"integration_tests_v2", runIntegrationTestsV2, "testnet_integration_tests", runTestnetIntegrationTests,
	)
}

// TestMain ... run main controller
func TestMain(m *testing.M) {
	ParseEnv()
	code := m.Run()
	os.Exit(code)
}

// requireDispersalRetrievalEigenDA ... ensure that blob was successfully dispersed/read to/from EigenDA
func requireDispersalRetrievalEigenDA(t *testing.T, cm *metrics.CountMap, mode commitments.CommitmentMode) {
	writeCount, err := cm.Get(string(mode), http.MethodPost)
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(string(mode), http.MethodGet)
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

// requireWriteReadSecondary ... ensure that secondary backend was successfully written/read to/from
func requireWriteReadSecondary(t *testing.T, cm *metrics.CountMap, bt common.BackendType) {
	writeCount, err := cm.Get(http.MethodPut, store.Success, bt.String())
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(http.MethodGet, store.Success, bt.String())
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

// requireStandardClientSetGet ... ensures that std proxy client can disperse and read a blob
func requireStandardClientSetGet(t *testing.T, ts e2e.TestSuite, blob []byte) {
	cfg := &standard_client.Config{
		URL: ts.Address(),
	}
	daClient := standard_client.New(cfg)

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, blob)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, blob, preimage)

}

// requireOPClientSetGet ... ensures that alt-da client can disperse and read a blob
func requireOPClientSetGet(t *testing.T, ts e2e.TestSuite, blob []byte, precompute bool) {
	daClient := altda.NewDAClient(ts.Address(), false, precompute)

	commit, err := daClient.SetInput(ts.Ctx, blob)
	require.NoError(t, err)

	preimage, err := daClient.GetInput(ts.Ctx, commit)
	require.NoError(t, err)
	require.Equal(t, blob, preimage)

}
