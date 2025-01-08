package encoder

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	GrpcPort                 string
	MaxConcurrentRequests    int
	RequestPoolSize          int
	RequestQueueSize         int
	EnableGnarkChunkEncoding bool
	PreventReencoding        bool
	Backend                  string
	GPUEnable                bool
	PprofHttpPort            string
	EnablePprof              bool
}
