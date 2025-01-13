package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		Expect(err).To(BeNil())
		Expect(key1).To(Not(BeNil()))
		Expect(blobStatus1).To(Not(BeNil()))
		Expect(*blobStatus1).To(Equal(dispv2.Queued))

		blobStatus2, key2, err := disp.DisperseBlob(ctx, paddedData2, 0, []uint8{0, 1}, 0)
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
		var signedBatch1 *disperserpb.SignedBatch
		var signedBatch2 *disperserpb.SignedBatch
		var blobVerification1 *disperserpb.BlobVerificationInfo
		var blobVerification2 *disperserpb.BlobVerificationInfo
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

				signedBatch1 = reply1.GetSignedBatch()
				batchHeader1 = signedBatch1.GetHeader()
				Expect(batchHeader1).To(Not(BeNil()))
				Expect(batchHeader1.GetBatchRoot()).To(Not(BeNil()))
				Expect(batchHeader1.GetReferenceBlockNumber()).To(BeNumerically(">", 0))
				attestation := reply1.GetSignedBatch().GetAttestation()
				Expect(attestation).To(Not(BeNil()))
				Expect(attestation.QuorumNumbers).To(ConsistOf([]uint32{0, 1}))
				Expect(len(attestation.NonSignerPubkeys)).To(Equal(0))
				Expect(attestation.ApkG2).To(Not(BeNil()))
				Expect(len(attestation.QuorumApks)).To(Equal(2))
				Expect(attestation.QuorumSignedPercentages).To(Equal([]byte{100, 100}))
				blobVerification1 = reply1.GetBlobVerificationInfo()
				Expect(blobVerification1).To(Not(BeNil()))
				Expect(blobVerification1.GetBlobCertificate()).To(Not(BeNil()))
				blobCert1, err = corev2.BlobCertificateFromProtobuf(blobVerification1.GetBlobCertificate())
				Expect(err).To(BeNil())
				inclusionProofBytes := blobVerification1.GetInclusionProof()
				blobIndex := blobVerification1.GetBlobIndex()
				proof, err := core.DeserializeMerkleProof(inclusionProofBytes, uint64(blobIndex))
				Expect(err).To(BeNil())
				certHash, err := blobCert1.Hash()
				Expect(err).To(BeNil())
				_, err = blobCert1.BlobHeader.BlobKey()
				Expect(err).To(BeNil())
				verified, err := merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader1.BatchRoot}, keccak256.New())
				Expect(err).To(BeNil())
				Expect(verified).To(BeTrue())

				signedBatch2 = reply2.GetSignedBatch()
				batchHeader2 = signedBatch2.GetHeader()
				Expect(batchHeader2).To(Not(BeNil()))
				Expect(batchHeader2.GetBatchRoot()).To(Not(BeNil()))
				Expect(batchHeader2.GetReferenceBlockNumber()).To(BeNumerically(">", 0))
				attestation = reply2.GetSignedBatch().GetAttestation()
				Expect(attestation).To(Not(BeNil()))

				attestation2 := reply2.GetSignedBatch().GetAttestation()
				Expect(attestation2).To(Not(BeNil()))
				Expect(attestation2.QuorumNumbers).To(Equal(attestation.QuorumNumbers))
				Expect(len(attestation2.NonSignerPubkeys)).To(Equal(len(attestation.NonSignerPubkeys)))
				Expect(attestation2.ApkG2).To(Equal(attestation.ApkG2))
				Expect(len(attestation2.QuorumApks)).To(Equal(len(attestation.QuorumApks)))
				Expect(attestation2.QuorumSignedPercentages).To(Equal(attestation.QuorumSignedPercentages))

				blobVerification2 = reply2.GetBlobVerificationInfo()
				Expect(blobVerification2).To(Not(BeNil()))
				Expect(blobVerification2.GetBlobCertificate()).To(Not(BeNil()))
				blobCert2, err = corev2.BlobCertificateFromProtobuf(blobVerification2.GetBlobCertificate())
				Expect(err).To(BeNil())
				inclusionProofBytes = blobVerification2.GetInclusionProof()
				blobIndex = blobVerification2.GetBlobIndex()
				proof, err = core.DeserializeMerkleProof(inclusionProofBytes, uint64(blobIndex))
				Expect(err).To(BeNil())
				certHash, err = blobCert2.Hash()
				Expect(err).To(BeNil())
				verified, err = merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader2.BatchRoot}, keccak256.New())
				Expect(err).To(BeNil())
				Expect(verified).To(BeTrue())
				loop = false
			}
		}

		// necessary to ensure that reference block number < current block number
		mineAnvilBlocks(1)

		// test onchain verification
		attestation, err := convertAttestation(signedBatch1.GetAttestation())
		Expect(err).To(BeNil())
		proof, err := convertBlobVerificationInfo(blobVerification1)
		Expect(err).To(BeNil())

		var batchRoot [32]byte
		copy(batchRoot[:], batchHeader1.BatchRoot)

		err = verifierContract.VerifyBlobV2FromSignedBatch(
			&bind.CallOpts{},
			verifierbindings.SignedBatch{
				BatchHeader: verifierbindings.BatchHeaderV2{
					BatchRoot:            batchRoot,
					ReferenceBlockNumber: uint32(batchHeader1.ReferenceBlockNumber),
				},
				Attestation: *attestation,
			},
			*proof,
		)
		Expect(err).To(BeNil())

		attestation, err = convertAttestation(signedBatch2.GetAttestation())
		Expect(err).To(BeNil())
		proof, err = convertBlobVerificationInfo(blobVerification2)
		Expect(err).To(BeNil())
		copy(batchRoot[:], batchHeader2.BatchRoot)
		err = verifierContract.VerifyBlobV2FromSignedBatch(
			&bind.CallOpts{},
			verifierbindings.SignedBatch{
				BatchHeader: verifierbindings.BatchHeaderV2{
					BatchRoot:            batchRoot,
					ReferenceBlockNumber: uint32(batchHeader2.ReferenceBlockNumber),
				},
				Attestation: *attestation,
			},
			*proof,
		)
		Expect(err).To(BeNil())

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

		// Test retrieval from DA network
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
		Expect(err).To(BeNil())
		Expect(restored).To(Equal(paddedData2))
	})
})

