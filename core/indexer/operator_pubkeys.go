package indexer

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/indexer"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

const (
	PubKeyAddedToQuorums     = "pubkey_added_to_quorums"
	PubKeyRemovedFromQuorums = "pubkey_removed_from_quorums"
	NewPubKeyRegistration    = "new_pubkey_registration"
)

type OperatorPubKeysPair struct {
	PubKeyG1 *bn254.G1Affine
	PubKeyG2 *bn254.G2Affine
}

type OperatorPubKeys struct {
	Operators    map[core.OperatorID]OperatorPubKeysPair
	QuorumTotals map[core.QuorumID]*bn254.G1Affine
}

type OperatorPubKeysAccumulator struct {
	Logger common.Logger
}

func NewOperatorPubKeysAccumulator(logger common.Logger) *OperatorPubKeysAccumulator {
	return &OperatorPubKeysAccumulator{
		Logger: logger,
	}
}

var _ indexer.Accumulator = (*OperatorPubKeysAccumulator)(nil)

func (a *OperatorPubKeysAccumulator) InitializeObject(header indexer.Header) (indexer.AccumulatorObject, error) {
	return &OperatorPubKeys{
		Operators:    make(map[core.OperatorID]OperatorPubKeysPair),
		QuorumTotals: make(map[core.QuorumID]*bn254.G1Affine),
	}, nil
}

func newFpElement(x *big.Int) fp.Element {
	var p fp.Element
	p.SetBigInt(x)
	return p
}

func (a *OperatorPubKeysAccumulator) UpdateObject(object indexer.AccumulatorObject, header *indexer.Header, event indexer.Event) (indexer.AccumulatorObject, error) {
	pubKeys, ok := object.(*OperatorPubKeys)
	if !ok {
		return object, ErrIncorrectObject
	}

	switch event.Type {
	case PubKeyAddedToQuorums:
		payload, ok := event.Payload.(PubKeyAddedEvent)
		if !ok {
			return object, ErrIncorrectEvent
		}

		pubKeysPair := OperatorPubKeysPair{
			PubKeyG1: &bn254.G1Affine{
				X: newFpElement(payload.RegEvent.PubkeyG1.X),
				Y: newFpElement(payload.RegEvent.PubkeyG1.Y),
			},
			PubKeyG2: &bn254.G2Affine{
				X: struct{ A0, A1 fp.Element }{
					A0: newFpElement(payload.RegEvent.PubkeyG2.X[1]),
					A1: newFpElement(payload.RegEvent.PubkeyG2.X[0]),
				},
				Y: struct{ A0, A1 fp.Element }{
					A0: newFpElement(payload.RegEvent.PubkeyG2.Y[1]),
					A1: newFpElement(payload.RegEvent.PubkeyG2.Y[0]),
				},
			},
		}

		p := core.G1Point{G1Affine: pubKeysPair.PubKeyG1}
		operatorID := p.GetOperatorID()

		for _, quorumID := range payload.AddedEvent.QuorumNumbers {

			totals, ok := pubKeys.QuorumTotals[core.QuorumID(quorumID)]
			if !ok {
				totals = &bn254.G1Affine{}
			}
			totals.Add(totals, pubKeysPair.PubKeyG1)

			pubKeys.QuorumTotals[core.QuorumID(quorumID)] = totals
		}

		pubKeys.Operators[operatorID] = pubKeysPair
	case PubKeyRemovedFromQuorums:
		// TODO: The operator ID is not available in the event payload, so this requires additional work.

		// payload, ok := event.Payload.(*blspubkeyreg.ContractBLSPubkeyRegistryPubkeyRemovedFromQuorums)
		// if !ok {
		// 	return object, ErrIncorrectEvent
		// }

		// operatorID := core.OperatorId(payload.Operator)
		// pubKeysPair, ok := pubKeys.Operators[operatorID]
		// if !ok {
		// 	return object, ErrOperatorNotFound
		// }

		// for _, quorumID := range payload.QuorumNumbers {

		// 	totals, ok := pubKeys.QuorumTotals[core.QuorumID(quorumID)]
		// 	if !ok {
		// 		totals = &bn254.G1Affine{}
		// 	}
		// 	totals.Sub(totals, pubKeysPair.PubKeyG1)
		// 	pubKeys.QuorumTotals[core.QuorumID(quorumID)] = totals
		// }

		// delete(pubKeys.Operators, operatorID)
	}

	return object, nil
}

// SerializeObject object takes the accummulator object, and serializes it using the rules for the specified fork.
func (a *OperatorPubKeysAccumulator) SerializeObject(object indexer.AccumulatorObject, fork indexer.UpgradeFork) ([]byte, error) {
	switch fork {
	case "genesis":
		obj, ok := object.(*OperatorPubKeys)
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

func (a *OperatorPubKeysAccumulator) DeserializeObject(data []byte, fork indexer.UpgradeFork) (indexer.AccumulatorObject, error) {
	switch fork {
	case "genesis":
		var (
			obj OperatorPubKeys
			buf = bytes.NewBuffer(data)
			dec = gob.NewDecoder(buf)
		)

		if err := dec.Decode(&obj); err != nil {
			return nil, err
		}

		return &obj, nil
	default:
		return nil, ErrUnrecognizedFork
	}
}
