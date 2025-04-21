package tracing

// TracingConfig contains configuration for tracing
type TracingConfig struct {
	Enabled     bool
	ServiceName string
	Endpoint    string
	SampleRatio float64
}
