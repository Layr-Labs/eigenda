package dataapi_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	subgraphmock "github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph/mock"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

var (
	subgraphOperatorRegistereds = []*subgraph.Operator{
		{
			Id:              "0x000763fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211",
			OperatorId:      "0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311",
			Operator:        "0x000563fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211",
			BlockTimestamp:  "1696975449000000000",
			BlockNumber:     "87",
			TransactionHash: "0x000163fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211",
		},
		{
			Id:              "0x000763fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f212",
			OperatorId:      "0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568310",
			Operator:        "0x000563fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f212",
			BlockTimestamp:  "1696975459000000000",
			BlockNumber:     "88",
			TransactionHash: "0x000163fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f212",
		},
	}

	subgraphOperatorDeregistereds = []*subgraph.Operator{
		{
			Id:              "0x000763fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f222",
			OperatorId:      "0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311",
			Operator:        "0x000223fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211",
			BlockTimestamp:  "1702666046",
			BlockNumber:     "22",
			TransactionHash: "0x000223fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211",
		},
	}

	subgraphBatches = []*subgraph.Batches{
		{
			Id:              "0x000763fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f207",
			BatchId:         "1",
			BatchHeaderHash: "0x890588400acb4f9f7f438c0d21734acb36a6c4c75df6560827e23b452bbdcc69",
			BlockTimestamp:  "1696975449000000000",
			BlockNumber:     "87",
			TxHash:          "0x63fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f207",
			GasFees: subgraph.GasFees{
				Id:       "0x0006afd9ce41ba0f3414ba2650a9cd2f47c0e22af21651f7fd902f71df678c5d9942",
				GasPrice: "1000045336",
				GasUsed:  "249815",
				TxFee:    "249826325612840",
			},
		},
		{
			Id:              "0x0007c601ff50ae500ec114a4430c1af872b14488a447f378c5c64adc36476e1101e1",
			BatchId:         "0",
			BatchHeaderHash: "0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568310",
			BlockTimestamp:  "1696975448000000000",
			BlockNumber:     "86",
			TxHash:          "0xc601ff50ae500ec114a4430c1af872b14488a447f378c5c64adc36476e1101e1",
			GasFees: subgraph.GasFees{
				Id:       "0x0006afd9ce41ba0f3414ba2650a9cd2f47c0e22af21651f7fd902f71df678c5d9942",
				GasPrice: "1000045336",
				GasUsed:  "249815",
				TxFee:    "249826325612840",
			},
		},
		{
			Id:              "0x0007de6f42234e643c6b427c349778cb41418f590ba899ac079c24427369d9c029aa",
			BatchId:         "2",
			BatchHeaderHash: "0x46c57a96296eb1b1d23f72b9ce3b2252fc5e2534c3008f5ce5e2afb06487a5eb",
			BlockTimestamp:  "1696975450000000000",
			BlockNumber:     "88",
			TxHash:          "0xde6f42234e643c6b427c349778cb41418f590ba899ac079c24427369d9c029aa",
			GasFees: subgraph.GasFees{
				Id:       "0x0006afd9ce41ba0f3414ba2650a9cd2f47c0e22af21651f7fd902f71df678c5d9942",
				GasPrice: "1000045336",
				GasUsed:  "249815",
				TxFee:    "249826325612840",
			},
		},
	}

	subgraphIndexedOperatorInfos = &subgraph.IndexedOperatorInfo{
		Id:         "0x000223fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f222",
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
		SocketUpdates: []subgraph.SocketUpdates{
			{
				Socket: "localhost:32006;32007",
			},
		},
	}
)

func TestQueryBatchesWithLimit(t *testing.T) {
	mockSubgraphApi := &subgraphmock.MockSubgraphApi{}
	subgraphClient := dataapi.NewSubgraphClient(mockSubgraphApi)
	mockSubgraphApi.On("QueryBatches").Return(subgraphBatches, nil)
	batches, err := subgraphClient.QueryBatchesWithLimit(context.Background(), 2, 0)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(batches))

	assert.Equal(t, []byte("0x0007de6f42234e643c6b427c349778cb41418f590ba899ac079c24427369d9c029aa"), batches[0].Id)
	assert.Equal(t, uint64(2), batches[0].BatchId)
	assert.Equal(t, []byte("0x46c57a96296eb1b1d23f72b9ce3b2252fc5e2534c3008f5ce5e2afb06487a5eb"), batches[0].BatchHeaderHash)
	assert.Equal(t, uint64(1696975450), batches[0].BlockTimestamp)
	assert.Equal(t, uint64(88), batches[0].BlockNumber)
	assert.Equal(t, []byte("0xde6f42234e643c6b427c349778cb41418f590ba899ac079c24427369d9c029aa"), batches[0].TxHash)
	assertGasFees(t, batches[0].GasFees)

	assert.Equal(t, []byte("0x000763fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f207"), batches[1].Id)
	assert.Equal(t, uint64(1), batches[1].BatchId)
	assert.Equal(t, []byte("0x890588400acb4f9f7f438c0d21734acb36a6c4c75df6560827e23b452bbdcc69"), batches[1].BatchHeaderHash)
	assert.Equal(t, uint64(1696975449), batches[1].BlockTimestamp)
	assert.Equal(t, uint64(87), batches[1].BlockNumber)
	assert.Equal(t, []byte("0x63fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f207"), batches[1].TxHash)
	assertGasFees(t, batches[1].GasFees)
}

