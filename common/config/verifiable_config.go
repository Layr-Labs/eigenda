package config

// VerifiableConfig is an interface for configurations that can be validated.
type VerifiableConfig interface {
	// Verify checks that the configuration is valid, returning an error if it is not.
	Verify() error
}

// Configuration that includes documentation metadata.
type DocumentedConfig interface {
	VerifiableConfig

	// Returns the name of the configuration. By convention, this should be in CamelCase.
	GetName() string

	// Returns the environment variable prefix for the configuration. By convention,
	// these should be in SCREAMING_SNAKE_CASE.
	GetEnvVarPrefix() string

	// Returns a list of packages that need to be loaded in order to fully resolve
	// the configuration and all nested types within the configuration. Used for scraping godocs.
	GetPackagePaths() []string
}
