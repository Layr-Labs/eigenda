package clients

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
)

// RequiredQuorumsStore provides a mapping from cert verifier address to the quorums required by the eigenDACertVerifier
// contract located at that address.
type RequiredQuorumsStore struct {
	requiredQuorumsCache sync.Map
	certVerifier         verification.ICertVerifier
}

// NewRequiredQuorumsStore creates a new RequiredQuorumsStore utility
func NewRequiredQuorumsStore(certVerifier verification.ICertVerifier) (*RequiredQuorumsStore, error) {
	return &RequiredQuorumsStore{
		requiredQuorumsCache: sync.Map{},
		certVerifier:         certVerifier,
	}, nil
}

// GetQuorumNumbersRequired returns the required quorums for a given cert verifier contract
//
// If the required quorums for the input certVerifierAddress are already known, they are returned immediately. If
// the required quorums are unknown, this method will attempt to fetch the required quorums from the contract. If the
// fetch is successful, the internal cache is updated with the result, and the result is returned. If the fetch
// is not successful, an error is returned.
//
// NOTE: it is UNSAFE to modify the returned list of quorums
func (rqs *RequiredQuorumsStore) GetQuorumNumbersRequired(
	ctx context.Context,
	certVerifierAddress string,
) ([]uint8, error) {
	requiredQuorums, keyAlreadyExists := rqs.requiredQuorumsCache.Load(certVerifierAddress)
	if keyAlreadyExists {
		return requiredQuorums.([]uint8), nil
	}

	requiredQuorums, err := rqs.certVerifier.GetQuorumNumbersRequired(ctx, certVerifierAddress)
	if err != nil {
		return nil, fmt.Errorf("retrieve required quorum numbers from cert verifier contract: %w", err)
	}

	rqs.requiredQuorumsCache.Store(certVerifierAddress, requiredQuorums)

	return requiredQuorums.([]uint8), nil
}
