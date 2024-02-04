package indexer

import (
	"bytes"
	"encoding/gob"

	"github.com/Layr-Labs/eigenda/common"
	regcoord "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/indexer"
)

const (
	OperatorSocketUpdate = "operator_socket_update"
)

type OperatorSockets map[core.OperatorID]string

type OperatorSocketsAccumulator struct {
	Logger common.Logger
}

func NewOperatorSocketsAccumulator(logger common.Logger) *OperatorSocketsAccumulator {
	return &OperatorSocketsAccumulator{
		Logger: logger,
	}
}

func (a *OperatorSocketsAccumulator) InitializeObject(header indexer.Header) (indexer.AccumulatorObject, error) {
	return make(OperatorSockets), nil
}

func (a *OperatorSocketsAccumulator) UpdateObject(object indexer.AccumulatorObject, header *indexer.Header, event indexer.Event) (indexer.AccumulatorObject, error) {
	sockets, ok := object.(OperatorSockets)
	if !ok {
		return object, ErrIncorrectObject
	}

	if event.Type != OperatorSocketUpdate {
		return object, ErrIncorrectEvent
	}

	payload, ok := event.Payload.(*regcoord.ContractRegistryCoordinatorOperatorSocketUpdate)
	if !ok {
		return object, ErrIncorrectEvent
	}

	sockets[payload.OperatorId] = payload.Socket

	return object, nil
}

func (a *OperatorSocketsAccumulator) SerializeObject(object indexer.AccumulatorObject, fork indexer.UpgradeFork) ([]byte, error) {
	switch fork {
	case "genesis":
		obj, ok := object.(OperatorSockets)
		if !ok {
			return nil, ErrIncorrectObject
		}

		var (
			buff bytes.Buffer
			enc  = gob.NewEncoder(&buff)
		)

		if err := enc.Encode(obj); err != nil {
			return nil, err
		}

		return buff.Bytes(), nil
	default:
		return nil, ErrUnrecognizedFork
	}
}

func (a *OperatorSocketsAccumulator) DeserializeObject(data []byte, fork indexer.UpgradeFork) (indexer.AccumulatorObject, error) {
	switch fork {
	case "genesis":
		var (
			obj OperatorSockets
			buf = bytes.NewBuffer(data)
			dec = gob.NewDecoder(buf)
		)

		if err := dec.Decode(&obj); err != nil {
			return nil, err
		}

		return obj, nil
	default:
		return nil, ErrUnrecognizedFork
	}
}
