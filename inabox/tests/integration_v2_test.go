package integration_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func TestEndToEndV2Scenario(t *testing.T) {
	/*
		This end to end test ensures that:
		1. a blob can be dispersed using the lower level disperser client to successfully produce a blob status response
		2. the blob certificate can be verified on chain using the immutable static EigenDACertVerifier and EigenDACertVerifierRouter contracts
		3. the blob can be retrieved from the disperser relay using the blob certificate
		4. the blob can be retrieved from the DA validator network using the blob certificate
		5. updates to the EigenDACertVerifierRouter contract can be made to add a new cert verifier with at a future activation block number
		6. the new cert verifier will be used to verify the blob certificate at the future activation block number

		TODO: Decompose this test into smaller tests that cover each of the above steps individually.
	*/
	// Create a fresh test harness for this test
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err, "Failed to create test harness")
	defer testHarness.Cleanup()

	ctx := t.Context()
	// mine finalization_delay # of blocks given sometimes registry coordinator updates can sometimes happen
	// in-between the current_block_number - finalization_block_delay. This ensures consistent test execution.
	integration.MineAnvilBlocks(t, testHarness.RPCClient, 6)

	payload1 := randomPayload(992)
	payload2 := randomPayload(123)

	// certificates are verified within the payload disperser client
	cert1, err := testHarness.PayloadDisperser.SendPayload(ctx, payload1)
	require.NoError(t, err)

	cert2, err := testHarness.PayloadDisperser.SendPayload(ctx, payload2)
	require.NoError(t, err)

	err = testHarness.StaticCertVerifier.CheckDACert(ctx, cert1)
	require.NoError(t, err)

	err = testHarness.RouterCertVerifier.CheckDACert(ctx, cert1)
	require.NoError(t, err)

	// test onchain verification using cert #2
	err = testHarness.StaticCertVerifier.CheckDACert(ctx, cert2)
	require.NoError(t, err)

	err = testHarness.RouterCertVerifier.CheckDACert(ctx, cert2)
	require.NoError(t, err)

	eigenDAV3Cert1, ok := cert1.(*coretypes.EigenDACertV3)
	require.True(t, ok)

	eigenDAV3Cert2, ok := cert2.(*coretypes.EigenDACertV3)
	require.True(t, ok)

	// test retrieval from disperser relay subnet
	actualPayload1, err := testHarness.RelayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert1)
	require.NoError(t, err)
	require.NotNil(t, actualPayload1)
	require.Equal(t, payload1, actualPayload1)

	actualPayload2, err := testHarness.RelayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert2)
	require.NoError(t, err)
	require.NotNil(t, actualPayload2)
	require.Equal(t, payload2, actualPayload2)

	// test distributed retrieval from DA network validator nodes
	actualPayload1, err = testHarness.ValidatorRetrievalClientV2.GetPayload(
		ctx,
		eigenDAV3Cert1,
	)
	require.NoError(t, err)
	require.NotNil(t, actualPayload1)
	require.Equal(t, payload1, actualPayload1)

	actualPayload2, err = testHarness.ValidatorRetrievalClientV2.GetPayload(
		ctx,
		eigenDAV3Cert2,
	)
	require.NoError(t, err)
	require.NotNil(t, actualPayload2)
	require.Equal(t, payload2, actualPayload2)

	/*
		enforce correct functionality of the EigenDACertVerifierRouter contract:
			1. ensure that a verifier can't be added at the latest block number
			2. ensure that a verifier can be added two blocks in the future
			3. ensure that the new verifier can be read from the contract when queried using a future rbn
			4. ensure that the old verifier can still be read from the contract when queried using the latest block number
			5. ensure that the new verifier is used to verify a cert at the future rbn after dispersal
	*/

	// ensure that a verifier can't be added at the latest block number
	latestBlock, err := testHarness.EthClient.BlockNumber(ctx)
	require.NoError(t, err)

	opts, err := testHarness.GetDeployerTransactOpts()
	require.NoError(t, err)
	_, err = testHarness.EigenDACertVerifierRouter.AddCertVerifier(
		opts,
		uint32(latestBlock),
		gethcommon.HexToAddress("0x0"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), getSolidityFunctionSig("ABNNotInFuture(uint32)"))

	// ensure that a verifier #2 can be added two blocks in the future where activation_block_number = latestBlock + 2
	opts, err = testHarness.GetDeployerTransactOpts()
	require.NoError(t, err)
	tx, err := testHarness.EigenDACertVerifierRouter.AddCertVerifier(
		opts,
		uint32(latestBlock)+2,
		gethcommon.HexToAddress("0x0"),
	)
	require.NoError(t, err)
	integration.MineAnvilBlocks(t, testHarness.RPCClient, 1)

	// ensure that tx successfully executed
	err = validateTxReceipt(ctx, testHarness, tx.Hash())
	require.NoError(t, err)

	// ensure that new verifier can be read from the contract at the future rbn
	verifier, err := testHarness.EigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock+2))
	require.NoError(t, err)
	require.Equal(t, gethcommon.HexToAddress("0x0"), verifier)

	// and that old one still lives at the latest block number - 1
	verifier, err = testHarness.EigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock-1))
	require.NoError(t, err)
	require.Equal(t, globalInfra.TestConfig.EigenDA.CertVerifier, verifier.String())

	// progress anvil chain 10 blocks
	integration.MineAnvilBlocks(t, testHarness.RPCClient, 10)

	// disperse blob #3 to trigger the new cert verifier which should fail
	// since the address is not a valid cert verifier and the GetQuorums call will fail
	payload3 := randomPayload(1234)
	cert3, err := testHarness.PayloadDisperser.SendPayload(ctx, payload3)
	require.Contains(t, err.Error(), "no contract code at given address")
	require.Nil(t, cert3)

	latestBlock, err = testHarness.EthClient.BlockNumber(ctx)
	require.NoError(t, err)

	opts, err = testHarness.GetDeployerTransactOpts()
	require.NoError(t, err)
	tx, err = testHarness.EigenDACertVerifierRouter.AddCertVerifier(
		opts,
		uint32(latestBlock)+2,
		gethcommon.HexToAddress(globalInfra.TestConfig.EigenDA.CertVerifier),
	)
	require.NoError(t, err)
	integration.MineAnvilBlocks(t, testHarness.RPCClient, 10)

	err = validateTxReceipt(ctx, testHarness, tx.Hash())
	require.NoError(t, err)

	// ensure that new verifier #3 can be used for successful verification
	// now disperse blob #4 to trigger the new cert verifier which should pass
	// ensure that a verifier can be added two blocks in the future
	payload4 := randomPayload(1234)
	cert4, err := testHarness.PayloadDisperser.SendPayload(ctx, payload4)
	require.NoError(t, err)
	err = testHarness.RouterCertVerifier.CheckDACert(ctx, cert4)
	require.NoError(t, err)

	err = testHarness.StaticCertVerifier.CheckDACert(ctx, cert4)
	require.NoError(t, err)

	// now force verification to fail by modifying the cert contents
	eigenDAV3Cert4, ok := cert4.(*coretypes.EigenDACertV3)
	require.True(t, ok)

	// modify the merkle root of the batch header and ensure verification fails
	// TODO: Test other cert verification failure cases as well
	eigenDAV3Cert4.BatchHeader.BatchRoot = gethcommon.Hash{0x1, 0x2, 0x3, 0x4}

	var certErr *verification.CertVerifierInvalidCertError
	err = testHarness.RouterCertVerifier.CheckDACert(ctx, eigenDAV3Cert4)
	require.IsType(t, &verification.CertVerifierInvalidCertError{}, err)
	require.True(t, errors.As(err, &certErr))
	// TODO(samlaf): after we update to CertVerifier 4.0.0 whose checkDACert will return error bytes,
	// we should check that extra bytes returned start with signature of the InvalidInclusionProof error
	require.Equal(t, verification.StatusInvalidCert, certErr.StatusCode)

	err = testHarness.StaticCertVerifier.CheckDACert(ctx, eigenDAV3Cert4)
	require.IsType(t, &verification.CertVerifierInvalidCertError{}, err)
	require.True(t, errors.As(err, &certErr))
	// TODO(samlaf): after we update to CertVerifier 4.0.0 whose checkDACert will return error bytes,
	// we should check that extra bytes returned start with signature of the InvalidInclusionProof error
	require.Equal(t, verification.StatusInvalidCert, certErr.StatusCode)
}

func validateTxReceipt(ctx context.Context, testHarness *integration.TestHarness, txHash gethcommon.Hash) error {
	receipt, err := testHarness.EthClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return err
	}
	if receipt == nil {
		return fmt.Errorf("transaction receipt not found for hash: %s", txHash.Hex())
	}
	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}
	return nil
}

func getSolidityFunctionSig(methodSig string) string {
	sig := []byte(methodSig)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(sig)
	selector := hash.Sum(nil)[:4] // take the first 4 bytes for the function selector
	return "0x" + hex.EncodeToString(selector)
}

func randomPayload(size int) coretypes.Payload {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	return coretypes.Payload(data)
}
