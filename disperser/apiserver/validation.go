package apiserver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
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

func (s *DispersalServer) getAccountRate(origin, authenticatedAddress string, quorumID core.QuorumID) (*PerUserRateInfo, string, error) {
	unauthRates, ok := s.rateConfig.QuorumRateInfos[quorumID]
	if !ok {
		return nil, "", fmt.Errorf("no configured rate exists for quorum %d", quorumID)
	}

	rates := &PerUserRateInfo{
		Throughput: unauthRates.PerUserUnauthThroughput,
		BlobRate:   unauthRates.PerUserUnauthBlobRate,
	}

	// Check if the address is in the allowlist
	if len(authenticatedAddress) > 0 {
		quorumRates, ok := s.rateConfig.Allowlist[authenticatedAddress]
		if ok {
			rateInfo, ok := quorumRates[quorumID]
			if ok {
				key := "address:" + authenticatedAddress
				if rateInfo.Throughput > 0 {
					rates.Throughput = rateInfo.Throughput
				}
				if rateInfo.BlobRate > 0 {
					rates.BlobRate = rateInfo.BlobRate
				}
				return rates, key, nil
			}
		}
	}

	// Check if the origin is in the allowlist

	key := "ip:" + origin

	for account, rateInfoByQuorum := range s.rateConfig.Allowlist {
		if !strings.Contains(origin, account) {
			continue
		}

		rateInfo, ok := rateInfoByQuorum[quorumID]
		if !ok {
			break
		}

		if rateInfo.Throughput > 0 {
			rates.Throughput = rateInfo.Throughput
		}

		if rateInfo.BlobRate > 0 {
			rates.BlobRate = rateInfo.BlobRate
		}

		break
	}

	return rates, key, nil

}

func (s *DispersalServer) checkRateLimitsAndAddRates(ctx context.Context, blob *core.Blob, origin, authenticatedAddress string) error {

	// TODO(robert): Remove these locks once we have resolved ratelimiting approach
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, param := range blob.RequestHeader.SecurityParams {

		rates, ok := s.rateConfig.QuorumRateInfos[param.QuorumID]
		quorumId := string(param.QuorumID)
		if !ok {
			return fmt.Errorf("no configured rate exists for quorum %d", param.QuorumID)
		}
		accountRates, accountKey, err := s.getAccountRate(origin, authenticatedAddress, param.QuorumID)
		if err != nil {
			return err
		}

		// Get the encoded blob size from the blob header. Calculation is done in a way that nodes can replicate
		blobSize := len(blob.Data)
		length := encoding.GetBlobLength(uint(blobSize))
		encodedLength := encoding.GetEncodedBlobLength(length, uint8(param.ConfirmationThreshold), uint8(param.AdversaryThreshold))
		encodedSize := encoding.GetBlobSize(encodedLength)

		s.logger.Debug("checking rate limits", "origin", origin, "address", authenticatedAddress, "quorum", param.QuorumID, "encodedSize", encodedSize, "blobSize", blobSize,
			"accountThroughput", accountRates.Throughput, "accountBlobRate", accountRates.BlobRate, "accountKey", accountKey)

		// Check System Ratelimit
		systemQuorumKey := fmt.Sprintf("%s:%d", systemAccountKey, param.QuorumID)
		allowed, err := s.ratelimiter.AllowRequest(ctx, systemQuorumKey, encodedSize, rates.TotalUnauthThroughput)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("system byte ratelimit exceeded", "systemQuorumKey", systemQuorumKey, "rate", rates.TotalUnauthThroughput)
			s.metrics.HandleSystemRateLimitedRequest(quorumId, blobSize, "DisperseBlob")
			return errSystemThroughputRateLimit
		}

		systemQuorumKey = fmt.Sprintf("%s:%d-blobrate", systemAccountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, systemQuorumKey, blobRateMultiplier, rates.TotalUnauthBlobRate)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("system blob ratelimit exceeded", "systemQuorumKey", systemQuorumKey, "rate", float32(rates.TotalUnauthBlobRate)/blobRateMultiplier)
			s.metrics.HandleSystemRateLimitedRequest(quorumId, blobSize, "DisperseBlob")
			return errSystemBlobRateLimit
		}

		// Check Account Ratelimit

		accountQuorumKey := fmt.Sprintf("%s:%d", accountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, accountQuorumKey, encodedSize, accountRates.Throughput)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account byte ratelimit exceeded", "accountQuorumKey", accountQuorumKey, "rate", accountRates.Throughput)
			s.metrics.HandleAccountRateLimitedRequest(quorumId, blobSize, "DisperseBlob")
			return errAccountThroughputRateLimit
		}

		accountQuorumKey = fmt.Sprintf("%s:%d-blobrate", accountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, accountQuorumKey, blobRateMultiplier, accountRates.BlobRate)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account blob ratelimit exceeded", "accountQuorumKey", accountQuorumKey, "rate", float32(accountRates.BlobRate)/blobRateMultiplier)
			s.metrics.HandleAccountRateLimitedRequest(quorumId, blobSize, "DisperseBlob")
			return errAccountBlobRateLimit
		}

		// Update the quorum rate
		blob.RequestHeader.SecurityParams[i].QuorumRate = accountRates.Throughput
	}
	return nil

}
