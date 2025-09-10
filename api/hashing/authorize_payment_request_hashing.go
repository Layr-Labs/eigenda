package hashing

import (
	"fmt"

	controller "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"golang.org/x/crypto/sha3"
)

// ControllerAuthorizePaymentRequestDomain is the domain for hashing AuthorizePaymentRequest messages.
const ControllerAuthorizePaymentRequestDomain = "controller.AuthorizePaymentRequest"

// HashAuthorizePaymentRequest hashes the given AuthorizePaymentRequest (excluding the disperser_signature field).
func HashAuthorizePaymentRequest(request *controller.AuthorizePaymentRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	hasher.Write([]byte(ControllerAuthorizePaymentRequestDomain))

	err := hashBlobHeader(hasher, request.GetBlobHeader())
	if err != nil {
		return nil, fmt.Errorf("hash blob header: %w", err)
	}

	// We intentionally do not hash the disperser_signature field, otherwise that signature would be self referential

	return hasher.Sum(nil), nil
}
