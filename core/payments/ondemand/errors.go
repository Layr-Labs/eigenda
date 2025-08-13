package ondemand

import (
	"errors"
)

// Sentinel errors for on-demand operations
var (
	// ErrQuorumNotSupported indicates that one or more requested quorums are not supported for on-demand payments.
	// Currently, only quorums 0 and 1 are supported for on-demand payments.
	ErrQuorumNotSupported = errors.New("quorum not supported for on-demand payments")
)