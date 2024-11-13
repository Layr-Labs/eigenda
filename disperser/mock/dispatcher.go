package mock

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/mock"
)

type Dispatcher struct {
	mock.Mock
	state *coremock.PrivateOperatorState
}

var _ disperser.Dispatcher = (*Dispatcher)(nil)

func NewDispatcher(state *coremock.PrivateOperatorState) *Dispatcher {
	return &Dispatcher{
		state: state,
	}
}

func (d *Dispatcher) DisperseBatch(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, header *core.BatchHeader) chan core.SigningMessage {
	args := d.Called()
	var nonSigners map[core.OperatorID]struct{}
	if args.Get(0) != nil {
		nonSigners = args.Get(0).(map[core.OperatorID]struct{})
	}
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
		for id := range state.IndexedOperators {
			info := d.state.PrivateOperators[id]
			if _, ok := nonSigners[id]; ok {
				update <- core.SigningMessage{
					Signature: nil,
					Operator:  id,
					Err:       errors.New("not a signer"),
				}
			} else {
				sig := info.KeyPair.SignMessage(message)

				update <- core.SigningMessage{
					Signature: sig,
					Operator:  id,
					Err:       nil,
				}
			}
		}
	}()

	return update
}
