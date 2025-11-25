package test

import (
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/stretchr/testify/require"
)

func TestConfigParsing(t *testing.T) {
	dir := t.TempDir()

	cfg := &StandardConfig{
		Foo: "example",
		Bar: 42,
		Baz: true,
	}

	err := config.DocumentConfig(
		func() config.DocumentedConfig {
			return cfg
		},
		dir,
		true)

	require.NoError(t, err)

	// It's tricky to verify the exact contents of the generated file, since it is designed for human consumption.
	// But we can look for a few strings that should definitely be there.

	content, err := os.ReadFile(dir + "/NameForStandardConfig.md")
	require.NoError(t, err)

	expectedStrings := []string{
		"NameForStandardConfig",
		// Foo
		"Foo",
		"TEST_FOO",
		"string",
		"This is variable Foo.",
		"example",
		// Bar
		"Bar",
		"TEST_BAR",
		"int",
		"This is variable Bar.",
		"Bar has more than one line of documentation.",
		"42",
		// Baz
		"Baz",
		"TEST_BAZ",
		"This is variable Baz. It has a '}' character, which used to cause a bug.",
		"bool",
		"true",
		// Nested.NestedField
		"Nested.NestedField",
		"TEST_NESTED_NESTED_FIELD",
		"string",
		"This is variable NestedField.",
		// Nested.DoublyNested.DoublyNestedField
		"Nested.DoublyNested.DoublyNestedField",
		"TEST_NESTED_DOUBLY_NESTED_DOUBLY_NESTED_FIELD",
		"int",
		"This is variable DoublyNestedField.",
	}

	for _, str := range expectedStrings {
		require.Contains(t, string(content), str)
	}

	// Look for some strings that should NOT be there.
	unexpectedStrings := []string{
		"privateIgnoredField",
		"// This field is unexported and should be ignored.",
	}

	for _, str := range unexpectedStrings {
		require.NotContains(t, string(content), str)
	}
}
