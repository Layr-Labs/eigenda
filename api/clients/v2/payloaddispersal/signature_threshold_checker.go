package payloaddispersal

import (
	"context"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
)

// signatureThresholdChecker is a utility that interprets the contents of a BlobStatusReply to determine whether
// all quorums have met the confirmation threshold.
type signatureThresholdChecker struct {
	confirmationThreshold   uint8
	requiredQuorums         []uint32
	quorumSignedPercentages []uint8
}

// newSignatureThresholdChecker constructs the utility struct
//
// This method may make a view-only eth contract call to fetch the confirmation threshold, if a new EigenDACertVerifier
// contract has become active.
func newSignatureThresholdChecker(
	ctx context.Context,
	certVerifier clients.ICertVerifier,
	blobStatusReply *dispgrpc.BlobStatusReply,
) (*signatureThresholdChecker, error) {
	if blobStatusReply == nil {
		return nil, fmt.Errorf("blobStatusReply is nil")
	}
	blobInclusionInfo := blobStatusReply.GetBlobInclusionInfo()
	if blobInclusionInfo == nil {
		return nil, fmt.Errorf("blobInclusionInfo is nil")
	}
	blobCertificate := blobInclusionInfo.GetBlobCertificate()
	if blobCertificate == nil {
		return nil, fmt.Errorf("blobCertificate is nil")
	}
	blobHeader := blobCertificate.GetBlobHeader()
	if blobHeader == nil {
		return nil, fmt.Errorf("blobHeader is nil")
	}
	quorumNumbers := blobHeader.GetQuorumNumbers()
	if quorumNumbers == nil {
		return nil, fmt.Errorf("quorumNumbers is nil")
	}
	if len(quorumNumbers) == 0 {
		return nil, fmt.Errorf("quorumNumbers is empty")
	}
	signedBatch := blobStatusReply.GetSignedBatch()
	if signedBatch == nil {
		return nil, fmt.Errorf("signedBatch is nil")
	}
	batchHeader := signedBatch.GetHeader()
	if batchHeader == nil {
		return nil, fmt.Errorf("batchHeader is nil")
	}
	referenceBlockNumber := batchHeader.GetReferenceBlockNumber()
	attestation := signedBatch.GetAttestation()
	if attestation == nil {
		return nil, fmt.Errorf("attestation is nil")
	}
	quorumSignedPercentages := attestation.GetQuorumSignedPercentages()
	if quorumSignedPercentages == nil {
		return nil, fmt.Errorf("quorumSignedPercentages is nil")
	}

	if len(quorumSignedPercentages) != len(quorumNumbers) {
		return nil, fmt.Errorf("expected number of quorum signed percentages to match number of quorums."+
			"quorum signed percentages count: %d. quorum number count: %d",
			len(quorumSignedPercentages), len(quorumNumbers))
	}

	confirmationThreshold, err := certVerifier.GetConfirmationThreshold(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get confirmation threshold: %w", err)
	}

	return &signatureThresholdChecker{
		confirmationThreshold:   confirmationThreshold,
		requiredQuorums:         quorumNumbers,
		quorumSignedPercentages: quorumSignedPercentages,
	}, nil
}

// signatureThresholdsMet returns true if each signedPercentage meets or exceeds the required confirmation threshold
func (stc *signatureThresholdChecker) signatureThresholdsMet() bool {
	for _, signedPercentage := range stc.quorumSignedPercentages {
		if signedPercentage < stc.confirmationThreshold {
			return false
		}
	}

	return true
}

// describeFailureToMeetThresholds returns a string describing a failure to meet all confirmation thresholds.
// This method should be called after it's already known that thresholds haven't been met, to produce a helpful
// error output.
func (stc *signatureThresholdChecker) describeFailureToMeetThresholds(blobKey string) string {
	stringBuilder := strings.Builder{}

	stringBuilder.WriteString("\nBlob Key: ")
	stringBuilder.WriteString(blobKey)
	stringBuilder.WriteString(fmt.Sprintf("\nConfirmation Threshold: %d%%", stc.confirmationThreshold))

	for index, quorum := range stc.requiredQuorums {
		signedPercentage := stc.quorumSignedPercentages[index]

		stringBuilder.WriteString(fmt.Sprintf("\n  Quorum %d: %d%% signed", quorum, signedPercentage))

		if signedPercentage < stc.confirmationThreshold {
			stringBuilder.WriteString(" (DOES NOT MEET THRESHOLD)")
		}
	}

	return stringBuilder.String()
}
