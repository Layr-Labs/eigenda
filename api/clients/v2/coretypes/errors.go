package coretypes

import "errors"

var (
	ErrBlobLengthSymbolsNotPowerOf2 = errors.New("blob length is not a power of 2")
)
