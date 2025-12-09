package v2

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
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

// TODO: make doc with explicit details of why the logic is the way that it is
func AuthenticateBlobRequest(
	header *core.BlobHeader,
	signature []byte,
	anchorSignature []byte,
	disperserId uint32,
	chainId *big.Int,
) error {
	hasLegacySignature := len(signature) > 0
	hasAnchorSignature := len(anchorSignature) > 0

	// At least one signature must be present
	if !hasLegacySignature && !hasAnchorSignature {
		return errors.New("no signatures provided")
	}

	blobKey, err := header.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %v", err)
	}

	// Validate legacy signature if present
	if hasLegacySignature {
		if len(signature) != 65 {
			return fmt.Errorf("signature length is unexpected: %d", len(signature))
		}

		sigPublicKeyECDSA, err := crypto.SigToPub(blobKey[:], signature)
		if err != nil {
			return fmt.Errorf("failed to recover public key from signature: %w", err)
		}

		if header.PaymentMetadata.AccountID.Cmp(crypto.PubkeyToAddress(*sigPublicKeyECDSA)) != 0 {
			return errors.New("signature doesn't match with provided public key")
		}
	}

	// Validate anchor signature if present
	if hasAnchorSignature {
		if len(anchorSignature) != 65 {
			return fmt.Errorf("anchor signature length is unexpected: %d", len(signature))
		}

		anchorHash, err := hashing.ComputeDispersalAnchorHash(chainId, disperserId, blobKey)
		if err != nil {
			return fmt.Errorf("compute dispersal anchor hash: %w", err)
		}

		anchorSignaturePublicKeyECDSA, err := crypto.SigToPub(anchorHash, anchorSignature)
		if err != nil {
			return fmt.Errorf("recover public key from anchor signature: %w", err)
		}

		if header.PaymentMetadata.AccountID.Cmp(crypto.PubkeyToAddress(*anchorSignaturePublicKeyECDSA)) != 0 {
			return errors.New("anchor signature doesn't match with provided public key")
		}
	}

	return nil
}

// AuthenticatePaymentStateRequest verifies the signature of the payment state request
// The signature is signed over the byte representation of the account ID and requestHash
// See implementation of BlobRequestSigner.SignPaymentStateRequest for more details
func (a *authenticator) AuthenticatePaymentStateRequest(accountAddr common.Address, request *pb.GetPaymentStateRequest) error {
	sig := request.GetSignature()
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(sig) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(sig))
	}

	requestHash, err := hashing.HashGetPaymentStateRequest(accountAddr, request.GetTimestamp())
	if err != nil {
		return fmt.Errorf("failed to hash request: %w", err)
	}
	hash := sha256.Sum256(requestHash)

	// Verify the signature
	sigPublicKeyECDSA, err := crypto.SigToPub(hash[:], sig)
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

	timestamp := request.GetTimestamp()
	if err := a.ReplayGuardian.VerifyRequest(requestHash, time.Unix(0, int64(timestamp))); err != nil {
		return fmt.Errorf("failed to verify request: %v", err)
	}

	return nil
}
