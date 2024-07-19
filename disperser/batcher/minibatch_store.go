package batcher

import (
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
	PutBatch(batch *BatchRecord) error
	GetBatch(batchID uuid.UUID) (*BatchRecord, error)
	PutMiniBatch(minibatch *MinibatchRecord) error
	GetMiniBatch(batchID uuid.UUID, minibatchIndex uint) (*MinibatchRecord, error)
	PutDispersalRequest(request *DispersalRequest) error
	GetDispersalRequest(batchID uuid.UUID, minibatchIndex uint) (*DispersalRequest, error)
	PutDispersalResponse(response *DispersalResponse) error
	GetDispersalResponse(batchID uuid.UUID, minibatchIndex uint) (*DispersalResponse, error)
	GetPendingBatch() (*BatchRecord, error)
}
