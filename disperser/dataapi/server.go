package dataapi

import (
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
)

type (
	Throughput struct {
		Throughput float64 `json:"throughput"`
		Timestamp  uint64  `json:"timestamp"`
	}

	Meta struct {
		Size      int    `json:"size"`
		NextToken string `json:"next_token,omitempty"`
	}

	OperatorNonsigningPercentageMetrics struct {
		OperatorId           string  `json:"operator_id"`
		OperatorAddress      string  `json:"operator_address"`
		QuorumId             uint8   `json:"quorum_id"`
		TotalUnsignedBatches int     `json:"total_unsigned_batches"`
		TotalBatches         int     `json:"total_batches"`
		Percentage           float64 `json:"percentage"`
		StakePercentage      float64 `json:"stake_percentage"`
	}

	OperatorsNonsigningPercentage struct {
		Meta Meta                                   `json:"meta"`
		Data []*OperatorNonsigningPercentageMetrics `json:"data"`
	}

	OperatorStake struct {
		QuorumId        string  `json:"quorum_id"`
		OperatorId      string  `json:"operator_id"`
		OperatorAddress string  `json:"operator_address"`
		StakePercentage float64 `json:"stake_percentage"`
		Rank            int     `json:"rank"`
		StakeAmount     float64 `json:"stake_amount"`
	}

	OperatorsStakeResponse struct {
		CurrentBlock         uint32                      `json:"current_block"`
		StakeRankedOperators map[string][]*OperatorStake `json:"stake_ranked_operators"`
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

	OperatorLiveness struct {
		OperatorId      string `json:"operator_id"`
		DispersalSocket string `json:"dispersal_socket"`
		DispersalOnline bool   `json:"dispersal_online"`
		DispersalStatus string `json:"dispersal_status"`
		RetrievalSocket string `json:"retrieval_socket"`
		RetrievalOnline bool   `json:"retrieval_online"`
		RetrievalStatus string `json:"retrieval_status"`
	}

	OperatorPortCheckResponse struct {
		OperatorId      string `json:"operator_id"`
		DispersalSocket string `json:"dispersal_socket"`
		DispersalOnline bool   `json:"dispersal_online"`
		DispersalStatus string `json:"dispersal_status"`
		RetrievalSocket string `json:"retrieval_socket"`
		RetrievalOnline bool   `json:"retrieval_online"`
		RetrievalStatus string `json:"retrieval_status"`
	}
	SemverReportResponse struct {
		Semver map[string]*semver.SemverMetrics `json:"semver"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

type ServerInterface interface {
	Start() error
	Shutdown() error
}
