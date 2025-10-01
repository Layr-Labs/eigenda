package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// ParseConfigFromCLI parses configuration, pulling config paths from command line arguments. Assumes that any command
// line argument, if present, is a config path. The resulting config is written into target. If there should be default
// values in the config, target should be initialized with those default values before calling this function.
func ParseConfigFromCLI[T VerifiableConfig](constructor func() T, envPrefix string) (T, error) {
	configPaths := make([]string, 0)
	configPaths = append(configPaths, os.Args...)

	cfg, err := ParseConfig(constructor, envPrefix, configPaths...)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse config from CLI: %w", err)
	}

	return cfg, nil
}

// ParseConfig parses the configuration from the given paths and environment variables. Configuration files are
// loaded in order, with later files overriding earlier ones. Environment variables are loaded last, and overrid values
// from all configuration files. If there are default values in the config, target should be initialized with those
// default values before calling this function.
//
// The envPrefix is used to prefix environment variables. May not be empty. Environment variables with this prefix
// are required to be bound to a configuration field, otherwise an error is returned. TODO expand docs
func ParseConfig[T VerifiableConfig](constructor func() T, envPrefix string, configPaths ...string) (T, error) {
	if envPrefix == "" {
		var zero T
		return zero, fmt.Errorf("envPrefix may not be empty")
	}

	target := constructor()

	// Configure viper.
	viperInstance := viper.New()
	viperInstance.SetEnvPrefix(envPrefix)
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperInstance.AutomaticEnv()

	// Load each config file in order.
	for i, path := range configPaths {
		err := loadConfigFile(viperInstance, path, i == 0)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("failed to load config file %q: %w", path, err)
		}
	}

	// Walk the struct and figure out what environment variables to bind to it.
	boundVars, err := bindEnvs(viperInstance, envPrefix, target)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to bind environment variables: %w", err)
	}

	// Make sure there aren't any invalid environment variables set.
	err = checkForInvalidEnvVars(boundVars, envPrefix)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("invalid environment variables: %w", err)
	}

	// Use viper to unmarshal environment variables into the target struct.
	decoderConfig := &mapstructure.DecoderConfig{
		ErrorUnused:      true,
		WeaklyTypedInput: true, // Allow automatic type conversion from strings (e.g., env vars)
		Result:           target,
		TagName:          "mapstructure",
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to create decoder: %w", err)
	}
	if err := decoder.Decode(viperInstance.AllSettings()); err != nil {
		var zero T
		return zero, fmt.Errorf("failed to decode settings: %w", err)
	}

	// Verify configuration invariants.
	err = target.Verify()
	if err != nil {
		var zero T
		return zero, fmt.Errorf("invalid configuration: %w", err)
	}

	return target, nil
}

func loadConfigFile(v *viper.Viper, path string, firstConfig bool) error {
	path, err := util.SanitizePath(path) // TODO extract sanitize path from LittDB
	if err != nil {
		return fmt.Errorf("failed to sanitize config path %q: %w", path, err)
	}

	exists, err := util.Exists(path)
	if err != nil {
		return fmt.Errorf("failed to check if config path %q exists: %w", path, err)
	}
	if !exists {
		return fmt.Errorf("config path %q does not exist", path)
	}

	v.SetConfigFile(path)
	if firstConfig {
		err = v.ReadInConfig()
		if err != nil {
			return fmt.Errorf("failed to read config file %q: %w", path, err)
		}
	} else {
		err = v.MergeInConfig()
		if err != nil {
			return fmt.Errorf("failed to merge config file %q: %w", path, err)
		}
	}

	return nil
}

// Walks a struct tree and automatically binds each field to an environment variable based on the given prefix
// and the field's path in the struct tree. For example, given a struct like:
//
//	type MyStruct struct {
//	    Foo string
//	    Bar struct {
//	        Baz int
//	    }
//	}
//
// and a prefix of "MYSTRUCT", this function will bind the following environment variables:
//
//	MYSTRUCT_FOO -> Foo
//	MYSTRUCT_BAR_BAZ -> Bar.Baz
//
// This function uses reflection to walk the struct tree, so it only works with exported fields. It also only works
// with basic types (string, int, bool, etc.) and nested structs. It does not work with slices, maps, or other complex
// types.
//
// This function is recursive, so it will walk nested structs to any depth.
//
// This function returns a set containing the names of all environment variables that were bound. This is used
// to detect unused environment variables (which are likely misconfigurations).
func bindEnvs(v *viper.Viper, prefix string, target any, path ...string) (map[string]struct{}, error) {
	boundVars := make(map[string]struct{})

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
	targetType := targetValue.Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		name := strings.ToLower(field.Name)
		keyPath := append(path, name)

		switch field.Type.Kind() {

		case reflect.Struct:
			// Recurse for nested structs
			tmp := reflect.New(field.Type).Elem().Interface()
			nestedBoundVars, err := bindEnvs(v, prefix, tmp, keyPath...)
			if err != nil {
				return nil, fmt.Errorf("failed to bind envs for field %s: %w", field.Name, err)
			}
			for k := range nestedBoundVars {
				boundVars[k] = struct{}{}
			}
		case reflect.Ptr:
			// Handle pointer to struct
			if field.Type.Elem().Kind() == reflect.Struct {
				tmp := reflect.New(field.Type.Elem()).Interface()
				nestedBoundVars, err := bindEnvs(v, prefix, tmp, keyPath...)
				if err != nil {
					return nil, fmt.Errorf("failed to bind envs for field %s: %w", field.Name, err)
				}
				for k := range nestedBoundVars {
					boundVars[k] = struct{}{}
				}
			} else {
				// Pointer to non-struct type, bind as regular field
				env := prefix + "_" + strings.ToUpper(strings.ReplaceAll(strings.Join(keyPath, "_"), ".", "_"))
				boundVars[env] = struct{}{}
				if err := v.BindEnv(strings.Join(keyPath, "."), env); err != nil {
					return nil, err
				}
			}
		default:
			env := prefix + "_" + strings.ToUpper(strings.ReplaceAll(strings.Join(keyPath, "_"), ".", "_"))
			boundVars[env] = struct{}{}
			if err := v.BindEnv(strings.Join(keyPath, "."), env); err != nil {
				return nil, err
			}
		}
	}

	return boundVars, nil
}

// checkForInvalidEnvVars checks for any environment variables with the given prefix that were not bound to any
// configuration fields. This is used to detect misconfigurations where an environment variable is set, but it does
// not correspond to any configuration field (e.g. due to a typo).
//
// This function returns an error if any invalid environment variables are found.
func checkForInvalidEnvVars(boundVars map[string]struct{}, envPrefix string) error {

	if envPrefix == "" {
		// Nothing we can do if there is no prefix.
		return nil
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		if !strings.HasPrefix(key, envPrefix+"_") {
			continue
		}
		if _, ok := boundVars[key]; !ok {
			return fmt.Errorf("environment variable %q is not bound to any configuration field", key)
		}
	}

	return nil
}
