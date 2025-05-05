package v2

import (
	"encoding/hex"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	disperserv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

// Base types shared acorss various API response types
type (
	OperatorIdentity struct {
		OperatorId      string `json:"operator_id"`
		OperatorAddress string `json:"operator_address"`
	}

	AttestationInfo struct {
		Attestation *corev2.Attestation          `json:"attestation"`
		Nonsigners  map[uint8][]OperatorIdentity `json:"nonsigners"`
		Signers     map[uint8][]OperatorIdentity `json:"signers"`
	}

	BatchHeader struct {
		BatchRoot            string `json:"batch_root"`
		ReferenceBlockNumber uint64 `json:"reference_block_number"`
	}

	BlobInclusionInfo struct {
		BatchHeader    *BatchHeader `json:"batch_header"`
		BlobKey        string       `json:"blob_key"`
		BlobIndex      uint32       `json:"blob_index"`
		InclusionProof string       `json:"inclusion_proof"`
	}

	BlobMetadata struct {
		BlobHeader    *corev2.BlobHeader `json:"blob_header"`
		Signature     string             `json:"signature"`
		BlobStatus    string             `json:"blob_status"`
		BlobSizeBytes uint64             `json:"blob_size_bytes"`
		RequestedAt   uint64             `json:"requested_at"`
		ExpiryUnixSec uint64             `json:"expiry_unix_sec"`
	}
)

// Operator types
type (
	OperatorDispersal struct {
		BatchHeaderHash string       `json:"batch_header_hash"`
		BatchHeader     *BatchHeader `json:"batch_header"`
		DispersedAt     uint64       `json:"dispersed_at"`
		Signature       string       `json:"signature"`
	}
	OperatorDispersalFeedResponse struct {
		OperatorIdentity OperatorIdentity     `json:"operator_identity"`
		OperatorSocket   string               `json:"operator_socket"`
		Dispersals       []*OperatorDispersal `json:"dispersals"`
	}

	OperatorSigningInfo struct {
		OperatorId              string  `json:"operator_id"`
		OperatorAddress         string  `json:"operator_address"`
		QuorumId                uint8   `json:"quorum_id"`
		TotalUnsignedBatches    int     `json:"total_unsigned_batches"`
		TotalResponsibleBatches int     `json:"total_responsible_batches"`
		TotalBatches            int     `json:"total_batches"`
		SigningPercentage       float64 `json:"signing_percentage"`
		StakePercentage         float64 `json:"stake_percentage"`
	}
	OperatorsSigningInfoResponse struct {
		StartBlock          uint32                 `json:"start_block"`
		EndBlock            uint32                 `json:"end_block"`
		StartTimeUnixSec    int64                  `json:"start_time_unix_sec"`
		EndTimeUnixSec      int64                  `json:"end_time_unix_sec"`
		OperatorSigningInfo []*OperatorSigningInfo `json:"operator_signing_info"`
	}

	OperatorStake struct {
		QuorumId        string  `json:"quorum_id"`
		OperatorId      string  `json:"operator_id"`
		OperatorAddress string  `json:"operator_address"`
		StakePercentage float64 `json:"stake_percentage"`
		Rank            int     `json:"rank"`
	}
	OperatorsStakeResponse struct {
		CurrentBlock         uint32                      `json:"current_block"`
		StakeRankedOperators map[string][]*OperatorStake `json:"stake_ranked_operators"`
	}

	OperatorDispersalResponse struct {
		Response *corev2.DispersalResponse `json:"operator_dispersal_response"`
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
	OperatorLivenessResponse struct {
		Operators []*OperatorLiveness `json:"operators"`
	}

	SemverReportResponse struct {
		Semver map[string]*semver.SemverMetrics `json:"semver"`
	}
)

// Blob types
type (
	BlobResponse struct {
		BlobKey       string             `json:"blob_key"`
		BlobHeader    *corev2.BlobHeader `json:"blob_header"`
		Status        string             `json:"status"`
		DispersedAt   uint64             `json:"dispersed_at"`
		BlobSizeBytes uint64             `json:"blob_size_bytes"`
	}

	BlobCertificateResponse struct {
		Certificate *corev2.BlobCertificate `json:"blob_certificate"`
	}

	BlobAttestationInfoResponse struct {
		BlobKey         string             `json:"blob_key"`
		BatchHeaderHash string             `json:"batch_header_hash"`
		InclusionInfo   *BlobInclusionInfo `json:"blob_inclusion_info"`
		AttestationInfo *AttestationInfo   `json:"attestation_info"`
	}

	BlobInfo struct {
		BlobKey      string        `json:"blob_key"`
		BlobMetadata *BlobMetadata `json:"blob_metadata"`
	}
	BlobFeedResponse struct {
		Blobs  []BlobInfo `json:"blobs"`
		Cursor string     `json:"cursor"`
	}
)

// Batch types
type (
	SignedBatch struct {
		BatchHeader     *BatchHeader     `json:"batch_header"`
		AttestationInfo *AttestationInfo `json:"attestation_info"`
	}

	BatchResponse struct {
		BatchHeaderHash    string                    `json:"batch_header_hash"`
		SignedBatch        *SignedBatch              `json:"signed_batch"`
		BlobKeys           []string                  `json:"blob_key"`
		BlobInclusionInfos []*BlobInclusionInfo      `json:"blob_inclusion_infos"`
		BlobCertificates   []*corev2.BlobCertificate `json:"blob_certificates"`
	}

	BatchInfo struct {
		BatchHeaderHash         string                  `json:"batch_header_hash"`
		BatchHeader             *BatchHeader            `json:"batch_header"`
		AttestedAt              uint64                  `json:"attested_at"`
		AggregatedSignature     *core.Signature         `json:"aggregated_signature"`
		QuorumNumbers           []core.QuorumID         `json:"quorum_numbers"`
		QuorumSignedPercentages map[core.QuorumID]uint8 `json:"quorum_signed_percentages"`
	}
	BatchFeedResponse struct {
		Batches []*BatchInfo `json:"batches"`
	}
)

// Account types
type (
	AccountBlobFeedResponse struct {
		AccountId string     `json:"account_id"`
		Blobs     []BlobInfo `json:"blobs"`
	}
)

// System types
type (
	MetricSummary struct {
		TotalBytesPosted      uint64  `json:"total_bytes_posted"`
		AverageBytesPerSecond float64 `json:"average_bytes_per_second"`
		StartTimestampSec     int64   `json:"start_timestamp_sec"`
		EndTimestampSec       int64   `json:"end_timestamp_sec"`
	}

	Metric struct {
		Throughput float64 `json:"throughput"`
	}

	Throughput struct {
		Throughput float64 `json:"throughput"`
		Timestamp  uint64  `json:"timestamp"`
	}

	SigningRateDataPoint struct {
		SigningRate float64 `json:"signing_rate"`
		Timestamp   uint64  `json:"timestamp"`
	}
	QuorumSigningRateData struct {
		QuorumId   string                 `json:"quorum_id"`
		DataPoints []SigningRateDataPoint `json:"data_points"`
	}
	NetworkSigningRateResponse struct {
		QuorumSigningRates []QuorumSigningRateData `json:"quorum_signing_rates"`
	}
)

func createBatchHeader(bh *corev2.BatchHeader) *BatchHeader {
	return &BatchHeader{
		BatchRoot:            hex.EncodeToString(bh.BatchRoot[:]),
		ReferenceBlockNumber: bh.ReferenceBlockNumber,
	}
}

func createBlobInclusionInfo(bi *corev2.BlobInclusionInfo) *BlobInclusionInfo {
	return &BlobInclusionInfo{
		BatchHeader:    createBatchHeader(bi.BatchHeader),
		BlobKey:        bi.BlobKey.Hex(),
		BlobIndex:      bi.BlobIndex,
		InclusionProof: hex.EncodeToString(bi.InclusionProof),
	}
}

func createBlobMetadata(bm *disperserv2.BlobMetadata) *BlobMetadata {
	return &BlobMetadata{
		BlobHeader:    bm.BlobHeader,
		Signature:     hex.EncodeToString(bm.Signature[:]),
		BlobStatus:    bm.BlobStatus.String(),
		BlobSizeBytes: bm.BlobSize,
		RequestedAt:   bm.RequestedAt,
		ExpiryUnixSec: bm.Expiry,
	}
}
