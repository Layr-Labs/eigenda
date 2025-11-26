package integration_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/common/geth"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestArbCustomDAGetSupportedHeaderBytesMethod(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	proxyConfig.EnabledServersConfig = &enablement.EnabledServersConfig{
		Metric:        false,
		ArbCustomDA:   true,
		RestAPIConfig: enablement.RestApisEnabled{},
	}
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	ethClient, err := geth.SafeDial(t.Context(), ts.ArbAddress())
	require.NoError(t, err)
	rpcClient := ethClient.Client()

	var supportedHeaderBytesResult *arbitrum_altda.SupportedHeaderBytesResult
	err = rpcClient.Call(&supportedHeaderBytesResult,
		arbitrum_altda.MethodGetSupportedHeaderBytes)
	require.NoError(t, err)
	require.Equal(t, supportedHeaderBytesResult.HeaderBytes[0][0], uint8(commitments.ArbCustomDAHeaderByte))
}

func TestArbCustomDAStoreAndRecoverMethods(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	proxyConfig.EnabledServersConfig = &enablement.EnabledServersConfig{
		Metric:        false,
		ArbCustomDA:   true,
		RestAPIConfig: enablement.RestApisEnabled{},
	}
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	ethClient, err := geth.SafeDial(t.Context(), ts.ArbAddress())
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
