package metrics

// Config provides configuration for a Metrics instance.
type Config struct {
	// Namespace is the namespace for the metrics.
	Namespace string

	// HTTPPort is the port to serve metrics on.
	HTTPPort int
}
