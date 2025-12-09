package dispersal

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/common/disperser"
	"github.com/Layr-Labs/eigenda/common/reputation"
	authv2 "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Contains the information needed for disperser selection.
type disperserInfo struct {
	id              uint32
	grpcUri         string
	reputationScore float64
}

// Supplies DisperserClients based on a dynamic set of eligible dispersers and their reputations.
//
// This struct is goroutine safe.
type DisperserClientMultiplexer struct {
	logger            logging.Logger
	config            *DisperserClientMultiplexerConfig
	disperserRegistry disperser.DisperserRegistry
	signer            *authv2.LocalBlobRequestSigner
	committer         *committer.Committer
	dispersalMetrics  metrics.DispersalMetricer
	// map from disperser ID to corresponding client that can communicate with that disperser
	clients map[uint32]*DisperserClient
	// map from disperser ID to its reputation tracker
	reputations map[uint32]*reputation.Reputation
	// chooses dispersers based on reputation
	reputationSelector *reputation.ReputationSelector[*disperserInfo]
	// indicates whether Close() has been called
	closed bool
	lock   sync.Mutex
}

func NewDisperserClientMultiplexer(
	logger logging.Logger,
	config *DisperserClientMultiplexerConfig,
	disperserRegistry disperser.DisperserRegistry,
	signer *authv2.LocalBlobRequestSigner,
	committer *committer.Committer,
	dispersalMetrics metrics.DispersalMetricer,
	random *rand.Rand,
) (*DisperserClientMultiplexer, error) {
	reputationSelector, err := reputation.NewReputationSelector(
		logger,
		&config.SelectorConfig,
		random,
		func(d *disperserInfo) float64 { return d.reputationScore },
	)
	if err != nil {
		return nil, fmt.Errorf("create reputation selector: %w", err)
	}

	return &DisperserClientMultiplexer{
		logger:             logger,
		config:             config,
		disperserRegistry:  disperserRegistry,
		signer:             signer,
		committer:          committer,
		dispersalMetrics:   dispersalMetrics,
		clients:            make(map[uint32]*DisperserClient),
		reputations:        make(map[uint32]*reputation.Reputation),
		reputationSelector: reputationSelector,
	}, nil
}

// Closes all underlying [DisperserClient]s
func (dcm *DisperserClientMultiplexer) Close() error {
	dcm.lock.Lock()
	defer dcm.lock.Unlock()

	if dcm.closed {
		return nil
	}
	dcm.closed = true

	var errs []error
	for id, client := range dcm.clients {
		if err := client.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close client %d: %w", id, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("close disperser clients: %w", errors.Join(errs...))
	}
	return nil
}

// Returns a client for the best available disperser based on the current reputations.
func (dcm *DisperserClientMultiplexer) GetDisperserClient(
	ctx context.Context,
	now time.Time,
	// if true, only consider dispersers that support on-demand payments
	onDemandPayment bool,
) (*DisperserClient, error) {
	// we could try to be more fine-grained about our locking, but it's probably not worth the complexity unless
	// contention actually becomes an issue
	dcm.lock.Lock()
	defer dcm.lock.Unlock()

	if dcm.closed {
		return nil, fmt.Errorf("disperser client multiplexer is closed")
	}

	eligibleDispersers, err := dcm.getEligibleDispersers(ctx, now, onDemandPayment)
	if err != nil {
		return nil, fmt.Errorf("get eligible dispersers: %w", err)
	}

	if len(eligibleDispersers) == 0 {
		return nil, fmt.Errorf("no eligible dispersers")
	}

	selectedDisperserInfo, err := dcm.reputationSelector.Select(eligibleDispersers)
	if err != nil {
		return nil, fmt.Errorf("select disperser: %w", err)
	}

	dcm.cleanupOutdatedClient(selectedDisperserInfo.id, selectedDisperserInfo.grpcUri)

	client, exists := dcm.clients[selectedDisperserInfo.id]
	if !exists {
		// create a new client for the selected disperser
		clientConfig := &DisperserClientConfig{
			GrpcUri:                  selectedDisperserInfo.grpcUri,
			UseSecureGrpcFlag:        dcm.config.UseSecureGrpcFlag,
			DisperserConnectionCount: dcm.config.DisperserConnectionCount,
			DisperserID:              selectedDisperserInfo.id,
			ChainID:                  dcm.config.ChainID,
		}

		client, err = NewDisperserClient(
			dcm.logger,
			clientConfig,
			dcm.signer,
			dcm.committer,
			dcm.dispersalMetrics,
		)
		if err != nil {
			return nil, fmt.Errorf("create disperser client for ID %d: %w", selectedDisperserInfo.id, err)
		}

		dcm.clients[selectedDisperserInfo.id] = client
	}

	return client, nil
}

