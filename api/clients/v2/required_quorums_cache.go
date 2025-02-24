package clients

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
)

// requiredQuorumsStore provides a mapping from cert verifier address to the quorums required by the eigenDACertVerifier
// contract located at that address.
type requiredQuorumsStore struct {
	requiredQuorumsCache sync.Map
	certVerifier         verification.ICertVerifier
}

// newRequiredQuorumsStore creates a new requiredQuorumsStore utility
func newRequiredQuorumsStore(certVerifier verification.ICertVerifier) (*requiredQuorumsStore, error) {
	return &requiredQuorumsStore{
		requiredQuorumsCache: sync.Map{},
		certVerifier:         certVerifier,
	}, nil
}

// getQuorumNumbersRequired returns the required quorums for a given cert verifier contract
//
// If the required quorums for the input certVerifierAddress are already known, they are returned immediately. If
// the required quorums are unknown, this method will attempt to fetch the required quorums from the contract. If the
// fetch is successful, the internal cache is updated with the result, and the result is returned. If the fetch
// is not successful, an error is returned.
func (rqs *requiredQuorumsStore) getQuorumNumbersRequired(ctx context.Context, certVerifierAddress string) ([]uint8, error) {
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
