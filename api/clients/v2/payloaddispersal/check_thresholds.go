package payloaddispersal

import (
	"context"
	"fmt"
	"strings"

	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
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
	confirmationThreshold byte,
	blobStatusReply *dispgrpc.BlobStatusReply,
	blobKey string,
) error {
	quorumNumbers := blobStatusReply.GetBlobInclusionInfo().GetBlobCertificate().GetBlobHeader().GetQuorumNumbers()
	if len(quorumNumbers) == 0 {
		return fmt.Errorf("expected >0 quorum numbers: %v", protoToString(blobStatusReply))
	}

	quorumSignedPercentages := blobStatusReply.GetSignedBatch().GetAttestation().GetQuorumSignedPercentages()
	if len(quorumSignedPercentages) != len(quorumNumbers) {
		return fmt.Errorf("expected number of signed percentages to match number of quorums. "+
			"signed percentages count: %d; quorum count: %d",
			len(quorumSignedPercentages), len(quorumNumbers))
	}

	batchHeader := blobStatusReply.GetSignedBatch().GetHeader()
	if batchHeader == nil {
		return fmt.Errorf("expected non-nil batch header: %v", protoToString(blobStatusReply))
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

func protoToString(protoMessage proto.Message) string {
	return prototext.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}.Format(protoMessage)
}
