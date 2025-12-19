package reservationvalidation

import "errors"

// ErrInsufficientBandwidth indicates a reservation debit was rejected due to lack of bandwidth capacity.
//
// This is a sentinel error intended to be wrapped (via fmt.Errorf("%w", ...)) so callers can reliably detect
// this condition via errors.Is, without relying on string matching.
var ErrInsufficientBandwidth = errors.New("insufficient bandwidth")
