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

type DisperserClientMultiplexer struct {
	logger                   logging.Logger
	config                   *DisperserClientMultiplexerConfig
	disperserRegistry        clients.DisperserRegistry
	signer                   corev2.BlobRequestSigner
	committer                *committer.Committer
	dispersalMetrics         metrics.DispersalMetricer
	disperserConnectionCount uint
	clients                  map[uint32]*DisperserClient
	reputations              map[uint32]*reputation.Reputation
	mu                       sync.Mutex
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
		config:                   config,
		clients:                  make(map[uint32]*DisperserClient),
		disperserRegistry:        disperserRegistry,
		reputations:              make(map[uint32]*reputation.Reputation),
		logger:                   logger,
		signer:                   signer,
		committer:                committer,
		dispersalMetrics:         dispersalMetrics,
		disperserConnectionCount: disperserConnectionCount,
	}
}

func (dcm *DisperserClientMultiplexer) Close() error {
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

func (dcm *DisperserClientMultiplexer) GetDisperserClient(
	ctx context.Context,
	now time.Time,
	onDemandPayment bool,
) (*DisperserClient, error) {
	// we could try to be more fine-grained about our locking, but it's probably not worth the complexity unless
	// contention actually becomes an issue
	dcm.mu.Lock()
	defer dcm.mu.Unlock()

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
		clientConfig := &DisperserClientConfig{
			Hostname:                 connectionInfo.Hostname,
			Port:                     fmt.Sprintf("%d", connectionInfo.Port),
			UseSecureGrpcFlag:        true,
			DisperserConnectionCount: dcm.disperserConnectionCount,
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

// Checks if the existing client for the given disperser ID is outdated based on the current connection info.
// If it is outdated, closes the existing client and removes it from the map.
func (dcm *DisperserClientMultiplexer) cleanupOutdatedClient(
	disperserID uint32,
	latestConnectionInfo *clients.DisperserConnectionInfo,
) {
	client, exists := dcm.clients[disperserID]
	if !exists {
		return
	}

	oldConfig := client.GetConfig()
	if oldConfig.Hostname != latestConnectionInfo.Hostname ||
		oldConfig.Port != fmt.Sprintf("%d", latestConnectionInfo.Port) {
		if err := client.Close(); err != nil {
			dcm.logger.Warn("failed to close outdated disperser client", "disperserID", disperserID, "err", err)
		}
		delete(dcm.clients, disperserID)
	}
}

func (dcm *DisperserClientMultiplexer) getEligibleDispersers(
	ctx context.Context,
	now time.Time,
	onDemandPayment bool,
) (map[uint32]*reputation.Reputation, error) {
	defaultDispersers, err := dcm.disperserRegistry.GetDefaultDispersers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get default dispersers: %w", err)
	}

	eligibleDispersers := make(map[uint32]*reputation.Reputation)
	for _, id := range defaultDispersers {
		if _, exists := dcm.reputations[id]; !exists {
			dcm.reputations[id] = reputation.NewReputation(dcm.config.ReputationConfig, now)
		}
		eligibleDispersers[id] = dcm.reputations[id]
	}

	for _, id := range dcm.config.AdditionalDispersers {
		if _, exists := dcm.reputations[id]; !exists {
			dcm.reputations[id] = reputation.NewReputation(dcm.config.ReputationConfig, now)
		}
		eligibleDispersers[id] = dcm.reputations[id]
	}

	for _, id := range dcm.config.DisperserBlacklist {
		delete(eligibleDispersers, id)
	}

	if onDemandPayment {
		onDemandDispersers, err := dcm.disperserRegistry.GetOnDemandDispersers(ctx)
		if err != nil {
			return nil, fmt.Errorf("get on-demand dispersers: %w", err)
		}

		onDemandSet := make(map[uint32]struct{})
		for _, id := range onDemandDispersers {
			onDemandSet[id] = struct{}{}
		}

		for id := range eligibleDispersers {
			if _, exists := onDemandSet[id]; !exists {
				delete(eligibleDispersers, id)
			}
		}
	}

	return eligibleDispersers, nil
}

func (dcm *DisperserClientMultiplexer) chooseDisperser(
	now time.Time,
	eligibleDispersers map[uint32]*reputation.Reputation,
) (uint32, error) {
	if len(eligibleDispersers) == 0 {
		return 0, fmt.Errorf("no eligible dispersers available")
	}

	// Apply forgiveness to all eligible dispersers
	for _, rep := range eligibleDispersers {
		rep.Forgive(now)
	}

	// Choose the disperser with the highest reputation
	var bestID uint32
	bestScore := -1.0
	for id, rep := range eligibleDispersers {
		if rep.ReputationScore > bestScore {
			bestScore = rep.ReputationScore
			bestID = id
		}
	}

	return bestID, nil
}
