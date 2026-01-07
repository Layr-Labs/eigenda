package e2e

import (
	"encoding/hex"
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	ethClient, err := geth.SafeDial(t.Context(), testSuite.ArbAddress())
	require.NoError(t, err)
	rpcClient := ethClient.Client()

	var supportedHeaderBytesResult *arbitrum_altda.SupportedHeaderBytesResult
	err = rpcClient.Call(&supportedHeaderBytesResult,
		arbitrum_altda.MethodGetSupportedHeaderBytes)
	require.NoError(t, err)
	require.Len(t, supportedHeaderBytesResult.HeaderBytes, 1)
	require.Equal(t, supportedHeaderBytesResult.HeaderBytes[0], uint8(commitments.ArbCustomDAHeaderByte))

}

func TestArbCustomDAGetMaxMessageSizeMethod(t *testing.T) {
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

	// Calculate the expected max payload size from the config
	expectedMaxPayloadSize, err := codec.BlobSymbolsToMaxPayloadSize(
		uint32(appCfg.StoreBuilderConfig.ClientConfigV2.MaxBlobSizeBytes / encoding.BYTES_PER_SYMBOL))
	require.NoError(t, err)

	ethClient, err := geth.SafeDial(t.Context(), testSuite.ArbAddress())
	require.NoError(t, err)
	rpcClient := ethClient.Client()

	// ensure that the max payload size value returned is correct
	var maxMessageSizeResult *arbitrum_altda.MaxMessageSizeResult
	err = rpcClient.Call(&maxMessageSizeResult,
		arbitrum_altda.MethodGetMaxMessageSize)
	require.NoError(t, err)
	require.NotNil(t, maxMessageSizeResult)
	require.Equal(t, expectedMaxPayloadSize, uint32(maxMessageSizeResult.MaxSize))

	// ensure that the max payload size value is respected as an upper limit for dispersal attempts

	var storeResult *arbitrum_altda.StoreResult
	seqMessageArg := "0x" + hex.EncodeToString(testutils.RandBytes(int(expectedMaxPayloadSize)+5))
	timeoutArg := hexutil.Uint(200)

	err = rpcClient.Call(&storeResult, arbitrum_altda.MethodStore,
		seqMessageArg,
		timeoutArg)

	require.Error(t, err)
	require.Equal(t, err.Error(), arbitrum_altda.ErrMessageTooLarge.Error())
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

	ethClient, err := geth.SafeDial(t.Context(), testSuite.ArbAddress())
	require.NoError(t, err)
	rpcClient := ethClient.Client()

	var storeResult *arbitrum_altda.StoreResult
	seqMessageArg := "0xDEADBEEF"
	timeoutArg := hexutil.Uint(200)

	err = rpcClient.Call(&storeResult, arbitrum_altda.MethodStore,
		seqMessageArg,
		timeoutArg)
	require.NoError(t, err)

	var recoverPayloadResult *arbitrum_altda.PayloadResult
	batchNum := hexutil.Uint(0)
	batchBlockHash := gethcommon.HexToHash("0x43")

	// pad 40 bytes for "message header"
	seqMessage := hexutil.Bytes(make([]byte, 40))
	seqMessage = append(seqMessage, storeResult.SerializedDACert...)

	err = rpcClient.Call(&recoverPayloadResult, arbitrum_altda.MethodRecoverPayload,
		batchNum,
		batchBlockHash,
		seqMessage,
	)
	require.NoError(t, err)

}
