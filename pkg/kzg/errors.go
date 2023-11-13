package kzg

import (
	"errors"
)

var ErrFrListTooLarge = errors.New("ErrFrListTooLarge")
var ErrG1ListTooLarge = errors.New("ErrG1ListTooLarge")
var ErrZeroPolyTooLarge = errors.New("ErrZeroPolyTooLarge")
var ErrDestNotPowerOfTwo = errors.New("ErrDestNotPowerOfTwo")
var ErrEmptyLeaves = errors.New("ErrEmptyLeaves")
var ErrEmptyPoly = errors.New("ErrEmptyPoly")
var ErrNotEnoughScratch = errors.New("ErrNotEnoughScratch")
var ErrInvalidDestinationLength = errors.New("ErrInvalidDestinationLength")
var ErrDomainTooSmall = errors.New("ErrDomainTooSmall")
var ErrLengthNotPowerOfTwo = errors.New("ErrLengthNotPowerOfTwo")
var ErrInvalidPolyLengthTooLarge = errors.New("ErrInvalidPolyLengthTooLarge")
var ErrInvalidPolyLengthPowerOfTwo = errors.New("ErrInvalidPolyLengthPowerOfTwo")
