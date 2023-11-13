package indexer

// HeaderService
type HeaderService interface {

	// GetHeaders returns a list of new headers since the indicated header. PullNewHeaders automatically handles batching and waiting for a specified period if it is already at head. PullNewHeaders sets the finalization status of the headers according to a finalization rule.
	PullNewHeaders(lastHeader *Header) (Headers, bool, error)

	// PullLatestHeader gets the latest header from the chain client
	PullLatestHeader(finalized bool) (*Header, error)
}
