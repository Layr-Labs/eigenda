package apiserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
)

func (s *DispersalServer) validateBlobRequest(ctx context.Context, blob *core.Blob) error {

	securityParams := blob.RequestHeader.SecurityParams
	if len(securityParams) == 0 {
		return fmt.Errorf("invalid request: security_params must not be empty")
	}
	if len(securityParams) > 256 {
		return fmt.Errorf("invalid request: security_params must not exceed 256")
	}

	seenQuorums := make(map[uint8]struct{})
	// The quorum ID must be in range [0, 254]. It'll actually be converted
	// to uint8, so it cannot be greater than 254.
	for _, param := range securityParams {
		if _, ok := seenQuorums[param.QuorumID]; ok {
			return fmt.Errorf("invalid request: security_params must not contain duplicate quorum_id")
		}
		seenQuorums[param.QuorumID] = struct{}{}

		if param.QuorumID >= s.quorumCount {
			err := s.updateQuorumCount(ctx)
			if err != nil {
				return fmt.Errorf("failed to get onchain quorum count: %w", err)
			}

			if param.QuorumID >= s.quorumCount {
				return fmt.Errorf("invalid request: the quorum_id must be in range [0, %d], but found %d", s.quorumCount-1, param.QuorumID)
			}
		}
	}

	blobSize := len(blob.Data)
	// The blob size in bytes must be in range [1, maxBlobSize].
	if blobSize > maxBlobSize {
		return fmt.Errorf("blob size cannot exceed 2 MiB")
	}
	if blobSize == 0 {
		return fmt.Errorf("blob size must be greater than 0")
	}

	if err := blob.RequestHeader.Validate(); err != nil {
		s.logger.Warn("invalid header", "err", err)
		for _, param := range securityParams {
			quorumId := string(param.QuorumID)
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return err
	}

	return nil

}

func (s *DispersalServer) checkRateLimitsAndAddRates(ctx context.Context, blob *core.Blob, origin, authenticatedAddress string) error {

	// TODO(robert): Remove these locks once we have resolved ratelimiting approach
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, param := range blob.RequestHeader.SecurityParams {

		rates, ok := s.rateConfig.QuorumRateInfos[param.QuorumID]
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
		encodedLength := encoding.GetEncodedBlobLength(length, uint8(param.QuorumThreshold), uint8(param.AdversaryThreshold))
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
			return errSystemThroughputRateLimit
		}

		systemQuorumKey = fmt.Sprintf("%s:%d-blobrate", systemAccountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, systemQuorumKey, blobRateMultiplier, rates.TotalUnauthBlobRate)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("system blob ratelimit exceeded", "systemQuorumKey", systemQuorumKey, "rate", float32(rates.TotalUnauthBlobRate)/blobRateMultiplier)
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
			return errAccountThroughputRateLimit
		}

		accountQuorumKey = fmt.Sprintf("%s:%d-blobrate", accountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, accountQuorumKey, blobRateMultiplier, accountRates.BlobRate)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account blob ratelimit exceeded", "accountQuorumKey", accountQuorumKey, "rate", float32(accountRates.BlobRate)/blobRateMultiplier)
			return errAccountBlobRateLimit
		}

		// Update the quorum rate
		blob.RequestHeader.SecurityParams[i].QuorumRate = accountRates.Throughput
	}
	return nil

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
