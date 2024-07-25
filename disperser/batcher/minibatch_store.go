package batcher

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type BatchStatus uint

// Pending: the batch has been created and minibatches are being added. There can be only one pending batch at a time.
// Formed: the batch has been formed and no more minibatches can be added. Implies that all minibatch records and dispersal request records have been created.
//
// Attested: the batch has been attested.
// Failed: the batch has failed.
//
// The batch lifecycle is as follows:
// Pending -> Formed -> Attested
// \              /
//  \-> Failed <-/
const (
	BatchStatusPending BatchStatus = iota
	BatchStatusFormed
	BatchStatusAttested
	BatchStatusFailed
)

type BatchRecord struct {
	ID                   uuid.UUID
	CreatedAt            time.Time
	ReferenceBlockNumber uint
	Status               BatchStatus
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
	GetBatchesByStatus(ctx context.Context, status BatchStatus) ([]*BatchRecord, error)
	UpdateBatchStatus(ctx context.Context, batchID uuid.UUID, status BatchStatus) error
	PutMinibatch(ctx context.Context, minibatch *MinibatchRecord) error
	GetMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*MinibatchRecord, error)
	GetMinibatches(ctx context.Context, batchID uuid.UUID) ([]*MinibatchRecord, error)
	PutDispersalRequest(ctx context.Context, request *DispersalRequest) error
	GetDispersalRequest(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*DispersalRequest, error)
	GetMinibatchDispersalRequests(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*DispersalRequest, error)
	PutDispersalResponse(ctx context.Context, response *DispersalResponse) error
	GetDispersalResponse(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*DispersalResponse, error)
	GetMinibatchDispersalResponses(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*DispersalResponse, error)

	// GetLatestFormedBatch returns the latest batch that has been formed.
	// If there is no formed batch, it returns nil.
	// It also returns the minibatches that belong to the batch in the ascending order of minibatch index.
	GetLatestFormedBatch(ctx context.Context) (batch *BatchRecord, minibatches []*MinibatchRecord, err error)
	BatchDispersed(ctx context.Context, batchID uuid.UUID) (bool, error)
}
