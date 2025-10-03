package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	String   string
	Int      int
	Int64    int64
	Int32    int32
	Int16    int16
	Int8     int8
	Uint     uint
	Uint64   uint64
	Uint32   uint32
	Uint16   uint16
	Uint8    uint8
	Float64  float64
	Float32  float32
	Duration time.Duration
	Bool     bool
	Bar      Bar
	Baz      *Baz
}

func DefaultFoo() *Foo {
	return &Foo{}
}

func (f *Foo) Verify() error {
	if f.String == "invalid" {
		return fmt.Errorf("String may not be 'invalid'")
	}

	return nil
}

type Bar struct {
	A   string
	B   int
	C   bool
	Baz *Baz
}

func (b *Bar) Verify() error {
	return nil
}

type Baz struct {
	X string
	Y int
	Z bool
}

func (b *Baz) Verify() error {
	return nil
}

func TestTOMLParsing(t *testing.T) {

	configFile := "testdata/config.toml"

	foo, err := ParseConfig(DefaultFoo, "FOO", configFile)
	require.NoError(t, err)

	// Top-level fields
	require.Equal(t, "this value came from config.toml", foo.String)
	require.Equal(t, 0, foo.Int)
	require.Equal(t, int64(1), foo.Int64)
	require.Equal(t, int32(3), foo.Int32)
	require.Equal(t, int16(4), foo.Int16)
	require.Equal(t, int8(5), foo.Int8)
	require.Equal(t, uint(6), foo.Uint)
	require.Equal(t, uint64(7), foo.Uint64)
	require.Equal(t, uint32(8), foo.Uint32)
	require.Equal(t, uint16(9), foo.Uint16)
	require.Equal(t, uint8(10), foo.Uint8)
	require.Equal(t, 11.11, foo.Float64)
	require.Equal(t, float32(12.12), foo.Float32)
	require.Equal(t, 5*time.Second, foo.Duration)
	require.Equal(t, false, foo.Bool)

	// Bar field
	require.Equal(t, "bar A", foo.Bar.A)
	require.Equal(t, 25, foo.Bar.B)
	require.Equal(t, true, foo.Bar.C)
	// Bar.Baz field
	require.NotNil(t, foo.Bar.Baz)
	require.Equal(t, "barD baz X", foo.Bar.Baz.X)
	require.Equal(t, 26, foo.Bar.Baz.Y)
	require.Equal(t, false, foo.Bar.Baz.Z)

	// Baz field
	require.NotNil(t, foo.Baz)
	require.Equal(t, "baz X", foo.Baz.X)
	require.Equal(t, 27, foo.Baz.Y)
	require.Equal(t, true, foo.Baz.Z)
}

func TestJSONParsing(t *testing.T) {

	configFile := "testdata/config.json"

	foo, err := ParseConfig(DefaultFoo, "FOO", configFile)
	require.NoError(t, err)

	// Top-level fields
	require.Equal(t, "this value came from config.json", foo.String)
	require.Equal(t, 100, foo.Int)
	require.Equal(t, int64(101), foo.Int64)
	require.Equal(t, int32(103), foo.Int32)
	require.Equal(t, int16(104), foo.Int16)
	require.Equal(t, int8(105), foo.Int8)
	require.Equal(t, uint(106), foo.Uint)
	require.Equal(t, uint64(107), foo.Uint64)
	require.Equal(t, uint32(108), foo.Uint32)
	require.Equal(t, uint16(109), foo.Uint16)
	require.Equal(t, uint8(110), foo.Uint8)
	require.Equal(t, 111.11, foo.Float64)
	require.Equal(t, float32(112.12), foo.Float32)
	require.Equal(t, 1*time.Hour, foo.Duration)
	require.Equal(t, true, foo.Bool)

	// Bar field
	require.Equal(t, "json bar A", foo.Bar.A)
	require.Equal(t, 125, foo.Bar.B)
	require.Equal(t, false, foo.Bar.C)

	// Bar.Baz field
	require.NotNil(t, foo.Bar.Baz)
	require.Equal(t, "json barD baz X", foo.Bar.Baz.X)
	require.Equal(t, 126, foo.Bar.Baz.Y)
	require.Equal(t, true, foo.Bar.Baz.Z)

	// Baz field
	require.NotNil(t, foo.Baz)
	require.Equal(t, "json baz X", foo.Baz.X)
	require.Equal(t, 127, foo.Baz.Y)
	require.Equal(t, false, foo.Baz.Z)
}

