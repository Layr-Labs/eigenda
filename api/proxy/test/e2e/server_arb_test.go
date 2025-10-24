package e2e

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func TestArbCustomDAGetSupportedHeaderBytesMethod(t *testing.T) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), common.V2EigenDABackend, nil)
	appCfg := testutils.BuildTestSuiteConfig(testCfg)
	appCfg.EnabledServersConfig = &enablement.EnabledServersConfig{
		Metric:        false,
		ArbCustomDA:   true,
		RestAPIConfig: enablement.RestApisEnabled{},
	}

	testSuite, teardown := testutils.CreateTestSuite(appCfg)
	defer teardown()

	rpcClient, err := rpc.Dial(testSuite.ArbAddress())
	require.NoError(t, err)

	var supportedHeaderBytesResult *arbitrum_altda.SupportedHeaderBytesResult
	err = rpcClient.Call(&supportedHeaderBytesResult,
		arbitrum_altda.MethodGetSupportedHeaderBytes)
	require.NoError(t, err)
	require.Equal(t, supportedHeaderBytesResult.HeaderBytes[0], uint8(commitments.ArbCustomDAHeaderByte))

}

func TestArbCustomDAStoreAndRecoverMethods(t *testing.T) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), common.V2EigenDABackend, nil)
	appCfg := testutils.BuildTestSuiteConfig(testCfg)
	appCfg.EnabledServersConfig = &enablement.EnabledServersConfig{
		Metric:        false,
		ArbCustomDA:   true,
		RestAPIConfig: enablement.RestApisEnabled{},
	}

	testSuite, teardown := testutils.CreateTestSuite(appCfg)
	defer teardown()

	rpcClient, err := rpc.Dial(testSuite.ArbAddress())
	require.NoError(t, err)

	var storeResult *arbitrum_altda.StoreResult
	seqMessageArg := "0xDEADBEEF"
	timeoutArg := hexutil.Uint(200)
	disableFallbackStoreDataOnChain := false

	err = rpcClient.Call(&storeResult, arbitrum_altda.MethodStore,
		seqMessageArg,
		timeoutArg,
		disableFallbackStoreDataOnChain)
	require.NoError(t, err)

	var recoverPayloadResult *arbitrum_altda.PayloadResult
	batchNum := hexutil.Uint(0)
	batchBlockHash := gethcommon.HexToHash("0x43")

	err = rpcClient.Call(&recoverPayloadResult, arbitrum_altda.MethodRecoverPayload,
		batchNum,
		batchBlockHash,
		storeResult.SerializedDACert,
	)
	require.NoError(t, err)

}
