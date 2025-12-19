package indexer

import "time"

type Config struct {
	// The frequency to pull data from The Graph.
	PullInterval time.Duration
}
