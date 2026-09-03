package v2

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"golang.org/x/sync/errgroup"
)

type signingInfoAccountingMode string

const (
	legacySigningInfoAccountingMode     signingInfoAccountingMode = "legacy"
	blobQuorumSigningInfoAccountingMode signingInfoAccountingMode = "blob_quorums"

	maxBlobQuorumSigningInfoIntervalSeconds = 12 * 60 * 60
	batchQuorumProfileLoadChunkSize         = 100
	batchQuorumProfileLoadConcurrency       = 8
)

// batchQuorumProfile contains the minimum batch data needed to determine which
// operators were capable of validating every blob in the batch.
type batchQuorumProfile struct {
	blobQuorums [][]core.QuorumID
	quorums     map[core.QuorumID]struct{}
}

func newBatchQuorumProfile(batch *corev2.Batch) (*batchQuorumProfile, error) {
	if batch == nil || batch.BatchHeader == nil {
		return nil, fmt.Errorf("batch header is required")
	}

	blobQuorums, err := batch.GetBlobQuorumNumbers()
	if err != nil {
		return nil, fmt.Errorf("get blob quorum numbers: %w", err)
	}

	return newBatchQuorumProfileFromBlobQuorums(blobQuorums)
}

func newBatchQuorumProfileFromBlobQuorums(blobQuorums [][]core.QuorumID) (*batchQuorumProfile, error) {
	if len(blobQuorums) == 0 {
		return nil, fmt.Errorf("batch must contain at least one blob quorum set")
	}
	profile := &batchQuorumProfile{
		blobQuorums: make([][]core.QuorumID, len(blobQuorums)),
		quorums:     make(map[core.QuorumID]struct{}),
	}
	for i, quorumNumbers := range blobQuorums {
		if len(quorumNumbers) == 0 {
			return nil, fmt.Errorf("blob quorum set %d has no quorums", i)
		}

		profile.blobQuorums[i] = append([]core.QuorumID(nil), quorumNumbers...)
		for _, quorum := range quorumNumbers {
			profile.quorums[quorum] = struct{}{}
		}
	}

	return profile, nil
}

// operatorIsEligible mirrors the current validator behavior: an operator can
// sign only when it has a chunk assignment for every blob in the batch.
func (p *batchQuorumProfile) operatorIsEligible(operatorQuorums []core.QuorumID) bool {
	for _, blobQuorums := range p.blobQuorums {
		assigned := false
		for _, blobQuorum := range blobQuorums {
			for _, operatorQuorum := range operatorQuorums {
				if blobQuorum == operatorQuorum {
					assigned = true
					break
				}
			}
			if assigned {
				break
			}
		}
		if !assigned {
			return false
		}
	}

	return true
}

func deduplicateAttestations(attestations []*corev2.Attestation) ([]*corev2.Attestation, error) {
	deduplicated := make([]*corev2.Attestation, 0, len(attestations))
	indices := make(map[string]int, len(attestations))

	for _, attestation := range attestations {
		if attestation == nil || attestation.BatchHeader == nil {
			return nil, fmt.Errorf("attestation batch header is required")
		}
		batchHeaderHash, err := attestation.BatchHeader.Hash()
		if err != nil {
			return nil, fmt.Errorf("hash batch header: %w", err)
		}
		key := hex.EncodeToString(batchHeaderHash[:])

		index, ok := indices[key]
		if !ok {
			indices[key] = len(deduplicated)
			deduplicated = append(deduplicated, attestation)
			continue
		}

		current := deduplicated[index]
		if attestation.AttestedAt > current.AttestedAt ||
			(attestation.AttestedAt == current.AttestedAt && len(attestation.QuorumNumbers) > len(current.QuorumNumbers)) {
			deduplicated[index] = attestation
		}
	}

	return deduplicated, nil
}

