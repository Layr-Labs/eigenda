package payloaddispersal

import (
	"context"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// thresholdNotMetError represents an error when signature thresholds are not met
type thresholdNotMetError struct {
	BlobKey               string
	ConfirmationThreshold uint8
	// these are the quorum numbers defined in the blob header
	BlobQuorumNumbers []uint32
	// map from quorumID to percent signed from the quorum
	SignedPercentagesMap map[uint32]uint8
}

// Error implements the error interface and returns a formatted error message
func (e *thresholdNotMetError) Error() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(fmt.Sprintf(
		"Blob Key: %s, Confirmation Threshold: %d%% [", e.BlobKey, e.ConfirmationThreshold))

	for index, quorumID := range e.BlobQuorumNumbers {
		signedPercentage := e.SignedPercentagesMap[quorumID]

		stringBuilder.WriteString(fmt.Sprintf("quorum_%d: %d%%", quorumID, signedPercentage))

		if signedPercentage < e.ConfirmationThreshold {
			stringBuilder.WriteString(" (DOES NOT MEET THRESHOLD)")
		}

		if index < len(e.BlobQuorumNumbers)-1 {
			stringBuilder.WriteString(", ")
		}
	}
	stringBuilder.WriteString("]")

	return stringBuilder.String()
}

// checkThresholds verifies if all quorums meet the confirmation threshold and returns a structured error if they don't
func checkThresholds(
	ctx context.Context,
	certVerifier *verification.CertVerifier,
	blobStatusReply *dispgrpc.BlobStatusReply,
	blobKey string,
) error {
	blobQuorumNumbers := blobStatusReply.GetBlobInclusionInfo().GetBlobCertificate().GetBlobHeader().GetQuorumNumbers()
	if len(blobQuorumNumbers) == 0 {
		return fmt.Errorf("expected >0 quorum numbers in blob header: %v", protoToString(blobStatusReply))
	}

	attestation := blobStatusReply.GetSignedBatch().GetAttestation()
	batchQuorumNumbers := attestation.GetQuorumNumbers()
	batchSignedPercentages := attestation.GetQuorumSignedPercentages()

	if len(batchQuorumNumbers) != len(batchSignedPercentages) {
		return fmt.Errorf("batch quorum number count and signed percentage count don't match")
	}

	// map from quorum ID to the percentage stake signed from that quorum
	signedPercentagesMap := make(map[uint32]uint8, len(batchQuorumNumbers))
	for index, quorumID := range batchQuorumNumbers {
		signedPercentagesMap[quorumID] = batchSignedPercentages[index]
	}

	batchHeader := blobStatusReply.GetSignedBatch().GetHeader()
	if batchHeader == nil {
		return fmt.Errorf("expected non-nil batch header: %v", protoToString(blobStatusReply))
	}

	confirmationThreshold, err := certVerifier.GetConfirmationThreshold(ctx, batchHeader.GetReferenceBlockNumber())
	if err != nil {
		return fmt.Errorf("get confirmation threshold: %w", err)
	}

	// Check if all thresholds are met for the quorums defined in the blob header
	for _, quorum := range blobQuorumNumbers {
		signedPercentage := signedPercentagesMap[quorum]
		if signedPercentage < confirmationThreshold {
			return &thresholdNotMetError{
				BlobKey:               blobKey,
				ConfirmationThreshold: confirmationThreshold,
				BlobQuorumNumbers:     blobQuorumNumbers,
				SignedPercentagesMap:  signedPercentagesMap,
			}
		}
	}

	return nil
}

func protoToString(protoMessage proto.Message) string {
	return prototext.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}.Format(protoMessage)
}
