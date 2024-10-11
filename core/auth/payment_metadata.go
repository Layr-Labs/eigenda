package auth

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

/* SUBJECT TO MODIFICATIONS */

// EIP712Signer handles EIP-712 domain specific signing operations over typed and structured data
type EIP712Signer struct {
	domain apitypes.TypedDataDomain
	types  apitypes.Types
}

// NewEIP712Signer creates a new EIP712Signer instance
func NewEIP712Signer(chainID *big.Int, verifyingContract common.Address) *EIP712Signer {
	return &EIP712Signer{
		domain: apitypes.TypedDataDomain{
			Name:              "EigenDA",
			Version:           "1",
			ChainId:           (*math.HexOrDecimal256)(chainID),
			VerifyingContract: verifyingContract.Hex(),
		},
		types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PaymentMetadata": []apitypes.Type{
				{Name: "accountID", Type: "string"},
				{Name: "binIndex", Type: "uint32"},
				{Name: "cumulativePayment", Type: "uint64"},
				{Name: "dataLength", Type: "uint32"},
				{Name: "quorumNumbers", Type: "uint8[]"},
			},
		},
	}
}

// SignPaymentMetadata signs a PaymentMetadata using EIP-712
func (s *EIP712Signer) SignPaymentMetadata(header *core.PaymentMetadata, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	typedData := apitypes.TypedData{
		Types:       s.types,
		PrimaryType: "PaymentMetadata",
		Domain:      s.domain,
		Message: apitypes.TypedDataMessage{
			"accountID":         header.AccountID,
			"binIndex":          fmt.Sprintf("%d", header.BinIndex),
			"cumulativePayment": fmt.Sprintf("%d", header.CumulativePayment),
			"dataLength":        fmt.Sprintf("%d", header.DataLength),
			"quorumNumbers":     convertUint8SliceToMap(header.QuorumNumbers),
		},
	}

	signature, err := s.signTypedData(typedData, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing payment metadata (header): %v", err)
	}

	return signature, nil
}

func convertUint8SliceToMap(params []uint8) []string {
	result := make([]string, len(params))
	for i, param := range params {
		result[i] = fmt.Sprintf("%d", param) // Converting uint32 to string
	}
	return result
}

// RecoverSender recovers the sender's address from a signed PaymentMetadata
func (s *EIP712Signer) RecoverSender(header *core.PaymentMetadata) (common.Address, error) {
	typedData := apitypes.TypedData{
		Types:       s.types,
		PrimaryType: "PaymentMetadata",
		Domain:      s.domain,
		Message: apitypes.TypedDataMessage{
			"accountID":         header.AccountID,
			"binIndex":          fmt.Sprintf("%d", header.BinIndex),
			"cumulativePayment": fmt.Sprintf("%d", header.CumulativePayment),
			"dataLength":        fmt.Sprintf("%d", header.DataLength),
			"quorumNumbers":     convertUint8SliceToMap(header.QuorumNumbers),
		},
	}

	return s.recoverTypedData(typedData, header.Signature)
}

func (s *EIP712Signer) signTypedData(typedData apitypes.TypedData, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("error hashing EIP712Domain: %v", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("error hashing primary type: %v", err)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	digest := crypto.Keccak256(rawData)

	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing digest: %v", err)
	}

	return signature, nil
}

func (s *EIP712Signer) recoverTypedData(typedData apitypes.TypedData, signature []byte) (common.Address, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return common.Address{}, fmt.Errorf("error hashing EIP712Domain: %v", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Address{}, fmt.Errorf("error hashing primary type: %v", err)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	digest := crypto.Keccak256(rawData)

	pubKey, err := crypto.SigToPub(digest, signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("error recovering public key: %v", err)
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

// ConstructPaymentMetadata creates a PaymentMetadata with a valid signature
func ConstructPaymentMetadata(
	signer *EIP712Signer,
	binIndex uint32,
	cumulativePayment uint64,
	dataLength uint32,
	quorumNumbers []uint8,
	privateKey *ecdsa.PrivateKey,
) (*core.PaymentMetadata, error) {
	accountID := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	header := &core.PaymentMetadata{
		AccountID:         accountID,
		BinIndex:          binIndex,
		CumulativePayment: cumulativePayment,
		QuorumNumbers:     quorumNumbers,
		DataLength:        dataLength,
	}

	signature, err := signer.SignPaymentMetadata(header, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing payment metadata (header): %v", err)
	}

	header.Signature = signature
	return header, nil
}
