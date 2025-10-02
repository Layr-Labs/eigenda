# Configuration Management

This configuration "framework" attempts to achieve maximal simplicity when it comes to creating, modifying, and
and maintaining configuration. Configuration is inherently a simple concept, and so the execution of configuration
should likewise be simple.

# Config is Struct

```
grug say, me want config in this struct. why config need be more than struct?
```

In order to define configuration, the user of this framework provides a simple struct that meets the following
requirements:

1. All variables must be exported.
2. Variables must all be "simple" types.
    - any primitive (`int`, `float`, `string`, etc.)
    - `time.Duration`
    - nested structs that themselves only contain simple types (recursive type nesting not permitted)
    - pointers to any of the above
3. The struct must implement the `config.VerifiableConfig` interface (see below).
4. The config must have a default constructor method.


```go
// VerifiableConfig is an interface for configurations that can be validated.
type VerifiableConfig interface {
	// Verify checks that the configuration is valid, returning an error if it is not.
	Verify() error
}
```

Although in theory the `Verify()` method can be a no-op, it is highly recommended to implement basic sanity checking.

The "constructor" for a config object must satisfy the interface `func() T` where `T` 
implements `config.VerifiableConfig`.

# How to Load Config

## ParseConfig()

There are two ways config can be loaded.

The first way is to pass in a list of zero or more configuration files to `ParseConfig()`.

```go
ParseConfig[T VerifiableConfig](constructor func() T, envPrefix string, configPaths ...string) (T, error)
```

`ParseConfig()` will load data from the configuration files in order (later files override values from earlier files).
After loading configuration files, `ParseConfig()` loads environment variables (overriding values set by config files).

Example:

```go
const MyAppPrefix = "MYAPP"

type MyConfig struct {
    // ...
}

func DefaultMyConfig() *MyConfig {
    return ...
}

cfg, err := config.ParseConfig(DefaultMyConfig, MyAppPrefix, "path/to/config1.toml", ..., "path/to/configN.toml")

// cfg is an instance of MyConfig that now contains all loaded config
```

## ParseConfigFromCLI()

An alternate way of parsing configuration is with the following method:

```
ParseConfigFromCLI[T VerifiableConfig](constructor func() T, envPrefix string) (T, error)
```

The primary difference between this method and `ParseConfig()` is that `ParseConfigFromCLI()` assumes that any
command line arguments provided to the process are configuration file paths. Using this method is functionally
equivalent to parsing for file paths from the CLI arguments, then passing those file paths into `ParseConfig()`.

Although use of `ParseConfigFromCLI()` is not required to use this framework, it is highly encouraged. Any time
configuration is sourced through multiple pathways, complexity grows. If configuration is large enough that the
config framework is needed, then it's best if all configuration flows through the config framework. A hodge-podge
of CLI arguments mixed with the configuration framework adds a lot of unnecessary complexity.

# Documenting Config

The proper place for documenting configuration is godocs in the configuration struct. A well documented struct should 
be understandable even by people who don't know how to read/write golang. By tightly coupling documentation with code,
it becomes less likely for documentation and implementation to drift apart.

# Configuration Files

The config framework supports configuration files in any format supported by the 
[viper](https://github.com/spf13/viper) framework.


- JSON
- YAML (including .yml)
- TOML
- HCL (HashiCorp Configuration Language)
- dotenv / envfiles (.env)
- Java Properties (.properties)

```
grug like toml, toml is simple. 
grug no like json, json no let grug comment things. 
grug no like yaml, yaml look simple sometimes, but grug know yaml only pretend be simple.
but grug not reach for club if not use toml.
```

In order to set a variable in a config file, simply mirror the struct and variable names in "the obvious way".
Below is an example using toml.

```go
type Foo struct {
    X Int
    Y Float
    Z String
    Bar Bar
}

type Bar struct {
    A String
    B Duration
    C Baz
}

type Baz struct {
    ValueStoredInAVariableWithALongName String
}
```

The following TOML file can be loaded into the structs above.

```toml
X = 1234
Y = 3.14159265359
Z = "this is a string"
Bar.A = "this is another string"
Bar.B = "5s"
Bar.C.ValueStoredInAVariableWithALongName = "yet another string"
```

## Mistyped Config

If there is a config value that does not have a corresponding entry in the config struct, the configuration framework
will return an error when it attempts to parse the config. This is very intentional. Unmatched config file entries
almost always signal a mistake in the configuration files. At the very least, returning an error for unmatched config
forces config files to be kept clean and well maintained.

# Environment Variables

The config framework supports loading config from environment variables. Although the primary intended use case for
environment variables is for loading secrets, there is nothing stopping configuration from being loaded entirely
from environment variables.

The configuration framework requires that a prefix be defined for environment variables. By convention, this prefix
should contain only upper case letters and underscores.

For each entry in a config struct, there is an environment variable that is mapped to that entry. The name of the
environment variable is `PREFIX_NAMEOFVARIABLE`. If the variable is in a nested struct, for each "parent variable",
add the name of the parent variable in uppercase, and separate parent variables with underscores. 

The following example shows the names of the environment variables that could be used to configure the following
struct.

```go
const MyPrefix = "PREFIX"

type Foo struct {
    X Int                                      // PREFIX_X
    Y Float                                    // PREFIX_Y
    Z String                                   // PREFIX_Z
    Bar Bar
}

type Bar struct {
    A String                                   // PREFIX_BAR_A
    B Duration                                 // PREFIX_BAR_B
    C Baz
}

type Baz struct {
    ValueStoredInAVariableWithALongName String // PREFIX_BAR_C_VALUESTOREDINAVARIABLEWITHALONGNAME
}
```

## Mistyped Environment Variables

The config framework looks at all environment variables that begin with the prefix. If it finds any environment
variable with the prefix that does not map to an entry in the config struct, it returns an error. This is intentional.
Similar to mistyped config, an environment variable that doesn't map to a config entry is likely to be a bug.

# Default Values

The purpose of the constructor is to set default values in the struct. The config API requires a constructor method
in order to strongly encourage users of this framework to set sane default values where possible. In general,
the fewer values that are required to be set, the easier it is to configure something.

# Required Values

If there are values that must be set by the end user, then return an error with an appropriate message in `Verify()` 
if any of those values are unset.