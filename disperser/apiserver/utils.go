package apiserver

import (
	"context"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
)

func (s *DispersalServer) updateQuorumCount(ctx context.Context) error {
	currentBlock, err := s.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}
	count, err := s.tx.GetQuorumCount(ctx, currentBlock)
	if err != nil {
		return err
	}

	s.logger.Debug("updating quorum count", "currentBlock", currentBlock, "count", count)
	s.mu.Lock()
	s.quorumCount = count
	s.mu.Unlock()
	return nil
}

func getResponseStatus(status disperser.BlobStatus) pb.BlobStatus {
	switch status {
	case disperser.Processing:
		return pb.BlobStatus_PROCESSING
	case disperser.Confirmed:
		return pb.BlobStatus_CONFIRMED
	case disperser.Failed:
		return pb.BlobStatus_FAILED
	case disperser.Finalized:
		return pb.BlobStatus_FINALIZED
	case disperser.InsufficientSignatures:
		return pb.BlobStatus_INSUFFICIENT_SIGNATURES
	default:
		return pb.BlobStatus_UNKNOWN
	}
}

func getBlobFromRequest(req *pb.DisperseBlobRequest) *core.Blob {
	params := make([]*core.SecurityParam, len(req.SecurityParams))

	for i, param := range req.GetSecurityParams() {
		params[i] = &core.SecurityParam{
			QuorumID:           core.QuorumID(param.QuorumId),
			AdversaryThreshold: uint8(param.AdversaryThreshold),
			QuorumThreshold:    uint8(param.QuorumThreshold),
		}
	}

	data := req.GetData()

	blob := &core.Blob{
		RequestHeader: core.BlobRequestHeader{
			BlobAuthHeader: core.BlobAuthHeader{
				AccountID: req.AccountId,
			},
			SecurityParams: params,
		},
		Data: data,
	}

	return blob
}
