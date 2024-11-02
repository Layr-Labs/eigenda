package auth

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type paymentSigner struct {
	PrivateKey *ecdsa.PrivateKey
}

var _ core.PaymentSigner = &paymentSigner{}

func NewPaymentSigner(privateKeyHex string) (*paymentSigner, error) {
	if len(privateKeyHex) == 0 {
		return nil, fmt.Errorf("private key cannot be empty")
	}
	privateKeyBytes := common.FromHex(privateKeyHex)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hex to ECDSA private key: %w", err)
	}

	return &paymentSigner{
		PrivateKey: privateKey,
	}, nil
}

// SignBlobPayment signs the payment header and returns the signature
func (s *paymentSigner) SignBlobPayment(pm *core.PaymentMetadata) ([]byte, error) {
	hash, err := pm.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash payment header: %w", err)
	}

	sig, err := crypto.Sign(hash[:], s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %w", err)
	}

	return sig, nil
}

type NoopPaymentSigner struct{}

func NewNoopPaymentSigner() *NoopPaymentSigner {
	return &NoopPaymentSigner{}
}

func (s *NoopPaymentSigner) SignBlobPayment(header *core.PaymentMetadata) ([]byte, error) {
	return nil, fmt.Errorf("noop signer cannot sign blob payment header")
}

func (s *NoopPaymentSigner) GetAccountID() string {
	return ""
}

// VerifyPaymentSignature verifies the signature against the payment metadata
func VerifyPaymentSignature(paymentHeader *core.PaymentMetadata, paymentSignature []byte) error {
	hash, err := paymentHeader.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash payment header: %w", err)
	}

	recoveredPubKey, err := crypto.SigToPub(hash[:], paymentSignature)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %w", err)
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)
	accountId := common.HexToAddress(paymentHeader.AccountID)
	if recoveredAddress != accountId {
		return fmt.Errorf("signature address %s does not match account id %s", recoveredAddress.Hex(), accountId.Hex())
	}

	ok := crypto.VerifySignature(
		crypto.FromECDSAPub(recoveredPubKey),
		hash[:],
		paymentSignature[:len(paymentSignature)-1], // Remove recovery ID
	)

	if !ok {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

// GetAccountID returns the Ethereum address of the signer
func (s *paymentSigner) GetAccountID() string {
	publicKey := crypto.FromECDSAPub(&s.PrivateKey.PublicKey)
	hash := crypto.Keccak256(publicKey[1:])

	return common.BytesToAddress(hash[12:]).Hex()
}