func convertBlobVerificationInfo(verificationInfo *disperserpb.BlobVerificationInfo) (*verifierbindings.BlobVerificationProofV2, error) {
	blobCertificate, err := corev2.BlobCertificateFromProtobuf(verificationInfo.GetBlobCertificate())
	if err != nil {
		return nil, err
	}
	paymentHeaderHash, err := blobCertificate.BlobHeader.PaymentMetadata.Hash()
	if err != nil {
		return nil, err
	}

	inclusionProof := verificationInfo.GetInclusionProof()
	blobIndex := verificationInfo.GetBlobIndex()

	commitX := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.Commitment.X.BigInt(commitX)
	commitY := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.Commitment.Y.BigInt(commitY)
	lengthCommitX0 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthCommitment.X.A0.BigInt(lengthCommitX0)
	lengthCommitX1 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthCommitment.X.A1.BigInt(lengthCommitX1)
	lengthCommitY0 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthCommitment.Y.A0.BigInt(lengthCommitY0)
	lengthCommitY1 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthCommitment.Y.A1.BigInt(lengthCommitY1)
	lengthProofX0 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthProof.X.A0.BigInt(lengthProofX0)
	lengthProofX1 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthProof.X.A1.BigInt(lengthProofX1)
	lengthProofY0 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthProof.Y.A0.BigInt(lengthProofY0)
	lengthProofY1 := big.NewInt(0)
	blobCertificate.BlobHeader.BlobCommitments.LengthProof.Y.A1.BigInt(lengthProofY1)
	return &verifierbindings.BlobVerificationProofV2{
		BlobCertificate: verifierbindings.BlobCertificate{
			BlobHeader: verifierbindings.BlobHeaderV2{
				Version:       uint16(blobCertificate.BlobHeader.BlobVersion),
				QuorumNumbers: blobCertificate.BlobHeader.QuorumNumbers,
				Commitment: verifierbindings.BlobCommitment{
					Commitment: verifierbindings.BN254G1Point{
						X: commitX,
						Y: commitY,
					},
					LengthCommitment: verifierbindings.BN254G2Point{
						X: [2]*big.Int{lengthCommitX0, lengthCommitX1},
						Y: [2]*big.Int{lengthCommitY0, lengthCommitY1},
					},
					LengthProof: verifierbindings.BN254G2Point{
						X: [2]*big.Int{lengthProofX0, lengthProofX1},
						Y: [2]*big.Int{lengthProofY0, lengthProofY1},
					},
					DataLength: uint32(blobCertificate.BlobHeader.BlobCommitments.Length),
				},
				PaymentHeaderHash: paymentHeaderHash,
			},
			RelayKeys: blobCertificate.RelayKeys,
		},
		InclusionProof: inclusionProof,
		BlobIndex:      blobIndex,
	}, nil
}

