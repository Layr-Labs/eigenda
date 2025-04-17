package payloaddispersal

import (
	"context"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
)

// thresholdNotMetError represents an error when signature thresholds are not met
type thresholdNotMetError struct {
	BlobKey               string
	ConfirmationThreshold uint8
	QuorumNumbers         []uint32
	SignedPercentages     []uint8
}

// Error implements the error interface and returns a formatted error message
func (e *thresholdNotMetError) Error() string {
	stringBuilder := strings.Builder{}

	stringBuilder.WriteString("\nBlob Key: ")
	stringBuilder.WriteString(e.BlobKey)
	stringBuilder.WriteString(fmt.Sprintf("\nConfirmation Threshold: %d%%", e.ConfirmationThreshold))

	for index, quorum := range e.QuorumNumbers {
		signedPercentage := e.SignedPercentages[index]

		stringBuilder.WriteString(fmt.Sprintf("\n  Quorum %d: %d%% signed", quorum, signedPercentage))

		if signedPercentage < e.ConfirmationThreshold {
			stringBuilder.WriteString(" (DOES NOT MEET THRESHOLD)")
		}
	}

	return stringBuilder.String()
}

// checkThresholds verifies if all quorums meet the confirmation threshold and returns a structured error if they don't
func checkThresholds(
	ctx context.Context,
	certVerifier clients.ICertVerifier,
	blobStatusReply *dispgrpc.BlobStatusReply,
	blobKey string,
) error {
	if blobStatusReply == nil {
		return fmt.Errorf("blobStatusReply is nil")
	}
	blobInclusionInfo := blobStatusReply.GetBlobInclusionInfo()
	if blobInclusionInfo == nil {
		return fmt.Errorf("blobInclusionInfo is nil")
	}
	blobCertificate := blobInclusionInfo.GetBlobCertificate()
	if blobCertificate == nil {
		return fmt.Errorf("blobCertificate is nil")
	}
	blobHeader := blobCertificate.GetBlobHeader()
	if blobHeader == nil {
		return fmt.Errorf("blobHeader is nil")
	}
	quorumNumbers := blobHeader.GetQuorumNumbers()
	if quorumNumbers == nil {
		return fmt.Errorf("quorumNumbers is nil")
	}
	if len(quorumNumbers) == 0 {
		return fmt.Errorf("quorumNumbers is empty")
	}
	signedBatch := blobStatusReply.GetSignedBatch()
	if signedBatch == nil {
		return fmt.Errorf("signedBatch is nil")
	}
	batchHeader := signedBatch.GetHeader()
	if batchHeader == nil {
		return fmt.Errorf("batchHeader is nil")
	}
	referenceBlockNumber := batchHeader.GetReferenceBlockNumber()
	attestation := signedBatch.GetAttestation()
	if attestation == nil {
		return fmt.Errorf("attestation is nil")
	}
	quorumSignedPercentages := attestation.GetQuorumSignedPercentages()
	if quorumSignedPercentages == nil {
		return fmt.Errorf("quorumSignedPercentages is nil")
	}

	if len(quorumSignedPercentages) != len(quorumNumbers) {
		return fmt.Errorf("expected number of quorum signed percentages to match number of quorums."+
			"quorum signed percentages count: %d. quorum number count: %d",
			len(quorumSignedPercentages), len(quorumNumbers))
	}

	confirmationThreshold, err := certVerifier.GetConfirmationThreshold(ctx, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("get confirmation threshold: %w", err)
	}

	// Check if all thresholds are met
	for _, signedPercentage := range quorumSignedPercentages {
		if signedPercentage < confirmationThreshold {
			return &thresholdNotMetError{
				BlobKey:               blobKey,
				ConfirmationThreshold: confirmationThreshold,
				QuorumNumbers:         quorumNumbers,
				SignedPercentages:     quorumSignedPercentages,
			}
		}
	}

	return nil
}
