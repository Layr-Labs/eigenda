package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
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

		// deploy payment dynamodb tables

		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))
		err = disp.PopulateAccountant(ctx)
		Expect(err).To(BeNil())

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
		var batchHeader1 *commonpb.BatchHeader
		var batchHeader2 *commonpb.BatchHeader
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

				batchHeader1 = reply1.GetSignedBatch().GetHeader()
				Expect(batchHeader1).To(Not(BeNil()))
				Expect(batchHeader1.GetBatchRoot()).To(Not(BeNil()))
				Expect(batchHeader1.GetReferenceBlockNumber()).To(BeNumerically(">", 0))
				attestation := reply1.GetSignedBatch().GetAttestation()
				Expect(attestation).To(Not(BeNil()))
				Expect(attestation.QuorumNumbers).To(Equal([]uint32{0, 1}))
				Expect(len(attestation.NonSignerPubkeys)).To(Equal(0))
				Expect(attestation.ApkG2).To(Not(BeNil()))
				Expect(len(attestation.QuorumApks)).To(Equal(2))
				Expect(attestation.QuorumSignedPercentages).To(Equal([]byte{100, 100}))
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
				verified, err := merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader1.BatchRoot}, keccak256.New())
				Expect(err).To(BeNil())
				Expect(verified).To(BeTrue())

				batchHeader2 = reply2.GetSignedBatch().GetHeader()
				Expect(batchHeader2).To(Not(BeNil()))
				Expect(batchHeader2.GetBatchRoot()).To(Not(BeNil()))
				Expect(batchHeader2.GetReferenceBlockNumber()).To(BeNumerically(">", 0))
				attestation = reply2.GetSignedBatch().GetAttestation()
				Expect(attestation).To(Not(BeNil()))
				Expect(attestation.QuorumNumbers).To(Equal([]uint32{0, 1}))
				Expect(len(attestation.NonSignerPubkeys)).To(Equal(0))
				Expect(attestation.ApkG2).To(Not(BeNil()))
				Expect(len(attestation.QuorumApks)).To(Equal(2))
				Expect(attestation.QuorumSignedPercentages).To(Equal([]byte{100, 100}))
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
				verified, err = merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader2.BatchRoot}, keccak256.New())
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

		b, err := retrievalClientV2.GetBlob(ctx, blobCert1.BlobHeader, batchHeader1.ReferenceBlockNumber, 0)
		Expect(err).To(BeNil())
		restored := bytes.TrimRight(b, "\x00")
		Expect(restored).To(Equal(paddedData1))
		b, err = retrievalClientV2.GetBlob(ctx, blobCert1.BlobHeader, batchHeader1.ReferenceBlockNumber, 1)
		restored = bytes.TrimRight(b, "\x00")
		Expect(err).To(BeNil())
		Expect(restored).To(Equal(paddedData1))
		b, err = retrievalClientV2.GetBlob(ctx, blobCert2.BlobHeader, batchHeader2.ReferenceBlockNumber, 0)
		restored = bytes.TrimRight(b, "\x00")
		Expect(err).To(BeNil())
		Expect(restored).To(Equal(paddedData2))
		b, err = retrievalClientV2.GetBlob(ctx, blobCert2.BlobHeader, batchHeader2.ReferenceBlockNumber, 1)
		restored = bytes.TrimRight(b, "\x00")
		Expect(err).NotTo(BeNil())
		Expect(restored).To(BeNil())
	})
})
