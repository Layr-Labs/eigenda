package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// EIP712Domain represents the EIP-712 domain for our blob headers
var EIP712Domain = apitypes.TypedDataDomain{
	Name:              "EigenDA",
	Version:           "1",
	ChainId:           (*math.HexOrDecimal256)(big.NewInt(17000)),
	VerifyingContract: common.HexToAddress("0x1234000000000000000000000000000000000000").Hex(),
}

// Protocol defines parameters: epoch length and rate-limit window interval
type ActiveReservation struct {
	dataRate    uint32 // bandwith being reserved
	startEpoch  uint32 // index of epoch where reservation begins
	endEpoch    uint32 // index of epoch where reservation ends
	quorumSplit []byte // each byte is a percentage at the corresponding quorum index
}

// Protocol defines parameters: FixedFeePerByte; fine to leave global rate-limit offchain atm
type OnDemandPayment struct {
	amountDeposited big.Int
	// amountCollected big.Int
}

// // Create the typed data for EIP-712 signature verification
// typedData := apitypes.TypedData{
// 	Types: apitypes.Types{
// 		"EIP712Domain": []apitypes.Type{
// 			{Name: "name", Type: "string"},
// 			{Name: "version", Type: "string"},
// 			{Name: "chainId", Type: "uint256"},
// 			{Name: "verifyingContract", Type: "address"},
// 		},
// 		"BlobHeader": []apitypes.Type{
// 			{Name: "version", Type: "uint32"},
// 			{Name: "accountID", Type: "string"},
// 			{Name: "nonce", Type: "uint32"},
// 			{Name: "binIndex", Type: "uint32"},
// 			{Name: "cumulativePayment", Type: "uint64"},
// 			{Name: "commitment", Type: "bytes"},
// 			{Name: "dataLength", Type: "uint32"},
// 			{Name: "blobQuorumParams", Type: "BlobQuorumParam[]"},
// 		},
// 		"BlobQuorumParam": []apitypes.Type{
// 			{Name: "quorumID", Type: "uint8"},
// 			{Name: "adversaryThreshold", Type: "uint32"},
// 			{Name: "quorumThreshold", Type: "uint32"},
// 		},
// 	},
// 	Domain:      EIP712Domain,
// 	PrimaryType: "BlobHeader",
// 	Message: apitypes.TypedDataMessage{
// 		"version":           header.Version,
// 		"accountID":         header.AccountID,
// 		"nonce":             header.Nonce,
// 		"binIndex":          header.BinIndex,
// 		"cumulativePayment": header.CumulativePayment,
// 		"commitment":        header.Commitment.Bytes(),
// 		"dataLength":        header.DataLength,
// 		"blobQuorumParams":  header.BlobQuorumParams,
// 	},
// }
