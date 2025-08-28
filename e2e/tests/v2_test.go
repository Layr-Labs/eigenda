package e2e_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func TestV2EndToEnd(t *testing.T) {
	/*
		This end to end test ensures that:
		1. a blob can be dispersed using the lower level disperser client to successfully produce a blob status response
		2. the blob certificate can be verified on chain using the immutable static EigenDACertVerifier and EigenDACertVerifierRouter contracts
		3. the blob can be retrieved from the disperser relay using the blob certificate
		4. the blob can be retrieved from the DA validator network using the blob certificate
		5. updates to the EigenDACertVerifierRouter contract can be made to add a new cert verifier with at a future activation block number
		6. the new cert verifier will be used to verify the blob certificate at the future activation block number
	*/

	ctx := context.Background()
	
	// mine finalization_delay # of blocks given sometimes registry coordinator updates can sometimes happen
	// in-between the current_block_number - finalization_block_delay. This ensures consistent test execution.
	mineAnvilBlocks(t, 6)

	payload1 := randomPayload(992)
	payload2 := randomPayload(123)

	// certificates are verified within the payload disperser client
	cert1, err := suite.payloadDisperser.SendPayload(ctx, payload1)
	require.NoError(t, err)

	cert2, err := suite.payloadDisperser.SendPayload(ctx, payload2)
	require.NoError(t, err)

	err = suite.staticCertVerifier.CheckDACert(ctx, cert1)
	require.NoError(t, err)

	err = suite.routerCertVerifier.CheckDACert(ctx, cert1)
	require.NoError(t, err)

	// test onchain verification using cert #2
	err = suite.staticCertVerifier.CheckDACert(ctx, cert2)
	require.NoError(t, err)

	err = suite.routerCertVerifier.CheckDACert(ctx, cert2)
	require.NoError(t, err)

	eigenDAV3Cert1, ok := cert1.(*coretypes.EigenDACertV3)
	require.True(t, ok)

	eigenDAV3Cert2, ok := cert2.(*coretypes.EigenDACertV3)
	require.True(t, ok)

	// test retrieval from disperser relay subnet
	actualPayload1, err := suite.relayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert1)
	require.NoError(t, err)
	require.NotNil(t, actualPayload1)
	assert.Equal(t, payload1, actualPayload1)

	actualPayload2, err := suite.relayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert2)
	require.NoError(t, err)
	require.NotNil(t, actualPayload2)
	assert.Equal(t, payload2, actualPayload2)

	// test distributed retrieval from DA network validator nodes
	actualPayload1, err = suite.validatorRetrievalClientV2.GetPayload(
		ctx,
		eigenDAV3Cert1,
	)
	require.NoError(t, err)
	require.NotNil(t, actualPayload1)
	assert.Equal(t, payload1, actualPayload1)

	actualPayload2, err = suite.validatorRetrievalClientV2.GetPayload(
		ctx,
		eigenDAV3Cert2,
	)
	require.NoError(t, err)
	require.NotNil(t, actualPayload2)
	assert.Equal(t, payload2, actualPayload2)

	/*
		enforce correct functionality of the EigenDACertVerifierRouter contract:
			1. ensure that a verifier can't be added at the latest block number
			2. ensure that a verifier can be added two blocks in the future
			3. ensure that the new verifier can be read from the contract when queried using a future rbn
			4. ensure that the old verifier can still be read from the contract when queried using the latest block number
			5. ensure that the new verifier is used to verify a cert at the future rbn after dispersal
	*/

	// ensure that a verifier can't be added at the latest block number
	latestBlock, err := suite.ethClient.BlockNumber(ctx)
	require.NoError(t, err)
	_, err = suite.eigenDACertVerifierRouter.AddCertVerifier(suite.deployerTransactorOpts, uint32(latestBlock), gethcommon.HexToAddress("0x0"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), getSolidityFunctionSig("ABNNotInFuture(uint32)"))

	// ensure that a verifier #2 can be added two blocks in the future where activation_block_number = latestBlock + 2
	tx, err := suite.eigenDACertVerifierRouter.AddCertVerifier(suite.deployerTransactorOpts, uint32(latestBlock)+2, gethcommon.HexToAddress("0x0"))
	require.NoError(t, err)
	mineAnvilBlocks(t, 1)

	// ensure that tx successfully executed
	err = validateTxReceipt(ctx, t, tx.Hash())
	require.NoError(t, err)

	// ensure that new verifier can be read from the contract at the future rbn
	verifier, err := suite.eigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock+2))
	require.NoError(t, err)
	assert.Equal(t, gethcommon.HexToAddress("0x0"), verifier)

	// and that old one still lives at the latest block number - 1
	verifier, err = suite.eigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock-1))
	require.NoError(t, err)
	assert.Equal(t, suite.eigenDAV2CertVerifierAddress, verifier.String())

	// progress anvil chain 10 blocks
	mineAnvilBlocks(t, 10)

	// disperse blob #3 to trigger the new cert verifier which should fail
	// since the address is not a valid cert verifier and the GetQuorums call will fail
	payload3 := randomPayload(1234)
	cert3, err := suite.payloadDisperser.SendPayload(ctx, payload3)
	assert.Contains(t, err.Error(), "no contract code at given address")
	assert.Nil(t, cert3)

	latestBlock, err = suite.ethClient.BlockNumber(ctx)
	require.NoError(t, err)

	tx, err = suite.eigenDACertVerifierRouter.AddCertVerifier(suite.deployerTransactorOpts, uint32(latestBlock)+2, gethcommon.HexToAddress(suite.eigenDAV2CertVerifierAddress))
	require.NoError(t, err)
	mineAnvilBlocks(t, 10)

	err = validateTxReceipt(ctx, t, tx.Hash())
	require.NoError(t, err)

	// ensure that new verifier #3 can be used for successful verification
	// now disperse blob #4 to trigger the new cert verifier which should pass
	// ensure that a verifier can be added two blocks in the future
	payload4 := randomPayload(1234)
	cert4, err := suite.payloadDisperser.SendPayload(ctx, payload4)
	require.NoError(t, err)
	err = suite.routerCertVerifier.CheckDACert(ctx, cert4)
	require.NoError(t, err)

	err = suite.staticCertVerifier.CheckDACert(ctx, cert4)
	require.NoError(t, err)

	// now force verification to fail by modifying the cert contents
	eigenDAV3Cert4, ok := cert4.(*coretypes.EigenDACertV3)
	require.True(t, ok)

	// modify the merkle root of the batch header and ensure verification fails
	// TODO: Test other cert verification failure cases as well
	eigenDAV3Cert4.BatchHeader.BatchRoot = gethcommon.Hash{0x1, 0x2, 0x3, 0x4}

	err = suite.routerCertVerifier.CheckDACert(ctx, eigenDAV3Cert4)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Merkle inclusion proof for blob batch is invalid")
	err = suite.staticCertVerifier.CheckDACert(ctx, eigenDAV3Cert4)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Merkle inclusion proof for blob batch is invalid")
}

