package disperser

import "time"

type ServerConfig struct {
	GrpcPort string

	// This timeout is used control the maximum age of a DisperseBlobAuthenticated() RPC call
	// (via a context with a timeout).
	GrpcTimeout time.Duration

	// The maximum permissible age of a GRPC connection before it is closed. If zero, then the server will not close
	// connections based on age.
	MaxConnectionAge time.Duration

	// When the server closes a connection due to MaxConnectionAge, it will wait for this grace period before
	// forcibly closing the connection. This allows in-flight requests to complete.
	MaxConnectionAgeGrace time.Duration

	// MaxIdleConnectionAge is the maximum time a connection can be idle before it is closed.
	MaxIdleConnectionAge time.Duration

	PprofHttpPort string
	EnablePprof   bool

	// DisableGetBlobCommitment, if true, causes the GetBlobCommitment gRPC endpoint to return
	// a deprecation error. This endpoint is deprecated and will be removed in a future release.
	DisableGetBlobCommitment bool

	// The amount of time to retain signing rate data.
	SigningRateRetentionPeriod time.Duration

	// The interval at which to poll for signing rate data from the controller.
	SigningRatePollInterval time.Duration
}
