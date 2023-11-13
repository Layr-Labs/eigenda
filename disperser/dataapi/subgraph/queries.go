package subgraph

import "github.com/shurcooL/graphql"

type (
	Batches struct {
		Id              graphql.String
		BatchId         graphql.String
		BatchHeaderHash graphql.String
		BlockTimestamp  graphql.String
		BlockNumber     graphql.String
		TxHash          graphql.String
		GasFees         GasFees
	}
	GasFees struct {
		Id       graphql.String
		GasUsed  graphql.String
		GasPrice graphql.String
		TxFee    graphql.String
	}
	OperatorRegistered struct {
		Id              graphql.String
		OperatorId      graphql.String
		Operator        graphql.String
		BlockTimestamp  graphql.String
		BlockNumber     graphql.String
		TransactionHash graphql.String
	}
	BatchNonSigningOperatorIds struct {
		NonSigning struct {
			NonSigners []struct {
				OperatorId graphql.String `graphql:"operatorId"`
			} `graphql:"nonSigners"`
		} `graphql:"nonSigning"`
	}
	queryBatches struct {
		Batches []*Batches `graphql:"batches(orderDirection: $orderDirection, orderBy: $orderBy, first: $first, skip: $skip)"`
	}
	queryOperatorRegistereds struct {
		OperatorRegistereds []*OperatorRegistered `graphql:"operatorRegistereds(first: $first)"`
	}
	queryBatchNonSigningOperatorIdsInInterval struct {
		BatchNonSigningOperatorIds []*BatchNonSigningOperatorIds `graphql:"batches(first: $first, skip: $skip, where: {blockTimestamp_gt: $blockTimestamp_gt})"`
	}
)
