package client

import (
	"fmt"
	"time"
)

type Config struct {
	// RPC is the HTTP provider URL for the Data Availability node.
	RPC string

	// The total amount of time that the client will spend waiting for EigenDA to confirm a blob
	StatusQueryTimeout time.Duration

	// The amount of time to wait between status queries of a newly dispersed blob
	StatusQueryRetryInterval time.Duration

	// The total amount of time that the client will waiting for a response from the EigenDA disperser
	ResponseTimeout time.Duration

	// The quorum IDs to write blobs to using this client. Should not include quorums 0 or 1.
	CustomQuorumIDs []uint

	// Signer private key in hex encoded format. This key should not be associated with an Ethereum address holding any funds.
	SignerPrivateKeyHex string

	// Whether to disable TLS for an insecure connection when connecting to a local EigenDA disperser instance.
	DisableTLS bool
}

var DefaultQuorums = map[uint]bool{0: true, 1: true}

func (c *Config) Check() error {
	for _, e := range c.CustomQuorumIDs {
		if DefaultQuorums[e] {
			return fmt.Errorf("EigenDA client config failed validation because CustomQuorumIDs includes a default quorum ID %d. Because it is included by default this quorum ID can be removed from the client configuration", e)
		}
	}
	if c.StatusQueryRetryInterval == 0 {
		c.StatusQueryRetryInterval = 5 * time.Second
	}
	if c.StatusQueryTimeout == 0 {
		c.StatusQueryTimeout = 25 * time.Minute
	}
	return nil
}
