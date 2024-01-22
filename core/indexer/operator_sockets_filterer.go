package indexer

import (
	"context"

	"github.com/Layr-Labs/eigenda/common"
	blsregcoord "github.com/Layr-Labs/eigenda/contracts/bindings/BLSRegistryCoordinatorWithIndices"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type OperatorSocketsFilterer interface {
	FilterHeaders(headers indexer.Headers) ([]indexer.HeaderAndEvents, error)
	GetSyncPoint(latestHeader *indexer.Header) (uint64, error)
	SetSyncPoint(latestHeader *indexer.Header) error
	FilterFastMode(headers indexer.Headers) (*indexer.Header, indexer.Headers, error)
	WatchOperatorSocketUpdate(ctx context.Context, operatorId core.OperatorID) (chan string, error)
}

type operatorSocketsFilterer struct {
	Filterer bind.ContractFilterer
	Address  gethcommon.Address

	FastMode bool
}

func NewOperatorSocketsFilterer(eigenDAServiceManagerAddr gethcommon.Address, client common.EthClient) (*operatorSocketsFilterer, error) {

	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(eigenDAServiceManagerAddr, client)
	if err != nil {
		return nil, err
	}

	blsRegAddress, err := contractEigenDAServiceManager.RegistryCoordinator(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return &operatorSocketsFilterer{
		Address:  blsRegAddress,
		Filterer: client,
		FastMode: false,
	}, nil
}

func (f *operatorSocketsFilterer) FilterHeaders(headers indexer.Headers) ([]indexer.HeaderAndEvents, error) {
	if err := headers.OK(); err != nil {
		return nil, err
	}

	filterer, err := blsregcoord.NewContractBLSRegistryCoordinatorWithIndicesFilterer(f.Address, f.Filterer)
	if err != nil {
		return nil, err
	}
	opts := &bind.FilterOpts{
		Start: headers.First().Number,
		End:   &headers.Last().Number,
	}

	it, err := filterer.FilterOperatorSocketUpdate(opts, [][32]byte{}) // todo: does this work
	if err != nil {
		return nil, err
	}

	var events []indexer.HeaderAndEvents

	for it.Next() {
		event := it.Event

		header, err := headers.GetHeaderByNumber(event.Raw.BlockNumber)
		if err != nil {
			return nil, err
		}
		if !header.BlockHashIs(event.Raw.BlockHash.Bytes()) {
			continue
		}

		events = append(events, indexer.HeaderAndEvents{
			Header: header,
			Events: []indexer.Event{{Type: OperatorSocketUpdate, Payload: event}},
		})
	}

	return events, nil
}

func (f *operatorSocketsFilterer) GetSyncPoint(latestHeader *indexer.Header) (uint64, error) {
	return 0, nil
}

func (f *operatorSocketsFilterer) SetSyncPoint(latestHeader *indexer.Header) error {
	f.FastMode = true
	return nil
}

func (f *operatorSocketsFilterer) FilterFastMode(headers indexer.Headers) (*indexer.Header, indexer.Headers, error) {
	if len(headers) == 0 {
		return nil, nil, nil
	}
	if f.FastMode {
		f.FastMode = false
		return headers.First(), headers, nil
	}
	return nil, headers, nil
}

func (f *operatorSocketsFilterer) WatchOperatorSocketUpdate(ctx context.Context, operatorId core.OperatorID) (chan string, error) {
	filterer, err := blsregcoord.NewContractBLSRegistryCoordinatorWithIndicesFilterer(f.Address, f.Filterer)
	if err != nil {
		return nil, err
	}

	sink := make(chan *blsregcoord.ContractBLSRegistryCoordinatorWithIndicesOperatorSocketUpdate)
	operatorID := [][32]byte{operatorId}
	_, err = filterer.WatchOperatorSocketUpdate(&bind.WatchOpts{Context: ctx}, sink, operatorID)
	if err != nil {
		return nil, err
	}
	socketChan := make(chan string)
	go func() {
		defer close(socketChan)
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-sink:
				socketChan <- event.Socket
			}
		}
	}()
	return socketChan, nil
}
