package v2

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
)

// ToBinding converts a SignedBatch into a contractEigenDABlobVerifier.SignedBatch
func (b *SignedBatch) ToBinding() (*verifierBindings.SignedBatch, error) {
	convertedBatchHeader, err := b.GetHeader().ToBinding()
	if err != nil {
		return nil, fmt.Errorf("convert batch header: %s", err)
	}

	convertedAttestation, err := b.GetAttestation().toBinding()
	if err != nil {
		return nil, fmt.Errorf("convert attestation: %s", err)
	}

	outputSignedBatch := &verifierBindings.SignedBatch{
		BatchHeader: *convertedBatchHeader,
		Attestation: *convertedAttestation,
	}

	return outputSignedBatch, nil
}

// toBinding converts an Attestation into a contractEigenDABlobVerifier.Attestation
func (a *Attestation) toBinding() (*verifierBindings.Attestation, error) {
	nonSignerPubkeys, err := repeatedBytesToG1Points(a.GetNonSignerPubkeys())
	if err != nil {
		return nil, fmt.Errorf("convert non signer pubkeys to g1 points: %s", err)
	}

	quorumApks, err := repeatedBytesToG1Points(a.GetQuorumApks())
	if err != nil {
		return nil, fmt.Errorf("convert quorum apks to g1 points: %s", err)
	}

	sigma, err := common.BytesToBN254G1Point(a.GetSigma())
	if err != nil {
		return nil, fmt.Errorf("convert sigma to g1 point: %s", err)
	}

	apkG2, err := common.BytesToBN254G2Point(a.GetApkG2())
	if err != nil {
		return nil, fmt.Errorf("convert apk g2 to g2 point: %s", err)
	}

	convertedAttestation := &verifierBindings.Attestation{
		NonSignerPubkeys: nonSignerPubkeys,
		QuorumApks:       quorumApks,
		Sigma:            *sigma,
		ApkG2:            *apkG2,
		QuorumNumbers:    a.GetQuorumNumbers(),
	}

	return convertedAttestation, nil
}

// ToBinding converts a BlobVerificationInfo into a contractEigenDABlobVerifier.BlobVerificationProofV2
func (i *BlobVerificationInfo) ToBinding() (*verifierBindings.BlobVerificationProofV2, error) {
	convertedBlobCertificate, err := i.GetBlobCertificate().ToBinding()

	if err != nil {
		return nil, fmt.Errorf("convert blob certificate: %s", err)
	}

	return &verifierBindings.BlobVerificationProofV2{
		BlobCertificate: *convertedBlobCertificate,
		BlobIndex:       i.GetBlobIndex(),
		InclusionProof:  i.GetInclusionProof(),
	}, nil
}

// repeatedBytesToG1Points accepts an array of byte arrays, and returns an array of contractEigenDABlobVerifier.BN254G1Point
func repeatedBytesToG1Points(repeatedBytes [][]byte) ([]verifierBindings.BN254G1Point, error) {
	var outputPoints []verifierBindings.BN254G1Point
	for _, bytes := range repeatedBytes {
		g1Point, err := common.BytesToBN254G1Point(bytes)
		if err != nil {
			return nil, fmt.Errorf("deserialize g1 point: %s", err)
		}

		outputPoints = append(outputPoints, *g1Point)
	}

	return outputPoints, nil
}
