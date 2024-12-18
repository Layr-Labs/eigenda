# What and Why?

This package provides simple light-weight configuration management library. Its goal is to reduce the amount
of boilerplate code required to manage configuration but without the complexities of a large framework.
It's essentially just a few wrapper methods around the standard golang json decoder.

## Simplicity, ~~Simplicity, Simplicity~~

The design ethos of this framework is that configuration should be super simple. There shouldn't be lots of different
ways of configuring things. This framework does not attempt to support any of the following features:

- Environment variables
- Command line flags
- Non-json configuration files
- Configuration files that change at runtime
- Configuration validation
- Required fields (this is a form of configuration validation)

# How To Use

## Basic Operation

First, define a struct that represents your configuration. For example:

```go
type BasicConfig struct {
	Foo string
	Bar int
	Baz bool
}
```

Put your configuration into a json file. For example, `config.json`:

```json
{
    "Foo": "Hello",
    "Bar": 42,
    "Baz": true
}
```

Then, at runtime, load the json file into the stuct via

```go

// Put default values here.
myConfig := BasicConfig{
    Foo: "default",
    Bar: 0,
    Baz: false,
}

err := configuration.Load(&myConfig, "config.json")

if err == nil {
    // Done! Configuration is now in the myConfig object.
}
```

## Default Values

If a field is defined in a struct but the field is not present in the json configuration file, then the field will
use its default value. The default value for a field is set by assigning a value to the struct that is passed
into the `Load()` function.

## Invalid Configuration

A configuration is considered to be invalid if there is a value in the json file that cannot be mapped to a field
in the struct. This is a safety feature. If there are unused values in the json file, either there is a typo
in the configuration file, or there is crud that should be removed.

## Layered Config Files

It's possible to provide multiple configuration files to the `Load()` function. The configuration files are loaded
in order, with the last file taking precedence. For example:

```go
myConfig := BasicConfig{
    Foo: "default",
    Bar: 0,
    Baz: false,
}

err := configuration.Load(&myConfig, "config.json", "config.local.json")

if err == nil {
    // Done! Configuration is now in the myConfig object.
}
```

In this example, the values in `config.json` override the default values, and the values in `config.local.json` override
both the default values and the values in `config.json`.

## Supported Data Types

The configuration framework technically supports all data types that the json decoder supports. However, the following
data types have unit tests. In general, it's probably wise to avoid using types not explicitly tested.

### Primitives

- `bool`: use the values `true` or `false`
- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `uint`
- `uint8`
- `uint16`
- `uint32`
- `uint64`
- `float32`: warning, precision may be lost
- `float64`: warning, precision may be lost
- `string`
- `time.Time`: use the format `RFC3339` (e.g. `2006-01-02T15:04:05Z07:00`)
- `time.Duration`: specified as an integer number of nanoseconds

### Structs

The configuration framework supports arbitrary structs containing primitives, maps, arrays,
and other arbitrarily nested structs.

### Maps and Lists

The configuration framework supports maps and lists.

Maps may use any primitive type as a key, and may use primitives, structs, maps, or arrays as values. 
Arrays may contain any primitive, struct, map, or list.
