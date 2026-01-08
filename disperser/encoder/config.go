package encoder

import "github.com/Layr-Labs/eigenda/encoding"

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	// MaxConcurrentRequestsDangerous limits the number of concurrent encoding requests the server will handle,
	// which also limits the number of concurrent GPU encodings if GPUEnable is true.
	// This is a dangerous setting because setting it too high may lead to out-of-memory panics on the GPU.
	MaxConcurrentRequestsDangerous int
	// RequestPoolSize is the maximum number of requests in the request pool.
	RequestPoolSize int
	// RequestQueueSize is the maximum number of requests in the request queue.
	RequestQueueSize int
	// EnableGnarkChunkEncoding if true, will produce chunks in Gnark, instead of Gob
	EnableGnarkChunkEncoding bool
	// PreventReencoding if true, will prevent reencoding of chunks by checking
	// if the chunk already exists in the chunk store
	PreventReencoding bool
	// Backend to use for encoding. Supported values are "gnark" and "icicle".
	Backend encoding.BackendType
	// GPUEnable enables GPU, falls back to CPU if not available
	GPUEnable bool
	// PprofHttpPort is the http port which the pprof server is listening
	PprofHttpPort string
	// EnablePprof starts the pprof server
	EnablePprof bool
}
