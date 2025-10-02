package config

// VerifiableConfig is an interface for configurations that can be validated.
type VerifiableConfig interface {
	// Verify checks that the configuration is valid, returning an error if it is not.
	Verify() error
}
