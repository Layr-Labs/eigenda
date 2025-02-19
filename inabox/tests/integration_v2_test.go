package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/docker/go-units"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
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
		signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
		Expect(err).To(BeNil())

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
		var blobCert1 *corev2.BlobCertificate
		var blobCert2 *corev2.BlobCertificate
		var batchHeader1 *commonpb.BatchHeader
		var batchHeader2 *commonpb.BatchHeader
		var signedBatch1 *disperserpb.SignedBatch
		var signedBatch2 *disperserpb.SignedBatch
		var blobInclusion1 *disperserpb.BlobInclusionInfo
		var blobInclusion2 *disperserpb.BlobInclusionInfo
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
				blobInclusion1 = reply1.GetBlobInclusionInfo()
				Expect(blobInclusion1).To(Not(BeNil()))
				Expect(blobInclusion1.GetBlobCertificate()).To(Not(BeNil()))
				blobCert1, err = corev2.BlobCertificateFromProtobuf(blobInclusion1.GetBlobCertificate())
				Expect(err).To(BeNil())
				inclusionProofBytes := blobInclusion1.GetInclusionProof()
				blobIndex := blobInclusion1.GetBlobIndex()
				proof, err := core.DeserializeMerkleProof(inclusionProofBytes, uint64(blobIndex))
				Expect(err).To(BeNil())
				certHash, err := blobCert1.Hash()
				Expect(err).To(BeNil())
				_, err = blobCert1.BlobHeader.BlobKey()
				Expect(err).To(BeNil())
				verified, err := merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader1.BatchRoot}, keccak256.New())
				Expect(err).To(BeNil())
				Expect(verified).To(BeTrue())
				Expect(blobCert1.Signature).To(HaveLen(65))
				Expect(len(blobCert1.RelayKeys)).To((BeNumerically(">", 0)))

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

				blobInclusion2 = reply2.GetBlobInclusionInfo()
				Expect(blobInclusion2).To(Not(BeNil()))
				Expect(blobInclusion2.GetBlobCertificate()).To(Not(BeNil()))
				blobCert2, err = corev2.BlobCertificateFromProtobuf(blobInclusion2.GetBlobCertificate())
				Expect(err).To(BeNil())
				inclusionProofBytes = blobInclusion2.GetInclusionProof()
				blobIndex = blobInclusion2.GetBlobIndex()
				proof, err = core.DeserializeMerkleProof(inclusionProofBytes, uint64(blobIndex))
				Expect(err).To(BeNil())
				certHash, err = blobCert2.Hash()
				Expect(err).To(BeNil())
				verified, err = merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{batchHeader2.BatchRoot}, keccak256.New())
				Expect(err).To(BeNil())
				Expect(verified).To(BeTrue())
				Expect(blobCert2.Signature).To(HaveLen(65))
				Expect(len(blobCert2.RelayKeys)).To((BeNumerically(">", 0)))
				loop = false
			}
		}

		// necessary to ensure that reference block number < current block number
		mineAnvilBlocks(1)

		// test onchain verification
		attestation, err := convertAttestation(signedBatch1.GetAttestation())
		Expect(err).To(BeNil())
		proof, err := convertBlobInclusionInfo(blobInclusion1)
		Expect(err).To(BeNil())

		var batchRoot [32]byte
		copy(batchRoot[:], batchHeader1.BatchRoot)

		err = verifierContract.VerifyDACertV2FromSignedBatch(
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
		proof, err = convertBlobInclusionInfo(blobInclusion2)
		Expect(err).To(BeNil())
		copy(batchRoot[:], batchHeader2.BatchRoot)
		err = verifierContract.VerifyDACertV2FromSignedBatch(
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
			Sockets:            relays,
			MaxGRPCMessageSize: units.GiB,
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

		blob1Key, err := blobCert1.BlobHeader.BlobKey()
		Expect(err).To(BeNil())

		blob2Key, err := blobCert2.BlobHeader.BlobKey()
		Expect(err).To(BeNil())

		// Test retrieval from DA network
		b, err := retrievalClientV2.GetBlob(
			ctx,
			blob1Key,
			blobCert1.BlobHeader.BlobVersion,
			blobCert1.BlobHeader.BlobCommitments,
			batchHeader1.ReferenceBlockNumber,
			0)
		Expect(err).To(BeNil())
		restored := bytes.TrimRight(b, "\x00")
		Expect(restored).To(Equal(paddedData1))
		b, err = retrievalClientV2.GetBlob(
			ctx,
			blob1Key,
			blobCert1.BlobHeader.BlobVersion,
			blobCert1.BlobHeader.BlobCommitments,
			batchHeader1.ReferenceBlockNumber,
			1)
		restored = bytes.TrimRight(b, "\x00")
		Expect(err).To(BeNil())
		Expect(restored).To(Equal(paddedData1))
		b, err = retrievalClientV2.GetBlob(
			ctx,
			blob2Key,
			blobCert2.BlobHeader.BlobVersion,
			blobCert2.BlobHeader.BlobCommitments,
			batchHeader2.ReferenceBlockNumber,
			0)
		restored = bytes.TrimRight(b, "\x00")
		Expect(err).To(BeNil())
		Expect(restored).To(Equal(paddedData2))
		b, err = retrievalClientV2.GetBlob(
			ctx,
			blob2Key,
			blobCert2.BlobHeader.BlobVersion,
			blobCert2.BlobHeader.BlobCommitments,
			batchHeader2.ReferenceBlockNumber,
			1)
		restored = bytes.TrimRight(b, "\x00")
		Expect(err).To(BeNil())
		Expect(restored).To(Equal(paddedData2))
	})
})

