package encoder

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	MaxConcurrentRequests    int
	RequestPoolSize          int
	RequestQueueSize         int
	EnableGnarkChunkEncoding bool
	PreventReencoding        bool
	Backend                  string
	GPUEnable                bool
	GPUConcurrentFrameGenerationDangerous int64
	PprofHttpPort                         string
	EnablePprof                           bool
}
