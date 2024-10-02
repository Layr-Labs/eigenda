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
	dataRate     uint32   // bandwith being reserved
	startEpoch   uint32   // index of epoch where reservation begins
	endEpoch     uint32   // index of epoch where reservation ends
	quoromNumber []uint32 // each byte is a percentage at the corresponding quorum index
	quorumSplit  []byte   // each byte is a percentage at the corresponding quorum index
}

// Protocol defines parameters: FixedFeePerByte; fine to leave global rate-limit offchain atm
type OnDemandPayment struct {
	amountDeposited big.Int
	// amountCollected big.Int
}
