package v2

import "github.com/Layr-Labs/eigenda/core"

type BlobStatus uint

const (
	Queued BlobStatus = iota
	Encoded
	Certified
	Failed
)

type BlobMetadata struct {
	core.BlobHeaderV2 `json:"blob_header"`

	BlobKey    core.BlobKey `json:"blob_key"`
	BlobStatus BlobStatus   `json:"blob_status"`
	// Expiry is Unix timestamp of the blob expiry in seconds from epoch
	Expiry      uint64 `json:"expiry"`
	NumRetries  uint   `json:"num_retries"`
	BlobSize    uint64 `json:"blob_size"`
	RequestedAt uint64 `json:"requested_at"`
}
