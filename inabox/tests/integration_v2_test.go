package integration_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	"github.com/Layr-Labs/eigenda/api/clients/v2"

	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
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

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
		signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
		Expect(err).To(BeNil())

		// TODO: update this to use the payload disperser client instead since it wraps the
		//       disperser client and provides additional functionality
		//       after doing an experiment with this its currently infeasible since the disperser chooses the latest block #
		//       for cert rbn which breaks an invariant on the registry coordinator that requires block.number > rbn
		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))

		payload1 := randomPayload(992)
		blob1, err := payload1.ToBlob(codecs.PolynomialFormEval)
		Expect(err).To(BeNil())

		payload2 := randomPayload(123)
		blob2, err := payload2.ToBlob(codecs.PolynomialFormEval)
		Expect(err).To(BeNil())

		reply1, err := disperseBlob(disp, blob1.Serialize())
		Expect(err).To(BeNil())

		reply2, err := disperseBlob(disp, blob2.Serialize())
		Expect(err).To(BeNil())

		// necessary to ensure that reference block number < current block number
		mineAnvilBlocks(1)

		// test onchain verification using cert #1
		eigenDACert, err := certBuilder.BuildCert(ctx, coretypes.VersionThreeCert, reply1)
		Expect(err).To(BeNil())

		err = staticCertVerifier.CheckDACert(ctx, eigenDACert)
		Expect(err).To(BeNil())

		err = routerCertVerifier.CheckDACert(ctx, eigenDACert)
		Expect(err).To(BeNil())

		// test onchain verification using cert #2
		eigenDACert2, err := certBuilder.BuildCert(ctx, coretypes.VersionThreeCert, reply2)
		Expect(err).To(BeNil())

		err = staticCertVerifier.CheckDACert(ctx, eigenDACert2)
		Expect(err).To(BeNil())

		err = routerCertVerifier.CheckDACert(ctx, eigenDACert2)
		Expect(err).To(BeNil())

		eigenDAV3Cert1, ok := eigenDACert.(*coretypes.EigenDACertV3)
		Expect(ok).To(BeTrue())

		eigenDAV3Cert2, ok := eigenDACert2.(*coretypes.EigenDACertV3)
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

		// ensure that a verifier can be added two blocks in the future
		tx, err := eigenDACertVerifierRouter.AddCertVerifier(deployerTransactorOpts, uint32(latestBlock)+2, gethcommon.HexToAddress("0x0"))
		Expect(err).To(BeNil())

		mineAnvilBlocks(1)
		latestBlock += 1

		// ensure that tx successfully executed
		receipt, err := ethClient.TransactionReceipt(ctx, tx.Hash())
		Expect(err).To(BeNil())
		Expect(receipt).To(Not(BeNil()))
		Expect(receipt.Status).To(Equal(uint64(1)))

		// ensure that new verifier can be read from the contract at the future rbn
		verifier, err := eigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock+1))
		Expect(err).To(BeNil())
		Expect(verifier).To(Equal(gethcommon.HexToAddress("0x0")))

		// and that old one still lives at the latest block number - 1
		verifier, err = eigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock-1))
		Expect(err).To(BeNil())
		Expect(verifier.String()).To(Equal(testConfig.EigenDA.CertVerifier))

		// progress anvil chain to ensure latest block number is set for rbn so
		// invalid verifier can be triggered
		mineAnvilBlocks(3)

		// disperse blob #3 to trigger the new cert verifier

		blob3, err := randomPayload(1000).ToBlob(codecs.PolynomialFormEval)
		Expect(err).To(BeNil())

		reply3, err := disperseBlob(disp, blob3.Serialize())
		Expect(err).To(BeNil())

		// necessary to ensure that reference block number < current block number
		mineAnvilBlocks(1)

		eigenDACert3, err := certBuilder.BuildCert(ctx, coretypes.VersionThreeCert, reply3)
		Expect(err).To(BeNil())

		// should fail since the rbn -> cert verifier is an empty address which should trigger a reversion
		err = routerCertVerifier.CheckDACert(ctx, eigenDACert3)
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("no contract code at given address"))

		// should pass using old verifier
		err = staticCertVerifier.CheckDACert(ctx, eigenDACert3)
		Expect(err).To(BeNil())

	})
})

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

func disperseBlob(disp clients.DisperserClient, blob []byte) (*disperserpb.BlobStatusReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	_, key, err := disp.DisperseBlob(blob, 0, []uint8{0, 1})
	if err != nil {
		return nil, err
	}

	var reply *disperserpb.BlobStatusReply

	for {
		select {
		case <-ctx.Done():
			Fail("timed out")
		case <-ticker.C:
			reply, err = disp.GetBlobStatus(context.Background(), key)
			if err != nil {
				return nil, err
			}

			status, err := dispv2.BlobStatusFromProtobuf(reply.GetStatus())
			if err != nil {
				return nil, err
			}

			if status != dispv2.Complete {
				continue
			}
			return reply, nil
		}
	}
}
