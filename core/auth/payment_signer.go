package auth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PaymentSigner struct {
	PrivateKey *ecdsa.PrivateKey
}

var _ core.PaymentSigner = &PaymentSigner{}

func NewPaymentSigner(privateKeyHex string) *PaymentSigner {

	privateKeyBytes := common.FromHex(privateKeyHex)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	return &PaymentSigner{
		PrivateKey: privateKey,
	}
}

func (s *PaymentSigner) SignBlobPayment(header *commonpb.PaymentHeader) ([]byte, error) {
	// Set the account id to the hex encoded public key of the signer
	header.AccountId = hex.EncodeToString(crypto.FromECDSAPub(&s.PrivateKey.PublicKey))
	pm := core.ConvertPaymentHeader(header)
	hash := pm.Hash()

	sig, err := crypto.Sign(hash.Bytes(), s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %v", err)
	}

	return sig, nil
}

func (s *PaymentSigner) SignAccountID(accountID string) ([]byte, error) {
	hash := crypto.Keccak256Hash([]byte(accountID))
	sig, err := crypto.Sign(hash.Bytes(), s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign account ID: %v", err)
	}

	return sig, nil
}

type NoopPaymentSigner struct{}

func NewNoopPaymentSigner() *NoopPaymentSigner {
	return &NoopPaymentSigner{}
}

func (s *NoopPaymentSigner) SignBlobPayment(header *commonpb.PaymentHeader) ([]byte, error) {
	return nil, fmt.Errorf("noop signer cannot sign blob payment header")
}

func (s *NoopPaymentSigner) GetAccountID() (string, error) {
	return "", fmt.Errorf("noop signer cannot get accountID")
}

// VerifyPaymentSignature verifies the signature against the payment metadata
func VerifyPaymentSignature(paymentHeader *commonpb.PaymentHeader, paymentSignature []byte) bool {
	pubKeyBytes, err := hex.DecodeString(paymentHeader.AccountId)
	if err != nil {
		log.Printf("Failed to decode AccountId: %v\n", err)
		return false
	}
	accountPubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		log.Printf("Failed to unmarshal public key: %v\n", err)
		return false
	}

	pm := core.ConvertPaymentHeader(paymentHeader)
	hash := pm.Hash()

	return crypto.VerifySignature(
		crypto.FromECDSAPub(accountPubKey),
		hash.Bytes(),
		paymentSignature[:len(paymentSignature)-1], // Remove recovery ID
	)
}

// VerifyAccountSignature verifies the signature against an account ID
func VerifyAccountSignature(accountID string, paymentSignature []byte) bool {
	pubKeyBytes, err := hex.DecodeString(accountID)
	if err != nil {
		log.Printf("Failed to decode AccountId: %v\n", err)
		return false
	}
	accountPubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		log.Printf("Failed to unmarshal public key: %v\n", err)
		return false
	}

	return crypto.VerifySignature(
		crypto.FromECDSAPub(accountPubKey),
		[]byte(accountID),
		paymentSignature[:len(paymentSignature)-1], // Remove recovery ID
	)
}

func (s *PaymentSigner) GetAccountID() (string, error) {
	return hex.EncodeToString(crypto.FromECDSAPub(&s.PrivateKey.PublicKey)), nil
}
