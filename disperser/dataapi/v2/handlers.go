package v2

import (
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
)

// Response Types
type ErrorResponse struct {
	Error string `json:"error"`
}

type BlobResponse struct {
	BlobHeader    *corev2.BlobHeader `json:"blob_header"`
	Status        string             `json:"status"`
	DispersedAt   uint64             `json:"dispersed_at"`
	BlobSizeBytes uint64             `json:"blob_size_bytes"`
}

type SignedBatch struct {
	BatchHeader *corev2.BatchHeader `json:"batch_header"`
	Attestation *corev2.Attestation `json:"attestation"`
}

type BatchResponse struct {
	BatchHeaderHash       string                         `json:"batch_header_hash"`
	SignedBatch           *SignedBatch                   `json:"signed_batch"`
	BlobVerificationInfos []*corev2.BlobVerificationInfo `json:"blob_verification_infos"`
}

type OperatorStake struct {
	QuorumId        string  `json:"quorum_id"`
	OperatorId      string  `json:"operator_id"`
	StakePercentage float64 `json:"stake_percentage"`
	Rank            int     `json:"rank"`
}

type OperatorsStakeResponse struct {
	StakeRankedOperators map[string][]*OperatorStake `json:"stake_ranked_operators"`
}

type SemverReportResponse struct {
	Semver map[string]*semver.SemverMetrics `json:"semver"`
}

type OperatorPortCheckResponse struct {
	OperatorId      string `json:"operator_id"`
	DispersalSocket string `json:"dispersal_socket"`
	RetrievalSocket string `json:"retrieval_socket"`
	DispersalOnline bool   `json:"dispersal_online"`
	RetrievalOnline bool   `json:"retrieval_online"`
}

// FetchBlobHandler godoc
//
//	@Summary	Fetch blob metadata by blob key
//	@Tags		Feed
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob key in hex string"
//	@Success	200			{object}	BlobResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blob/{blob_key} [get]
func FetchBlobHandler() {}

// FetchBatchHandler godoc
//
//	@Summary	Fetch batch by the batch header hash
//	@Tags		Feed
//	@Produce	json
//	@Param		batch_header_hash	path		string	true	"Batch header hash in hex string"
//	@Success	200					{object}	BlobResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/batch/{batch_header_hash} [get]
func FetchBatchHandler() {}

// FetchOperatorsStake godoc
//
//	@Summary	Operator stake distribution query
//	@Tags		Operators
//	@Produce	json
//	@Param		operator_id	query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200			{object}	OperatorsStakeResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/stake [get]
func FetchOperatorsStake() {}

// FetchOperatorsNodeInfo godoc
//
//	@Summary	Active operator semver
//	@Tags		Operators
//	@Produce	json
//	@Success	200	{object}	SemverReportResponse
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/nodeinfo [get]
func FetchOperatorsNodeInfo() {}

// CheckOperatorsReachability godoc
//
//	@Summary	Operator node reachability check
//	@Tags		Operators
//	@Produce	json
//	@Param		operator_id	query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200			{object}	OperatorPortCheckResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/reachability [get]
func CheckOperatorsReachability() {}
