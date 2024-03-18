package eth

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	head "github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// block is finalized if its distance from HEAD is greater than some configurable number.
const DistanceFromHead = 100

type HeaderService struct {
	rpcEthClient common.RPCEthClient
	logger       logging.Logger
}

func NewHeaderService(logger logging.Logger, rpcEthClient common.RPCEthClient) *HeaderService {
	return &HeaderService{logger: logger, rpcEthClient: rpcEthClient}
}

// GetHeaders returns a list of new headers since the indicated header.
func (h *HeaderService) PullNewHeaders(lastHeader *head.Header) (head.Headers, bool, error) {
	ctx := context.Background()
	latestHeader, err := h.getHeaderByNumber(ctx, nil)
	if err != nil {
		h.logger.Error("Error. Cannot get latest header:", "err", err)
		return nil, false, err
	}

	lastHeaderNum := lastHeader.Number
	latestHeaderNum := latestHeader.Number.Uint64()

	if latestHeaderNum == lastHeaderNum {
		return []*head.Header{lastHeader}, true, nil
	}

	starting := lastHeaderNum + 1
	count := latestHeaderNum - starting + 1

	newHeaders, err := h.headersByRange(ctx, starting, int(count))
	if err != nil {
		h.logger.Error("Error. Cannot get latest header: ", "err", err)
		return nil, false, err
	}

	headers := make(head.Headers, 0, len(newHeaders))
	for _, header := range newHeaders {
		headerNum := header.Number.Uint64()
		finalized := latestHeaderNum-headerNum > DistanceFromHead

		headers = append(headers, &head.Header{
			BlockHash:     header.Hash(),
			PrevBlockHash: header.ParentHash,
			Number:        headerNum,
			Finalized:     finalized,
			CurrentFork:   "",
			IsUpgrade:     false,
		})
	}
	return headers, false, nil
}

// PullLatestHeader gets the latest header from the chain client
func (h *HeaderService) PullLatestHeader(finalized bool) (*head.Header, error) {
	ctx := context.Background()

	header, err := h.getHeaderByNumber(ctx, nil)
	if err != nil {
		h.logger.Error("Error. Cannot get latest header", "err", err)
		return nil, err
	}

	diff := header.Number.Int64() - DistanceFromHead
	if finalized && diff >= DistanceFromHead {
		latestFinalized, err := h.getHeaderByNumber(ctx, big.NewInt(diff))
		if err != nil {
			h.logger.Error("Error. Cannot get finalized header", "err", err)
			return nil, err
		}
		return &head.Header{
			BlockHash:     latestFinalized.Hash(),
			PrevBlockHash: latestFinalized.ParentHash,
			Number:        latestFinalized.Number.Uint64(),
			Finalized:     true,
			CurrentFork:   "",
			IsUpgrade:     false,
		}, nil
	}

	return &head.Header{
		BlockHash:     header.Hash(),
		PrevBlockHash: header.ParentHash,
		Number:        header.Number.Uint64(),
		Finalized:     false,
		CurrentFork:   "",
		IsUpgrade:     false,
	}, nil
}

func (h *HeaderService) headersByRange(ctx context.Context, startHeight uint64, count int) ([]*types.Header, error) {
	height := startHeight
	batchElems := make([]rpc.BatchElem, count)
	for i := 0; i < count; i++ {
		batchElems[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args: []interface{}{
				toBlockNumArg(new(big.Int).SetUint64(height + uint64(i))),
				false,
			},
			Result: new(types.Header),
			Error:  nil,
		}
	}

	if err := h.rpcEthClient.BatchCallContext(ctx, batchElems); err != nil {
		return nil, err
	}

	out := make([]*types.Header, count)
	for i := 0; i < len(batchElems); i++ {
		if batchElems[i].Error != nil {
			return nil, batchElems[i].Error
		}
		out[i] = batchElems[i].Result.(*types.Header)
	}

	return out, nil
}

func (h *HeaderService) getHeaderByNumber(ctx context.Context, number *big.Int) (types.Header, error) {
	var header = types.Header{}
	if err := h.rpcEthClient.CallContext(ctx, &header, "eth_getBlockByNumber", toBlockNumArg(number), false); err != nil {
		return types.Header{}, err
	}
	return header, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
