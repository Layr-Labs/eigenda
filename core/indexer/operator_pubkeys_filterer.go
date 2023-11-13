package indexer

import (
	"context"
	"errors"
	"sort"

	"github.com/Layr-Labs/eigenda/common"
	blspubkeyreg "github.com/Layr-Labs/eigenda/contracts/bindings/BLSPubkeyRegistry"
	blspubkeycompendium "github.com/Layr-Labs/eigenda/contracts/bindings/BLSPublicKeyCompendium"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/indexer"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type PubKeyAddedEvent struct {
	AddedEvent *blspubkeyreg.ContractBLSPubkeyRegistryOperatorAddedToQuorums
	RegEvent   *blspubkeycompendium.ContractBLSPublicKeyCompendiumNewPubkeyRegistration
}

type operatorPubKeysEvent struct {
	Header      *indexer.Header
	BlockHash   gethcommon.Hash
	BlockNumber uint64
	Index       uint
	Type        string
	Payload     any
}

type operatorPubKeysEventFilterer struct {
	f  *blspubkeyreg.ContractBLSPubkeyRegistryFilterer
	cf *pubkeyRegistrationEventFilterer
}

func newOperatorPubKeysEventFilterer(
	addr gethcommon.Address,
	filterer bind.ContractFilterer,
	regFilterer *pubkeyRegistrationEventFilterer,
) (*operatorPubKeysEventFilterer, error) {
	f, err := blspubkeyreg.NewContractBLSPubkeyRegistryFilterer(addr, filterer)
	if err != nil {
		return nil, err
	}
	return &operatorPubKeysEventFilterer{
		f:  f,
		cf: regFilterer,
	}, nil
}

func (f operatorPubKeysEventFilterer) FilterEvents(
	headers indexer.Headers, opts *bind.FilterOpts,
) ([]operatorPubKeysEvent, error) {
	pubKeyAddedEvts, err := f.filterPubKeyAddedToQuorums(headers, opts)
	if err != nil {
		return nil, err
	}

	pubKeyRemovedEvts, err := f.filterPubKeyRemovedFromQuorums(headers, opts)
	if err != nil {
		return nil, err
	}

	events := append(pubKeyAddedEvts, pubKeyRemovedEvts...)
	sort.Slice(events, func(i, j int) bool {
		if events[i].BlockNumber != events[j].BlockNumber {
			return events[i].BlockNumber < events[j].BlockNumber
		}
		return events[i].Index < events[j].Index
	})
	return events, nil
}

