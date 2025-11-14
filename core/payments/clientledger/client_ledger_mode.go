package clientledger

import "fmt"

// ClientLedgerMode represents the mode of operation for the client ledger, indicating which types of payment should
// be active.
type ClientLedgerMode string

const (
	// Only reservation payments are active
	ClientLedgerModeReservationOnly ClientLedgerMode = "reservation-only"

	// Only on-demand payments are active
	ClientLedgerModeOnDemandOnly ClientLedgerMode = "on-demand-only"

	// Both reservation and on-demand payments are active
	ClientLedgerModeReservationAndOnDemand ClientLedgerMode = "reservation-and-on-demand"
)

// Converts a string to ClientLedgerMode. Panics if an unrecognized mode string is provided.
func ParseClientLedgerMode(mode string) ClientLedgerMode {
	switch mode {
	case string(ClientLedgerModeReservationOnly):
		return ClientLedgerModeReservationOnly
	case string(ClientLedgerModeOnDemandOnly):
		return ClientLedgerModeOnDemandOnly
	case string(ClientLedgerModeReservationAndOnDemand):
		return ClientLedgerModeReservationAndOnDemand
	default:
		panic(fmt.Sprintf("unrecognized client ledger mode: %s", mode))
	}
}
