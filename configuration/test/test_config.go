package test

import "time"

// BasicConfig is a config struct with simple types.
type BasicConfig struct {
	Foo string
	Bar int
	Baz bool
}

func DefaultBasicConfig() *BasicConfig {
	return &BasicConfig{
		Foo: "this is a default value",
		Bar: 1337,
		Baz: false,
	}
}

// NestedConfig is a config struct with a nested config structs.
type NestedConfig struct {
	RecursiveConfig *NestedConfig
	BasicConfig     BasicConfig
}

// AllPrimitiveTypes contains all supported primitive types.
type AllPrimitiveTypes struct {
	Bool     bool
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	Uint     uint
	Uint8    uint8
	Uint16   uint16
	Uint32   uint32
	Uint64   uint64
	Float32  float32
	Float64  float64
	String   string
	Time     time.Time
	Duration time.Duration
}
