package reservation

import (
	"errors"
)

// Sentinel errors for reservation operations
var (
	// ErrQuorumNotPermitted indicates that one or more requested quorums are not permitted by the reservation.
	ErrQuorumNotPermitted = errors.New("quorum not permitted")

	// ErrTimeOutOfRange indicates the dispersal time is outside the reservation's valid time range.
	ErrTimeOutOfRange = errors.New("time outside reservation range")

	// ErrLockAcquisition indicates failure to acquire the internal reservation lock.
	ErrLockAcquisition = errors.New("acquire reservation lock")

	// ErrTimeMovedBackward indicates a timestamp was observed that is before a previously observed timestamp.
	// This is unlikely to happen in practice, but could be the result of clock drift and NTP adjustments.
	ErrTimeMovedBackward = errors.New("time moved backward")
)
