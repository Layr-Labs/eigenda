package thegraph_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ethcomm "github.com/ethereum/go-ethereum/common"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

var (
	quorums = []core.QuorumID{0}
)

type mockGraphQLQuerier struct {
	QueryFn func(ctx context.Context, q any, variables map[string]any) error
}

func (m mockGraphQLQuerier) Query(ctx context.Context, q any, variables map[string]any) error {
	return m.QueryFn(ctx, q, variables)
}

type mockChainState struct {
	GetCurrentBlockNumberFn func() (uint, error)
}

func (m mockChainState) GetCurrentBlockNumber() (uint, error) {
	return m.GetCurrentBlockNumberFn()
}

func (m *mockChainState) GetOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.OperatorState, error) {
	return nil, nil
}

func (m *mockChainState) GetOperatorStateByOperator(ctx context.Context, blockNumber uint, operator core.OperatorID) (*core.OperatorState, error) {
	return nil, nil
}

func TestIndexedChainState_GetIndexedOperatorState(t *testing.T) {
	logger := logging.NewNoopLogger()
	operatorsQueryCalled := false
	querier := &mockGraphQLQuerier{}
	querier.QueryFn = func(ctx context.Context, q any, variables map[string]any) error {
		switch res := q.(type) {
		case *thegraph.QueryQuorumAPKGql:
			pubKey := thegraph.AggregatePubkeyKeyGql{
				Apk_X: "3829803941453902453085939595934570464887466392754984985219704448765546217155",
				Apk_Y: "7864472681234874546092094912246874347602747071877011905183009416740980374479",
			}
			res.QuorumAPK = append(res.QuorumAPK, pubKey)
			return nil
		case *thegraph.QueryOperatorsGql:
			if operatorsQueryCalled {
				return nil
			}
			res.Operators = []thegraph.IndexedOperatorInfoGql{
				{
					Id:         "0x3eb7d5df61c48ec2718d8c8ad52304effc970ae92f19138e032dae07b7c0d629",
					PubkeyG1_X: "3336192159512049190945679273141887248666932624338963482128432381981287252980",
					PubkeyG1_Y: "15195175002875833468883745675063986308012687914999552116603423331534089122704",
					PubkeyG2_X: []graphql.String{
						"21597023645215426396093421944506635812143308313031252511177204078669540440732",
						"11405255666568400552575831267661419473985517916677491029848981743882451844775",
					},
					PubkeyG2_Y: []graphql.String{
						"9416989242565286095121881312760798075882411191579108217086927390793923664442",
						"13612061731370453436662267863740141021994163834412349567410746669651828926551",
					},
					SocketUpdates: []thegraph.SocketUpdates{{Socket: "localhost:32006;32007"}},
				},
			}
			operatorsQueryCalled = true
			return nil
		default:
			return nil
		}
	}

	chainState := &mockChainState{}
	chainState.GetCurrentBlockNumberFn = func() (uint, error) {
		return 1, nil
	}

	cs := thegraph.NewIndexedChainState(chainState, querier, logger)
	err := cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	state, err := cs.GetIndexedOperatorState(context.Background(), headerNum, quorums)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(state.IndexedOperators))
}

func TestIndexedChainState_GetIndexedOperatorInfoByOperatorId(t *testing.T) {
	logger := logging.NewNoopLogger()

	chainState := &mockChainState{}
	chainState.GetCurrentBlockNumberFn = func() (uint, error) {
		return 1, nil
	}

	querier := &mockGraphQLQuerier{}
	querier.QueryFn = func(ctx context.Context, q any, variables map[string]any) error {
		switch res := q.(type) {
		case *thegraph.QueryOperatorByIdGql:
			res.Operator = thegraph.IndexedOperatorInfoGql{
				Id:         "0x3eb7d5df61c48ec2718d8c8ad52304effc970ae92f19138e032dae07b7c0d629",
				PubkeyG1_X: "3336192159512049190945679273141887248666932624338963482128432381981287252980",
				PubkeyG1_Y: "15195175002875833468883745675063986308012687914999552116603423331534089122704",
				PubkeyG2_X: []graphql.String{
					"21597023645215426396093421944506635812143308313031252511177204078669540440732",
					"11405255666568400552575831267661419473985517916677491029848981743882451844775",
				},
				PubkeyG2_Y: []graphql.String{
					"9416989242565286095121881312760798075882411191579108217086927390793923664442",
					"13612061731370453436662267863740141021994163834412349567410746669651828926551",
				},
				SocketUpdates: []thegraph.SocketUpdates{{Socket: "localhost:32006;32007"}},
			}
			return nil
		default:
			return nil
		}
	}

	cs := thegraph.NewIndexedChainState(chainState, querier, logger)
	err := cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	opID := ethcomm.HexToHash("0x3eb7d5df61c48ec2718d8c8ad52304effc970ae92f19138e032dae07b7c0d629")
	info, err := cs.GetIndexedOperatorInfoByOperatorId(context.Background(), core.OperatorID(opID.Bytes()), uint32(headerNum))
	assert.NoError(t, err)
	assert.Equal(t, "3336192159512049190945679273141887248666932624338963482128432381981287252980", info.PubkeyG1.X.String())
	assert.Equal(t, "15195175002875833468883745675063986308012687914999552116603423331534089122704", info.PubkeyG1.Y.String())
}
