package batcher

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type BatchRecord struct {
	ID                   uuid.UUID
	CreatedAt            time.Time
	ReferenceBlockNumber uint
	HeaderHash           [32]byte
	AggregatePubKey      *core.G2Point
	AggregateSignature   *core.Signature
}

type MinibatchRecord struct {
	BatchID              uuid.UUID
	MinibatchIndex       uint
	BlobHeaderHashes     [][32]byte
	BatchSize            uint64 // in bytes
	ReferenceBlockNumber uint
}

type DispersalRequest struct {
	BatchID        uuid.UUID
	MinibatchIndex uint
	core.OperatorID
	OperatorAddress gcommon.Address
	Socket          string
	NumBlobs        uint
	RequestedAt     time.Time
}

type DispersalResponse struct {
	DispersalRequest
	Signatures  []*core.Signature
	RespondedAt time.Time
	Error       error
}

type MinibatchStore interface {
	PutBatch(ctx context.Context, batch *BatchRecord) error
	GetBatch(ctx context.Context, batchID uuid.UUID) (*BatchRecord, error)
	PutMinibatch(ctx context.Context, minibatch *MinibatchRecord) error
	GetMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*MinibatchRecord, error)
	PutDispersalRequest(ctx context.Context, request *DispersalRequest) error
	GetDispersalRequest(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*DispersalRequest, error)
	GetDispersalRequests(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*DispersalRequest, error)
	PutDispersalResponse(ctx context.Context, response *DispersalResponse) error
	GetDispersalResponse(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*DispersalResponse, error)
	GetDispersalResponses(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*DispersalResponse, error)
	GetPendingBatch(ctx context.Context) (*BatchRecord, error)
}
