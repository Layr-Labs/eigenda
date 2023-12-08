package indexer

import (
	"fmt"

	dacommon "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/indexer"
	indexereth "github.com/Layr-Labs/eigenda/indexer/eth"
	inmemstore "github.com/Layr-Labs/eigenda/indexer/inmem"
	"github.com/ethereum/go-ethereum/common"
)

func SetupNewIndexer(
	config *indexer.Config,
	gethClient dacommon.EthClient,
	rpcClient dacommon.RPCEthClient,
	eigenDAServiceManagerAddr string,
	logger dacommon.Logger,
) (indexer.Indexer, error) {

	eigenDAServiceManager := common.HexToAddress(eigenDAServiceManagerAddr)

	pubKeyFilterer, err := NewOperatorPubKeysFilterer(eigenDAServiceManager, gethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create new operator pubkeys filter: %w", err)
	}

	socketsFilterer, err := NewOperatorSocketsFilterer(eigenDAServiceManager, gethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create new operator sockets filter: %w", err)
	}

	handlers := []indexer.AccumulatorHandler{
		{
			Acc:      NewOperatorPubKeysAccumulator(logger),
			Filterer: pubKeyFilterer,
			Status:   indexer.Good,
		},
		{
			Acc:      NewOperatorSocketsAccumulator(logger),
			Filterer: socketsFilterer,
			Status:   indexer.Good,
		},
	}

	var (
		upgrader    = &Upgrader{}
		headerStore = inmemstore.NewHeaderStore()
		headerSrvc  = indexereth.NewHeaderService(logger, rpcClient)
	)
	return indexer.New(
		config,
		handlers,
		headerSrvc,
		headerStore,
		upgrader,
		logger,
	), nil
}
