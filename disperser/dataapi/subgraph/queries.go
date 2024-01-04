package subgraph

import (
	"github.com/shurcooL/graphql"
)

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
	Operator struct {
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
	SocketUpdates struct {
		Socket graphql.String
	}
	IndexedOperatorInfo struct {
		Id         graphql.String
		PubkeyG1_X graphql.String   `graphql:"pubkeyG1_X"`
		PubkeyG1_Y graphql.String   `graphql:"pubkeyG1_Y"`
		PubkeyG2_X []graphql.String `graphql:"pubkeyG2_X"`
		PubkeyG2_Y []graphql.String `graphql:"pubkeyG2_Y"`
		// Socket is the socket address of the operator, in the form "host:port"
		SocketUpdates []SocketUpdates `graphql:"socketUpdates(first: 1, orderBy: blockNumber, orderDirection: desc)"`
	}
	queryBatches struct {
		Batches []*Batches `graphql:"batches(orderDirection: $orderDirection, orderBy: $orderBy, first: $first, skip: $skip)"`
	}
	queryBatchesByBlockTimestampRange struct {
		Batches []*Batches `graphql:"batches(first: $first, orderBy: blockTimestamp, where: {and: [{ blockTimestamp_gte: $blockTimestamp_gte}, {blockTimestamp_lte: $blockTimestamp_lte}]})"`
	}
	queryOperatorRegistereds struct {
		OperatorRegistereds []*Operator `graphql:"operatorRegistereds(first: $first)"`
	}
	queryBatchNonSigningOperatorIdsInInterval struct {
		BatchNonSigningOperatorIds []*BatchNonSigningOperatorIds `graphql:"batches(first: $first, skip: $skip, where: {blockTimestamp_gt: $blockTimestamp_gt})"`
	}
	queryOperatorRegisteredsGTBlockTimestamp struct {
		OperatorRegistereds []*Operator `graphql:"operatorRegistereds(orderBy: blockTimestamp, where: {blockTimestamp_gt: $blockTimestamp_gt})"`
	}
	queryOperatorDeregisteredsGTBlockTimestamp struct {
		OperatorDeregistereds []*Operator `graphql:"operatorDeregistereds(orderBy: blockTimestamp, where: {blockTimestamp_gt: $blockTimestamp_gt})"`
	}
	queryOperatorById struct {
		Operator IndexedOperatorInfo `graphql:"operator(id: $id)"`
	}
)
