package v2

import (
	"encoding/hex"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
)

type BlobStatus uint

const (
	Queued BlobStatus = iota
	Encoded
	Certified
	Failed
)

type BlobVersion uint32

type BlobKey [32]byte

func (b BlobKey) Hex() string {
	return hex.EncodeToString(b[:])
}

func HexToBlobKey(h string) (BlobKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return BlobKey{}, err
	}
	return BlobKey(b), nil
}

type BlobHeader struct {
	BlobVersion     BlobVersion              `json:"version"`
	BlobQuorumInfos []*core.BlobQuorumInfo   `json:"blob_quorum_infos"`
	BlobCommitment  encoding.BlobCommitments `json:"commitments"`

	core.PaymentMetadata `json:"payment_metadata"`
}

type BlobMetadata struct {
	BlobHeader `json:"blob_header"`

	BlobStatus  BlobStatus `json:"blob_status"`
	Expiry      uint64     `json:"expiry"`
	NumRetries  uint       `json:"num_retries"`
	BlobSize    uint64     `json:"blob_size"`
	RequestedAt uint64     `json:"requested_at"`
}
