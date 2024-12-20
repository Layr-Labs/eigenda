package dataapi

import (
	"errors"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
)

var errNotFound = errors.New("not found")

type (
	ErrorResponse struct {
		Error string `json:"error"`
	}

	MetricSummary struct {
		AvgThroughput float64 `json:"avg_throughput"`
	}

	OperatorStake struct {
		QuorumId        string  `json:"quorum_id"`
		OperatorId      string  `json:"operator_id"`
		StakePercentage float64 `json:"stake_percentage"`
		Rank            int     `json:"rank"`
	}

	OperatorsStakeResponse struct {
		StakeRankedOperators map[string][]*OperatorStake `json:"stake_ranked_operators"`
	}

	OperatorPortCheckResponse struct {
		OperatorId      string `json:"operator_id"`
		DispersalSocket string `json:"dispersal_socket"`
		RetrievalSocket string `json:"retrieval_socket"`
		DispersalOnline bool   `json:"dispersal_online"`
		RetrievalOnline bool   `json:"retrieval_online"`
	}

	QueriedOperatorEjections struct {
		OperatorId      string  `json:"operator_id"`
		OperatorAddress string  `json:"operator_address"`
		Quorum          uint8   `json:"quorum"`
		BlockNumber     uint64  `json:"block_number"`
		BlockTimestamp  string  `json:"block_timestamp"`
		TransactionHash string  `json:"transaction_hash"`
		StakePercentage float64 `json:"stake_percentage"`
	}

	SemverReportResponse struct {
		Semver map[string]*semver.SemverMetrics `json:"semver"`
	}

	Metric struct {
		Throughput float64 `json:"throughput"`
		CostInGas  float64 `json:"cost_in_gas"`
		// deprecated: use TotalStakePerQuorum instead. Remove when the frontend is updated.
		TotalStake          *big.Int                   `json:"total_stake"`
		TotalStakePerQuorum map[core.QuorumID]*big.Int `json:"total_stake_per_quorum"`
	}

	Throughput struct {
		Throughput float64 `json:"throughput"`
		Timestamp  uint64  `json:"timestamp"`
	}

	Meta struct {
		Size      int    `json:"size"`
		NextToken string `json:"next_token,omitempty"`
	}
)
