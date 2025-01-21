package contractEigenDABlobVerifier

import (
	"fmt"
	"math"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserv2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

func SignedBatchProtoToBinding(inputBatch *disperserv2.SignedBatch) (*SignedBatch, error) {
	convertedBatchHeader, err := BatchHeaderProtoToBinding(inputBatch.GetHeader())
	if err != nil {
		return nil, fmt.Errorf("convert batch header: %s", err)
	}

	convertedAttestation, err := attestationProtoToBinding(inputBatch.GetAttestation())
	if err != nil {
		return nil, fmt.Errorf("convert attestation: %s", err)
	}

	outputSignedBatch := &SignedBatch{
		BatchHeader: *convertedBatchHeader,
		Attestation: *convertedAttestation,
	}

	return outputSignedBatch, nil
}

func BatchHeaderProtoToBinding(inputHeader *commonv2.BatchHeader) (*BatchHeaderV2, error) {
	var outputBatchRoot [32]byte

	inputBatchRoot := inputHeader.GetBatchRoot()
	if len(inputBatchRoot) != 32 {
		return nil, fmt.Errorf("BatchRoot must be 32 bytes (length was %d)", len(inputBatchRoot))
	}
	copy(outputBatchRoot[:], inputBatchRoot[:])

	inputReferenceBlockNumber := inputHeader.GetReferenceBlockNumber()
	if inputReferenceBlockNumber > math.MaxUint32 {
		return nil, fmt.Errorf(
			"ReferenceBlockNumber overflow: value was %d, but max allowable value is %d",
			inputReferenceBlockNumber,
			math.MaxUint32)
	}

	convertedHeader := &BatchHeaderV2{
		BatchRoot:            outputBatchRoot,
		ReferenceBlockNumber: uint32(inputReferenceBlockNumber),
	}

	return convertedHeader, nil
}

func attestationProtoToBinding(inputAttestation *disperserv2.Attestation) (*Attestation, error) {
	nonSignerPubkeys, err := repeatedBytesToBN254G1Points(inputAttestation.GetNonSignerPubkeys())
	if err != nil {
		return nil, fmt.Errorf("convert non signer pubkeys to g1 points: %s", err)
	}

	quorumApks, err := repeatedBytesToBN254G1Points(inputAttestation.GetQuorumApks())
	if err != nil {
		return nil, fmt.Errorf("convert quorum apks to g1 points: %s", err)
	}

	sigma, err := bytesToBN254G1Point(inputAttestation.GetSigma())
	if err != nil {
		return nil, fmt.Errorf("convert sigma to g1 point: %s", err)
	}

	apkG2, err := bytesToBN254G2Point(inputAttestation.GetApkG2())
	if err != nil {
		return nil, fmt.Errorf("convert apk g2 to g2 point: %s", err)
	}

	convertedAttestation := &Attestation{
		NonSignerPubkeys: nonSignerPubkeys,
		QuorumApks:       quorumApks,
		Sigma:            *sigma,
		ApkG2:            *apkG2,
		QuorumNumbers:    inputAttestation.GetQuorumNumbers(),
	}

	return convertedAttestation, nil
}

func VerificationProofProtoToBinding(inputVerificationInfo *disperserv2.BlobVerificationInfo) (*BlobVerificationProofV2, error) {
	convertedBlobCertificate, err := blobCertificateProtoToBinding(inputVerificationInfo.GetBlobCertificate())

	if err != nil {
		return nil, fmt.Errorf("convert blob certificate: %s", err)
	}

	return &BlobVerificationProofV2{
		BlobCertificate: *convertedBlobCertificate,
		BlobIndex:       inputVerificationInfo.GetBlobIndex(),
		InclusionProof:  inputVerificationInfo.GetInclusionProof(),
	}, nil
}

func blobCertificateProtoToBinding(inputCertificate *commonv2.BlobCertificate) (*BlobCertificate, error) {
	convertedBlobHeader, err := blobHeaderProtoToBinding(inputCertificate.GetBlobHeader())
	if err != nil {
		return nil, fmt.Errorf("convert blob header: %s", err)
	}

	return &BlobCertificate{
		BlobHeader: *convertedBlobHeader,
		RelayKeys:  inputCertificate.GetRelays(),
	}, nil
}

func blobHeaderProtoToBinding(inputHeader *commonv2.BlobHeader) (*BlobHeaderV2, error) {
	inputVersion := inputHeader.GetVersion()
	if inputVersion > math.MaxUint16 {
		return nil, fmt.Errorf(
			"version overflow: value was %d, but max allowable value is %d",
			inputVersion,
			math.MaxUint16)
	}

	var quorumNumbers []byte
	for _, quorumNumber := range inputHeader.GetQuorumNumbers() {
		if quorumNumber > math.MaxUint8 {
			return nil, fmt.Errorf(
				"quorum number overflow: value was %d, but max allowable value is %d",
				quorumNumber,
				uint8(math.MaxUint8))
		}

		quorumNumbers = append(quorumNumbers, byte(quorumNumber))
	}

	convertedBlobCommitment, err := blobCommitmentProtoToBinding(inputHeader.GetCommitment())
	if err != nil {
		return nil, fmt.Errorf("convert blob commitment: %s", err)
	}

	paymentHeaderHash, err := core.ConvertToPaymentMetadata(inputHeader.GetPaymentHeader()).Hash()
	if err != nil {
		return nil, fmt.Errorf("hash payment header: %s", err)
	}

	return &BlobHeaderV2{
		Version:           uint16(inputVersion),
		QuorumNumbers:     quorumNumbers,
		Commitment:        *convertedBlobCommitment,
		PaymentHeaderHash: paymentHeaderHash,
	}, nil
}

func blobCommitmentProtoToBinding(inputCommitment *common.BlobCommitment) (*BlobCommitment, error) {
	convertedCommitment, err := bytesToBN254G1Point(inputCommitment.GetCommitment())
	if err != nil {
		return nil, fmt.Errorf("convert commitment to g1 point: %s", err)
	}

	convertedLengthCommitment, err := bytesToBN254G2Point(inputCommitment.GetLengthCommitment())
	if err != nil {
		return nil, fmt.Errorf("convert length commitment to g2 point: %s", err)
	}

	convertedLengthProof, err := bytesToBN254G2Point(inputCommitment.GetLengthProof())
	if err != nil {
		return nil, fmt.Errorf("convert length proof to g2 point: %s", err)
	}

	return &BlobCommitment{
		Commitment:       *convertedCommitment,
		LengthCommitment: *convertedLengthCommitment,
		LengthProof:      *convertedLengthProof,
		DataLength:       inputCommitment.GetLength(),
	}, nil
}

// BlobCommitmentBindingToProto converts a BlobCommitment binding into a common.BlobCommitment protobuf
func BlobCommitmentBindingToProto(inputCommitment *BlobCommitment) *common.BlobCommitment {
	return &common.BlobCommitment{
		Commitment:       bn254G1PointToBytes(&inputCommitment.Commitment),
		LengthCommitment: bn254G2PointToBytes(&inputCommitment.LengthCommitment),
		LengthProof:      bn254G2PointToBytes(&inputCommitment.LengthProof),
		Length:           inputCommitment.DataLength,
	}
}

func bytesToBN254G1Point(bytes []byte) (*BN254G1Point, error) {
	var g1Point bn254.G1Affine
	_, err := g1Point.SetBytes(bytes)

	if err != nil {
		return nil, fmt.Errorf("deserialize g1 point: %s", err)
	}

	return &BN254G1Point{
		X: g1Point.X.BigInt(new(big.Int)),
		Y: g1Point.Y.BigInt(new(big.Int)),
	}, nil
}

func bn254G1PointToBytes(inputPoint *BN254G1Point) []byte {
	var x fp.Element
	x.SetBigInt(inputPoint.X)
	var y fp.Element
	y.SetBigInt(inputPoint.Y)

	g1Point := &bn254.G1Affine{X: x, Y: y}

	bytes := g1Point.Bytes()
	return bytes[:]
}

func bytesToBN254G2Point(bytes []byte) (*BN254G2Point, error) {
	var g2Point bn254.G2Affine

	// SetBytes checks that the result is in the correct subgroup
	_, err := g2Point.SetBytes(bytes)

	if err != nil {
		return nil, fmt.Errorf("deserialize g2 point: %s", err)
	}

	var x, y [2]*big.Int
	// Order is intentionally reversed when constructing BN254G2Point
	// (see https://github.com/Layr-Labs/eigenlayer-middleware/blob/512ce7326f35e8060b9d46e23f9c159c0000b546/src/libraries/BN254.sol#L43)
	x[0] = g2Point.X.A1.BigInt(new(big.Int))
	x[1] = g2Point.X.A0.BigInt(new(big.Int))

	y[0] = g2Point.Y.A1.BigInt(new(big.Int))
	y[1] = g2Point.Y.A0.BigInt(new(big.Int))

	return &BN254G2Point{
		X: x,
		Y: y,
	}, nil
}

func bn254G2PointToBytes(inputPoint *BN254G2Point) []byte {
	var g2Point bn254.G2Affine
	
	// Order is intentionally reversed when converting here
	// (see https://github.com/Layr-Labs/eigenlayer-middleware/blob/512ce7326f35e8060b9d46e23f9c159c0000b546/src/libraries/BN254.sol#L43)

	var xa0, xa1, ya0, ya1 fp.Element
	g2Point.X.A0 = *(xa0.SetBigInt(inputPoint.X[1]))
	g2Point.X.A1 = *(xa1.SetBigInt(inputPoint.X[0]))

	g2Point.Y.A0 = *(ya0.SetBigInt(inputPoint.Y[1]))
	g2Point.Y.A1 = *(ya1.SetBigInt(inputPoint.Y[0]))

	pointBytes := g2Point.Bytes()
	return pointBytes[:]
}

func repeatedBytesToBN254G1Points(repeatedBytes [][]byte) ([]BN254G1Point, error) {
	var outputPoints []BN254G1Point
	for _, bytes := range repeatedBytes {
		g1Point, err := bytesToBN254G1Point(bytes)
		if err != nil {
			return nil, fmt.Errorf("deserialize g1 point: %s", err)
		}

		outputPoints = append(outputPoints, *g1Point)
	}

	return outputPoints, nil
}
