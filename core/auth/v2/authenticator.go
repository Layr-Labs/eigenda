package v2

import (
	"crypto/sha256"
	"errors"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"time"

	"github.com/Layr-Labs/eigenda/common/replay"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type authenticator struct {
	replayGuardian replay.ReplayGuardian
}

func NewAuthenticator(maxTimeInPast, maxTimeInFuture time.Duration) *authenticator {
	return &authenticator{
		replayGuardian: replay.NewReplayGuardian(time.Now, maxTimeInPast, maxTimeInFuture),
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
// The signature is signed over the byte representation of the account ID and nonce
// See implementation of BlobRequestSigner.SignPaymentStateRequest for more details
func (a *authenticator) AuthenticatePaymentStateRequest(accountAddr common.Address, request *pb.GetPaymentStateRequest) error {
	sig := request.GetSignature()
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(sig) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(sig))
	}

	timestamp := request.GetTimestamp()
	nonce := request.GetNonce()

	// Verify the signature
	accountAddrWithNonce := append(accountAddr.Bytes(), nonce...)
	hash := sha256.Sum256(accountAddrWithNonce)
	sigPublicKeyECDSA, err := crypto.SigToPub(hash[:], sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	pubKeyAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	if accountAddr.Cmp(pubKeyAddr) != 0 {
		return errors.New("signature doesn't match with provided public key")
	}

	if err := a.replayGuardian.VerifyRequest(nonce, time.Unix(int64(timestamp), 0)); err != nil {
		return fmt.Errorf("invalid request: the provided nonce has already been used or the request has expired (%v)", err)
	}

	return nil
}
