package hashing

import (
	"fmt"

	controller "github.com/Layr-Labs/eigenda/api/grpc/controller/v1"
	"golang.org/x/crypto/sha3"
)

// ControllerAuthorizePaymentRequestDomain is the domain for hashing AuthorizePaymentRequest messages.
const ControllerAuthorizePaymentRequestDomain = "controller.AuthorizePaymentRequest"

func HashAuthorizePaymentRequest(request *controller.AuthorizePaymentRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	hasher.Write([]byte(ControllerAuthorizePaymentRequestDomain))

	err := hashBlobHeader(hasher, request.GetBlobHeader())
	if err != nil {
		return nil, fmt.Errorf("hash blob header: %w", err)
	}

	hasher.Write(request.GetClientSignature())

	return hasher.Sum(nil), nil
}
