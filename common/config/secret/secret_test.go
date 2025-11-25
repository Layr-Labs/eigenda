package secret

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGetAndSet(t *testing.T) {
	s := NewSecret("this is my secret A")

	require.Equal(t, "this is my secret A", s.Get())

	oldValue := s.Set("this is my secret B")
	require.Equal(t, "this is my secret A", oldValue)
	require.Equal(t, "this is my secret B", s.Get())
}

func TestSecretNotExposedViaPrintf(t *testing.T) {
	secretValue := "super-secret-password"
	s := NewSecret(secretValue)

	testCases := []struct {
		name   string
		format string
	}{
		{"default format", "%v"},
		{"string format", "%s"},
		{"quoted string", "%q"},
		{"go-syntax", "%#v"},
		{"type and value", "%T %v"},
		{"pointer", "%p"},
		{"detailed struct", "%+v"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := fmt.Sprintf(tc.format, s)
			require.NotContains(t, output, secretValue, "Secret value should not be exposed in format: %s", tc.format)
		})
	}
}

func TestSecretNotExposedViaJSON(t *testing.T) {
	secretValue := "super-secret-api-key"
	type Config struct {
		APIKey  *Secret[string] `json:"api_key"`
		Timeout int             `json:"timeout"`
	}

	config := Config{
		APIKey:  NewSecret(secretValue),
		Timeout: 30,
	}

	// Test JSON marshaling
	jsonBytes, err := json.Marshal(config)
	require.NoError(t, err)
	jsonStr := string(jsonBytes)

	require.NotContains(t, jsonStr, secretValue, "Secret value should not be exposed in JSON")

	// Test JSON with indent
	jsonIndentBytes, err := json.MarshalIndent(config, "", "  ")
	require.NoError(t, err)
	jsonIndentStr := string(jsonIndentBytes)

	require.NotContains(t, jsonIndentStr, secretValue, "Secret value should not be exposed in indented JSON")
}

func TestSecretNotExposedViaYAML(t *testing.T) {
	secretValue := "super-secret-token"
	type Config struct {
		Token   *Secret[string] `yaml:"token"`
		Enabled bool            `yaml:"enabled"`
	}

	config := Config{
		Token:   NewSecret(secretValue),
		Enabled: true,
	}

	// Test YAML marshaling
	yamlBytes, err := yaml.Marshal(config)
	require.NoError(t, err)
	yamlStr := string(yamlBytes)

	require.NotContains(t, yamlStr, secretValue, "Secret value should not be exposed in YAML")
}
