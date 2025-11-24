package test

import (
	"github.com/Layr-Labs/eigenda/common/config"
)

var _ config.DocumentedConfig = (*StandardConfig)(nil)

// This is a test config used by config_document_generator_test.go. Can't be in a test file since we need to import it.
type StandardConfig struct {
	// This is variable Foo.
	Foo string

	// This is variable Bar.
	// Bar has more than one line of documentation.
	Bar int

	// This is variable Baz. It has a '}' character, which used to cause a bug.
	Baz bool

	// This is a nested config, does not use a pointer.
	Nested NestedConfig

	// This field is unexported and should be ignored.
	privateIgnoredField string
}

type NestedConfig struct {
	// This is variable NestedField.
	NestedField string

	// This is a doubly nested config. Uses a pointer to a struct.
	DoublyNested *DoublyNestedConfig
}

type DoublyNestedConfig struct {
	// This is variable DoublyNestedField.
	DoublyNestedField int
}

func (s *StandardConfig) GetEnvVarPrefix() string {
	return "TEST"
}

func (s *StandardConfig) GetName() string {
	return "NameForStandardConfig"
}

func (s *StandardConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/common/config/test",
	}
}

func (s *StandardConfig) Verify() error {
	return nil
}
