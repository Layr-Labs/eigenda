package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigenda/api/clients/v2"

	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inabox v2 Integration", func() {
	/*
		This end to end test ensures that:
		1. a blob can be dispersed using the lower level disperser client to successfully produce a blob status response
		2. the blob certificate can be verified on chain using the EigenDACertVerifier and EigenDACertVerifierRouter contracts
		3. the blob can be retrieved from the disperser relay using the blob certificate
		4. the blob can be retrieved from the DA validator network using the blob certificate

		TODO: Decompose this test into smaller tests that cover each of the above steps individually.
	*/
	It("test end to end scenario", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
		signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
		Expect(err).To(BeNil())

		// TODO: update this to use the payload disperser client instead since it wraps the 
		//       disperser client and provides additional functionality
		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))

		data1 := make([]byte, 992)
		_, err = rand.Read(data1)
		Expect(err).To(BeNil())
		data2 := make([]byte, 123)
		_, err = rand.Read(data2)
		Expect(err).To(BeNil())

		paddedData1 := codec.ConvertByPaddingEmptyByte(data1)
		paddedData2 := codec.ConvertByPaddingEmptyByte(data2)

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1})
		Expect(err).To(BeNil())
		Expect(key1).To(Not(BeNil()))
		Expect(blobStatus1).To(Not(BeNil()))
		Expect(*blobStatus1).To(Equal(dispv2.Queued))

		blobStatus2, key2, err := disp.DisperseBlob(ctx, paddedData2, 0, []uint8{0, 1})
		Expect(err).To(BeNil())
		Expect(key2).To(Not(BeNil()))
		Expect(blobStatus2).To(Not(BeNil()))
		Expect(*blobStatus2).To(Equal(dispv2.Queued))

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		var reply1 *disperserpb.BlobStatusReply
		var reply2 *disperserpb.BlobStatusReply
		for loop := true; loop; {
			select {
			case <-ctx.Done():
				Fail("timed out")
			case <-ticker.C:
				reply1, err = disp.GetBlobStatus(context.Background(), key1)
				Expect(err).To(BeNil())
				Expect(reply1).To(Not(BeNil()))
				status1, err := dispv2.BlobStatusFromProtobuf(reply1.GetStatus())
				Expect(err).To(BeNil())

				reply2, err = disp.GetBlobStatus(context.Background(), key2)
				Expect(err).To(BeNil())
				Expect(reply2).To(Not(BeNil()))
				status2, err := dispv2.BlobStatusFromProtobuf(reply2.GetStatus())
				Expect(err).To(BeNil())

				if status1 != dispv2.Complete || status2 != dispv2.Complete {
					continue
				}
				loop = false
			}
		}

		// necessary to ensure that reference block number < current block number
		mineAnvilBlocks(1)

		// test onchain verification using disperser blob status reply #1
		eigenDACert, err := certBuilder.BuildCert(ctx, coretypes.VersionThreeCert, reply1)
		Expect(err).To(BeNil())

		err = staticCertVerifier.CheckDACert(ctx, eigenDACert)
		Expect(err).To(BeNil())

		err = routerCertVerifier.CheckDACert(ctx, eigenDACert)
		Expect(err).To(BeNil())

		// test onchain verification using disperser blob status reply #2
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


		// test retrieval from disperser relay

		// TODO: fix this test to use the payload retrieval client instead of the relay retrieval client
		// payload1, err := relayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert1)
		// Expect(err).To(BeNil())
		// Expect(payload1).To(Not(BeNil()))
		// restored := bytes.TrimRight(payload1.Serialize(), "\x00")
		// Expect(restored).To(Equal(payload1))

		// payload2, err := relayRetrievalClientV2.GetPayload(ctx, eigenDAV3Cert2)
		// Expect(err).To(BeNil())
		// Expect(payload2).To(Not(BeNil()))
		// restored = bytes.TrimRight(payload2.Serialize(), "\x00")
		// Expect(restored).To(Equal(payload2))


		// test retrieval from DA network

		blob1HeaderWithoutPayment, err := eigenDAV3Cert1.BlobHeader()
		Expect(err).To(BeNil())

		blob2HeaderWithoutPayment, err := eigenDAV3Cert2.BlobHeader()
		Expect(err).To(BeNil())

		b, err := validatorRetrievalClientV2.GetBlob(
			ctx,
			blob1HeaderWithoutPayment,
			uint64(eigenDAV3Cert1.ReferenceBlockNumber()),
		)
		Expect(err).To(BeNil())
		restored := bytes.TrimRight(b, "\x00")
		Expect(restored).To(Equal(paddedData1))
		b, err = validatorRetrievalClientV2.GetBlob(
			ctx,
			blob2HeaderWithoutPayment,
			uint64(eigenDAV3Cert2.ReferenceBlockNumber()),
		)
		restored = bytes.TrimRight(b, "\x00")
		Expect(err).To(BeNil())
		Expect(restored).To(Equal(paddedData2))

		// TODO: Figure out how to advance the disperser's reference block number 
		//       currently the disperser isn't respecting the latest anvil chain head when determining RBN
		latestBlock, err := ethClient.BlockNumber(ctx)
		Expect(err).To(BeNil())


		println("latest block number: ", latestBlock)
		tx, err := eigenDACertVerifierRouter.AddCertVerifier(deployerTransactorOpts, uint32(latestBlock) + 2, gethcommon.HexToAddress("0x0"))
		Expect(err).To(BeNil())

		mineAnvilBlocks(1000)
		receipt, err := ethClient.TransactionReceipt(ctx, tx.Hash())
		Expect(err).To(BeNil())
		Expect(receipt).To(Not(BeNil()))
		time.Sleep(30 * time.Second)

		var reply3 *disperserpb.BlobStatusReply
		for loop := true; loop; {
			select {
			case <-ctx.Done():
				Fail("timed out")
			case <-ticker.C:
				reply3, err = disp.GetBlobStatus(context.Background(), key1)
				Expect(err).To(BeNil())
				Expect(reply3).To(Not(BeNil()))
				status, err := dispv2.BlobStatusFromProtobuf(reply3.GetStatus())
				Expect(err).To(BeNil())
				if status != dispv2.Complete {
					continue
				}
				loop = false
			}
		}

		latestBlock, err = ethClient.BlockNumber(ctx)
		Expect(err).To(BeNil())
		println("latest block number after mining: ", latestBlock)

		eigenDACert3, err := certBuilder.BuildCert(ctx, coretypes.VersionThreeCert, reply1)
		eigenDAV3Cert3, ok := eigenDACert3.(*coretypes.EigenDACertV3)
		Expect(ok).To(BeTrue())

		println("cert 3 reference block #: ", eigenDAV3Cert3.ReferenceBlockNumber())
		Expect(err).To(BeNil())

		err = routerCertVerifier.CheckDACert(ctx, eigenDACert3)
		Expect(err).To(Not(BeNil()))


	})
})
