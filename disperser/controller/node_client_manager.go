package controller

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	lru "github.com/hashicorp/golang-lru/v2"
)

type NodeClientManager interface {
	GetClient(host, port string) (clients.NodeClient, error)
	// GetClientForValidator returns a client for a specific validator, using the appropriate disperser ID
	// if dual-signer mode is enabled. This method requires dual-signer support to be configured.
	GetClientForValidator(host, port string, validatorID core.OperatorID) (clients.NodeClient, error)
}

type nodeClientManager struct {
	// nodeClients is a cache of node clients keyed by socket address with disperser ID
	nodeClients   *lru.Cache[string, clients.NodeClient]
	requestSigner clients.DispersalRequestSigner
	disperserID   uint32
	logger        logging.Logger

	// Multi-signer support (optional, nil if not enabled)
	multiSigner *clients.MultiDispersalRequestSigner
	idTracker   *ValidatorDisperserIDTracker
}

var _ NodeClientManager = (*nodeClientManager)(nil)

func NewNodeClientManager(
	cacheSize int,
	requestSigner clients.DispersalRequestSigner,
	disperserID uint32,
	logger logging.Logger) (NodeClientManager, error) {

	closeClient := func(socket string, value clients.NodeClient) {

		if err := value.Close(); err != nil {
			logger.Error("failed to close node client", "err", err)
		}
	}
	nodeClients, err := lru.NewWithEvict(cacheSize, closeClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}

	return &nodeClientManager{
		nodeClients:   nodeClients,
		requestSigner: requestSigner,
		disperserID:   disperserID,
		logger:        logger,
		multiSigner:   nil,
		idTracker:     nil,
	}, nil
}

// NewNodeClientManagerWithMultiSigner creates a NodeClientManager with multi-signer support.
// This enables the manager to use different disperser IDs for different validators.
func NewNodeClientManagerWithMultiSigner(
	cacheSize int,
	multiSigner *clients.MultiDispersalRequestSigner,
	idTracker *ValidatorDisperserIDTracker,
	logger logging.Logger) (NodeClientManager, error) {

	closeClient := func(socket string, value clients.NodeClient) {
		if err := value.Close(); err != nil {
			logger.Error("failed to close node client", "err", err)
		}
	}
	nodeClients, err := lru.NewWithEvict(cacheSize, closeClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}

	return &nodeClientManager{
		nodeClients:   nodeClients,
		requestSigner: nil, // Not used in multi-signer mode
		disperserID:   0,   // Not used in multi-signer mode
		logger:        logger,
		multiSigner:   multiSigner,
		idTracker:     idTracker,
	}, nil
}

func (m *nodeClientManager) GetClient(host, port string) (clients.NodeClient, error) {
	socket := fmt.Sprintf("%s:%s", host, port)
	client, ok := m.nodeClients.Get(socket)
	if !ok {
		var err error
		client, err = clients.NewNodeClient(
			&clients.NodeClientConfig{
				Hostname:    host,
				Port:        port,
				DisperserID: m.disperserID,
			},
			m.requestSigner)
		if err != nil {
			return nil, fmt.Errorf("failed to create node client at %s: %w", socket, err)
		}

		m.nodeClients.Add(socket, client)
	}

	return client, nil
}

// GetClientForValidator returns a client for a specific validator.
// In multi-signer mode, it looks up the appropriate disperser ID for the validator
// and returns a client configured with that ID. In single-signer mode, it falls back
// to GetClient behavior.
func (m *nodeClientManager) GetClientForValidator(
	host, port string,
	validatorID core.OperatorID,
) (clients.NodeClient, error) {
	// If multi-signer mode is not enabled, fall back to single-signer behavior
	if m.multiSigner == nil {
		return m.GetClient(host, port)
	}

	// Get the disperser ID for this validator
	disperserID := m.idTracker.GetDisperserID(validatorID)

	// Create a signer wrapper for this specific disperser ID
	signer := &multiSignerWrapper{
		multiSigner: m.multiSigner,
		disperserID: disperserID,
	}

	// Cache key includes disperser ID to allow separate clients per ID
	cacheKey := fmt.Sprintf("%s:%s:%d", host, port, disperserID)
	client, ok := m.nodeClients.Get(cacheKey)
	if !ok {
		var err error
		client, err = clients.NewNodeClient(
			&clients.NodeClientConfig{
				Hostname:    host,
				Port:        port,
				DisperserID: disperserID,
			},
			signer)
		if err != nil {
			return nil, fmt.Errorf("failed to create node client at %s:%s with disperser ID %d: %w",
				host, port, disperserID, err)
		}

		m.nodeClients.Add(cacheKey, client)
	}

	return client, nil
}

// multiSignerWrapper wraps a MultiDispersalRequestSigner to sign with a specific disperser ID.
type multiSignerWrapper struct {
	multiSigner *clients.MultiDispersalRequestSigner
	disperserID uint32
}

func (w *multiSignerWrapper) SignStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
) ([]byte, error) {
	//nolint:wrapcheck // Thin wrapper delegates to underlying signer
	return w.multiSigner.SignStoreChunksRequest(ctx, request, w.disperserID)
}