// Helper functions

func mineAnvilBlocks(t *testing.T, numBlocks int) {
	for i := 0; i < numBlocks; i++ {
		err := suite.rpcClient.CallContext(context.Background(), nil, "evm_mine")
		require.NoError(t, err)
	}
}

func validateTxReceipt(ctx context.Context, t *testing.T, txHash gethcommon.Hash) error {
	receipt, err := suite.ethClient.TransactionReceipt(ctx, txHash)
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

func RandomG1Point() (encoding.G1Commitment, error) {
	// 1) pick r ← 𝐹ᵣ at random
	var r fr.Element
	if _, err := r.SetRandom(); err != nil {
		return encoding.G1Commitment{}, err
	}

	// 2) compute P = r·G₁ in Jacobian form
	G1Jac, _, _, _ := bn254.Generators()
	var Pjac bn254.G1Jac

	var rBigInt big.Int
	r.BigInt(&rBigInt)
	Pjac.ScalarMultiplication(&G1Jac, &rBigInt)

	// 3) convert to affine (x, y)
	var Paff bn254.G1Affine
	Paff.FromJacobian(&Pjac)
	return encoding.G1Commitment(Paff), nil
}

func RandomG2Point() (encoding.G2Commitment, error) {

	// 1) pick r ← 𝐹ᵣ at random
	var r fr.Element
	if _, err := r.SetRandom(); err != nil {
		return encoding.G2Commitment{}, err
	}

	// 2) compute P = r·G₂ in Jacobian form
	_, g2Jac, _, _ := bn254.Generators()
	var Pjac bn254.G2Jac

	var rBigInt big.Int
	r.BigInt(&rBigInt)
	Pjac.ScalarMultiplication(&g2Jac, &rBigInt)

	// 3) convert to affine (x, y)
	var Paff bn254.G2Affine
	Paff.FromJacobian(&Pjac)
	return encoding.G2Commitment(Paff), nil
}