func (f operatorPubKeysEventFilterer) filterPubKeyAddedToQuorums(
	headers indexer.Headers, opts *bind.FilterOpts,
) ([]operatorPubKeysEvent, error) {
	it, err := f.f.FilterOperatorAddedToQuorums(opts)
	if err != nil {
		return nil, err
	}

	events, err := f.filterEvents(headers, it, func(it any) operatorPubKeysEvent {
		event := it.(*blspubkeyreg.ContractBLSPubkeyRegistryOperatorAddedToQuorumsIterator).Event
		return operatorPubKeysEvent{
			BlockHash:   event.Raw.BlockHash,
			BlockNumber: event.Raw.BlockNumber,
			Index:       event.Raw.Index,
			Type:        PubKeyAddedToQuorums,
			Payload: PubKeyAddedEvent{
				AddedEvent: event,
				RegEvent:   nil,
			},
		}
	})
	if err != nil {
		return nil, err
	}

	events, err = f.cf.addPubkeyRegistration(events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (f operatorPubKeysEventFilterer) filterPubKeyRemovedFromQuorums(
	headers indexer.Headers, opts *bind.FilterOpts,
) ([]operatorPubKeysEvent, error) {
	it, err := f.f.FilterOperatorRemovedFromQuorums(opts)
	if err != nil {
		return nil, err
	}
	return f.filterEvents(headers, it, func(it any) operatorPubKeysEvent {
		event := it.(*blspubkeyreg.ContractBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator).Event
		return operatorPubKeysEvent{
			BlockHash:   event.Raw.BlockHash,
			BlockNumber: event.Raw.BlockNumber,
			Index:       event.Raw.Index,
			Type:        PubKeyRemovedFromQuorums,
			Payload:     event,
		}
	})
}

func (f operatorPubKeysEventFilterer) filterEvents(
	headers indexer.Headers,
	iter any,
	fn func(it any) operatorPubKeysEvent,
) ([]operatorPubKeysEvent, error) {
	var events []operatorPubKeysEvent

	it := iter.(interface {
		Next() bool
	})

	for it.Next() {
		event := fn(it)

		header, err := headers.GetHeaderByNumber(event.BlockNumber)
		if err != nil {
			return nil, err
		}
		if !header.BlockHashIs(event.BlockHash.Bytes()) {
			continue
		}

		event.Header = header
		events = append(events, event)
	}

	return events, nil
}

type pubkeyRegistrationEventFilterer struct {
	addr     gethcommon.Address
	f        *blspubkeycompendium.ContractBLSPublicKeyCompendiumFilterer
	filterer bind.ContractFilterer
}

func newPubkeyRegistrationEventFilterer(
	addr gethcommon.Address,
	filterer bind.ContractFilterer,
) (*pubkeyRegistrationEventFilterer, error) {
	f, err := blspubkeycompendium.NewContractBLSPublicKeyCompendiumFilterer(addr, filterer)
	if err != nil {
		return nil, err
	}
	return &pubkeyRegistrationEventFilterer{
		addr:     addr,
		f:        f,
		filterer: filterer,
	}, nil
}

func (f pubkeyRegistrationEventFilterer) addPubkeyRegistration(events []operatorPubKeysEvent) ([]operatorPubKeysEvent, error) {

	if len(events) == 0 {
		return events, nil
	}

	ctx := context.Background()

	operators := make([]interface{}, len(events))
	for i, event := range events {
		operators[i] = event.Payload.(PubKeyAddedEvent).AddedEvent.Operator
	}

	// TODO(robert): Properly set the topic0
	query := [][]interface{}{
		// {"NewPubkeyRegistration(indexed address,(uint256,uint256),(uint256[2],uint256[2]))"},
		{},
		operators,
	}

	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, err
	}

	q := ethereum.FilterQuery{
		Addresses: []gethcommon.Address{f.addr},
		Topics:    topics,
	}

	vLogs, err := f.filterer.FilterLogs(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(vLogs) == 0 {
		return nil, errors.New("no pubkey registration events found")
	}

	eventMap := make(map[gethcommon.Address]*blspubkeycompendium.ContractBLSPublicKeyCompendiumNewPubkeyRegistration, len(vLogs))
	for _, vLog := range vLogs {
		event, err := f.f.ParseNewPubkeyRegistration(vLog)
		if err != nil {
			return nil, err
		}
		eventMap[event.Operator] = event
	}

	for i, event := range events {
		regEvent, ok := eventMap[event.Payload.(PubKeyAddedEvent).AddedEvent.Operator]
		if !ok {
			return nil, errors.New("no pubkey event found for registration event")
		}
		payload := event.Payload.(PubKeyAddedEvent)
		payload.RegEvent = regEvent
		events[i].Payload = payload
	}

	return events, nil
}

type OperatorPubKeysFilterer struct {
	Logger                  common.Logger
	Filterer                bind.ContractFilterer
	BlsRegAddress           gethcommon.Address
	PubKeyCompendiumAddress gethcommon.Address

	FastMode bool
}

func NewOperatorPubKeysFilterer(eigenDAServiceManagerAddr gethcommon.Address, client common.EthClient) (*OperatorPubKeysFilterer, error) {

	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(eigenDAServiceManagerAddr, client)
	if err != nil {
		return nil, err
	}

	blsRegAddress, err := contractEigenDAServiceManager.BlsPubkeyRegistry(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	blsRegistry, err := blspubkeyreg.NewContractBLSPubkeyRegistry(blsRegAddress, client)
	if err != nil {
		return nil, err
	}

	pubkeyCompendiumAddress, err := blsRegistry.PubkeyCompendium(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return &OperatorPubKeysFilterer{
		Filterer:                client,
		BlsRegAddress:           blsRegAddress,
		PubKeyCompendiumAddress: pubkeyCompendiumAddress,
	}, nil
}

var _ indexer.Filterer = (*OperatorPubKeysFilterer)(nil)

func (f *OperatorPubKeysFilterer) FilterHeaders(headers indexer.Headers) ([]indexer.HeaderAndEvents, error) {
	if err := headers.OK(); err != nil {
		return nil, err
	}

	regFilterer, err := newPubkeyRegistrationEventFilterer(f.PubKeyCompendiumAddress, f.Filterer)
	if err != nil {
		return nil, err
	}

	filterer, err := newOperatorPubKeysEventFilterer(f.BlsRegAddress, f.Filterer, regFilterer)
	if err != nil {
		return nil, err
	}

	opts := &bind.FilterOpts{
		Start: headers.First().Number,
		End:   &headers.Last().Number,
	}

	events, err := filterer.FilterEvents(headers, opts)
	if err != nil {
		return nil, err
	}

	var res []indexer.HeaderAndEvents

	for _, event := range events {
		res = append(res, indexer.HeaderAndEvents{
			Header: event.Header,
			Events: []indexer.Event{{Type: event.Type, Payload: event.Payload}},
		})
	}

	return res, nil
}

// GetSyncPoint determines the BlockNumber at which it needs to start syncing from based on both 1) its ability to full its entire state from the chain and 2) its indexing duration requirements.
func (f *OperatorPubKeysFilterer) GetSyncPoint(latestHeader *indexer.Header) (uint64, error) {
	return 0, nil
}

// SetSyncPoint sets the Accumulator to operate in fast mode.
func (f *OperatorPubKeysFilterer) SetSyncPoint(latestHeader *indexer.Header) error {
	f.FastMode = true
	return nil
}

// HandleFastMode handles the fast mode operation of the accumulator. In this mode, it will ignore all headers until it reaching the BlockNumber associated with GetSyncPoint. Upon reaching this BlockNumber, it will pull its entire state from the chain and then proceed with normal syncing.
func (f *OperatorPubKeysFilterer) FilterFastMode(headers indexer.Headers) (*indexer.Header, indexer.Headers, error) {
	if len(headers) == 0 {
		return nil, nil, nil
	}
	if f.FastMode {
		f.FastMode = false
		return headers[0], headers, nil
	}
	return nil, headers, nil
}
