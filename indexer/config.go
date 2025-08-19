package indexer

import "time"

type Config struct {
	PullInterval            time.Duration
	ContractDeploymentBlock uint64 // The block number at which the contract was deployed
}
