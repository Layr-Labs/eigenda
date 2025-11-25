package dispersal

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/common/reputation"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Supplies DisperserClients based on a dynamic set of eligible dispersers and their reputations.
//
// This struct is goroutine safe.
type DisperserClientMultiplexer struct {
	logger            logging.Logger
	config            *DisperserClientMultiplexerConfig
	disperserRegistry clients.DisperserRegistry
	signer            corev2.BlobRequestSigner
	committer         *committer.Committer
	dispersalMetrics  metrics.DispersalMetricer
	// number of grpc connections to each disperser
	disperserConnectionCount uint
	// map from disperser ID to corresponding client that can communicate with that disperser
	clients map[uint32]*DisperserClient
	// map from disperser ID to its reputation tracker
	reputations map[uint32]*reputation.Reputation
	// indicates whether Close() has been called
	closed bool
	lock   sync.Mutex
}

func NewDisperserClientMultiplexer(
	logger logging.Logger,
	config *DisperserClientMultiplexerConfig,
	disperserRegistry clients.DisperserRegistry,
	signer corev2.BlobRequestSigner,
	committer *committer.Committer,
	dispersalMetrics metrics.DispersalMetricer,
	disperserConnectionCount uint,
) *DisperserClientMultiplexer {
	return &DisperserClientMultiplexer{
		logger:                   logger,
		config:                   config,
		disperserRegistry:        disperserRegistry,
		signer:                   signer,
		committer:                committer,
		dispersalMetrics:         dispersalMetrics,
		disperserConnectionCount: disperserConnectionCount,
		clients:                  make(map[uint32]*DisperserClient),
		reputations:              make(map[uint32]*reputation.Reputation),
	}
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

	chosenDisperser, err := dcm.chooseDisperser(now, eligibleDispersers)
	if err != nil {
		return nil, fmt.Errorf("choose disperser: %w", err)
	}

	connectionInfo, err := dcm.disperserRegistry.GetDisperserConnectionInfo(ctx, chosenDisperser)
	if err != nil {
		return nil, fmt.Errorf("get disperser connection info for ID %d: %w", chosenDisperser, err)
	}

	dcm.cleanupOutdatedClient(chosenDisperser, connectionInfo)

	client, exists := dcm.clients[chosenDisperser]
	if !exists {
		// create a new client for the chosen disperser
		clientConfig := &DisperserClientConfig{
			Hostname:                 connectionInfo.Hostname,
			Port:                     fmt.Sprintf("%d", connectionInfo.Port),
			UseSecureGrpcFlag:        true,
			DisperserConnectionCount: dcm.disperserConnectionCount,
			DisperserID:              chosenDisperser,
		}

		client, err = NewDisperserClient(
			dcm.logger,
			clientConfig,
			dcm.signer,
			dcm.committer,
			dcm.dispersalMetrics,
		)
		if err != nil {
			return nil, fmt.Errorf("create disperser client for ID %d: %w", chosenDisperser, err)
		}

		dcm.clients[chosenDisperser] = client
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

// Checks if the existing client for the given disperser ID is outdated based on the current connection info.
// If it is outdated, closes the existing client and removes it from the map.
func (dcm *DisperserClientMultiplexer) cleanupOutdatedClient(
	disperserID uint32,
	latestConnectionInfo *clients.DisperserConnectionInfo,
) {
	client, exists := dcm.clients[disperserID]
	if !exists {
		// nothing to clean up, if the client doesn't exist
		return
	}

	// check if the latest connection info matches the existing client's config
	// if not, the existing client is outdated and should be closed and removed
	oldConfig := client.GetConfig()
	if oldConfig.Hostname != latestConnectionInfo.Hostname ||
		oldConfig.Port != fmt.Sprintf("%d", latestConnectionInfo.Port) {
		if err := client.Close(); err != nil {
			dcm.logger.Errorf("failed to close outdated disperser client for disperserID %d: %v", disperserID, err)
		}
		// remove the outdated client from the map, but don't delete the reputation. reputation is presumed to remain
		// relevant for a given disperser ID, even if the connection info changes
		delete(dcm.clients, disperserID)
	}
}

// Returns the IDs of all eligible dispersers, along with their reputations.
func (dcm *DisperserClientMultiplexer) getEligibleDispersers(
	ctx context.Context,
	now time.Time,
	onDemandPayment bool,
) (map[uint32]*reputation.Reputation, error) {
	defaultDispersers, err := dcm.disperserRegistry.GetDefaultDispersers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get default dispersers: %w", err)
	}

	// start by assuming that all default dispersers are eligible
	eligibleDispersers := make(map[uint32]*reputation.Reputation)
	for id := range defaultDispersers {
		if _, exists := dcm.reputations[id]; !exists {
			dcm.reputations[id] = reputation.NewReputation(dcm.config.ReputationConfig, now)
		}
		eligibleDispersers[id] = dcm.reputations[id]
	}

	// add any additional dispersers specified in the config
	for _, id := range dcm.config.AdditionalDispersers {
		if _, exists := dcm.reputations[id]; !exists {
			dcm.reputations[id] = reputation.NewReputation(dcm.config.ReputationConfig, now)
		}
		eligibleDispersers[id] = dcm.reputations[id]
	}

	// remove any dispersers that are blacklisted
	for _, id := range dcm.config.DisperserBlacklist {
		delete(eligibleDispersers, id)
	}

	// if on-demand payment support is required, filter out dispersers that don't support it
	if onDemandPayment {
		onDemandDispersers, err := dcm.disperserRegistry.GetOnDemandDispersers(ctx)
		if err != nil {
			return nil, fmt.Errorf("get on-demand dispersers: %w", err)
		}

		// Rebuild eligibleDispersers with only on-demand dispersers
		filtered := make(map[uint32]*reputation.Reputation, len(onDemandDispersers))
		for id := range onDemandDispersers {
			if reputation, exists := eligibleDispersers[id]; exists {
				filtered[id] = reputation
			}
		}
		eligibleDispersers = filtered
	}

	return eligibleDispersers, nil
}

// Chooses the best disperser from the eligible set based on their reputations.
func (dcm *DisperserClientMultiplexer) chooseDisperser(
	now time.Time,
	eligibleDispersers map[uint32]*reputation.Reputation,
) (uint32, error) {
	if len(eligibleDispersers) == 0 {
		return 0, fmt.Errorf("no eligible dispersers")
	}

	// Choose the disperser with the highest reputation
	//
	// TODO(litt3): At some point, we might consider adding some randomness here
	var bestID uint32
	bestScore := -1.0
	for disperserId, disperserReputation := range eligibleDispersers {
		score := disperserReputation.Score(now)
		if score > bestScore {
			bestScore = score
			bestID = disperserId
		}
	}

	return bestID, nil
}