func TestYAMLParsing(t *testing.T) {

	configFile := "testdata/config.yaml"

	foo, err := ParseConfig(DefaultFoo, "FOO", configFile)
	require.NoError(t, err)

	// Top-level fields
	require.Equal(t, "this value came from config.yaml", foo.String)
	require.Equal(t, 200, foo.Int)
	require.Equal(t, int64(201), foo.Int64)
	require.Equal(t, int32(203), foo.Int32)
	require.Equal(t, int16(204), foo.Int16)
	require.Equal(t, int8(105), foo.Int8)
	require.Equal(t, uint(206), foo.Uint)
	require.Equal(t, uint64(207), foo.Uint64)
	require.Equal(t, uint32(208), foo.Uint32)
	require.Equal(t, uint16(209), foo.Uint16)
	require.Equal(t, uint8(210), foo.Uint8)
	require.Equal(t, 211.11, foo.Float64)
	require.Equal(t, float32(212.12), foo.Float32)
	require.Equal(t, 33*time.Minute, foo.Duration)
	require.Equal(t, false, foo.Bool)

	// Bar field
	require.Equal(t, "yaml bar A", foo.Bar.A)
	require.Equal(t, 225, foo.Bar.B)
	require.Equal(t, true, foo.Bar.C)

	// Bar.Baz field
	require.NotNil(t, foo.Bar.Baz)
	require.Equal(t, "yaml barD baz X", foo.Bar.Baz.X)
	require.Equal(t, 226, foo.Bar.Baz.Y)
	require.Equal(t, false, foo.Bar.Baz.Z)

	// Baz field
	require.NotNil(t, foo.Baz)
	require.Equal(t, "yaml baz X", foo.Baz.X)
	require.Equal(t, 227, foo.Baz.Y)
	require.Equal(t, true, foo.Baz.Z)
}

func TestTOMLConfigOverride(t *testing.T) {

	configFile := "testdata/config.toml"
	overrideFile := "testdata/config_override.toml"

	foo, err := ParseConfig(DefaultFoo, "FOO", configFile, overrideFile)
	require.NoError(t, err)

	// Top-level fields - mix of base and override
	require.Equal(t, "this value came from config.toml", foo.String) // from base
	require.Equal(t, -1, foo.Int)                                    // from override
	require.Equal(t, int64(1), foo.Int64)                            // from base
	require.Equal(t, int32(-3), foo.Int32)                           // from override
	require.Equal(t, int16(4), foo.Int16)                            // from base
	require.Equal(t, int8(-5), foo.Int8)                             // from override
	require.Equal(t, uint(6), foo.Uint)                              // from base
	require.Equal(t, uint64(10007), foo.Uint64)                      // from override
	require.Equal(t, uint32(8), foo.Uint32)                          // from base
	require.Equal(t, uint16(10009), foo.Uint16)                      // from override
	require.Equal(t, uint8(10), foo.Uint8)                           // from base
	require.Equal(t, -11.11, foo.Float64)                            // from override
	require.Equal(t, 5*time.Second, foo.Duration)                    // from base
	require.Equal(t, float32(12.12), foo.Float32)                    // from base
	require.Equal(t, true, foo.Bool)                                 // from override

	// Bar field - mix of base and override
	require.Equal(t, "bar A", foo.Bar.A) // from base
	require.Equal(t, -25, foo.Bar.B)     // from override
	require.Equal(t, true, foo.Bar.C)    // from base

	// Baz field - mix of base and override
	require.NotNil(t, foo.Baz)
	require.Equal(t, "toml baz partial X", foo.Baz.X) // from override
	require.Equal(t, 27, foo.Baz.Y)                   // from base
	require.Equal(t, false, foo.Baz.Z)                // from override
}

