package encoder

const (
	Localhost = "0.0.0.0"
)

type ServerConfig struct {
	GrpcPort                 string
	MaxConcurrentRequests    int
	RequestPoolSize          int
	EnableGnarkChunkEncoding bool
	PreventReencoding        bool
}
