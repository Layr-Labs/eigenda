package clientledger

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
