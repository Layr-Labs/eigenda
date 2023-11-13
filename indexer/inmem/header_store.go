package inmem

import (
	"errors"

	"github.com/Layr-Labs/eigenda/indexer"
)

var (
	ErrObjectNotFound        = errors.New("object not found")
	ErrHeaderNotFound        = errors.New("header with number not found")
	ErrInconsistentHash      = errors.New("header at number does not match")
	ErrPrevBlockHashNotFound = errors.New("previous block hash not found")
)

type Payloads map[indexer.Accumulator][]byte

type Header struct {
	*indexer.Header
	Payloads Payloads
}

func AddPayloads(headers indexer.Headers, payloads Payloads) []*Header {

	copyPayloads := func(payloads Payloads) Payloads {
		payloadCopy := make(Payloads)
		for k, v := range payloads {
			payloadCopy[k] = v
		}
		return payloadCopy
	}

	newHeaders := make([]*Header, len(headers))
	for ind := range headers {
		newHeaders[ind] = new(Header)
		newHeaders[ind].Header = headers[ind]
		newHeaders[ind].Payloads = copyPayloads(payloads)
	}
	return newHeaders
}

type HeaderStore struct {
	Chain          []*Header
	IndOffset      int
	FinalizedIndex int
}

var _ indexer.HeaderStore = (*HeaderStore)(nil)

func NewHeaderStore() *HeaderStore {
	return &HeaderStore{
		Chain:          make([]*Header, 0),
		IndOffset:      0,
		FinalizedIndex: 0,
	}
}

func (h *HeaderStore) getHeaderByNumber(number uint64) (*Header, int, bool) {

	ind := int(number) - h.IndOffset

	if ind < 0 || ind >= len(h.Chain) {
		return nil, 0, false
	}

	return h.Chain[ind], ind, true

}

func (h *HeaderStore) getHeader(header *indexer.Header) (*Header, int, error) {
	myHeader, ind, found := h.getHeaderByNumber(header.Number)
	if !found {
		return nil, 0, ErrHeaderNotFound
	}
	if header.BlockHash != [32]byte{} && myHeader.BlockHash != header.BlockHash {
		return nil, 0, ErrInconsistentHash
	}

	return myHeader, ind, nil
}

func (h *HeaderStore) updateFinalizedIndex() {

	finalizedIndex := h.FinalizedIndex
	for ind := h.FinalizedIndex; ind < len(h.Chain); ind++ {
		if h.Chain[ind].Finalized {
			finalizedIndex = ind
		} else {
			break
		}
	}
	h.FinalizedIndex = finalizedIndex

}

// Addheaders finds the header  It then crawls along this list of headers until it finds the point of divergence with its existing chain. All new headers from this point of divergence onward are returned.
func (h *HeaderStore) AddHeaders(headers indexer.Headers) (indexer.Headers, error) {

	if len(headers) == 0 {
		return headers, nil
	}

	if !headers.IsOrdered() {
		return nil, indexer.ErrHeadersUnordered
	}

	if len(h.Chain) == 0 {
		h.IndOffset = int(headers[0].Number)
		h.Chain = AddPayloads(headers, make(Payloads))
		h.updateFinalizedIndex()
		return headers, nil
	}

	myHeader, _, found := h.getHeaderByNumber(headers[len(headers)-1].Number)
	if found && myHeader.BlockHash == headers[len(headers)-1].BlockHash {
		return nil, nil
	}

	ind, myInd, err := func() (int, int, error) {

		for ind := len(headers) - 1; ind >= 0; ind-- {
			myHeader, myInd, found := h.getHeaderByNumber(headers[ind].Number - 1)
			if found {
				if myHeader.BlockHash == headers[ind].PrevBlockHash {
					return ind, myInd, nil
				}
			}
		}
		return 0, 0, ErrPrevBlockHashNotFound
	}()
	if err != nil {
		return nil, err
	}

	newHeaders := AddPayloads(headers[ind:], h.Chain[myInd].Payloads)
	h.Chain = append(h.Chain[:myInd+1], newHeaders...)
	h.updateFinalizedIndex()

	return headers[ind:], nil

}

// GetLatestHeader returns the most recent header that the HeaderService has previously pulled
func (h *HeaderStore) GetLatestHeader(finalized bool) (*indexer.Header, error) {
	if len(h.Chain) == 0 {
		return nil, indexer.ErrNoHeaders
	}
	var index int
	if finalized {
		index = h.FinalizedIndex
	} else {
		index = len(h.Chain) - 1
	}
	if index < 0 && index >= len(h.Chain) {
		return nil, ErrHeaderNotFound
	}
	return h.Chain[index].Header, nil
}

// AttachObject takes an accumulator object and attaches it to a header so that it can be retrieved using GetObject
func (h *HeaderStore) AttachObject(object indexer.AccumulatorObject, header *indexer.Header, acc indexer.Accumulator,
) error {

	_, ind, err := h.getHeader(header)
	if err != nil {
		return err
	}

	data, err := acc.SerializeObject(object, indexer.UpgradeFork(header.CurrentFork))
	if err != nil {
		return err
	}

	h.Chain[ind].Payloads[acc] = data

	return nil
}

// GetObject takes in a header and retrieves the accumulator object attached to the latest header prior to the supplied header having the requested object type.
func (h *HeaderStore) GetObject(header *indexer.Header, acc indexer.Accumulator) (indexer.AccumulatorObject, *indexer.Header, error) {

	data, myHeader, found := func() (data []byte, myHeader *Header, found bool) {
		for ind := int(header.Number); ind >= 0; ind-- {

			queryHeader := &indexer.Header{
				Number: uint64(ind),
			}

			myHeader, _, err := h.getHeader(queryHeader)
			if err != nil {
				return nil, nil, false
			}

			var ok bool
			data, ok = myHeader.Payloads[acc]
			if ok {
				return data, myHeader, true
			}
		}
		return nil, nil, false
	}()

	if !found {
		return nil, nil, ErrObjectNotFound
	}

	obj, err := acc.DeserializeObject(data, indexer.UpgradeFork(myHeader.CurrentFork))
	if err != nil {
		return nil, nil, err
	}

	return obj, myHeader.Header, nil
}

// GetObject retrieves the accumulator object attached to the latest header having the requested object type.
func (h *HeaderStore) GetLatestObject(acc indexer.Accumulator, finalized bool) (indexer.AccumulatorObject, *indexer.Header, error) {
	header, err := h.GetLatestHeader(finalized)
	if err != nil {
		return nil, nil, err
	}
	return h.GetObject(header, acc)
}

// GetObject retrieves the accumulator object attached to the latest header having the requested object type.
func (h *HeaderStore) FastForward() {
	h.Chain = make([]*Header, 0)
}
