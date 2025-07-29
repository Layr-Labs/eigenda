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
//
// TODO(samlaf): I don't think we need this interface. We probably shouldn't have the separate relay
// and validator retrieval clients that implement this interface. Instead,
// we should have a single PayloadRetriever that knows how to retrieve blobs from either
// relays or validators, and then decodes them to (encoded) payloads.
type PayloadRetriever interface {
	// GetPayload retrieves a payload from some backend, using the provided certificate
	// GetPayload should return a [coretypes.ErrBlobDecodingFailedDerivationError] if the blob cannot be decoding according
	// to one of the encodings available via [codecs.PayloadEncodingVersion]s.
	GetPayload(ctx context.Context, eigenDACert coretypes.RetrievableEigenDACert) (*coretypes.Payload, error)

	// GetEncodedPayload retrieves an encoded payload from some backend, using the provided certificate.
	// This method performs the same operations as GetPayload but stops before decoding the payload,
	// returning the encoded form instead.
	GetEncodedPayload(ctx context.Context, eigenDACert coretypes.RetrievableEigenDACert) (*coretypes.EncodedPayload, error)
}