func TestJSONConfigOverride(t *testing.T) {

	configFile := "testdata/config.json"
	overrideFile := "testdata/config_override.json"

	foo, err := ParseConfig(DefaultFoo, "FOO", configFile, overrideFile)
	require.NoError(t, err)

	// Top-level fields - mix of base and override
	require.Equal(t, "this value came from config.json", foo.String) // from base
	require.Equal(t, -100, foo.Int)                                  // from override
	require.Equal(t, int64(101), foo.Int64)                          // from base
	require.Equal(t, int32(-103), foo.Int32)                         // from override
	require.Equal(t, int16(104), foo.Int16)                          // from base
	require.Equal(t, int8(-15), foo.Int8)                            // from override
	require.Equal(t, uint(106), foo.Uint)                            // from base
	require.Equal(t, uint64(100007), foo.Uint64)                     // from override
	require.Equal(t, uint32(108), foo.Uint32)                        // from base
	require.Equal(t, uint16(10009), foo.Uint16)                      // from override
	require.Equal(t, uint8(110), foo.Uint8)                          // from base
	require.Equal(t, -111.11, foo.Float64)                           // from override
	require.Equal(t, float32(112.12), foo.Float32)                   // from base
	require.Equal(t, 1*time.Hour, foo.Duration)                      // from base
	require.Equal(t, false, foo.Bool)                                // from override

	// Bar field - mix of base and override
	require.Equal(t, "json bar A", foo.Bar.A) // from base
	require.Equal(t, -125, foo.Bar.B)         // from override
	require.Equal(t, false, foo.Bar.C)        // from base

	// Bar.Baz field - from base (not overridden)
	require.NotNil(t, foo.Bar.Baz)
	require.Equal(t, "json barD baz X", foo.Bar.Baz.X) // from base
	require.Equal(t, 126, foo.Bar.Baz.Y)               // from base
	require.Equal(t, true, foo.Bar.Baz.Z)              // from base

	// Baz field - mix of base and override
	require.NotNil(t, foo.Baz)
	require.Equal(t, "json baz partial X", foo.Baz.X) // from override
	require.Equal(t, 127, foo.Baz.Y)                  // from base
	require.Equal(t, true, foo.Baz.Z)                 // from override
}

func TestYAMLConfigOverride(t *testing.T) {

	configFile := "testdata/config.yaml"
	overrideFile := "testdata/config_override.yaml"

	foo, err := ParseConfig(DefaultFoo, "FOO", configFile, overrideFile)
	require.NoError(t, err)

	// Top-level fields - mix of base and override
	require.Equal(t, "this value came from config.yaml", foo.String) // from base
	require.Equal(t, -200, foo.Int)                                  // from override
	require.Equal(t, int64(201), foo.Int64)                          // from base
	require.Equal(t, int32(-203), foo.Int32)                         // from override
	require.Equal(t, int16(204), foo.Int16)                          // from base
	require.Equal(t, int8(-15), foo.Int8)                            // from override
	require.Equal(t, uint(206), foo.Uint)                            // from base
	require.Equal(t, uint64(200007), foo.Uint64)                     // from override
	require.Equal(t, uint32(208), foo.Uint32)                        // from base
	require.Equal(t, uint16(20009), foo.Uint16)                      // from override
	require.Equal(t, uint8(210), foo.Uint8)                          // from base
	require.Equal(t, -211.11, foo.Float64)                           // from override
	require.Equal(t, float32(212.12), foo.Float32)                   // from base
	require.Equal(t, 33*time.Minute, foo.Duration)                   // from base
	require.Equal(t, true, foo.Bool)                                 // from override

	// Bar field - mix of base and override
	require.Equal(t, "yaml bar A", foo.Bar.A) // from base
	require.Equal(t, -225, foo.Bar.B)         // from override
	require.Equal(t, true, foo.Bar.C)         // from base

	// Bar.Baz field - from base (not overridden)
	require.NotNil(t, foo.Bar.Baz)
	require.Equal(t, "yaml barD baz X", foo.Bar.Baz.X) // from base
	require.Equal(t, 226, foo.Bar.Baz.Y)               // from base
	require.Equal(t, false, foo.Bar.Baz.Z)             // from base

	// Baz field - mix of base and override
	require.NotNil(t, foo.Baz)
	require.Equal(t, "yaml baz partial X", foo.Baz.X) // from override
	require.Equal(t, 227, foo.Baz.Y)                  // from base
	require.Equal(t, false, foo.Baz.Z)                // from override
}

