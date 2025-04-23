package validator

import "time"

// ValidatorClientConfig contains the configuration for the validator retrieval client.
type ValidatorClientConfig struct {

	// If 1.0, then the validator retrieval client will attempt to download the exact number of chunks
	// needed to reconstruct the blob. If higher than 1.0, then the validator retrieval client will
	// pessimistically assume that some operators will not respond in time, and will download
	// additional chunks from other operators. For example, at 2.0, the validator retrieval client
	// will download twice the number of chunks needed to reconstruct the blob. Setting this to below
	// 1.0 is not supported.
	//
	// The default value is 2.0.
	DownloadPessimism float64

	// If 1.0, then the validator retrieval client will attempt to verify the exact number of chunks
	// needed to reconstruct the blob. If higher than 1.0, then the validator retrieval client will
	// pessimistically assume that some operators sent invalid chunks, and will verify additional chunks
	// from other operators. For example, at 2.0, the validator retrieval client will verify twice the number of
	// chunks needed to reconstruct the blob. Setting this to below 1.0 is not supported.
	//
	// The default value is 1.0.
	VerificationPessimism float64

	// After this amount of time passes, the validator retrieval client will assume that the operator is not
	// responding, and will start downloading from a different operator. The download is not terminated when
	// this timeout is reached.
	//
	// The default value is 10 seconds.
	PessimisticTimeout time.Duration

	// The absolute limit on the time to wait for a download to complete. If this timeout is reached, the
	// download will be terminated.
	//
	// The default value is 30 seconds.
	DownloadTimeout time.Duration

	// The control loop periodically wakes up to do work. This is the period of that control loop.
	//
	// The default value is 1 second.
	ControlLoopPeriod time.Duration

	// If true, then the validator retrieval client will log detailed information about the download process
	// (at debug level).
	//
	// The default value is false.
	DetailedLogging bool
}

// DefaultClientConfig returns the default configuration for the validator retrieval client.
func DefaultClientConfig() *ValidatorClientConfig {
	return &ValidatorClientConfig{
		DownloadPessimism:     2.0,
		VerificationPessimism: 1.0,
		PessimisticTimeout:    10 * time.Second,
		DownloadTimeout:       30 * time.Second,
		ControlLoopPeriod:     1 * time.Second,
		DetailedLogging:       false,
	}
}
