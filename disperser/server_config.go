package disperser

import "time"

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	GrpcPort    string
	GrpcTimeout time.Duration

	PprofHttpPort string
	EnablePprof   bool
}