func TestQueryOperators(t *testing.T) {
	mockSubgraphApi := &subgraphmock.MockSubgraphApi{}
	mockSubgraphApi.On("QueryOperators").Return(subgraphOperatorRegistereds, nil)
	subgraphClient := dataapi.NewSubgraphClient(mockSubgraphApi)
	operators, err := subgraphClient.QueryOperatorsWithLimit(context.Background(), 2)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(operators))

	assert.NotNil(t, operators[0])
	assert.Equal(t, []byte("0x000763fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211"), operators[0].Id)
	assert.Equal(t, []byte("0x000563fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211"), operators[0].Operator)
	assert.Equal(t, []byte("0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311"), operators[0].OperatorId)
	assert.Equal(t, uint64(1696975449), operators[0].BlockTimestamp)
	assert.Equal(t, uint64(87), operators[0].BlockNumber)
	assert.Equal(t, []byte("0x000163fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211"), operators[0].TransactionHash)

	assert.NotNil(t, operators[1])
	assert.Equal(t, []byte("0x000763fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f212"), operators[1].Id)
	assert.Equal(t, []byte("0x000563fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f212"), operators[1].Operator)
	assert.Equal(t, []byte("0xe1cdae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568310"), operators[1].OperatorId)
	assert.Equal(t, uint64(1696975459), operators[1].BlockTimestamp)
	assert.Equal(t, uint64(88), operators[1].BlockNumber)
	assert.Equal(t, []byte("0x000163fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f212"), operators[1].TransactionHash)
}

func TestQueryIndexedDeregisteredOperatorsInTheLast14Days(t *testing.T) {
	mockSubgraphApi := &subgraphmock.MockSubgraphApi{}
	mockSubgraphApi.On("QueryDeregisteredOperatorsGreaterThanBlockTimestamp").Return(subgraphOperatorDeregistereds, nil)
	mockSubgraphApi.On("QueryOperatorInfoByOperatorIdAtBlockNumber").Return(subgraphIndexedOperatorInfos, nil)
	subgraphClient := dataapi.NewSubgraphClient(mockSubgraphApi)
	indexedDeregisteredOperatorState, err := subgraphClient.QueryIndexedDeregisteredOperatorsInTheLast14Days(context.Background())
	assert.NoError(t, err)

	operators := indexedDeregisteredOperatorState.Operators
	assert.Equal(t, 1, len(operators))

	var operatorId [32]byte
	copy(operatorId[:], []byte("0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311"))
	operator := operators[operatorId]

	assert.NotNil(t, operator)

	expectedIndexedOperatorInfo, err := dataapi.ConvertOperatorInfoGqlToIndexedOperatorInfo(subgraphIndexedOperatorInfos)
	assert.NoError(t, err)

	assert.Equal(t, expectedIndexedOperatorInfo.PubkeyG1, operator.PubkeyG1)
	assert.Equal(t, expectedIndexedOperatorInfo.PubkeyG2, operator.PubkeyG2)
	assert.Equal(t, "localhost:32006;32007", operator.Socket)
	assert.Equal(t, uint64(22), uint64(operator.BlockNumber))
	assert.Equal(t, []byte("0xe22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311"), operator.Metadata.OperatorId)
	assert.Equal(t, []byte("0x000223fb86a79eda47c891d8826474d80b6a935ad2a2b5de921933e05c67f320f211"), operator.Metadata.TransactionHash)
	assert.Equal(t, uint64(22), uint64(operator.Metadata.BlockNumber))
}

func assertGasFees(t *testing.T, gasFees *dataapi.GasFees) {
	assert.NotNil(t, gasFees)
	assert.Equal(t, []byte("0x0006afd9ce41ba0f3414ba2650a9cd2f47c0e22af21651f7fd902f71df678c5d9942"), gasFees.Id)
	assert.Equal(t, uint64(249815), gasFees.GasUsed)
	assert.Equal(t, uint64(1000045336), gasFees.GasPrice)
	assert.Equal(t, uint64(249826325612840), gasFees.TxFee)
}
