package v1

import (
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigenda/encoding"
)

// Response Types
type ErrorResponse struct {
	Error string `json:"error"`
}

type BlobMetadataResponse struct {
	BlobKey                 string                    `json:"blob_key"`
	BatchHeaderHash         string                    `json:"batch_header_hash"`
	BlobIndex               uint32                    `json:"blob_index"`
	SignatoryRecordHash     string                    `json:"signatory_record_hash"`
	ReferenceBlockNumber    uint32                    `json:"reference_block_number"`
	BatchRoot               string                    `json:"batch_root"`
	BlobInclusionProof      string                    `json:"blob_inclusion_proof"`
	BlobCommitment          *encoding.BlobCommitments `json:"blob_commitment"`
	BatchId                 uint32                    `json:"batch_id"`
	ConfirmationBlockNumber uint32                    `json:"confirmation_block_number"`
	ConfirmationTxnHash     string                    `json:"confirmation_txn_hash"`
	Fee                     string                    `json:"fee"`
	SecurityParams          []*core.SecurityParam     `json:"security_params"`
	RequestAt               uint64                    `json:"requested_at"`
	BlobStatus              disperser.BlobStatus      `json:"blob_status"`
}

type BlobsResponse struct {
	Meta Meta                    `json:"meta"`
	Data []*BlobMetadataResponse `json:"data"`
}

type Meta struct {
	Size      int    `json:"size"`
	NextToken string `json:"next_token,omitempty"`
}

type OperatorNonsigningPercentageMetrics struct {
	OperatorId           string  `json:"operator_id"`
	OperatorAddress      string  `json:"operator_address"`
	QuorumId             uint8   `json:"quorum_id"`
	TotalUnsignedBatches int     `json:"total_unsigned_batches"`
	TotalBatches         int     `json:"total_batches"`
	Percentage           float64 `json:"percentage"`
	StakePercentage      float64 `json:"stake_percentage"`
}

