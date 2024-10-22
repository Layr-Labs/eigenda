package clients

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
)

type EigenDAClientConfig struct {
	// RPC is the HTTP provider URL for the EigenDA Disperser
	RPC string

	// Timeout used when making dispersals to the EigenDA Disperser
	// TODO: we should change this param as its name is quite confusing
	ResponseTimeout time.Duration

	// The total amount of time that the client will spend waiting for EigenDA
	// to confirm a blob after it has been dispersed
	// Note that reasonable values for this field will depend on the value of WaitForFinalization.
	StatusQueryTimeout time.Duration

	// The amount of time to wait between status queries of a newly dispersed blob
	StatusQueryRetryInterval time.Duration

	// If true, will wait for the blob to finalize, if false, will wait only for the blob to confirm.
	WaitForFinalization bool

	// The quorum IDs to write blobs to using this client. Should not include default quorums 0 or 1.
	// TODO: should we change this to core.QuorumID instead? https://github.com/Layr-Labs/eigenda/blob/style--improve-api-clients-comments/core/data.go#L18
	CustomQuorumIDs []uint

	// Signer private key in hex encoded format. This key is currently purely used for authn/authz on the disperser.
	// For security, it should not be associated with an Ethereum address holding any funds.
	// This might change once we introduce payments.
	// OPTIONAL: this value is optional, and if set to "", will result in a read-only eigenDA client,
	// that can retrieve blobs but cannot disperse blobs.
	SignerPrivateKeyHex string

	// Whether to disable TLS for an insecure connection when connecting to a local EigenDA disperser instance.
	DisableTLS bool

	// The blob encoding version to use when writing blobs from the high level interface.
	PutBlobEncodingVersion codecs.BlobEncodingVersion

	// Point verification mode does an IFFT on data before it is written, and does an FFT on data after it is read.
	// This makes it possible to open points on the KZG commitment to prove that the field elements correspond to
	// the commitment. With this mode disabled, you will need to supply the entire blob to perform a verification
	// that any part of the data matches the KZG commitment.
	DisablePointVerificationMode bool
}

func (c *EigenDAClientConfig) CheckAndSetDefaults() error {
	if c.StatusQueryRetryInterval == 0 {
		c.StatusQueryRetryInterval = 5 * time.Second
	}
	if c.StatusQueryTimeout == 0 {
		c.StatusQueryTimeout = 25 * time.Minute
	}
	if c.ResponseTimeout == 0 {
		c.ResponseTimeout = 30 * time.Second
	}
	if len(c.SignerPrivateKeyHex) > 0 && len(c.SignerPrivateKeyHex) != 64 {
		return fmt.Errorf("a valid length SignerPrivateKeyHex needs to have 64 bytes")
	}

	if len(c.RPC) == 0 {
		return fmt.Errorf("EigenDAClientConfig.RPC not set")
	}
	return nil
}
