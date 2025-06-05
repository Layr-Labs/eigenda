package integration_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/docker/go-units"
	"google.golang.org/grpc"

	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	nodegrpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	verifierbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	"github.com/Layr-Labs/eigenda/core"

	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	caws "github.com/Layr-Labs/eigenda/common/aws"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"golang.org/x/crypto/sha3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
				PaymentMetadata: core.PaymentMetadata{
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

		publicKey, err := caws.LoadPublicKeyKMS(ctx, keyManager, *keyID)
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

	_, key, err := disp.DisperseBlob(ctx, blob, 0, []uint8{0, 1})
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
