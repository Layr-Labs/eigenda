package encoding

type BackendType string

const (
	BackendDefault BackendType = "default"
	BackendIcicle  BackendType = "icicle"
)

type Config struct {
	NumWorker   uint64
	BackendType BackendType
	EnableGPU   bool
	Verbose     bool
}

func WithIcicleBackend() Config {
	return Config{
		BackendType: BackendIcicle,
	}
}