// Reports the outcome of a dispersal attempt to the reputation system.
// If success is true, the disperser's reputation is improved; otherwise, it is degraded.
// Returns an error if the disperserID is not found in the reputation system.
func (dcm *DisperserClientMultiplexer) ReportDispersalOutcome(
	disperserID uint32,
	success bool,
	now time.Time,
) error {
	dcm.lock.Lock()
	defer dcm.lock.Unlock()

	if dcm.closed {
		return fmt.Errorf("disperser client multiplexer is closed")
	}

	reputation, exists := dcm.reputations[disperserID]
	if !exists {
		return fmt.Errorf("disperser ID %d not found in reputation system", disperserID)
	}

	if success {
		reputation.ReportSuccess(now)
	} else {
		reputation.ReportFailure(now)
	}

	return nil
}

// Checks if the existing client for the given disperser ID is outdated based on the current network address.
// If it is outdated, closes the existing client and removes it from the map.
//
// NOTE: This method has an edge case where clients that have already been returned to callers
// via GetDisperserClient() may be closed while still in use. This will cause those in-flight operations
// to fail.
//
// This is an acceptable trade-off because:
//  1. gRPC URI changes for dispersers are rare in practice
//  2. When they do occur, the affected dispersals will fail gracefully with errors
//  3. Failed dispersals during a disperser's gRPC URI transition are tolerable
//  4. The alternative (reference counting) adds significant complexity for a rare edge case
func (dcm *DisperserClientMultiplexer) cleanupOutdatedClient(
	disperserID uint32,
	latestGrpcUri string,
) {
	client, exists := dcm.clients[disperserID]
	if !exists {
		// nothing to clean up, if the client doesn't exist
		return
	}

	// check if the latest gRPC URI matches the existing client's config
	// if not, the existing client is outdated and should be closed and removed
	oldConfig := client.GetConfig()
	if oldConfig.GrpcUri != latestGrpcUri {
		if err := client.Close(); err != nil {
			dcm.logger.Errorf("failed to close outdated disperser client for disperserID %d: %v", disperserID, err)
		}
		// remove the outdated client from the map, but don't delete the reputation. reputation is presumed to remain
		// relevant for a given disperser ID, even if the gRPC URI changes
		delete(dcm.clients, disperserID)
	}
}

// Returns the list of all eligible dispersers, along with their reputations scores and URIs.
//
// All dispersers returned by this function will have corresponding entries in dcm.reputations, since new reputations
// are created internally as needed.
func (dcm *DisperserClientMultiplexer) getEligibleDispersers(
	ctx context.Context,
	now time.Time,
	onDemandPayment bool,
) ([]*disperserInfo, error) {
	defaultDispersers, err := dcm.disperserRegistry.GetDefaultDispersers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get default dispersers: %w", err)
	}

	// Combine default dispersers and additional dispersers
	potentiallyEligibleDispersers := make([]uint32, 0, len(defaultDispersers)+len(dcm.config.AdditionalDispersers))
	potentiallyEligibleDispersers = append(potentiallyEligibleDispersers, defaultDispersers...)
	potentiallyEligibleDispersers = append(potentiallyEligibleDispersers, dcm.config.AdditionalDispersers...)

	eligibleDispersers := make([]*disperserInfo, 0, len(potentiallyEligibleDispersers))
	for _, disperserId := range potentiallyEligibleDispersers {
		if slices.Contains(dcm.config.DisperserBlacklist, disperserId) {
			continue
		}

		// Skip if on-demand payment is required and disperser doesn't support it
		if onDemandPayment {
			supportsOnDemand, err := dcm.disperserRegistry.IsOnDemandDisperser(ctx, disperserId)
			if err != nil {
				dcm.logger.Errorf(
					"failed to check if disperser ID %d supports on-demand, excluding: %v", disperserId, err)
				continue
			}
			if !supportsOnDemand {
				continue
			}
		}

		grpcUri, err := dcm.disperserRegistry.GetDisperserGrpcUri(ctx, disperserId)
		if err != nil {
			dcm.logger.Errorf("failed to get URI for disperser ID %d, excluding from eligible dispersers: %v",
				disperserId, err)
			continue
		}

		// Initialize reputation if it doesn't exist
		if _, exists := dcm.reputations[disperserId]; !exists {
			dcm.reputations[disperserId] = reputation.NewReputation(dcm.config.ReputationConfig, now)
		}

		score := dcm.reputations[disperserId].Score(now)
		dcm.dispersalMetrics.RecordDisperserReputationScore(disperserId, score)
		eligibleDispersers = append(eligibleDispersers, &disperserInfo{
			id:              disperserId,
			grpcUri:         grpcUri,
			reputationScore: score,
		})
	}

	return eligibleDispersers, nil
}
