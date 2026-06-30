// Package service wires together the chainstate indexer and its HTTP API
// server into a single runnable unit. It exists as a separate package (rather
// than living in package chainstate) because the api package imports
// chainstate, so the combined wiring must sit above both to avoid an import
// cycle. Both the chainstate-indexer binary and the e2e tests construct the
// service through this package so they exercise the same startup path.
package service

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/chainstate"
	"github.com/Layr-Labs/eigenda/chainstate/api"
	"github.com/Layr-Labs/eigenda/chainstate/store"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Service bundles a running indexer with its API server.
type Service struct {
	indexer   *chainstate.Indexer
	apiServer *api.Server
	ethClient *ethclient.Client
	logger    logging.Logger

	// errChan receives a fatal error from the API server goroutine, if any.
	errChan chan error
}

// New constructs a Service: it dials the configured Ethereum RPC, builds the
// indexer (resolving contract addresses from the EigenDADirectory), and creates
// the API server. It does not start any background work; call Start for that.
func New(
	ctx context.Context,
	config *chainstate.IndexerConfig,
	secret *chainstate.IndexerSecretConfig,
	logger logging.Logger,
) (*Service, error) {
	if len(secret.EthRpcUrls) == 0 {
		return nil, fmt.Errorf("no Ethereum RPC URLs configured")
	}

	ethClient, err := ethclient.Dial(secret.EthRpcUrls[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}

	indexer, err := chainstate.NewIndexer(ctx, config, ethClient, logger)
	if err != nil {
		ethClient.Close()
		return nil, fmt.Errorf("failed to create indexer: %w", err)
	}

	apiServer := api.NewServer(config, indexer.GetStore(), logger)

	return &Service{
		indexer:   indexer,
		apiServer: apiServer,
		ethClient: ethClient,
		logger:    logger,
		errChan:   make(chan error, 1),
	}, nil
}

// Start begins indexing and launches the API server in a background goroutine.
// A fatal API server error is delivered on the channel returned by Errors.
func (s *Service) Start(ctx context.Context) error {
	if err := s.indexer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start indexer: %w", err)
	}

	go func() {
		if err := s.apiServer.Start(ctx); err != nil {
			s.errChan <- err
		}
	}()

	return nil
}

// Errors returns a channel that receives a fatal API server error, if one
// occurs. At most one error is ever sent.
func (s *Service) Errors() <-chan error {
	return s.errChan
}

// Store returns the underlying store, primarily for tests and introspection.
func (s *Service) Store() store.Store {
	return s.indexer.GetStore()
}

// Wait blocks until the indexer's background goroutines have stopped (which
// happens after the context passed to Start is cancelled), then releases the
// Ethereum client. Waiting on the indexer ensures the persister's final state
// save completes before this returns.
func (s *Service) Wait() {
	s.indexer.Wait()
	s.ethClient.Close()
}
