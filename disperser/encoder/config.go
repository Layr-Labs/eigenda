package encoder

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	// MaxConcurrentRequestsDangerous limits the number of concurrent encoding requests the server will handle,
	// which also limits the number of concurrent GPU encodings if GPUEnable is true.
	// This is a dangerous setting because setting it too high may lead to out-of-memory panics on the GPU.
	MaxConcurrentRequestsDangerous int
	RequestPoolSize                int
	RequestQueueSize               int
	EnableGnarkChunkEncoding       bool
	PreventReencoding              bool
	Backend                        string
	GPUEnable                      bool
	PprofHttpPort                  string
	EnablePprof                    bool
}
