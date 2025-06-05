package integration_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inabox v2 Integration", func() {
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
	It("test end to end scenario", func() {
		ctx := context.Background()
		// mine finalization_delay # of blocks given sometimes registry coordinator updates can sometimes happen
		// in-between the current_block_number - finalization_block_delay. This ensures consistent test execution.
		mineAnvilBlocks(6)

		payload1 := randomPayload(992)
		payload2 := randomPayload(123)

		// certificates are verified within the payload disperser client
		cert1, err := payloadDisperser.SendPayload(ctx, payload1)
		Expect(err).To(BeNil())

		cert2, err := payloadDisperser.SendPayload(ctx, payload2)
		Expect(err).To(BeNil())

		err = staticCertVerifier.CheckDACert(ctx, cert1)
		Expect(err).To(BeNil())

		err = routerCertVerifier.CheckDACert(ctx, cert1)
		Expect(err).To(BeNil())

		// test onchain verification using cert #2
		err = staticCertVerifier.CheckDACert(ctx, cert2)
		Expect(err).To(BeNil())

		err = routerCertVerifier.CheckDACert(ctx, cert2)
		Expect(err).To(BeNil())

		eigenDAV3Cert1, ok := cert1.(*coretypes.EigenDACertV3)
		Expect(ok).To(BeTrue())

		eigenDAV3Cert2, ok := cert2.(*coretypes.EigenDACertV3)
		Expect(ok).To(BeTrue())

		// test retrieval from disperser relay subnet
		actualPayload1, err := relayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert1)
		Expect(err).To(BeNil())
		Expect(actualPayload1).To(Not(BeNil()))
		Expect(actualPayload1).To(Equal(payload1))

		actualPayload2, err := relayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert2)
		Expect(err).To(BeNil())
		Expect(actualPayload2).To(Not(BeNil()))
		Expect(actualPayload2).To(Equal(payload2))

		// test distributed retrieval from DA network validator nodes
		actualPayload1, err = validatorRetrievalClientV2.GetPayload(
			ctx,
			eigenDAV3Cert1,
		)
		Expect(err).To(BeNil())
		Expect(actualPayload1).To(Not(BeNil()))
		Expect(actualPayload1).To(Equal(payload1))

		actualPayload2, err = validatorRetrievalClientV2.GetPayload(
			ctx,
			eigenDAV3Cert2,
		)
		Expect(err).To(BeNil())
		Expect(actualPayload2).To(Not(BeNil()))
		Expect(actualPayload2).To(Equal(payload2))

		/*
			enforce correct functionality of the EigenDACertVerifierRouter contract:
				1. ensure that a verifier can't be added at the latest block number
				2. ensure that a verifier can be added two blocks in the future
				3. ensure that the new verifier can be read from the contract when queried using a future rbn
				4. ensure that the old verifier can still be read from the contract when queried using the latest block number
				5. ensure that the new verifier is used to verify a cert at the future rbn after dispersal
		*/

		// ensure that a verifier can't be added at the latest block number
		latestBlock, err := ethClient.BlockNumber(ctx)
		Expect(err).To(BeNil())
		_, err = eigenDACertVerifierRouter.AddCertVerifier(deployerTransactorOpts, uint32(latestBlock), gethcommon.HexToAddress("0x0"))
		Expect(err).Error()
		Expect(err.Error()).To(ContainSubstring(getSolidityFunctionSig("ABNNotInFuture(uint32)")))

		// ensure that a verifier can be added two blocks in the future where activation_block_number = latestBlock + 2
		tx, err := eigenDACertVerifierRouter.AddCertVerifier(deployerTransactorOpts, uint32(latestBlock)+2, gethcommon.HexToAddress("0x0"))
		Expect(err).To(BeNil())
		mineAnvilBlocks(1)

		// ensure that tx successfully executed
		err = validateTxReceipt(ctx, tx.Hash())
		Expect(err).To(BeNil())

		// ensure that new verifier can be read from the contract at the future rbn
		verifier, err := eigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock+2))
		Expect(err).To(BeNil())
		Expect(verifier).To(Equal(gethcommon.HexToAddress("0x0")))

		// and that old one still lives at the latest block number - 1
		verifier, err = eigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock-2))
		Expect(err).To(BeNil())
		Expect(verifier.String()).To(Equal(testConfig.EigenDA.CertVerifier))

		// progress anvil chain 10 blocks
		mineAnvilBlocks(10)

		// disperse blob #3 to trigger the new cert verifier which should fail
		// since the address is not a valid cert verifier and the GetQuorums call will fail
		payload3 := randomPayload(1234)
		cert3, err := payloadDisperser.SendPayload(ctx, payload3)
		Expect(err.Error()).To(ContainSubstring("no contract code at given address"))
		Expect(cert3).To(BeNil())

		latestBlock, err = ethClient.BlockNumber(ctx)
		Expect(err).To(BeNil())

		// now disperse blob #4 to trigger the new cert verifier which should pass
		// ensure that a verifier can be added two blocks in the future
		tx, err = eigenDACertVerifierRouter.AddCertVerifier(deployerTransactorOpts, uint32(latestBlock)+2, gethcommon.HexToAddress(testConfig.EigenDA.CertVerifier))
		Expect(err).To(BeNil())
		mineAnvilBlocks(10)

		err = validateTxReceipt(ctx, tx.Hash())
		Expect(err).To(BeNil())

		// ensure that new verifier can be used for successful verification
		payload4 := randomPayload(1234)
		cert4, err := payloadDisperser.SendPayload(ctx, payload4)
		Expect(err).To(BeNil())
		err = routerCertVerifier.CheckDACert(ctx, cert4)
		Expect(err).To(BeNil())

		err = staticCertVerifier.CheckDACert(ctx, cert4)
		Expect(err).To(BeNil())

		// now force verification to fail by modifying the cert contents
		eigenDAV3Cert4, ok := cert4.(*coretypes.EigenDACertV3)
		Expect(ok).To(BeTrue())

		// modify the merkle root of the batch header and ensure verification fails
		// TODO: Test other cert verification failure cases as well
		eigenDAV3Cert4.BatchHeader.BatchRoot = gethcommon.Hash{0x1, 0x2, 0x3, 0x4}

		err = routerCertVerifier.CheckDACert(ctx, eigenDAV3Cert4)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(ContainSubstring("Merkle inclusion proof for blob batch is invalid"))
		err = staticCertVerifier.CheckDACert(ctx, eigenDAV3Cert4)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(ContainSubstring("Merkle inclusion proof for blob batch is invalid"))
	})
})

func validateTxReceipt(ctx context.Context, txHash gethcommon.Hash) error {
	receipt, err := ethClient.TransactionReceipt(ctx, txHash)
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

func randomPayload(size int) *coretypes.Payload {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	return coretypes.NewPayload(data)
}
