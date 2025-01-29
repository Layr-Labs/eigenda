package clients

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
)

// PayloadRetriever represents something that knows how to retrieve a payload from some backend using a verification.EigenDACert
type PayloadRetriever interface {
	// GetPayload retrieves a payload from some backend, using the provided certificate
	GetPayload(ctx context.Context, eigenDACert *verification.EigenDACert) ([]byte, error)
}