type OperatorsNonsigningPercentage struct {
	Meta Meta                                   `json:"meta"`
	Data []*OperatorNonsigningPercentageMetrics `json:"data"`
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

type QueriedStateOperatorMetadata struct {
	OperatorId           string `json:"operator_id"`
	BlockNumber          uint   `json:"block_number"`
	Socket               string `json:"socket"`
	IsOnline             bool   `json:"is_online"`
	OperatorProcessError string `json:"operator_process_error"`
}

type QueriedStateOperatorsResponse struct {
	Meta Meta                            `json:"meta"`
	Data []*QueriedStateOperatorMetadata `json:"data"`
}

type QueriedOperatorEjections struct {
	OperatorId      string  `json:"operator_id"`
	OperatorAddress string  `json:"operator_address"`
	Quorum          uint8   `json:"quorum"`
	BlockNumber     uint64  `json:"block_number"`
	BlockTimestamp  string  `json:"block_timestamp"`
	TransactionHash string  `json:"transaction_hash"`
	StakePercentage float64 `json:"stake_percentage"`
}

type QueriedOperatorEjectionsResponse struct {
	Ejections []*QueriedOperatorEjections `json:"ejections"`
}

type ServiceAvailability struct {
	ServiceName   string `json:"service_name"`
	ServiceStatus string `json:"service_status"`
}

type ServiceAvailabilityResponse struct {
	Meta Meta                   `json:"meta"`
	Data []*ServiceAvailability `json:"data"`
}

type OperatorPortCheckResponse struct {
	OperatorId      string `json:"operator_id"`
	DispersalSocket string `json:"dispersal_socket"`
	RetrievalSocket string `json:"retrieval_socket"`
	DispersalOnline bool   `json:"dispersal_online"`
	RetrievalOnline bool   `json:"retrieval_online"`
}

type SemverReportResponse struct {
	Semver map[string]*semver.SemverMetrics `json:"semver"`
}

type Metric struct {
	Throughput          float64                    `json:"throughput"`
	CostInGas           float64                    `json:"cost_in_gas"`
	TotalStake          *big.Int                   `json:"total_stake"`
	TotalStakePerQuorum map[core.QuorumID]*big.Int `json:"total_stake_per_quorum"`
}

type Throughput struct {
	Throughput float64 `json:"throughput"`
	Timestamp  uint64  `json:"timestamp"`
}

type NonSigner struct {
	OperatorId      string  `json:"operator_id"`
	OperatorAddress string  `json:"operator_address"`
	QuorumId        uint8   `json:"quorum_id"`
	ErrorRate       float64 `json:"error_rate"`
	TotalBatches    int     `json:"total_batches"`
}

//	@Summary	Fetch blob metadata by blob key
//	@Tags		Feed
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob Key"
//	@Success	200			{object}	BlobMetadataResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/feed/blobs/{blob_key} [get]
func FetchBlob() {}

//	@Summary	Fetch blob metadata by batch header hash
//	@Tags		Feed
//	@Produce	json
//	@Param		batch_header_hash	path		string	true	"Batch Header Hash"
//	@Param		limit				query		int		false	"Limit [default: 10]"
//	@Param		next_token			query		string	false	"Next page token"
//	@Success	200					{object}	BlobsResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/feed/batches/{batch_header_hash}/blobs [get]
func FetchBlobsFromBatchHeaderHash() {}

//	@Summary	Fetch blobs metadata list
//	@Tags		Feed
//	@Produce	json
//	@Param		limit	query		int	false	"Limit [default: 10]"
//	@Success	200		{object}	BlobsResponse
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/feed/blobs [get]
func FetchBlobs() {}

//	@Summary	Operator stake distribution query
//	@Tags		OperatorsStake
//	@Produce	json
//	@Param		operator_id	query		string	true	"Operator ID"
//	@Success	200			{object}	OperatorsStakeResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators-info/operators-stake [get]
func OperatorsStake() {}

//	@Summary	Fetch list of operators that have been deregistered for days
//	@Tags		OperatorsInfo
//	@Produce	json
//	@Success	200	{object}	QueriedStateOperatorsResponse
//	@Failure	400	{object}	ErrorResponse	"error: Bad request"
//	@Failure	404	{object}	ErrorResponse	"error: Not found"
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/operators-info/deregistered-operators [get]
func FetchDeregisteredOperators() {}

//	@Summary	Fetch list of operators that have been registered for days
//	@Tags		OperatorsInfo
//	@Produce	json
//	@Success	200	{object}	QueriedStateOperatorsResponse
//	@Failure	400	{object}	ErrorResponse	"error: Bad request"
//	@Failure	404	{object}	ErrorResponse	"error: Not found"
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/operators-info/registered-operators [get]
func FetchRegisteredOperators() {}

//	@Summary	Fetch list of operator ejections over last N days
//	@Tags		OperatorsInfo
//	@Produce	json
//	@Param		days		query		int		false	"Lookback in days [default: 1]"
//	@Param		operator_id	query		string	false	"Operator ID filter [default: all operators]"
//	@Param		first		query		int		false	"Return first N ejections [default: 1000]"
//	@Param		skip		query		int		false	"Skip first N ejections [default: 0]"
//	@Success	200			{object}	QueriedOperatorEjectionsResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators-info/operator-ejections [get]
func FetchOperatorEjections() {}

//	@Summary	Operator node reachability port check
//	@Tags		OperatorsInfo
//	@Produce	json
//	@Param		operator_id	query		string	true	"Operator ID"
//	@Success	200			{object}	OperatorPortCheckResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators-info/port-check [get]
func OperatorPortCheck() {}

//	@Summary	Active operator semver scan
//	@Tags		OperatorsInfo
//	@Produce	json
//	@Success	200	{object}	SemverReportResponse
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/operators-info/semver-scan [get]
func SemverScan() {}

//	@Summary	Fetch metrics
//	@Tags		Metrics
//	@Produce	json
//	@Param		start	query		int	false	"Start unix timestamp [default: 1 hour ago]"
//	@Param		end		query		int	false	"End unix timestamp [default: unix time now]"
//	@Param		limit	query		int	false	"Limit [default: 10]"
//	@Success	200		{object}	Metric
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics [get]
func FetchMetrics() {}

// FetchMetricsThroughputHandler godoc
//
//	@Summary	Fetch throughput time series
//	@Tags		Metrics
//	@Produce	json
//	@Param		start	query		int	false	"Start unix timestamp [default: 1 hour ago]"
//	@Param		end		query		int	false	"End unix timestamp [default: unix time now]"
//	@Success	200		{object}	[]Throughput
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/throughput  [get]
func FetchMetricsThroughputHandler() {}

// FetchNonSigners godoc
//
//	@Summary	Fetch non signers
//	@Tags		Metrics
//	@Produce	json
//	@Param		interval	query		int	false	"Interval to query for non signers in seconds [default: 3600]"
//	@Success	200			{object}	[]NonSigner
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/non-signers  [get]
func FetchNonSigners() {}

// FetchOperatorsNonsigningPercentageHandler godoc
//
//	@Summary	Fetch operators non signing percentage
//	@Tags		Metrics
//	@Produce	json
//	@Param		interval	query		int		false	"Interval to query for operators nonsigning percentage [default: 3600]"
//	@Param		end			query		string	false	"End time (2006-01-02T15:04:05Z) to query for operators nonsigning percentage [default: now]"
//	@Param		live_only	query		string	false	"Whether return only live nonsigners [default: true]"
//	@Success	200			{object}	OperatorsNonsigningPercentage
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/operator-nonsigning-percentage  [get]
func FetchOperatorsNonsigningPercentageHandler() {}

//	@Summary	Get status of EigenDA Disperser service
//	@Tags		ServiceAvailability
//	@Produce	json
//	@Success	200	{object}	ServiceAvailabilityResponse
//	@Failure	400	{object}	ErrorResponse	"error: Bad request"
//	@Failure	404	{object}	ErrorResponse	"error: Not found"
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/disperser-service-availability [get]
func FetchDisperserServiceAvailability() {}

//	@Summary	Get status of EigenDA churner service
//	@Tags		ServiceAvailability
//	@Produce	json
//	@Success	200	{object}	ServiceAvailabilityResponse
//	@Failure	400	{object}	ErrorResponse	"error: Bad request"
//	@Failure	404	{object}	ErrorResponse	"error: Not found"
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/churner-service-availability [get]
func FetchChurnerServiceAvailability() {}

//	@Summary	Get status of EigenDA batcher
//	@Tags		ServiceAvailability
//	@Produce	json
//	@Success	200	{object}	ServiceAvailabilityResponse
//	@Failure	400	{object}	ErrorResponse	"error: Bad request"
//	@Failure	404	{object}	ErrorResponse	"error: Not found"
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/batcher-service-availability [get]
func FetchBatcherAvailability() {}
