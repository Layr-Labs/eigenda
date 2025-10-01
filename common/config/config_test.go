package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	String    string
	Int       int
	Int64     int64
	Int32     int32
	Int16     int16
	Int8      int8
	Uint      uint
	Uint64    uint64
	Uint32    uint32
	Uint16    uint16
	Uint8     uint8
	Float64   float64
	Float32   float32
	Bool      bool
	Recursive *Foo
	Bar       Bar
	Baz       *Baz
}

func (f *Foo) Verify() error {
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

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile)
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
	require.Equal(t, false, foo.Bool)

	// Recursive field
	require.NotNil(t, foo.Recursive)
	require.Equal(t, "this value also came from config.toml", foo.Recursive.String)
	require.Equal(t, 13, foo.Recursive.Int)
	require.Equal(t, int64(14), foo.Recursive.Int64)
	require.Equal(t, int32(15), foo.Recursive.Int32)
	require.Equal(t, int16(16), foo.Recursive.Int16)
	require.Equal(t, int8(17), foo.Recursive.Int8)
	require.Equal(t, uint(18), foo.Recursive.Uint)
	require.Equal(t, uint64(19), foo.Recursive.Uint64)
	require.Equal(t, uint32(20), foo.Recursive.Uint32)
	require.Equal(t, uint16(21), foo.Recursive.Uint16)
	require.Equal(t, uint8(22), foo.Recursive.Uint8)
	require.Equal(t, 23.23, foo.Recursive.Float64)
	require.Equal(t, float32(24.24), foo.Recursive.Float32)
	require.Equal(t, true, foo.Recursive.Bool)

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

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile)
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
	require.Equal(t, true, foo.Bool)

	// Recursive field
	require.NotNil(t, foo.Recursive)
	require.Equal(t, "this value also came from config.json", foo.Recursive.String)
	require.Equal(t, 113, foo.Recursive.Int)
	require.Equal(t, int64(114), foo.Recursive.Int64)
	require.Equal(t, int32(115), foo.Recursive.Int32)
	require.Equal(t, int16(116), foo.Recursive.Int16)
	require.Equal(t, int8(117), foo.Recursive.Int8)
	require.Equal(t, uint(118), foo.Recursive.Uint)
	require.Equal(t, uint64(119), foo.Recursive.Uint64)
	require.Equal(t, uint32(120), foo.Recursive.Uint32)
	require.Equal(t, uint16(121), foo.Recursive.Uint16)
	require.Equal(t, uint8(122), foo.Recursive.Uint8)
	require.Equal(t, 123.23, foo.Recursive.Float64)
	require.Equal(t, float32(124.24), foo.Recursive.Float32)
	require.Equal(t, false, foo.Recursive.Bool)

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

	configFile := "testdata/config.yml"

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile)
	require.NoError(t, err)

	// Top-level fields
	require.Equal(t, "this value came from config.yml", foo.String)
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
	require.Equal(t, false, foo.Bool)

	// Recursive field
	require.NotNil(t, foo.Recursive)
	require.Equal(t, "this value also came from config.yml", foo.Recursive.String)
	require.Equal(t, 213, foo.Recursive.Int)
	require.Equal(t, int64(214), foo.Recursive.Int64)
	require.Equal(t, int32(215), foo.Recursive.Int32)
	require.Equal(t, int16(216), foo.Recursive.Int16)
	require.Equal(t, int8(117), foo.Recursive.Int8)
	require.Equal(t, uint(218), foo.Recursive.Uint)
	require.Equal(t, uint64(219), foo.Recursive.Uint64)
	require.Equal(t, uint32(220), foo.Recursive.Uint32)
	require.Equal(t, uint16(221), foo.Recursive.Uint16)
	require.Equal(t, uint8(222), foo.Recursive.Uint8)
	require.Equal(t, 223.23, foo.Recursive.Float64)
	require.Equal(t, float32(224.24), foo.Recursive.Float32)
	require.Equal(t, true, foo.Recursive.Bool)

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

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile, overrideFile)
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
	require.Equal(t, float32(12.12), foo.Float32)                    // from base
	require.Equal(t, true, foo.Bool)                                 // from override

	// Recursive field - mix of base and override
	require.NotNil(t, foo.Recursive)
	require.Equal(t, "this value also came from config.toml", foo.Recursive.String) // from base
	require.Equal(t, -13, foo.Recursive.Int)                                        // from override
	require.Equal(t, int64(14), foo.Recursive.Int64)                                // from base
	require.Equal(t, int32(-15), foo.Recursive.Int32)                               // from override
	require.Equal(t, int16(16), foo.Recursive.Int16)                                // from base
	require.Equal(t, int8(-17), foo.Recursive.Int8)                                 // from override
	require.Equal(t, uint(18), foo.Recursive.Uint)                                  // from base
	require.Equal(t, uint64(100019), foo.Recursive.Uint64)                          // from override
	require.Equal(t, uint32(20), foo.Recursive.Uint32)                              // from base
	require.Equal(t, uint16(10021), foo.Recursive.Uint16)                              // from base
	require.Equal(t, uint8(22), foo.Recursive.Uint8)                                // from base
	require.Equal(t, -23.23, foo.Recursive.Float64)                                 // from override
	require.Equal(t, float32(24.24), foo.Recursive.Float32)                         // from base
	require.Equal(t, false, foo.Recursive.Bool)                                     // from override

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

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile, overrideFile)
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
	require.Equal(t, false, foo.Bool)                                // from override

	// Recursive field - mix of base and override
	require.NotNil(t, foo.Recursive)
	require.Equal(t, "this value also came from config.json", foo.Recursive.String) // from base
	require.Equal(t, -113, foo.Recursive.Int)                                       // from override
	require.Equal(t, int64(114), foo.Recursive.Int64)                               // from base
	require.Equal(t, int32(-115), foo.Recursive.Int32)                              // from override
	require.Equal(t, int16(116), foo.Recursive.Int16)                               // from base
	require.Equal(t, int8(-17), foo.Recursive.Int8)                                 // from override
	require.Equal(t, uint(118), foo.Recursive.Uint)                                 // from base
	require.Equal(t, uint64(1000019), foo.Recursive.Uint64)                         // from override
	require.Equal(t, uint32(120), foo.Recursive.Uint32)                             // from base
	require.Equal(t, uint16(10021), foo.Recursive.Uint16)                           // from override
	require.Equal(t, uint8(122), foo.Recursive.Uint8)                               // from base
	require.Equal(t, -123.23, foo.Recursive.Float64)                                // from override
	require.Equal(t, float32(124.24), foo.Recursive.Float32)                        // from base
	require.Equal(t, true, foo.Recursive.Bool)                                      // from override

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

	configFile := "testdata/config.yml"
	overrideFile := "testdata/config_override.yml"

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile, overrideFile)
	require.NoError(t, err)

	// Top-level fields - mix of base and override
	require.Equal(t, "this value came from config.yml", foo.String) // from base
	require.Equal(t, -200, foo.Int)                                 // from override
	require.Equal(t, int64(201), foo.Int64)                         // from base
	require.Equal(t, int32(-203), foo.Int32)                        // from override
	require.Equal(t, int16(204), foo.Int16)                         // from base
	require.Equal(t, int8(-15), foo.Int8)                           // from override
	require.Equal(t, uint(206), foo.Uint)                           // from base
	require.Equal(t, uint64(200007), foo.Uint64)                    // from override
	require.Equal(t, uint32(208), foo.Uint32)                       // from base
	require.Equal(t, uint16(20009), foo.Uint16)                     // from override
	require.Equal(t, uint8(210), foo.Uint8)                         // from base
	require.Equal(t, -211.11, foo.Float64)                          // from override
	require.Equal(t, float32(212.12), foo.Float32)                  // from base
	require.Equal(t, true, foo.Bool)                                // from override

	// Recursive field - mix of base and override
	require.NotNil(t, foo.Recursive)
	require.Equal(t, "this value also came from config.yml", foo.Recursive.String) // from base
	require.Equal(t, -213, foo.Recursive.Int)                                      // from override
	require.Equal(t, int64(214), foo.Recursive.Int64)                              // from base
	require.Equal(t, int32(-215), foo.Recursive.Int32)                             // from override
	require.Equal(t, int16(216), foo.Recursive.Int16)                              // from base
	require.Equal(t, int8(-17), foo.Recursive.Int8)                                // from override
	require.Equal(t, uint(218), foo.Recursive.Uint)                                // from base
	require.Equal(t, uint64(2000019), foo.Recursive.Uint64)                        // from override
	require.Equal(t, uint32(220), foo.Recursive.Uint32)                            // from base
	require.Equal(t, uint16(20021), foo.Recursive.Uint16)                          // from override
	require.Equal(t, uint8(222), foo.Recursive.Uint8)                              // from base
	require.Equal(t, -223.23, foo.Recursive.Float64)                               // from override
	require.Equal(t, float32(224.24), foo.Recursive.Float32)                       // from base
	require.Equal(t, false, foo.Recursive.Bool)                                    // from override

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

	var foo Foo
	err := ParseConfig(&foo, "FOO", configFile)
	require.Error(t, err)
}

