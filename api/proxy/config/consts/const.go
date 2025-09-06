package consts

const (
	EigenDAClientCategory   = "EigenDA V1 Client"
	EigenDAV2ClientCategory = "EigenDA V2 Client"
	LoggingFlagsCategory    = "Logging"
	MetricsFlagCategory     = "Metrics"
	MemstoreFlagsCategory   = "Memstore (for testing purposes - replaces EigenDA backend)"
	StorageFlagsCategory    = "Storage"
	S3Category              = "S3 Cache/Fallback"
	VerifierCategory        = "Cert Verifier (V1 only)"
	KZGCategory             = "KZG"
	ProxyServerCategory     = "Proxy Server"
	PaymentsCategory        = "Payments"

	DeprecatedRedisCategory = "Redis Cache/Fallback"

	// EnvVar prefix added in front of all environment variables accepted by the binary.
	// This acts as a namespace to avoid collisions with other binaries.
	GlobalEnvVarPrefix = "EIGENDA_PROXY"
)
