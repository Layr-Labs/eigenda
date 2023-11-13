package indexer

import "errors"

var (
	ErrNoHeaders = errors.New("no headers")
)

// HeaderStore is a stateful component that maintains a chain of headers and their finalization status.
type HeaderStore interface {

	// AddHeaders finds the header It then crawls along this list of headers until it finds the point of divergence with its existing chain. All new headers from this point of divergence onward are returned.
	AddHeaders(headers Headers) (Headers, error)

	// GetLatestHeader returns the most recent header that the HeaderService has previously pulled
	GetLatestHeader(finalized bool) (*Header, error)

	// AttachObject takes an accumulator object and attaches it to a header so that it can be retrieved using GetObject
	AttachObject(object AccumulatorObject, header *Header, acc Accumulator) error

	// GetObject takes in a header and retrieves the accumulator object attached to the latest header prior
	// to the supplied header having the requested object type.
	GetObject(header *Header, acc Accumulator) (AccumulatorObject, *Header, error)

	// GetLatestObject retrieves the accumulator object attached to the latest header having the requested object type.
	GetLatestObject(acc Accumulator, finalized bool) (AccumulatorObject, *Header, error)

	FastForward()
}
