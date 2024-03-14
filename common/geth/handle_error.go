package geth

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
)

type ChainConnFaults int64

const (
	SenderFault ChainConnFaults = iota
	RPCFault
	TooManyRequest
	Ok
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
		return SenderFault
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
		return SenderFault
	}

	// custom error parsing. If the error message does not contain any error code, which is 3 digit at minimum
	errMsg := err.Error()
	if len(errMsg) < 3 {
		return RPCFault
	}

	// prevent ddos if error message is too large
	if len(errMsg) > 1000 {
		return RPCFault
	}

	// for example, https://docs.alchemy.com/reference/error-reference 500 errors
	if strings.Contains(errMsg, "500") ||
		strings.Contains(errMsg, "503") {
		return RPCFault
	}

	// 400 errors
	// too many requests
	if strings.Contains(errMsg, "429") {
		return TooManyRequest
	}

	if strings.Contains(errMsg, "400") ||
		strings.Contains(errMsg, "401") ||
		strings.Contains(errMsg, "403") {
		return SenderFault
	}

	// by default it is rpc's fault
	return RPCFault
}
