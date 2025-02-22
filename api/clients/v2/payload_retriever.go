package clients

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
)

// PayloadRetriever represents something that knows how to retrieve a payload from some backend using a verification.EigenDACert
//
// This interface may be implemented to provide alternate retrieval methods, for example payload retrieval from an S3
// bucket instead of from EigenDA relays or nodes.
type PayloadRetriever interface {
	// GetPayload retrieves a payload from some backend, using the provided certificate
	GetPayload(ctx context.Context, eigenDACert *verification.EigenDACert) (*codecs.Payload, error)
}
