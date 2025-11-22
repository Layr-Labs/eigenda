package dataapi

import "time"

type Config struct {
	SocketAddr         string
	ServerMode         string
	AllowOrigins       []string
	DisperserHostname  string
	ChurnerHostname    string
	BatcherHealthEndpt string
	FeedDelay          time.Duration
}

type DataApiVersion uint

const (
	V1 DataApiVersion = 1
	V2 DataApiVersion = 2
)
