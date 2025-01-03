package verification

import (
	"context"
	"fmt"
	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
)

// BlobVerifier is responsible for making eth calls to verify blobs that have been received by the client
//
// Blob verification is not threadsafe.
type BlobVerifier struct {
	// the eth client that calls will be made to
	ethClient *ethclient.Client
	//
	blobVerifierCaller *verifierBindings.ContractEigenDABlobVerifierCaller
}

func NewBlobVerifier(
	ethRpcUrl string,
	verifyBlobV2FromSignedBatchAddress string) (*BlobVerifier, error) {

	client, err := ethclient.Dial(ethRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ETH RPC node: %s", err)
	}

	verifierCaller, err := verifierBindings.NewContractEigenDABlobVerifierCaller(ethcommon.HexToAddress(verifyBlobV2FromSignedBatchAddress), client)

	if err != nil {
		client.Close()

		return nil, fmt.Errorf("bind to verifier contract at %s: %s", verifyBlobV2FromSignedBatchAddress, err)
	}

	return &BlobVerifier{
		ethClient:          client,
		blobVerifierCaller: verifierCaller,
	}, nil
}

func (v *BlobVerifier) VerifyBlobV2FromSignedBatch(
	ctx context.Context,
	signedBatch *disperser.SignedBatch,
	blobVerificationProof *disperser.BlobVerificationInfo,
) (bool, error) {

	convertedSignedBatch, err := convertSignedBatch(signedBatch)
	if err != nil {
		return false, fmt.Errorf("convert signed batch: %s", err)
	}

	convertedBlobVerificationProof, err := convertBlobVerificationProof(blobVerificationProof)
	if err != nil {
		return false, fmt.Errorf("convert blob verification proof: %s", err)
	}

	err = v.blobVerifierCaller.VerifyBlobV2FromSignedBatch(
		&bind.CallOpts{Context: ctx},
		*convertedSignedBatch,
		*convertedBlobVerificationProof)

	if err != nil {
		return false, fmt.Errorf("verify blob v2 from signed batch: %s", err)
	}

	return true, nil
}

// Close closes the eth client. This method is threadsafe.
func (v *BlobVerifier) Close() {
	v.ethClient.Close()
}

func convertSignedBatch(inputSignedBatch *disperser.SignedBatch) (*verifierBindings.SignedBatch, error) {

	convertedBatchHeader, err := convertBatchHeader(inputSignedBatch.GetHeader())
	if err != nil {
		return nil, fmt.Errorf("convert batch header: %s", err)
	}

	convertedAttestation, err := convertAttestation(inputSignedBatch.Attestation)
	if err != nil {
		return nil, fmt.Errorf("convert attestation: %s", err)
	}

	outputSignedBatch := &verifierBindings.SignedBatch{
		BatchHeader: *convertedBatchHeader,
		Attestation: *convertedAttestation,
	}

	return outputSignedBatch, nil
}

func convertBatchHeader(inputBatchHeader *commonv2.BatchHeader) (*verifierBindings.BatchHeaderV2, error) {
	var outputBatchRoot [32]byte

	inputBatchRoot := inputBatchHeader.GetBatchRoot()
	if len(inputBatchRoot) != 32 {
		return nil, fmt.Errorf("BatchRoot must be 32 bytes (length was %d)", len(inputBatchRoot))
	}
	copy(outputBatchRoot[:], inputBatchRoot[:])

	inputReferenceBlockNumber := inputBatchHeader.GetReferenceBlockNumber()
	if inputReferenceBlockNumber > math.MaxUint32 {
		return nil, fmt.Errorf(
			"ReferenceBlockNumber overflow: value was %d, but max allowable value is %d",
			inputReferenceBlockNumber,
			math.MaxUint32)
	}

	convertedHeader := &verifierBindings.BatchHeaderV2{
		BatchRoot:            outputBatchRoot,
		ReferenceBlockNumber: uint32(inputReferenceBlockNumber),
	}

	return convertedHeader, nil
}

func convertAttestation(inputAttestation *disperser.Attestation) (*verifierBindings.Attestation, error) {
	nonSignerPubkeys, err := repeatedBytesToG1Points(inputAttestation.NonSignerPubkeys)
	if err != nil {
		return nil, fmt.Errorf("convert non signer pubkeys to g1 points: %s", err)
	}

	quorumApks, err := repeatedBytesToG1Points(inputAttestation.QuorumApks)
	if err != nil {
		return nil, fmt.Errorf("convert quorum apks to g1 points: %s", err)
	}

	sigma, err := bytesToBN254G1Point(inputAttestation.Sigma)
	if err != nil {
		return nil, fmt.Errorf("convert sigma to g1 point: %s", err)
	}

	apkG2, err := bytesToBN254G2Point(inputAttestation.ApkG2)
	if err != nil {
		return nil, fmt.Errorf("convert apk g2 to g2 point: %s", err)
	}

	convertedAttestation := &verifierBindings.Attestation{
		NonSignerPubkeys: nonSignerPubkeys,
		QuorumApks:       quorumApks,
		Sigma:            *sigma,
		ApkG2:            *apkG2,
		QuorumNumbers:    inputAttestation.QuorumNumbers,
	}

	return convertedAttestation, nil
}

