package thegraph_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ethcomm "github.com/ethereum/go-ethereum/common"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

var (
	quorums = []corev2.QuorumID{0}
)

type mockGraphQLQuerier struct {
	QueryFn func(ctx context.Context, q any, variables map[string]any) error
}

func (m mockGraphQLQuerier) Query(ctx context.Context, q any, variables map[string]any) error {
	return m.QueryFn(ctx, q, variables)
}

func TestIndexedChainState_GetIndexedOperatorState(t *testing.T) {
	logger := logging.NewNoopLogger()

	chainState, _ := mock.MakeChainDataMock(map[uint8]int{
		0: 1,
		1: 1,
		2: 1,
	})
	chainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	state, err := chainState.GetOperatorState(context.Background(), 1, quorums)
	assert.NoError(t, err)
	id := ""
	for key := range state.Operators[0] {
		id = key.Hex()
	}

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
					Id:         graphql.String(id),
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

	cs := thegraph.NewIndexedChainState(chainState, querier, logger)
	err = cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	indexedState, err := cs.GetIndexedOperatorState(context.Background(), headerNum, quorums)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(indexedState.IndexedOperators))
}

func TestIndexedChainState_GetIndexedOperatorStateMissingOperator(t *testing.T) {
	logger := logging.NewNoopLogger()

	chainState, _ := mock.MakeChainDataMock(map[uint8]int{
		0: 2,
		1: 2,
		2: 2,
	})
	chainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	state, err := chainState.GetOperatorState(context.Background(), 1, quorums)
	assert.NoError(t, err)
	id := ""
	for key := range state.Operators[0] {
		id = key.Hex()
		break
	}

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
					Id:         graphql.String(id),
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

	cs := thegraph.NewIndexedChainState(chainState, querier, logger)
	err = cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	_, err = cs.GetIndexedOperatorState(context.Background(), headerNum, quorums)
	assert.ErrorContains(t, err, "not found in indexed state")
}

func TestIndexedChainState_GetIndexedOperatorStateExtraOperator(t *testing.T) {
	logger := logging.NewNoopLogger()

	chainState, _ := mock.MakeChainDataMock(map[uint8]int{
		0: 1,
		1: 1,
		2: 1,
	})
	chainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	state, err := chainState.GetOperatorState(context.Background(), 1, quorums)
	assert.NoError(t, err)
	id := ""
	for key := range state.Operators[0] {
		id = key.Hex()
		break
	}

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
					Id:         graphql.String(id),
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

	cs := thegraph.NewIndexedChainState(chainState, querier, logger)
	err = cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	indexedState, err := cs.GetIndexedOperatorState(context.Background(), headerNum, quorums)
	assert.NoError(t, err)
	assert.Len(t, indexedState.IndexedOperators, 1)

}

func TestIndexedChainState_GetIndexedOperatorInfoByOperatorId(t *testing.T) {
	logger := logging.NewNoopLogger()

	chainState, _ := mock.MakeChainDataMock(map[uint8]int{
		0: 1,
		1: 1,
		2: 1,
	})
	chainState.On("GetCurrentBlockNumber").Return(uint(1), nil)

	state, err := chainState.GetOperatorState(context.Background(), 1, quorums)
	assert.NoError(t, err)
	id := ""
	for key := range state.Operators[0] {
		id = key.Hex()
	}

	querier := &mockGraphQLQuerier{}
	querier.QueryFn = func(ctx context.Context, q any, variables map[string]any) error {
		switch res := q.(type) {
		case *thegraph.QueryOperatorByIdGql:
			res.Operator = thegraph.IndexedOperatorInfoGql{
				Id:         graphql.String(id),
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
	err = cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	opID := ethcomm.HexToHash(id)
	info, err := cs.GetIndexedOperatorInfoByOperatorId(context.Background(), corev2.OperatorID(opID.Bytes()), uint32(headerNum))
	assert.NoError(t, err)
	assert.Equal(t, "3336192159512049190945679273141887248666932624338963482128432381981287252980", info.PubkeyG1.X.String())
	assert.Equal(t, "15195175002875833468883745675063986308012687914999552116603423331534089122704", info.PubkeyG1.Y.String())
}
