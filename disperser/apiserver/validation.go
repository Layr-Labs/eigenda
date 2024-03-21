package apiserver

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
)

func (s *DispersalServer) validateRequestAndGetBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*core.Blob, error) {

	data := req.GetData()
	blobSize := len(data)
	// The blob size in bytes must be in range [1, maxBlobSize].
	if blobSize > maxBlobSize {
		return nil, fmt.Errorf("blob size cannot exceed 2 MiB")
	}
	if blobSize == 0 {
		return nil, fmt.Errorf("blob size must be greater than 0")
	}

	if len(req.GetCustomQuorumNumbers()) > 256 {
		return nil, errors.New("invalid request: number of custom_quorum_numbers must not exceed 256")
	}

	quorumConfig, err := s.updateQuorumConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum config: %w", err)
	}

	if len(req.GetCustomQuorumNumbers()) > int(quorumConfig.QuorumCount) {
		return nil, errors.New("invalid request: number of custom_quorum_numbers must not exceed number of quorums")
	}

	seenQuorums := make(map[uint8]struct{})
	// The quorum ID must be in range [0, 254]. It'll actually be converted
	// to uint8, so it cannot be greater than 254.
	for i := range req.GetCustomQuorumNumbers() {

		if req.GetCustomQuorumNumbers()[i] > 254 {
			return nil, fmt.Errorf("invalid request: quorum_numbers must be in range [0, 254], but found %d", req.GetCustomQuorumNumbers()[i])
		}

		quorumID := uint8(req.GetCustomQuorumNumbers()[i])
		if quorumID >= quorumConfig.QuorumCount {
			return nil, fmt.Errorf("invalid request: the quorum_numbers must be in range [0, %d], but found %d", s.quorumConfig.QuorumCount-1, quorumID)
		}

		if _, ok := seenQuorums[quorumID]; ok {
			return nil, fmt.Errorf("invalid request: quorum_numbers must not contain duplicates")
		}
		seenQuorums[quorumID] = struct{}{}

	}

	// Add the required quorums to the list of quorums to check
	for _, quorumID := range quorumConfig.RequiredQuorums {
		if _, ok := seenQuorums[quorumID]; ok {
			return nil, fmt.Errorf("invalid request: quorum_numbers should not include the required quorums, but required quorum %d was found", quorumID)
		}
		seenQuorums[quorumID] = struct{}{}
	}

	if len(seenQuorums) == 0 {
		return nil, fmt.Errorf("invalid request: the blob must be sent to at least one quorum")
	}

	params := make([]*core.SecurityParam, len(seenQuorums))
	i := 0
	for quorumID := range seenQuorums {
		params[i] = &core.SecurityParam{
			QuorumID:              core.QuorumID(quorumID),
			AdversaryThreshold:    quorumConfig.SecurityParams[quorumID].AdversaryThreshold,
			ConfirmationThreshold: quorumConfig.SecurityParams[quorumID].ConfirmationThreshold,
		}
		err = params[i].Validate()
		if err != nil {
			return nil, fmt.Errorf("invalid request: %w", err)
		}
		i++
	}

	header := core.BlobRequestHeader{
		BlobAuthHeader: core.BlobAuthHeader{
			AccountID: req.AccountId,
		},
		SecurityParams: params,
	}

	blob := &core.Blob{
		RequestHeader: header,
		Data:          data,
	}

	return blob, nil
}