func TestInvalidTOML(t *testing.T) {
	configFile := "testdata/invalid_config.toml"

	_, err := ParseConfig(DefaultFoo, "FOO", configFile)
	require.Error(t, err)
}

func TestDefaultValues(t *testing.T) {

	configFile := "testdata/config_override.toml"

	constructor := func() *Foo {
		return &Foo{
			String:  "default string",
			Int:     42,
			Float64: 3.14,
			Bar: Bar{
				A: "default bar A",
				B: 84,
				C: true,
				Baz: &Baz{
					X: "default baz X",
					Y: 168,
					Z: false,
				},
			},
			Baz: &Baz{
				X: "default top-level baz X",
				Y: 336,
				Z: true,
			},
		}
	}

	foo, err := ParseConfig(constructor, "FOO", configFile)
	require.NoError(t, err)

	// Fields that are overridden by config_override.toml
	require.Equal(t, -1, foo.Int)               // overridden
	require.Equal(t, int32(-3), foo.Int32)      // overridden
	require.Equal(t, int8(-5), foo.Int8)        // overridden
	require.Equal(t, uint64(10007), foo.Uint64) // overridden
	require.Equal(t, uint16(10009), foo.Uint16) // overridden
	require.Equal(t, -11.11, foo.Float64)       // overridden
	require.Equal(t, true, foo.Bool)            // overridden

	// Fields that keep default values (not in override file)
	require.Equal(t, "default string", foo.String) // default
	require.Equal(t, int64(0), foo.Int64)          // default (zero value since not in override or default)
	require.Equal(t, int16(0), foo.Int16)          // default (zero value)
	require.Equal(t, uint(0), foo.Uint)            // default (zero value)
	require.Equal(t, uint32(0), foo.Uint32)        // default (zero value)
	require.Equal(t, uint8(0), foo.Uint8)          // default (zero value)
	require.Equal(t, float32(0), foo.Float32)      // default (zero value)

	// Bar field
	require.Equal(t, "default bar A", foo.Bar.A) // default
	require.Equal(t, -25, foo.Bar.B)             // overridden
	require.Equal(t, true, foo.Bar.C)            // default
	require.NotNil(t, foo.Bar.Baz)               // default (nested struct)
	require.Equal(t, "default baz X", foo.Bar.Baz.X)
	require.Equal(t, 168, foo.Bar.Baz.Y)
	require.Equal(t, false, foo.Bar.Baz.Z)

	// Baz field - mix of override and default
	require.NotNil(t, foo.Baz)
	require.Equal(t, "toml baz partial X", foo.Baz.X) // overridden
	require.Equal(t, 336, foo.Baz.Y)                  // default
	require.Equal(t, false, foo.Baz.Z)                // overridden
}

