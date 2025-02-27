package eth

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// ChainState implements core.ChainState interface to query on-chain states
// and filter socket updates through event logs to track the socket addresses
// of operators in a map.
type SocketStateCache struct {
	Reader core.Reader
	logger logging.Logger
	// A cache map of the operator registry, key: operator id, value: socket string
	SocketMap core.OperatorSockets
	// Mutex to access socket map
	socketMu sync.Mutex
	// The previous block number the socket map was updated at, inclusive
	socketPrevBlockNumber uint32
}

var _ core.SocketStateCache = (*SocketStateCache)(nil)

func NewSocketStateCache(ctx context.Context, reader core.Reader, logger logging.Logger) (*SocketStateCache, error) {
	currentBlockNumber, err := reader.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	cs := &SocketStateCache{
		Reader:    reader,
		logger:    logger,
		SocketMap: make(map[core.OperatorID]*core.OperatorSocket),
		// Set initial block number to current block number
		socketPrevBlockNumber: currentBlockNumber,
	}
	return cs, nil
}

// GetOperatorSocket returns the socket address for a given operator at the current block number,
// and it takes blockNumber due to the core.ChainState interface.
func (cs *SocketStateCache) GetOperatorSocket(ctx context.Context, operator core.OperatorID) (string, error) {
	socket, err := cs.Reader.GetOperatorSocket(ctx, operator)
	if err != nil {
		cs.logger.Warn("failed to get operator socket from Eth Client reader", "operator", operator, "error", err)
		return "", fmt.Errorf("failed to get socket for operator %x: %w", operator, err)
	}
	return socket, nil
}

// indexSocketMap filters event logs from the previously checked block number to the current block,
// to identify all socket update events in that block range, and update the socket map accordingly
func (cs *SocketStateCache) indexSocketMap(ctx context.Context) error {
	currentBlockNumber, err := cs.Reader.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	updates, err := cs.Reader.GetSocketUpdates(ctx, uint64(cs.socketPrevBlockNumber), uint64(currentBlockNumber))
	if err != nil {
		return fmt.Errorf("failed to get socket updates: %w", err)
	}

	for _, update := range updates {
		socket, err := convertSocketStringToOperatorSocket(update.Socket)
		if err != nil {
			cs.logger.Warn("failed to convert socket string to operator socket", "operator", update.OperatorId, "error", err)
			continue
		}
		cs.SocketMap[core.OperatorID(update.OperatorId)] = &socket
	}

	// Store current block as a processed block
	cs.socketPrevBlockNumber = currentBlockNumber

	return nil
}

// refreshSocketMap refresh the socket map for the given operators at the current block.
func (cs *SocketStateCache) refreshSocketMap(ctx context.Context, operators []core.OperatorID) error {
	for _, operator := range operators {
		_, ok := cs.SocketMap[operator]

		if !ok {
			socketString, err := cs.Reader.GetOperatorSocket(ctx, operator)
			if err != nil {
				return fmt.Errorf("failed to get socket for operator %x: %w", operator, err)
			}
			socket, err := convertSocketStringToOperatorSocket(socketString)
			if err != nil {
				cs.logger.Warn("failed to convert socket string to operator socket", "operator", operator, "error", err)
				continue
			}
			cs.SocketMap[operator] = &socket
		}
	}

	// Index for recent socket updates
	if err := cs.indexSocketMap(ctx); err != nil {
		return err
	}
	return nil
}

func (cs *SocketStateCache) GetOperatorSockets(ctx context.Context, operators []core.OperatorID) (core.OperatorSockets, error) {
	cs.socketMu.Lock()
	defer cs.socketMu.Unlock()

	if err := cs.refreshSocketMap(ctx, operators); err != nil {
		return core.OperatorSockets{}, fmt.Errorf("failed to index socket map: %w", err)
	}

	operatorSockets := make(map[core.OperatorID]*core.OperatorSocket)
	for operatorID, socket := range cs.SocketMap {
		operatorSockets[operatorID] = socket
	}
	return operatorSockets, nil
}

func convertSocketStringToOperatorSocket(socket string) (core.OperatorSocket, error) {
	host, v1Dispersal, v1Retrieval, v2Dispersal, v2Retrieval, err := core.ParseOperatorSocket(socket)
	if err != nil {
		return "", fmt.Errorf("failed to parse operator socket: %w", err)
	}
	return core.MakeOperatorSocket(host, v1Dispersal, v1Retrieval, v2Dispersal, v2Retrieval), nil
}
