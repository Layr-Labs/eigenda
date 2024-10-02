package meterer

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

/* SUBJECT TO BIG MODIFICATIONS */

// BlobHeader represents the header information for a blob
type BlobHeader struct {
	// Existing fields
	Commitment    core.G1Point
	DataLength    uint32
	QuorumNumbers []uint8

	// New fields
	Version   uint32
	AccountID string
	Nonce     uint32 // use nonce to prevent duplicate payments in the same reservation window
	BinIndex  uint32

	Signature []byte
	BlobSize  uint32
	// TODO: we are thinking the contract can use uint128 for cumulative payment,
	// but the definition on v2 uses uint64. Double check with team.
	CumulativePayment uint64
}

// // EIP712Domain represents the EIP-712 domain for our blob headers
// var EIP712Domain = apitypes.TypedDataDomain{
// 	Name:              "EigenDA",
// 	Version:           "1",
// 	ChainId:           (*math.HexOrDecimal256)(big.NewInt(17000)),
// 	VerifyingContract: common.HexToAddress("0x1234000000000000000000000000000000000000").Hex(),
// }

// Protocol defines parameters: epoch length and rate-limit window interval
type Reservation struct {
	dataRate    uint32 // bandwith being reserved
	startEpoch  uint32 // index of epoch where reservation begins
	endEpoch    uint32 // index of epoch where reservation ends
	quorumSplit []byte // each byte is a percentage at the corresponding quorum index
}

// Protocol defines parameters: FixedFeePerByte; fine to leave global rate-limit offchain atm
type OnDemand struct {
	amountDeposited big.Int
	amountCollected big.Int
}

// EIP712Signer handles EIP-712 signing operations
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
			"BlobHeader": []apitypes.Type{
				{Name: "version", Type: "uint32"},
				{Name: "accountID", Type: "string"},
				{Name: "nonce", Type: "uint32"},
				{Name: "binIndex", Type: "uint32"},
				{Name: "cumulativePayment", Type: "uint64"},
				{Name: "commitment", Type: "bytes"},
				{Name: "dataLength", Type: "uint32"},
				{Name: "quorumNumbers", Type: "uint8[]"},
			},
		},
	}
}

// SignBlobHeader signs a BlobHeader using EIP-712
func (s *EIP712Signer) SignBlobHeader(header *BlobHeader, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	commitment := header.Commitment.Serialize()
	typedData := apitypes.TypedData{
		Types:       s.types,
		PrimaryType: "BlobHeader",
		Domain:      s.domain,
		Message: apitypes.TypedDataMessage{
			"version":           fmt.Sprintf("%d", header.Version),
			"accountID":         header.AccountID,
			"nonce":             fmt.Sprintf("%d", header.Nonce),
			"binIndex":          fmt.Sprintf("%d", header.BinIndex),
			"cumulativePayment": fmt.Sprintf("%d", header.CumulativePayment),
			"commitment":        hexutil.Encode(commitment),
			"dataLength":        fmt.Sprintf("%d", header.DataLength),
			"quorumNumbers":     convertUint8SliceToMap(header.QuorumNumbers),
		},
	}

	signature, err := s.signTypedData(typedData, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing blob header: %v", err)
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

// RecoverSender recovers the sender's address from a signed BlobHeader
func (s *EIP712Signer) RecoverSender(header *BlobHeader) (common.Address, error) {
	typedData := apitypes.TypedData{
		Types:       s.types,
		PrimaryType: "BlobHeader",
		Domain:      s.domain,
		Message: apitypes.TypedDataMessage{
			"version":           fmt.Sprintf("%d", header.Version),
			"accountID":         header.AccountID,
			"nonce":             fmt.Sprintf("%d", header.Nonce),
			"binIndex":          fmt.Sprintf("%d", header.BinIndex),
			"cumulativePayment": fmt.Sprintf("%d", header.CumulativePayment),
			"commitment":        hexutil.Encode(header.Commitment.Serialize()),
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

// ConstructBlobHeader creates a BlobHeader with a valid signature
func ConstructBlobHeader(
	signer *EIP712Signer,
	version uint32,
	nonce uint32,
	binIndex uint32,
	cumulativePayment uint64,
	commitment core.G1Point,
	blobSize uint32,
	quorumNumbers []uint8,
	privateKey *ecdsa.PrivateKey,
) (*BlobHeader, error) {
	accountID := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	header := &BlobHeader{
		Version:           version,
		AccountID:         accountID,
		Nonce:             nonce,
		BinIndex:          binIndex,
		CumulativePayment: cumulativePayment,
		Commitment:        commitment,
		QuorumNumbers:     quorumNumbers,
		BlobSize:          blobSize,
		DataLength:        blobSize,
	}

	signature, err := signer.SignBlobHeader(header, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing blob header: %v", err)
	}

	header.Signature = signature
	return header, nil
}
