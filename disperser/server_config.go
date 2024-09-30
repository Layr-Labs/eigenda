package disperser

import "time"

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	HttpPort    string
	GrpcPort    string
	GrpcTimeout time.Duration
}
