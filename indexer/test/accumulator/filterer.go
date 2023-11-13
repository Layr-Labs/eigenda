package accumulator

import (
	"bytes"

	"github.com/Layr-Labs/eigenda/indexer"
	weth "github.com/Layr-Labs/eigenda/indexer/test/accumulator/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Filterer struct {
	Filterer bind.ContractFilterer
	Address  common.Address
	Accounts []common.Address

	FastMode bool
}

func (f *Filterer) FilterHeaders(headers indexer.Headers) ([]indexer.HeaderAndEvents, error) {

	if !headers.IsOrdered() {
		return nil, ErrHeadersNotOrdered
	}

	wethFilterer, err := weth.NewWethFilterer(f.Address, f.Filterer)
	if err != nil {
		return nil, err
	}

	opts := &bind.FilterOpts{
		Start: headers[0].Number,
		End:   &headers[len(headers)-1].Number,
	}

	iter, err := wethFilterer.FilterDeposit(opts, f.Accounts)
	if err != nil {
		return nil, err
	}

	headerAndEvents := make([]indexer.HeaderAndEvents, 0)

	for iter.Next() {

		event := *iter.Event

		header, err := headers.GetHeaderByNumber(event.Raw.BlockNumber)
		if err != nil {
			continue
		}

		if !bytes.Equal(header.BlockHash[:], event.Raw.BlockHash.Bytes()) {
			continue
		}

		headerAndEvents = append(headerAndEvents, indexer.HeaderAndEvents{
			Header: header,
			Events: []indexer.Event{
				{
					Type:    "Deposit",
					Payload: event,
				},
			},
		})

	}

	return headerAndEvents, nil

}

// GetSyncPoint determines the blockNumber at which it needs to start syncing from based on both 1) its ability to full its entire state from the chain and 2) its indexing duration requirements.
func (f *Filterer) GetSyncPoint(latestHeader indexer.Header) (uint64, error) {
	return 0, nil
}

// SetSyncPoint sets the Accumulator to operate in fast mode.
func (f *Filterer) SetSyncPoint(latestHeader indexer.Header) error {
	f.FastMode = true
	return nil
}

// HandleFastMode handles the fast mode operation of the accumulator. In this mode, it will ignore all headers until it reaching the blockNumber associated with GetSyncPoint. Upon reaching this blockNumber, it will pull its entire state from the chain and then proceed with normal syncing.
func (f *Filterer) FilterFastMode(headers []indexer.Header) (*indexer.Header, []indexer.Header, error) {

	if len(headers) == 0 {
		return nil, nil, nil
	}

	if f.FastMode {
		return &headers[0], headers, nil
	}
	return nil, headers, nil
}
