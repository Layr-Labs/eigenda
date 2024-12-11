package integration_test

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
)

var _ = Describe("Inabox v2 Integration", func() {
	It("test end to end scenario", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		disp, err := clients.NewDisperserClientV2(&clients.DisperserClientV2Config{
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

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		Expect(err).To(BeNil())
		Expect(key1).To(Not(BeNil()))
		Expect(blobStatus1).To(Not(BeNil()))
		Expect(*blobStatus1).To(Equal(dispv2.Queued))

		blobStatus2, key2, err := disp.DisperseBlob(ctx, paddedData2, 0, []uint8{0}, 0)
		Expect(err).To(BeNil())
		Expect(key2).To(Not(BeNil()))
		Expect(blobStatus2).To(Not(BeNil()))
		Expect(*blobStatus2).To(Equal(dispv2.Queued))

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		var reply1 *disperserpb.BlobStatusReply
		var reply2 *disperserpb.BlobStatusReply
		var blobCert1 *corev2.BlobCertificate
		var blobCert2 *corev2.BlobCertificate
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

				if status1 != dispv2.Certified || status2 != dispv2.Certified {
					continue
				}

				batchHeader := reply1.GetSignedBatch().GetHeader()
				Expect(batchHeader).To(Not(BeNil()))
				Expect(batchHeader.GetBatchRoot()).To(Not(BeNil()))
				Expect(batchHeader.GetReferenceBlockNumber()).To(BeNumerically(">", 0))
				attestation := reply1.GetSignedBatch().GetAttestation()
				Expect(attestation).To(Not(BeNil()))
				blobVerification := reply1.GetBlobVerificationInfo()
				Expect(blobVerification).To(Not(BeNil()))
				Expect(blobVerification.GetBlobCertificate()).To(Not(BeNil()))
				blobCert1, err = corev2.BlobCertificateFromProtobuf(blobVerification.GetBlobCertificate())
				Expect(err).To(BeNil())
				inclusionProofBytes := blobVerification.GetInclusionProof()
				blobIndex := blobVerification.GetBlobIndex()
				proof, err := core.DeserializeMerkleProof(inclusionProofBytes, uint64(blobIndex))
				Expect(err).To(BeNil())
				certHash, err := blobCert1.Hash()
				Expect(err).To(BeNil())
				_, err = blobCert1.BlobHeader.BlobKey()
				Expect(err).To(BeNil())
				verified, err := merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader.BatchRoot}, keccak256.New())
				Expect(err).To(BeNil())
				Expect(verified).To(BeTrue())

				batchHeader = reply2.GetSignedBatch().GetHeader()
				Expect(batchHeader).To(Not(BeNil()))
				Expect(batchHeader.GetBatchRoot()).To(Not(BeNil()))
				Expect(batchHeader.GetReferenceBlockNumber()).To(BeNumerically(">", 0))
				attestation = reply2.GetSignedBatch().GetAttestation()
				Expect(attestation).To(Not(BeNil()))
				blobVerification = reply2.GetBlobVerificationInfo()
				Expect(blobVerification).To(Not(BeNil()))
				Expect(blobVerification.GetBlobCertificate()).To(Not(BeNil()))
				blobCert2, err = corev2.BlobCertificateFromProtobuf(blobVerification.GetBlobCertificate())
				Expect(err).To(BeNil())
				inclusionProofBytes = blobVerification.GetInclusionProof()
				blobIndex = blobVerification.GetBlobIndex()
				proof, err = core.DeserializeMerkleProof(inclusionProofBytes, uint64(blobIndex))
				Expect(err).To(BeNil())
				certHash, err = blobCert2.Hash()
				Expect(err).To(BeNil())
				verified, err = merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader.BatchRoot}, keccak256.New())
				Expect(err).To(BeNil())
				Expect(verified).To(BeTrue())
				// TODO(ian-shim): verify the blob onchain using a mock rollup contract
				loop = false
			}
		}

		// Test retrieval from relay
		relayClient, err := clients.NewRelayClient(&clients.RelayClientConfig{
			Sockets: relays,
		}, logger)
		Expect(err).To(BeNil())
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		blob1Relays := make(map[corev2.RelayKey]struct{}, 0)
		blob2Relays := make(map[corev2.RelayKey]struct{}, 0)
		for _, k := range blobCert1.RelayKeys {
			blob1Relays[corev2.RelayKey(k)] = struct{}{}
		}
		for _, k := range blobCert2.RelayKeys {
			blob2Relays[corev2.RelayKey(k)] = struct{}{}
		}
		for relayKey := range relays {
			blob1, err := relayClient.GetBlob(ctx, relayKey, key1)
			if _, ok := blob1Relays[corev2.RelayKey(relayKey)]; ok {
				Expect(err).To(BeNil())
				Expect(blob1).To(Equal(paddedData1))
			} else {
				Expect(err).NotTo(BeNil())
			}

			blob2, err := relayClient.GetBlob(ctx, relayKey, key2)
			if _, ok := blob2Relays[corev2.RelayKey(relayKey)]; ok {
				Expect(err).To(BeNil())
				Expect(blob2).To(Equal(paddedData2))
			} else {
				Expect(err).NotTo(BeNil())
			}
		}

		// TODO(ian-shim): test retrieval from DA nodes via retrieval client
	})
})
