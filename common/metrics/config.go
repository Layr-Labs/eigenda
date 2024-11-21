package metrics

// Config provides configuration for a Metrics instance.
type Config struct {
	// Namespace is the namespace for the metrics.
	Namespace string

	// HTTPPort is the port to serve metrics on.
	HTTPPort int

	// MetricsBlacklist is a list of metrics to blacklist. To determine the fully qualified metric name
	// for this list, use the format "metricName:metricLabel" if the metric has a label, or just "metricLabel"
	// if the metric does not have a label. Any fully qualified metric name that matches exactly with an entry
	// in this list will be blacklisted (i.e. it will not be reported).
	MetricsBlacklist []string

	// MetricsFuzzyBlacklist is a list of metrics to blacklist. To determine the fully qualified metric name
	// for this list, use the format "metricName:metricLabel" if the metric has a label, or just "metricLabel"
	// if the metric does not have a label. Any fully qualified metric that contains one of these strings
	// in any position to be blacklisted (i.e. it will not be reported).
	MetricsFuzzyBlacklist []string
}
