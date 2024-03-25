package disperser

import "time"

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	GrpcPort    string
	GrpcTimeout time.Duration

	// Feature flags
	// Whether enable the dual quorums.
	// If false, only quorum 0 will be used as required quorum.
	EnableDualQuorums bool
}
