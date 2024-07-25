package mock

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
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

func (d *Dispatcher) SendBlobsToOperator(ctx context.Context, blobs []*core.BlobMessage, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) ([]*core.Signature, error) {
	args := d.Called(ctx, blobs, batchHeader, op)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*core.Signature), args.Error(1)
}

func (d *Dispatcher) AttestBatch(ctx context.Context, state *core.IndexedOperatorState, blobHeaderHashes [][32]byte, batchHeader *core.BatchHeader) (chan core.SigningMessage, error) {
	args := d.Called(ctx, state, blobHeaderHashes, batchHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(chan core.SigningMessage), args.Error(1)
}

func (d *Dispatcher) SendAttestBatchRequest(ctx context.Context, nodeDispersalClient node.DispersalClient, blobHeaderHashes [][32]byte, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) (*core.Signature, error) {
	args := d.Called(ctx, nodeDispersalClient, blobHeaderHashes, batchHeader, op)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Signature), args.Error(1)
}
