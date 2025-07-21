package clients

import (
	"context"

	_ "github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// PayloadRetriever represents something that knows how to retrieve a payload from some backend using a verification.EigenDACert
//
// This interface may be implemented to provide alternate retrieval methods, for example payload retrieval from an S3
// bucket instead of from EigenDA relays or nodes.
type PayloadRetriever interface {
	// GetPayload retrieves a payload from some backend, using the provided certificate
	// GetPayload should return a [coretypes.ErrBlobDecodingFailedDerivationError] if the blob cannot be decoding according
	// to one of the encodings available via [codecs.PayloadEncodingVersion]s.
	GetPayload(ctx context.Context, eigenDACert coretypes.RetrievableEigenDACert) (*coretypes.Payload, error)
}