func (s *ServerV2) getBatchQuorumProfiles(
	ctx context.Context,
	attestations []*corev2.Attestation,
) (map[string]*batchQuorumProfile, error) {
	profiles := make(map[string]*batchQuorumProfile, len(attestations))
	missingHashes := make([][32]byte, 0)
	missingKeys := make(map[string]struct{})

	for _, attestation := range attestations {
		// The controller initially persists an empty attestation without a
		// non-signer set. Preserve the existing behavior by excluding it.
		if len(attestation.QuorumNumbers) == 0 {
			continue
		}

		batchHeaderHash, err := attestation.BatchHeader.Hash()
		if err != nil {
			return nil, fmt.Errorf("hash batch header: %w", err)
		}
		key := hex.EncodeToString(batchHeaderHash[:])
		if len(attestation.BlobQuorumNumbers) > 0 {
			profile, err := newBatchQuorumProfileFromBlobQuorums(attestation.BlobQuorumNumbers)
			if err != nil {
				return nil, fmt.Errorf("create persisted batch quorum profile for batch %s: %w", key, err)
			}
			profiles[key] = profile
			continue
		}
		if profile, ok := s.batchQuorumProfileCache.Get(key); ok {
			profiles[key] = profile
			continue
		}
		if _, ok := missingKeys[key]; ok {
			continue
		}
		missingKeys[key] = struct{}{}
		missingHashes = append(missingHashes, batchHeaderHash)
	}

	if len(missingHashes) > 0 {
		var profilesMu sync.Mutex
		group, groupCtx := errgroup.WithContext(ctx)
		group.SetLimit(batchQuorumProfileLoadConcurrency)

		for start := 0; start < len(missingHashes); start += batchQuorumProfileLoadChunkSize {
			end := min(start+batchQuorumProfileLoadChunkSize, len(missingHashes))
			hashes := append([][32]byte(nil), missingHashes[start:end]...)
			group.Go(func() error {
				batches, err := s.blobMetadataStore.GetBatches(groupCtx, hashes)
				if err != nil {
					return fmt.Errorf("get batches: %w", err)
				}

				loadedProfiles := make(map[string]*batchQuorumProfile, len(batches))
				for _, batch := range batches {
					profile, err := newBatchQuorumProfile(batch)
					if err != nil {
						return fmt.Errorf("create batch quorum profile: %w", err)
					}
					batchHeaderHash, err := batch.BatchHeader.Hash()
					if err != nil {
						return fmt.Errorf("hash batch header: %w", err)
					}
					key := hex.EncodeToString(batchHeaderHash[:])
					if _, ok := missingKeys[key]; !ok {
						return fmt.Errorf("received unexpected batch %s", key)
					}
					loadedProfiles[key] = profile
				}

				profilesMu.Lock()
				defer profilesMu.Unlock()
				for key, profile := range loadedProfiles {
					profiles[key] = profile
					s.batchQuorumProfileCache.Add(key, profile)
				}
				return nil
			})
		}

		if err := group.Wait(); err != nil {
			return nil, fmt.Errorf("load historical batch quorum profiles: %w", err)
		}
	}

	for key := range missingKeys {
		if _, ok := profiles[key]; !ok {
			return nil, fmt.Errorf("batch quorum profile not found for batch %s", key)
		}
	}

	return profiles, nil
}

func computeBlobQuorumSigningStats(
	attestations []*corev2.Attestation,
	profiles map[string]*batchQuorumProfile,
	operatorQuorumIntervals dataapi.OperatorQuorumIntervals,
) (map[string]map[uint8]int, map[string]map[uint8]int, map[uint8]int, error) {
	numFailed := make(map[string]map[uint8]int)
	numResponsible := make(map[string]map[uint8]int)
	totalNumBatchesPerQuorum := make(map[uint8]int)

	for _, attestation := range attestations {
		if len(attestation.QuorumNumbers) == 0 {
			continue
		}

		batchHeaderHash, err := attestation.BatchHeader.Hash()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("hash batch header: %w", err)
		}
		key := hex.EncodeToString(batchHeaderHash[:])
		profile, ok := profiles[key]
		if !ok {
			return nil, nil, nil, fmt.Errorf("batch quorum profile not found for batch %s", key)
		}

		for quorum := range profile.quorums {
			totalNumBatchesPerQuorum[quorum]++
		}

		nonSigners := make(map[string]struct{}, len(attestation.NonSignerPubKeys))
		for _, pubkey := range attestation.NonSignerPubKeys {
			if pubkey == nil {
				return nil, nil, nil, fmt.Errorf("attestation contains a nil non-signer public key")
			}
			nonSigners[pubkey.GetOperatorID().Hex()] = struct{}{}
		}

		for operatorID := range operatorQuorumIntervals {
			operatorQuorums := operatorQuorumIntervals.GetQuorums(
				operatorID,
				uint32(attestation.ReferenceBlockNumber),
			)
			if !profile.operatorIsEligible(operatorQuorums) {
				continue
			}

			for _, quorum := range operatorQuorums {
				if _, ok := profile.quorums[quorum]; !ok {
					continue
				}
				if _, ok := numResponsible[operatorID]; !ok {
					numResponsible[operatorID] = make(map[uint8]int)
				}
				numResponsible[operatorID][quorum]++

				if _, failed := nonSigners[operatorID]; failed {
					if _, ok := numFailed[operatorID]; !ok {
						numFailed[operatorID] = make(map[uint8]int)
					}
					numFailed[operatorID][quorum]++
				}
			}
		}
	}

	return numFailed, numResponsible, totalNumBatchesPerQuorum, nil
}