func TestDefaultValues(t *testing.T) {

	configFile := "testdata/config_override.toml"

	defaultValue := &Foo{
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

	err := ParseConfig(defaultValue, "FOO", configFile)
	require.NoError(t, err)

	// Fields that are overridden by config_override.toml
	require.Equal(t, -1, defaultValue.Int)                 // overridden
	require.Equal(t, int32(-3), defaultValue.Int32)        // overridden
	require.Equal(t, int8(-5), defaultValue.Int8)          // overridden
	require.Equal(t, uint64(10007), defaultValue.Uint64)   // overridden
	require.Equal(t, uint16(10009), defaultValue.Uint16)   // overridden
	require.Equal(t, -11.11, defaultValue.Float64)         // overridden
	require.Equal(t, true, defaultValue.Bool)              // overridden

	// Fields that keep default values (not in override file)
	require.Equal(t, "default string", defaultValue.String) // default
	require.Equal(t, int64(0), defaultValue.Int64)          // default (zero value since not in override or default)
	require.Equal(t, int16(0), defaultValue.Int16)          // default (zero value)
	require.Equal(t, uint(0), defaultValue.Uint)            // default (zero value)
	require.Equal(t, uint32(0), defaultValue.Uint32)        // default (zero value)
	require.Equal(t, uint8(0), defaultValue.Uint8)          // default (zero value)
	require.Equal(t, float32(0), defaultValue.Float32)      // default (zero value)

	// Recursive field - mix of override and defaults
	require.NotNil(t, defaultValue.Recursive)
	require.Equal(t, -13, defaultValue.Recursive.Int)                 // overridden
	require.Equal(t, int32(-15), defaultValue.Recursive.Int32)        // overridden
	require.Equal(t, int8(-17), defaultValue.Recursive.Int8)          // overridden
	require.Equal(t, uint64(100019), defaultValue.Recursive.Uint64)   // overridden
	require.Equal(t, uint16(10021), defaultValue.Recursive.Uint16)    // overridden
	require.Equal(t, -23.23, defaultValue.Recursive.Float64)          // overridden
	require.Equal(t, false, defaultValue.Recursive.Bool)              // overridden

	// Bar field
	require.Equal(t, "default bar A", defaultValue.Bar.A) // default
	require.Equal(t, -25, defaultValue.Bar.B)             // overridden
	require.Equal(t, true, defaultValue.Bar.C)            // default
	require.NotNil(t, defaultValue.Bar.Baz)               // default (nested struct)
	require.Equal(t, "default baz X", defaultValue.Bar.Baz.X)
	require.Equal(t, 168, defaultValue.Bar.Baz.Y)
	require.Equal(t, false, defaultValue.Bar.Baz.Z)

	// Baz field - mix of override and default
	require.NotNil(t, defaultValue.Baz)
	require.Equal(t, "toml baz partial X", defaultValue.Baz.X) // overridden
	require.Equal(t, 336, defaultValue.Baz.Y)                  // default
	require.Equal(t, false, defaultValue.Baz.Z)                // overridden
}

func TestEnvironmentVariables(t *testing.T) {

	configFile := "testdata/config.toml"

	// Set environment variables to override some config values.
	os.Setenv("FOO_STRING", "value from env var")
	os.Setenv("FOO_INT", "-999")
	os.Setenv("FOO_RECURSIVE_INT", "-888")
	os.Setenv("FOO_BAR_B", "-777")
	os.Setenv("FOO_BAZ_PARTIAL_X", "env var baz X")
	os.Setenv("FOO_BAZ_PARTIAL_Y", "555")
	os.Setenv("FOO_BAZ_PARTIAL_Z", "true")
	os.Setenv("FOO_BAR_BAZ_X", "env var bar baz X")
	os.Setenv("FOO_BAR_BAZ_Y", "444")
	os.Setenv("FOO_BAR_BAZ_Z", "false")
	os.Setenv("FOO_INT64", "0")    // zero value
	os.Setenv("FOO_INT32", "0")    // zero value

	foo := &Foo{

	}
	err := ParseConfig(foo, "FOO", configFile)
	require.NoError(t, err)

	// Verify that environment variables have overridden the config file values.
	require.Equal(t, "value from env var", foo.String)    // from env
	require.Equal(t, -999, foo.Int)                       // from env
	require.Equal(t, int64(0), foo.Int64)                 // from env (zero value)
	require.Equal(t, int32(0), foo.Int32)                 // from env (zero value)
	require.Equal(t, int16(4), foo.Int16)                 // from config
	require.Equal(t, int8(5), foo.Int8)                   // from config
	require.Equal(t, uint(6), foo.Uint)                   // from config
	require.Equal(t, uint64(7), foo.Uint64)               // from config
	require.Equal(t, uint32(8), foo.Uint32)               // from config
	require.Equal(t, uint16(9), foo.Uint16)               // from config
	require.Equal(t, uint8(10), foo.Uint8)                // from config
	require.Equal(t, 11.11, foo.Float64)                  // from config
	require.Equal(t, float32(12.12), foo.Float32)         // from config
	require.Equal(t, false, foo.Bool)                     // from config

	// Recursive field
	require.NotNil(t, foo.Recursive)
	require.Equal(t, "this value also came from config.toml", foo.Recursive.String) // from config
	require.Equal(t, -888, foo.Recursive.Int)                                       // from env
	require.Equal(t, int64(14), foo.Recursive.Int64)                                // from config
	require.Equal(t, int32(15), foo.Recursive.Int32)                                // from config
	require.Equal(t, int16(16), foo.Recursive.Int16)                                // from config
	require.Equal(t, int8(17), foo.Recursive.Int8)                                  // from config
	require.Equal(t, uint(18), foo.Recursive.Uint)                                  // from config
	require.Equal(t, uint64(19), foo.Recursive.Uint64)                              // from config
	require.Equal(t, uint32(20), foo.Recursive.Uint32)                              // from config
	require.Equal(t, uint16(21), foo.Recursive.Uint16)                              // from config
	require.Equal(t, uint8(22), foo.Recursive.Uint8)                                // from config
	require.Equal(t, 23.23, foo.Recursive.Float64)                                  // from config
	require.Equal(t, float32(24.24), foo.Recursive.Float32)                         // from config
	require.Equal(t, true, foo.Recursive.Bool)                                      // from config

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