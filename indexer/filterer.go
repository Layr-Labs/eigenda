package indexer

type Event struct {
	Type    string
	Payload interface{}
}

type HeaderAndEvents struct {
	Header *Header
	Events []Event
}

type Filterer interface {

	// FilterHeaders accepts a list of incoming headers. Will throw an error is the accumulator does not have an existing header which can form a chain with the incoming headers. The Accumulator will discard any orphaned headers.
	FilterHeaders(headers Headers) ([]HeaderAndEvents, error)

	// GetSyncPoint determines the blockNumber at which it needs to start syncing from based on both 1) its ability to full its entire state from the chain and 2) its indexing duration requirements.
	GetSyncPoint(latestHeader *Header) (uint64, error)

	// SetSyncPoint sets the Accumulator to operate in fast mode.
	SetSyncPoint(latestHeader *Header) error

	// HandleFastMode handles the fast mode operation of the accumulator. In this mode, it will ignore all headers until it reaching the blockNumber associated with GetSyncPoint. Upon reaching this blockNumber, it will pull its entire state from the chain and then proceed with normal syncing.
	FilterFastMode(headers Headers) (*Header, Headers, error)
}
