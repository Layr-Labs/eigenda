package integration_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/docker/go-units"
	"google.golang.org/grpc"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	nodegrpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/core"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/core/payment"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/sha3"
)

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

var _ = Describe("Inabox v2 blacklisting Integration test", func() {
	It("test end to end scenario of blacklisting", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// random G1 point
		g1Commitment, err := RandomG1Point()
		if err != nil {
			Fail("failed to generate random G1 point")
		}

		g2Commitment, err := RandomG2Point()
		if err != nil {
			Fail("failed to generate random G2 point")
		}

		// random data blob certificate
		blobCert := &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion: 2,
				BlobCommitments: encoding.BlobCommitments{
					Commitment:       &g1Commitment,
					LengthCommitment: &g2Commitment,
					LengthProof:      &g2Commitment,
					Length:           100,
				},
				QuorumNumbers: []core.QuorumID{0, 1},
				PaymentMetadata: payment.PaymentMetadata{
					AccountID:         gethcommon.HexToAddress("0x1234567890123456789012345678901234567890"),
					Timestamp:         time.Now().UnixNano(),
					CumulativePayment: big.NewInt(100),
				},
			},
			Signature: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65},
			RelayKeys: []corev2.RelayKey{0, 1},
		}

		blobCertProto, err := blobCert.ToProtobuf()
		if err != nil {
			Fail("failed to convert blob certificate to protobuf")
		}
		fmt.Println("blobCertProto size", len(blobCertProto.String()))

		mineAnvilBlocks(1)
		// println("latest block number", deploy.GetLatestBlockNumber("http://localhost:8545"))

		request := &nodegrpc.StoreChunksRequest{
			Batch: &commonpb.Batch{
				Header: &commonpb.BatchHeader{
					BatchRoot:            []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					ReferenceBlockNumber: 70,
				},
				BlobCertificates: []*commonpb.BlobCertificate{
					blobCertProto,
				},
			},
			DisperserID: api.EigenLabsDisperserID,
			Timestamp:   uint32(time.Now().Unix()),
			Signature:   []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65},
		}

		hash, err := hashing.HashStoreChunksRequest(request)
		if err != nil {
			Fail("failed to hash request")
		}

		keyManager := kms.New(kms.Options{
			Region:       "us-east-1",
			BaseEndpoint: aws.String("http://localhost:4570"),
		})

		// pick the first key in the key manager
		keys, err := keyManager.ListKeys(ctx, &kms.ListKeysInput{})
		if err != nil {
			Fail("failed to list keys")
		}
		keyID := keys.Keys[0].KeyId

		publicKey, err := aws2.LoadPublicKeyKMS(ctx, keyManager, *keyID)
		if err != nil {
			Fail("failed to load public key")
		}

		signature, err := aws2.SignKMS(ctx, keyManager, *keyID, publicKey, hash)
		if err != nil {
			Fail("failed to sign request")
		}

		request.Signature = signature

		addr := fmt.Sprintf("%v:%v", "localhost", "32017")
		dialOptions := clients.GetGrpcDialOptions(false, 4*units.MiB)
		// conn, err := grpc.NewClient(addr, dialOptions...)
		// if err != nil {
		// 	Fail("failed to create grpc connection")
		// }
		conn, err := grpc.NewClient(addr, dialOptions...)
		if err != nil {
			Fail("failed to create grpc connection")
		}
		dispersalClient := nodegrpc.NewDispersalClient(conn)

		// after this request, the disperser should be blacklisted
		_, err = dispersalClient.StoreChunks(ctx, request)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(ContainSubstring("failed to validate blob request"))

		// should get error saying disperser is blacklisted
		_, err = dispersalClient.StoreChunks(ctx, request)
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(ContainSubstring("disperser is blacklisted"))

	})
})

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

		// ensure that a verifier #2 can be added two blocks in the future where activation_block_number = latestBlock + 2
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
		verifier, err = eigenDACertVerifierRouterCaller.GetCertVerifierAt(&bind.CallOpts{}, uint32(latestBlock-1))
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

		tx, err = eigenDACertVerifierRouter.AddCertVerifier(deployerTransactorOpts, uint32(latestBlock)+2, gethcommon.HexToAddress(testConfig.EigenDA.CertVerifier))
		Expect(err).To(BeNil())
		mineAnvilBlocks(10)

		err = validateTxReceipt(ctx, tx.Hash())
		Expect(err).To(BeNil())

		// ensure that new verifier #3 can be used for successful verification
		// now disperse blob #4 to trigger the new cert verifier which should pass
		// ensure that a verifier can be added two blocks in the future
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
