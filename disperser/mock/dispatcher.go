package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
)

type Dispatcher struct {
	state *mock.PrivateOperatorState
}

var _ disperser.Dispatcher = (*Dispatcher)(nil)

func NewDispatcher(state *mock.PrivateOperatorState) disperser.Dispatcher {
	return &Dispatcher{
		state: state,
	}
}

func (d *Dispatcher) DisperseBatch(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, header *core.BatchHeader) chan core.SigningMessage {
	update := make(chan core.SigningMessage)
	message, err := header.GetBatchHeaderHash()
	if err != nil {
		for id := range d.state.PrivateOperators {
			update <- core.SigningMessage{
				Signature: nil,
				Operator:  id,
				Err:       err,
			}
		}
	}

	go func() {
		for id, op := range d.state.PrivateOperators {
			sig := op.KeyPair.SignMessage(message)

			update <- core.SigningMessage{
				Signature: sig,
				Operator:  id,
				Err:       nil,
			}
		}
	}()

	return update
}
