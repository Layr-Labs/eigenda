package clients

import "time"

// PayloadDisperserConfig contains an embedded PayloadClientConfig, plus all additional configuration values needed
// by a PayloadDisperser
type PayloadDisperserConfig struct {
	PayloadClientConfig

	// DisperseBlobTimeout is the duration after which the PayloadDisperser will time out, when trying to disperse a
	// blob
	DisperseBlobTimeout time.Duration

	// BlobCertifiedTimeout is the duration after which the PayloadDisperser will time out, while polling
	// the disperser for blob status, waiting for BlobStatus_CERTIFIED
	BlobCertifiedTimeout time.Duration

	// BlobStatusPollInterval is the tick rate for the PayloadDisperser to use, while polling the disperser with
	// GetBlobStatus.
	BlobStatusPollInterval time.Duration
}
