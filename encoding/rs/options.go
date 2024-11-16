package rs

import "github.com/Layr-Labs/eigenda/encoding"

// EncoderOption defines the function signature for encoder options
type EncoderOption func(*Encoder)

// Option Definitions
func WithBackend(backend encoding.BackendType) EncoderOption {
	return func(e *Encoder) {
		e.Config.BackendType = backend
	}
}

func WithGPU(enable bool) EncoderOption {
	return func(e *Encoder) {
		e.Config.GPUEnable = enable
	}
}

func WithNumWorkers(workers uint64) EncoderOption {
	return func(e *Encoder) {
		e.Config.NumWorker = workers
	}
}

func WithVerbose(verbose bool) EncoderOption {
	return func(e *Encoder) {
		e.Config.Verbose = verbose
	}
}
