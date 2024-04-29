package accumulator

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/Layr-Labs/eigenda/indexer"
	weth "github.com/Layr-Labs/eigenda/indexer/test/accumulator/bindings"
)

var (
	ErrNotImplemented    = errors.New("not implemented")
	ErrIncorrectObject   = errors.New("incorrect object")
	ErrUnrecognizedFork  = errors.New("unrecognized fork")
	ErrHeadersNotOrdered = errors.New("headers not ordered")
)

type Accumulator struct {
}

type AccountBalanceV1 struct {
	Balance uint64
}

func (a *Accumulator) InitializeObject(header indexer.Header) (indexer.AccumulatorObject, error) {

	return AccountBalanceV1{
		Balance: 0,
	}, nil
}

func (a *Accumulator) UpdateObject(object indexer.AccumulatorObject, event indexer.Event) indexer.AccumulatorObject {

	deposit := event.Payload.(weth.WethDeposit)
	obj := object.(AccountBalanceV1)
	obj.Balance += deposit.Wad.Uint64()

	return obj
}

// Serialize object takes the accumulator object, and serializes it using the rules for the specified fork.
func (a *Accumulator) SerializeObject(object indexer.AccumulatorObject, fork indexer.UpgradeFork) ([]byte, error) {

	switch fork {
	case "genesis":

		obj, ok := object.(*AccountBalanceV1)
		if !ok {
			return nil, ErrIncorrectObject
		}

		var buff bytes.Buffer
		enc := gob.NewEncoder(&buff)

		// Encode the value.
		err := enc.Encode(obj)
		if err != nil {
			return nil, err
		}
		return buff.Bytes(), nil

	}

	return nil, ErrUnrecognizedFork

}

func (a *Accumulator) DeserializeObject(data []byte, fork indexer.UpgradeFork) (indexer.AccumulatorObject, error) {

	switch fork {
	case "genesis":

		obj := &AccountBalanceV1{}

		buff := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buff)

		// Encode the value.
		err := dec.Decode(obj)
		if err != nil {
			return nil, err
		}
		return obj, nil

	}

	return nil, ErrUnrecognizedFork

}