func TestEnvironmentVariables(t *testing.T) {

	configFile := "testdata/config.toml"

	// Set environment variables to override some config values.
	require.NoError(t, os.Setenv("PREFIX_STRING", "value from env var"))
	require.NoError(t, os.Setenv("PREFIX_INT", "-999"))
	require.NoError(t, os.Setenv("PREFIX_BAR_B", "-777"))
	require.NoError(t, os.Setenv("PREFIX_BAR_BAZ_X", "env var bar baz X"))
	require.NoError(t, os.Setenv("PREFIX_BAR_BAZ_Y", "444"))
	require.NoError(t, os.Setenv("PREFIX_BAR_BAZ_Z", "false"))
	require.NoError(t, os.Setenv("PREFIX_INT64", "0")) // zero value
	require.NoError(t, os.Setenv("PREFIX_INT32", "0")) // zero value

	require.NoError(t, os.Setenv("A_VARIABLE_THAT_DOES_NOT_HAVE_PREFIX", "should be ignored"))

	foo, err := ParseConfig(DefaultFoo, "PREFIX", configFile)
	require.NoError(t, err)

	// Verify that environment variables have overridden the config file values.
	require.Equal(t, "value from env var", foo.String) // from env
	require.Equal(t, -999, foo.Int)                    // from env
	require.Equal(t, int64(0), foo.Int64)              // from env (zero value)
	require.Equal(t, int32(0), foo.Int32)              // from env (zero value)
	require.Equal(t, int16(4), foo.Int16)              // from config
	require.Equal(t, int8(5), foo.Int8)                // from config
	require.Equal(t, uint(6), foo.Uint)                // from config
	require.Equal(t, uint64(7), foo.Uint64)            // from config
	require.Equal(t, uint32(8), foo.Uint32)            // from config
	require.Equal(t, uint16(9), foo.Uint16)            // from config
	require.Equal(t, uint8(10), foo.Uint8)             // from config
	require.Equal(t, 11.11, foo.Float64)               // from config
	require.Equal(t, float32(12.12), foo.Float32)      // from config
	require.Equal(t, 5*time.Second, foo.Duration)      // from config
	require.Equal(t, false, foo.Bool)                  // from config

	// Bar field
	require.Equal(t, "bar A", foo.Bar.A) // from config
	require.Equal(t, -777, foo.Bar.B)    // from env
	require.Equal(t, true, foo.Bar.C)    // from config

	// Bar.Baz field
	require.NotNil(t, foo.Bar.Baz)
	require.Equal(t, "env var bar baz X", foo.Bar.Baz.X) // from env
	require.Equal(t, 444, foo.Bar.Baz.Y)                 // from env
	require.Equal(t, false, foo.Bar.Baz.Z)               // from env

	// Baz field - the env vars use FOO_BAZ_PARTIAL_* which doesn't match foo.Baz,
	// so these should come from config
	require.NotNil(t, foo.Baz)
	require.Equal(t, "baz X", foo.Baz.X) // from config
	require.Equal(t, 27, foo.Baz.Y)      // from config
	require.Equal(t, true, foo.Baz.Z)    // from config
}

func TestInvalidEnvironmentVariable(t *testing.T) {
	configFile := "testdata/config.toml"

	// Set environment variables to override some config values.
	require.NoError(t, os.Setenv("PREFIX_STRING", "value from env var"))
	require.NoError(t, os.Setenv("PREFIX_THIS_VARIABLE_WAS_MISTYPED", "should not be ignored"))

	_, err := ParseConfig(DefaultFoo, "PREFIX", configFile)
	require.Error(t, err)

	require.NoError(t, os.Unsetenv("PREFIX_THIS_VARIABLE_WAS_MISTYPED"))
}

func TestVerificaitonFailure(t *testing.T) {
	configFile := "testdata/config.toml"

	// Set environment variables to override some config values.
	require.NoError(t, os.Setenv("PREFIX_STRING", "invalid")) // will cause verification to fail

	_, err := ParseConfig(DefaultFoo, "PREFIX", configFile)
	require.Error(t, err)
	require.Contains(t, err.Error(), "String may not be 'invalid'")
}

func TestDocGeneration(t *testing.T) {

	packagePaths := []string{
		"github.com/Layr-Labs/eigenda/node",
		"github.com/Layr-Labs/eigensdk-go/signer/bls/types",
		"github.com/Layr-Labs/eigenda/common/geth",
		"github.com/Layr-Labs/eigenda/common",
		"github.com/Layr-Labs/eigensdk-go/logging",
		"github.com/Layr-Labs/eigenda/encoding/kzg",
		"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation",
	}

	err := DocumentConfig(
		func() node.Config {
			return node.Config{}
		},
		"NODE",
		packagePaths,
		"config.md")
	require.NoError(t, err)

}
