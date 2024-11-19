package controller

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigensdk-go/logging"
	lru "github.com/hashicorp/golang-lru/v2"
)

type NodeClientManager interface {
	GetClient(host, port string) (clients.NodeClientV2, error)
}

type nodeClientManager struct {
	// nodeClients is a cache of node clients keyed by socket address
	nodeClients *lru.Cache[string, clients.NodeClientV2]
	logger      logging.Logger
}

var _ NodeClientManager = (*nodeClientManager)(nil)

func NewNodeClientManager(cacheSize int, logger logging.Logger) (*nodeClientManager, error) {
	closeClient := func(socket string, value clients.NodeClientV2) {
		if err := value.Close(); err != nil {
			logger.Error("failed to close node client", "err", err)
		}
	}
	nodeClients, err := lru.NewWithEvict(cacheSize, closeClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}

	return &nodeClientManager{
		nodeClients: nodeClients,
		logger:      logger,
	}, nil
}

func (m *nodeClientManager) GetClient(host, port string) (clients.NodeClientV2, error) {
	socket := fmt.Sprintf("%s:%s", host, port)
	client, ok := m.nodeClients.Get(socket)
	if !ok {
		var err error
		client, err = clients.NewNodeClientV2(&clients.NodeClientV2Config{
			Hostname: host,
			Port:     port,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create node client at %s: %w", socket, err)
		}

		m.nodeClients.Add(socket, client)
	}

	return client, nil
}