func convertAttestation(attestation *disperserpb.Attestation) (*verifierbindings.Attestation, error) {
	if attestation == nil {
		return nil, fmt.Errorf("attestation is nil")
	}
	nonSignerPubkeys := make([]verifierbindings.BN254G1Point, 0)
	for _, pubkey := range attestation.GetNonSignerPubkeys() {
		pk, err := convertG1Point(pubkey)
		if err != nil {
			return nil, err
		}
		nonSignerPubkeys = append(nonSignerPubkeys, *pk)
	}

	quorumApks := make([]verifierbindings.BN254G1Point, 0)
	for _, apk := range attestation.GetQuorumApks() {
		apk, err := convertG1Point(apk)
		if err != nil {
			return nil, err
		}
		quorumApks = append(quorumApks, *apk)
	}

	if attestation.GetSigma() == nil {
		return nil, fmt.Errorf("attestation sigma is nil")
	}
	sigma, err := convertG1Point(attestation.GetSigma())
	if err != nil {
		return nil, err
	}

	if attestation.GetApkG2() == nil {
		return nil, fmt.Errorf("attestation apkG2 is nil")
	}
	apkg2, err := convertG2Point(attestation.GetApkG2())
	if err != nil {
		return nil, err
	}

	return &verifierbindings.Attestation{
		NonSignerPubkeys: nonSignerPubkeys,
		QuorumApks:       quorumApks,
		Sigma:            *sigma,
		ApkG2:            *apkg2,
		QuorumNumbers:    attestation.GetQuorumNumbers(),
	}, nil
}

func convertG1Point(data []byte) (*verifierbindings.BN254G1Point, error) {
	point, err := new(core.G1Point).Deserialize(data)
	if err != nil {
		return nil, err
	}
	x := big.NewInt(0)
	y := big.NewInt(0)

	point.X.BigInt(x)
	point.Y.BigInt(y)
	return &verifierbindings.BN254G1Point{
		X: x,
		Y: y,
	}, nil
}

func convertG2Point(data []byte) (*verifierbindings.BN254G2Point, error) {
	point, err := new(core.G2Point).Deserialize(data)
	if err != nil {
		return nil, err
	}
	x0 := big.NewInt(0)
	x1 := big.NewInt(0)
	y0 := big.NewInt(0)
	y1 := big.NewInt(0)

	point.X.A0.BigInt(x0)
	point.X.A1.BigInt(x1)
	point.Y.A0.BigInt(y0)
	point.Y.A1.BigInt(y1)
	return &verifierbindings.BN254G2Point{
		X: [2]*big.Int{x1, x0},
		Y: [2]*big.Int{y1, y0},
	}, nil
}
