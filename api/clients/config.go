package clients

import (
	"fmt"
	"time"
)

type EigenDAClientConfig struct {
	// RPC is the HTTP provider URL for the Data Availability node.
	RPC string

	// The total amount of time that the client will spend waiting for EigenDA to confirm a blob
	StatusQueryTimeout time.Duration

	// The amount of time to wait between status queries of a newly dispersed blob
	StatusQueryRetryInterval time.Duration

	// The total amount of time that the client will waiting for a response from the EigenDA disperser
	ResponseTimeout time.Duration

	// The quorum IDs to write blobs to using this client. Should not include default quorums 0 or 1.
	CustomQuorumIDs []uint

	// Signer private key in hex encoded format. This key should not be associated with an Ethereum address holding any funds.
	SignerPrivateKeyHex string

	// Whether to disable TLS for an insecure connection when connecting to a local EigenDA disperser instance.
	DisableTLS bool

	// The blob encoding version to use when writing blobs from the high level interface.
	PutBlobEncodingVersion BlobEncodingVersion
}

func (c *EigenDAClientConfig) Check() error {
	if c.StatusQueryRetryInterval == 0 {
		c.StatusQueryRetryInterval = 5 * time.Second
	}
	if c.StatusQueryTimeout == 0 {
		c.StatusQueryTimeout = 25 * time.Minute
	}
	if c.ResponseTimeout == 0 {
		c.ResponseTimeout = 30 * time.Second
	}
	if len(c.SignerPrivateKeyHex) != 64 {
		return fmt.Errorf("EigenDAClientConfig.SignerPrivateKeyHex should be 64 hex characters long, should not have 0x prefix")
	}
	if len(c.RPC) == 0 {
		return fmt.Errorf("EigenDAClientConfig.RPC not set")
	}
	return nil
}
