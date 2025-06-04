package v2

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common/replay"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type authenticator struct {
	ReplayGuardian replay.ReplayGuardian
}

// NewBlobRequestAuthenticator creates an authenticator for blob requests.
// ReplayGuardian is not used for blob requests.
func NewBlobRequestAuthenticator() *authenticator {
	return &authenticator{
		ReplayGuardian: nil, // Not needed for blob requests
	}
}

// NewPaymentStateAuthenticator creates an authenticator for payment state requests,
// which requires replay protection.
func NewPaymentStateAuthenticator(maxTimeInPast, maxTimeInFuture time.Duration) *authenticator {
	return &authenticator{
		ReplayGuardian: replay.NewReplayGuardian(time.Now, maxTimeInPast, maxTimeInFuture),
	}
}

var _ core.BlobRequestAuthenticator = &authenticator{}

func (*authenticator) AuthenticateBlobRequest(header *core.BlobHeader, signature []byte) error {
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(signature) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(signature))
	}

	blobKey, err := header.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %v", err)
	}

	// Recover public key from signature
	sigPublicKeyECDSA, err := crypto.SigToPub(blobKey[:], signature)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	accountAddr := header.PaymentMetadata.AccountID
	pubKeyAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	if accountAddr.Cmp(pubKeyAddr) != 0 {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil
}

// AuthenticatePaymentStateRequest verifies the signature of the payment state request
// The signature is signed over the byte representation of the account ID and requestHash
// See implementation of BlobRequestSigner.SignPaymentStateRequest for more details
// Deprecated: Use AuthenticatePaymentStateForAllQuorumsRequest instead.
func (a *authenticator) AuthenticatePaymentStateRequest(accountAddr common.Address, request *pb.GetPaymentStateRequest) error {
	return a.authenticatePaymentRequest(accountAddr, request.GetSignature(), request.GetTimestamp())
}

// AuthenticatePaymentStateForAllQuorumsRequest verifies the signature of the quorum specific payment state request
// The signature is signed over the byte representation of the account ID and requestHash
// See implementation of BlobRequestSigner.SignPaymentStateRequest for more details
func (a *authenticator) AuthenticatePaymentStateForAllQuorumsRequest(accountAddr common.Address, request *pb.GetPaymentStateForAllQuorumsRequest) error {
	return a.authenticatePaymentRequest(accountAddr, request.GetSignature(), request.GetTimestamp())
}

// authenticatePaymentRequest is a helper method that handles the common authentication logic for payment requests
func (a *authenticator) authenticatePaymentRequest(accountAddr common.Address, signature []byte, timestamp uint64) error {
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(signature) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(signature))
	}

	requestHash, err := hashing.HashGetPaymentStateRequest(accountAddr, timestamp)
	if err != nil {
		return fmt.Errorf("failed to hash request: %w", err)
	}
	hash := sha256.Sum256(requestHash)

	// Verify the signature
	sigPublicKeyECDSA, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	pubKeyAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	if accountAddr.Cmp(pubKeyAddr) != 0 {
		return errors.New("signature doesn't match with provided public key")
	}

	if a.ReplayGuardian == nil {
		return errors.New("replay guardian is not configured for payment state requests")
	}

	if err := a.ReplayGuardian.VerifyRequest(requestHash, time.Unix(0, int64(timestamp))); err != nil {
		return fmt.Errorf("failed to verify request: %v", err)
	}

	return nil
}
