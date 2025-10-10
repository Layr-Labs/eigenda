package hashing

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	"golang.org/x/crypto/sha3"
)

// ControllerRefundPaymentRequestDomain is the domain for hashing RefundPaymentRequest messages.
const ControllerRefundPaymentRequestDomain = "controller.RefundPaymentRequest"

func HashRefundPaymentRequest(request *controller.RefundPaymentRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	hasher.Write([]byte(ControllerRefundPaymentRequestDomain))

	err := hashBlobHeader(hasher, request.GetBlobHeader())
	if err != nil {
		return nil, fmt.Errorf("hash blob header: %w", err)
	}

	return hasher.Sum(nil), nil
}