func convertBlobInclusionInfo(inclusionInfo *disperserpb.BlobInclusionInfo) (*verifierbindings.BlobInclusionInfo, error) {
	blobCertificate, err := corev2.BlobCertificateFromProtobuf(inclusionInfo.GetBlobCertificate())
	if err != nil {
		return nil, err
	}
	paymentHeaderHash, err := blobCertificate.BlobHeader.PaymentMetadata.Hash()
	if err != nil {
		return nil, err
	}

	inclusionProof := inclusionInfo.GetInclusionProof()
	blobIndex := inclusionInfo.GetBlobIndex()

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
	return &verifierbindings.BlobInclusionInfo{
		BlobCertificate: verifierbindings.BlobCertificate{
			BlobHeader: verifierbindings.BlobHeaderV2{
				Version:       uint16(blobCertificate.BlobHeader.BlobVersion),
				QuorumNumbers: blobCertificate.BlobHeader.QuorumNumbers,
				Commitment: verifierbindings.BlobCommitment{
					Commitment: verifierbindings.BN254G1Point{
						X: commitX,
						Y: commitY,
					},
					// Most crypptography library serializes a G2 point by having
					// A0 followed by A1 for both X, Y field of G2. However, ethereum
					// precompile assumes an ordering of A1, A0. We choose
					// to conform with Ethereum order when serializing a blobHeaderV2
					// for instance, gnark, https://github.com/Consensys/gnark-crypto/blob/de0d77f2b4d520350bc54c612828b19ce2146eee/ecc/bn254/marshal.go#L1078
					// Ethereum, https://eips.ethereum.org/EIPS/eip-197#definition-of-the-groups
					LengthCommitment: verifierbindings.BN254G2Point{
						X: [2]*big.Int{lengthCommitX1, lengthCommitX0},
						Y: [2]*big.Int{lengthCommitY1, lengthCommitY0},
					},
					LengthProof: verifierbindings.BN254G2Point{
						X: [2]*big.Int{lengthProofX1, lengthProofX0},
						Y: [2]*big.Int{lengthProofY1, lengthProofY0},
					},
					Length: uint32(blobCertificate.BlobHeader.BlobCommitments.Length),
				},
				PaymentHeaderHash: paymentHeaderHash,
			},
			Signature: blobCertificate.Signature,
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
