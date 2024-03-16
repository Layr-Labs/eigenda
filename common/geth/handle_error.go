package geth

import (
	"errors"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
)

type ChainConnFaults int64

const (
	Ok ChainConnFaults = iota
	ConnectionFault
	EVMFault
)

// this function accepts error message returned from API
// decide if it is servers fault
func HandleError(err error) ChainConnFaults {
	if err == nil {
		return Ok
	}

	// All geth core error codes from , https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.13.14/core
	// return false, as it is not RPC issue
	if errors.Is(err, core.ErrNonceTooLow) ||
		errors.Is(err, core.ErrNonceTooHigh) ||
		errors.Is(err, core.ErrNonceMax) ||
		errors.Is(err, core.ErrGasLimitReached) ||
		errors.Is(err, core.ErrInsufficientFundsForTransfer) ||
		errors.Is(err, core.ErrMaxInitCodeSizeExceeded) ||
		errors.Is(err, core.ErrInsufficientFunds) ||
		errors.Is(err, core.ErrGasUintOverflow) ||
		errors.Is(err, core.ErrIntrinsicGas) ||
		errors.Is(err, core.ErrTxTypeNotSupported) ||
		errors.Is(err, core.ErrTipAboveFeeCap) ||
		errors.Is(err, core.ErrTipVeryHigh) ||
		errors.Is(err, core.ErrFeeCapVeryHigh) ||
		errors.Is(err, core.ErrFeeCapTooLow) ||
		errors.Is(err, core.ErrSenderNoEOA) ||
		errors.Is(err, core.ErrBlobFeeCapTooLow) {
		return EVMFault
	}

	// All geth txpool error, https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.13.14/core
	// return false, as it is not RPC issue
	if errors.Is(err, txpool.ErrAlreadyKnown) ||
		errors.Is(err, txpool.ErrInvalidSender) ||
		errors.Is(err, txpool.ErrUnderpriced) ||
		errors.Is(err, txpool.ErrReplaceUnderpriced) ||
		errors.Is(err, txpool.ErrAccountLimitExceeded) ||
		errors.Is(err, txpool.ErrGasLimit) ||
		errors.Is(err, txpool.ErrNegativeValue) ||
		errors.Is(err, txpool.ErrOversizedData) ||
		errors.Is(err, txpool.ErrFutureReplacePending) {
		return EVMFault
	}

	return ConnectionFault
}
