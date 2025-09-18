package e2e

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func TestArbCustomDAIsValidHeaderByte(t *testing.T) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), common.V2EigenDABackend, nil)

	appCfg := testutils.BuildTestSuiteConfig(testCfg)
	testSuite, teardown := testutils.CreateTestSuite(appCfg)
	defer teardown()

	rpcClient, err := rpc.Dial(testSuite.ArbAddress())
	require.NoError(t, err)

	var validHeaderByteResult *arbitrum_altda.IsValidHeaderByteResult
	err = rpcClient.Call(&validHeaderByteResult, arbitrum_altda.MethodIsValidHeaderByte, arbitrum_altda.EigenDAV2MessageHeaderByte)
	require.NoError(t, err)

	var validHeaderByteFalseResult *arbitrum_altda.IsValidHeaderByteResult
	err = rpcClient.Call(&validHeaderByteFalseResult, arbitrum_altda.MethodIsValidHeaderByte, 0x69)
	require.NoError(t, err)

}

func TestArbCustomDAStoreAndRecoverPayloadFromBatch(t *testing.T) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), common.V2EigenDABackend, nil)

	appCfg := testutils.BuildTestSuiteConfig(testCfg)
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

	var recoverPayloadResult *arbitrum_altda.RecoverPayloadFromBatchResult
	batchNum := hexutil.Uint(0)
	batchBlockHash := gethcommon.HexToHash("0x43")
	preimageMap := new(arbitrum_altda.PreimagesMap)
	validateSeqMsg := false

	err = rpcClient.Call(&recoverPayloadResult, arbitrum_altda.MethodRecoverBatchFromPayload,
		batchNum,
		batchBlockHash,
		storeResult.SerializedDACert,
		preimageMap,
		validateSeqMsg,
	)
	require.NoError(t, err)

}
