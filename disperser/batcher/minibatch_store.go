package batcher

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type BatchStatus uint

// Pending: the batch has been created and minibatches are being added. There can be only one pending batch at a time.
// Formed: the batch has been formed and no more minibatches can be added. Implies that all minibatch records and dispersal request records have been created.
//
// Attested: the batch has been attested.
// Failed: the batch has failed. Retries will be attempted at the blob level.
//
// The batch lifecycle is as follows:
// Pending -> Formed -> Attested
// \              /
//
//	\-> Failed <-/
const (
	BatchStatusPending BatchStatus = iota
	BatchStatusFormed
	BatchStatusAttested
	BatchStatusFailed
)

type BatchRecord struct {
	ID                   uuid.UUID `dynamodbav:"-"`
	CreatedAt            time.Time
	ReferenceBlockNumber uint
	Status               BatchStatus `dynamodbav:"BatchStatus"`

	NumMinibatches uint
}

type DispersalResponse struct {
	Signatures  []*core.Signature
	RespondedAt time.Time
	Error       error
}

type MinibatchDispersal struct {
	BatchID         uuid.UUID `dynamodbav:"-"`
	MinibatchIndex  uint
	core.OperatorID `dynamodbav:"-"`
	OperatorAddress gcommon.Address
	Socket          string
	NumBlobs        uint
	RequestedAt     time.Time
	DispersalResponse
}

type BlobMinibatchMapping struct {
	BlobKey        *disperser.BlobKey `dynamodbav:"-"`
	BatchID        uuid.UUID          `dynamodbav:"-"`
	MinibatchIndex uint
	BlobIndex      uint
	BlobHeaderHash [32]byte
	core.BlobHeader
}

type MinibatchStore interface {
	PutBatch(ctx context.Context, batch *BatchRecord) error
	GetBatch(ctx context.Context, batchID uuid.UUID) (*BatchRecord, error)
	GetBatchesByStatus(ctx context.Context, status BatchStatus) ([]*BatchRecord, error)
	MarkBatchFormed(ctx context.Context, batchID uuid.UUID, numMinibatches uint) error
	UpdateBatchStatus(ctx context.Context, batchID uuid.UUID, status BatchStatus) error
	PutDispersal(ctx context.Context, dispersal *MinibatchDispersal) error
	UpdateDispersalResponse(ctx context.Context, dispersal *MinibatchDispersal, response *DispersalResponse) error
	GetDispersal(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*MinibatchDispersal, error)
	GetDispersalsByBatchID(ctx context.Context, batchID uuid.UUID) ([]*MinibatchDispersal, error)
	GetDispersalsByMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*MinibatchDispersal, error)
	GetBlobMinibatchMappings(ctx context.Context, blobKey disperser.BlobKey) ([]*BlobMinibatchMapping, error)
	GetBlobMinibatchMappingsByBatchID(ctx context.Context, batchID uuid.UUID) ([]*BlobMinibatchMapping, error)
	PutBlobMinibatchMappings(ctx context.Context, blobMinibatchMappings []*BlobMinibatchMapping) error
	// GetLatestFormedBatch returns the latest batch that has been formed.
	// If there is no formed batch, it returns nil.
	// It also returns the minibatches that belong to the batch in the ascending order of minibatch index.
	GetLatestFormedBatch(ctx context.Context) (batch *BatchRecord, err error)
	BatchDispersed(ctx context.Context, batchID uuid.UUID, numMinibatches uint) (bool, error)
}
