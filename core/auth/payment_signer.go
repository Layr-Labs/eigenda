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

// SignBlobPayment signs the payment header and returns the signature
func (s *PaymentSigner) SignBlobPayment(header *commonpb.PaymentHeader) ([]byte, error) {
	header.AccountId = s.GetAccountID()
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

func (s *NoopPaymentSigner) GetAccountID() string {
	return ""
}

// VerifyPaymentSignature verifies the signature against the payment metadata
func VerifyPaymentSignature(paymentHeader *commonpb.PaymentHeader, paymentSignature []byte) bool {
	pm := core.ConvertPaymentHeader(paymentHeader)
	hash := pm.Hash()

	recoveredPubKey, err := crypto.SigToPub(hash.Bytes(), paymentSignature)
	if err != nil {
		log.Printf("Failed to recover public key from signature: %v\n", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)
	accountId := common.HexToAddress(paymentHeader.AccountId)
	if recoveredAddress != accountId {
		log.Printf("Signature address %s does not match account id %s\n", recoveredAddress.Hex(), accountId.Hex())
		return false
	}

	return crypto.VerifySignature(
		crypto.FromECDSAPub(recoveredPubKey),
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

// GetAccountID returns the Ethereum address of the signer
func (s *PaymentSigner) GetAccountID() string {
	publicKey := crypto.FromECDSAPub(&s.PrivateKey.PublicKey)
	hash := crypto.Keccak256(publicKey[1:])

	return common.BytesToAddress(hash[12:]).Hex()
}