func convertBlobVerificationProof(inputBlobVerificationProof *disperser.BlobVerificationInfo) (*verifierBindings.BlobVerificationProofV2, error) {
	convertedBlobCertificate, err := convertBlobCertificate(inputBlobVerificationProof.BlobCertificate)
	if err != nil {
		return nil, fmt.Errorf("convert blob certificate: %s", err)
	}

	return &verifierBindings.BlobVerificationProofV2{
		BlobCertificate: *convertedBlobCertificate,
		BlobIndex:       inputBlobVerificationProof.BlobIndex,
		InclusionProof:  inputBlobVerificationProof.InclusionProof,
	}, nil
}

func convertBlobCertificate(inputBlobCertificate *commonv2.BlobCertificate) (*verifierBindings.BlobCertificate, error) {
	convertedBlobHeader, err := convertBlobHeader(inputBlobCertificate.GetBlobHeader())
	if err != nil {
		return nil, fmt.Errorf("convert blob header: %s", err)
	}

	return &verifierBindings.BlobCertificate{
		BlobHeader: *convertedBlobHeader,
		RelayKeys:  inputBlobCertificate.GetRelays(),
	}, nil
}

func convertBlobHeader(inputBlobHeader *commonv2.BlobHeader) (*verifierBindings.BlobHeaderV2, error) {
	inputVersion := inputBlobHeader.Version
	if inputVersion > math.MaxUint16 {
		return nil, fmt.Errorf(
			"version overflow: value was %d, but max allowable value is %d",
			inputVersion,
			math.MaxUint16)
	}

	var quorumNumbers []byte
	for _, quorumNumber := range inputBlobHeader.QuorumNumbers {
		if quorumNumber > math.MaxUint8 {
			return nil, fmt.Errorf("quorum number overflow: value was %d, but max allowable value is %d", quorumNumber, uint8(math.MaxUint8))
		}

		quorumNumbers = append(quorumNumbers, byte(quorumNumber))
	}

	convertedBlobCommitment, err := convertBlobCommitment(inputBlobHeader.Commitment)
	if err != nil {
		return nil, fmt.Errorf("convert blob commitment: %s", err)
	}

	paymentHeaderHash, err := core.ConvertToPaymentMetadata(inputBlobHeader.PaymentHeader).Hash()
	if err != nil {
		return nil, fmt.Errorf("hash payment header: %s", err)
	}

	return &verifierBindings.BlobHeaderV2{
		Version:           uint16(inputVersion),
		QuorumNumbers:     quorumNumbers,
		Commitment:        *convertedBlobCommitment,
		PaymentHeaderHash: paymentHeaderHash,
	}, nil
}

func convertBlobCommitment(inputCommitment *commonv1.BlobCommitment) (*verifierBindings.BlobCommitment, error) {
	convertedCommitment, err := bytesToBN254G1Point(inputCommitment.Commitment)
	if err != nil {
		return nil, fmt.Errorf("convert commitment to g1 point: %s", err)
	}

	convertedLengthCommitment, err := bytesToBN254G2Point(inputCommitment.LengthCommitment)
	if err != nil {
		return nil, fmt.Errorf("convert length commitment to g2 point: %s", err)
	}

	convertedLengthProof, err := bytesToBN254G2Point(inputCommitment.LengthProof)
	if err != nil {
		return nil, fmt.Errorf("convert length proof to g2 point: %s", err)
	}

	return &verifierBindings.BlobCommitment{
		Commitment:       *convertedCommitment,
		LengthCommitment: *convertedLengthCommitment,
		LengthProof:      *convertedLengthProof,
		DataLength:       inputCommitment.Length,
	}, nil
}

func repeatedBytesToG1Points(repeatedBytes [][]byte) ([]verifierBindings.BN254G1Point, error) {
	var outputPoints []verifierBindings.BN254G1Point
	for _, bytes := range repeatedBytes {
		g1Point, err := bytesToBN254G1Point(bytes)
		if err != nil {
			return nil, fmt.Errorf("deserialize g1 point: %s", err)
		}

		outputPoints = append(outputPoints, *g1Point)
	}

	return outputPoints, nil
}

func bytesToBN254G1Point(bytes []byte) (*verifierBindings.BN254G1Point, error) {
	g1Point, err := new(core.G1Point).Deserialize(bytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize g1 point: %s", err)
	}

	return &verifierBindings.BN254G1Point{
		X: g1Point.X.BigInt(new(big.Int)),
		Y: g1Point.Y.BigInt(new(big.Int)),
	}, nil
}

func bytesToBN254G2Point(bytes []byte) (*verifierBindings.BN254G2Point, error) {
	g2Point, err := new(core.G2Point).Deserialize(bytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize g2 point: %s", err)
	}

	var x [2]*big.Int
	x[0] = g2Point.X.A0.BigInt(new(big.Int))
	x[1] = g2Point.X.A1.BigInt(new(big.Int))

	var y [2]*big.Int
	y[0] = g2Point.Y.A0.BigInt(new(big.Int))
	y[1] = g2Point.Y.A1.BigInt(new(big.Int))

	return &verifierBindings.BN254G2Point{
		X: x,
		Y: y,
	}, nil
}
