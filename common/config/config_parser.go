package config

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/Layr-Labs/eigenda/common/config/secret"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// ParseConfig parses the configuration from the given paths and environment variables. Configuration files are
// loaded in order, with later files overriding earlier ones. Environment variables are loaded last, and override values
// from all configuration files. If there are default values in the config, those values should be in the provided cfg.
func ParseConfig[T VerifiableConfig](
	// Used to log debug information about environment variables if something goes wrong.
	logger logging.Logger,
	// The configuration to populate, should already contain any default values.
	cfg T,
	// The prefix to use for environment variables. If empty, then environment variables are not read.
	envPrefix string,
	// A map of environment variable aliases. The keys are environment variables that should be aliased to something
	// else, and the values are the environment variables they should be aliased to.
	//
	// Environment variables in this map should be fully qualified, including any prefixes.
	//
	// If nil, then no aliasing is performed.
	aliasedEnvVars map[string]string,
	// A list of environment variables that should be ignored when sanity checking environment variables.
	// Useful for situations where external systems set environment variables that would otherwise cause problems.
	//
	// Environment variables in this list should be fully qualified, including any prefixes.
	//
	// If nil, then no environment variables are ignored during sanity checking.
	ignoredEnvVars []string,
	// A list of zero or more paths to configuration files. Later files override earlier ones.
	// If environment variables are read, they override all configuration files.
	configPaths ...string,
) (T, error) {
	viperInstance := viper.New()

	// Load each config file in order.
	for i, path := range configPaths {
		err := loadConfigFile(viperInstance, path, i == 0)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("failed to load config file %q: %w", path, err)
		}
	}

	if envPrefix != "" {
		err := aliasEnvVars(logger, aliasedEnvVars)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("failed to alias environment variables: %w", err)
		}

		viperInstance.SetEnvPrefix(envPrefix)
		viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viperInstance.AutomaticEnv()

		// Walk the struct and figure out what environment variables to bind to it.
		boundVars, err := bindEnvs(viperInstance, envPrefix, cfg)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("failed to bind environment variables: %w", err)
		}

		// Make sure there aren't any invalid environment variables set.
		err = checkForInvalidEnvVars(logger, boundVars, envPrefix, aliasedEnvVars, ignoredEnvVars)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("invalid environment variables: %w", err)
		}
	}

	decoderConfig := &mapstructure.DecoderConfig{
		ErrorUnused:      true,
		WeaklyTypedInput: true, // Allow automatic type conversion from strings (e.g., env vars)
		Result:           cfg,
		TagName:          "mapstructure",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(), // Support time.Duration parsing from strings
			secret.DecodeHook,        // Support Secret parsing
			basicTypeSliceDecodeHook, // Support slices of basic types
		),
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
	err = cfg.Verify()
	if err != nil {
		var zero T
		return zero, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Applies environment variable aliases by copying the value of each aliased variable to its target variable.
// This function sets new environment variables using os.Setenv if old environment variables in need of
// aliasing are set.
func aliasEnvVars(logger logging.Logger, aliasedEnvVars map[string]string) error {
	for oldVar, newVar := range aliasedEnvVars {
		value, oldVarExists := os.LookupEnv(oldVar)

		if oldVarExists {
			_, newVarExists := os.LookupEnv(newVar)
			if newVarExists {
				return fmt.Errorf("cannot alias environment variable %q to %q: both are set", oldVar, newVar)
			}

			logger.Warnf("Deprecated environment variable %q is set; please use %q instead. "+
				"Support for this environment variable may be removed in a future release.", oldVar, newVar)

			err := os.Setenv(newVar, value)
			if err != nil {
				return fmt.Errorf("failed to set aliased environment variable %q: %w", newVar, err)
			}
		}
	}
	return nil
}

func loadConfigFile(v *viper.Viper, path string, firstConfig bool) error {
	path, err := util.SanitizePath(path)
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
// with basic types (string, int, bool, etc.), slices of basic types, nested structs, secret.Secret, and slices of 
// secret.Secret. It does not work with maps or other complex types.
//
// This function is recursive, so it will walk nested structs to any depth.
//
// This function returns a set containing the names of all environment variables that were bound. This is used
// to detect unused environment variables (which are likely misconfigurations).
func bindEnvs(
	// The viper instance that will parse environment variables.
	v *viper.Viper,
	// The prefix to use for environment variables.
	prefix string,
	// The struct to walk.
	target any,
	// The "path" to the current struct in the tree. This should be empty when calling this function initially.
	// Each step in the path is the name of a field in the config struct.
	path ...string,
) (map[string]struct{}, error) {

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

		// Get the mapstructure tag, or use field name if tag is not present
		fieldKey := field.Name
		if tag := field.Tag.Get("mapstructure"); tag != "" {
			fieldKey = tag
		}

		keyPath := append(path, fieldKey)

		switch field.Type.Kind() { //nolint:exhaustive // only handling struct, pointer, and slice types

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
		case reflect.Slice:
			// Handle slices
			elemType := field.Type.Elem()
			// Check if this is a slice of pointers to Secret
			if elemType.Kind() == reflect.Ptr &&
				elemType.Elem().Kind() == reflect.Struct &&
				elemType.Elem() == reflect.TypeOf((*secret.Secret)(nil)).Elem() {
				// Slice of *secret.Secret, bind as leaf value
				env := buildEnvVarName(prefix, keyPath...)
				boundVars[env] = struct{}{}
				if err := v.BindEnv(strings.Join(keyPath, "."), env); err != nil {
					return nil, fmt.Errorf("failed to bind env %s: %w", env, err)
				}
			} else if isBasicType(elemType) {
				// Slice of basic types (int, string, bool, float, etc.)
				// Bind as leaf value - mapstructure will handle comma-separated conversion
				env := buildEnvVarName(prefix, keyPath...)
				boundVars[env] = struct{}{}
				if err := v.BindEnv(strings.Join(keyPath, "."), env); err != nil {
					return nil, fmt.Errorf("failed to bind env %s: %w", env, err)
				}
			}
			// Other slice types (e.g., slices of structs) are not currently supported
			// for environment variable binding and are silently ignored.
		case reflect.Ptr:
			// Handle pointer to struct
			if field.Type.Elem().Kind() == reflect.Struct {
				// Check if this is a Secret type - if so, treat it as a leaf value
				elemType := field.Type.Elem()
				isSecretType := elemType == reflect.TypeOf((*secret.Secret)(nil)).Elem()
				if isSecretType {
					// Secret types should be bound as leaf values, not recursed into
					env := buildEnvVarName(prefix, keyPath...)
					boundVars[env] = struct{}{}
					if err := v.BindEnv(strings.Join(keyPath, "."), env); err != nil {
						return nil, fmt.Errorf("failed to bind env %s: %w", env, err)
					}
				} else {
					// Regular struct, recurse into it
					tmp := reflect.New(field.Type.Elem()).Interface()
					nestedBoundVars, err := bindEnvs(v, prefix, tmp, keyPath...)
					if err != nil {
						return nil, fmt.Errorf("failed to bind envs for field %s: %w", field.Name, err)
					}
					for k := range nestedBoundVars {
						boundVars[k] = struct{}{}
					}
				}
			} else {
				// Pointer to non-struct type, bind as regular field
				env := buildEnvVarName(prefix, keyPath...)
				boundVars[env] = struct{}{}
				if err := v.BindEnv(strings.Join(keyPath, "."), env); err != nil {
					return nil, fmt.Errorf("failed to bind env %s: %w", env, err)
				}
			}
		default:
			env := buildEnvVarName(prefix, keyPath...)
			boundVars[env] = struct{}{}
			if err := v.BindEnv(strings.Join(keyPath, "."), env); err != nil {
				return nil, fmt.Errorf("failed to bind env %s: %w", env, err)
			}
		}
	}

	return boundVars, nil
}

// Derive the name of an environment variable from the given prefix and path.
func buildEnvVarName(prefix string, path ...string) string {
	sb := strings.Builder{}
	sb.WriteString(prefix)

	for _, p := range path {
		sb.WriteString("_")
		sb.WriteString(toScreamingSnakeCase(p))
	}
	return sb.String()
}

// isBasicType checks if a type is a basic type that can be parsed from environment variables.
// This includes primitives (int, uint, float, bool), strings, and pointers to these types.
func isBasicType(t reflect.Type) bool {
	// Handle pointer to basic type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() { //nolint:exhaustive // only handling basic types, default handles all others
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

// basicTypeSliceDecodeHook is a mapstructure decode hook that handles slices of basic types.
// It converts string inputs from config files or environment variables into slices of basic types
// by splitting on commas. This allows environment variables to represent slices using a
// comma-separated format (e.g., "value1,value2,value3").
var basicTypeSliceDecodeHook mapstructure.DecodeHookFunc = func(
	from reflect.Type, to reflect.Type, data any) (any, error) {
	// Only handle string sources
	if from.Kind() != reflect.String {
		return data, nil
	}

	// Only handle slice targets
	if to.Kind() != reflect.Slice {
		return data, nil
	}

	// Only handle slices of basic types
	if !isBasicType(to.Elem()) {
		return data, nil
	}

	// Get the source data as a string
	sourceStr, ok := data.(string)
	if !ok {
		return data, nil
	}

	// If the source string is empty, return an empty slice
	if len(sourceStr) == 0 {
		return reflect.MakeSlice(to, 0, 0).Interface(), nil
	}

	// Split the string by commas
	parts := strings.Split(sourceStr, ",")

	// Create a slice of the target type
	result := reflect.MakeSlice(to, len(parts), len(parts))

	// Convert each part to the target element type using WeakDecode
	// which handles type conversion automatically
	for i, part := range parts {
		trimmedPart := strings.TrimSpace(part)

		// Create a pointer to a new instance of the target element type
		elemPtr := reflect.New(to.Elem())

		// Use WeakDecode directly - it's more efficient than creating a decoder each time
		if err := mapstructure.WeakDecode(trimmedPart, elemPtr.Interface()); err != nil {
			return nil, fmt.Errorf("failed to decode element %d (%q): %w", i, trimmedPart, err)
		}

		// Set the element in the result slice
		result.Index(i).Set(elemPtr.Elem())
	}

	return result.Interface(), nil
}

// checkForInvalidEnvVars checks for any environment variables with the given prefix that were not bound to any
// configuration fields. This is used to detect misconfigurations where an environment variable is set, but it does
// not correspond to any configuration field (e.g. due to a typo).
//
// This function returns an error if any invalid environment variables are found.
func checkForInvalidEnvVars(
	logger logging.Logger,
	boundVars map[string]struct{},
	envPrefix string,
	aliasedEnvVars map[string]string,
	ignoredEnvVars []string,
) error {
	if envPrefix == "" {
		// Nothing we can do if there is no prefix.
		return nil
	}

	ignoredSet := make(map[string]struct{}, len(ignoredEnvVars))
	for _, v := range ignoredEnvVars {
		ignoredSet[v] = struct{}{}
	}
	// The config parser will return an error if it discovers an environment variable that doesn't map to a struct
	// value. Since the aliased environment variables indirectly map to struct values, we need to instruct the config
	// parser to ignore them when it's checking for un-mapped environment variables.
	for k := range aliasedEnvVars {
		ignoredSet[k] = struct{}{}
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

		if _, ok := ignoredSet[key]; ok {
			// ignore this variable
			continue
		}

		if _, ok := boundVars[key]; !ok {
			sb := strings.Builder{}
			sb.WriteString("environment variable ")
			sb.WriteString(key)
			sb.WriteString(" is not bound to any configuration field. Legal environment variables are:\n")

			sortedVars := make([]string, 0, len(boundVars))
			for k := range boundVars {
				sortedVars = append(sortedVars, k)
			}
			slices.Sort(sortedVars)

			for _, k := range sortedVars {
				sb.WriteString("  - ")
				sb.WriteString(k)
				sb.WriteString("\n")
			}
			logger.Error(sb.String())

			return fmt.Errorf("environment variable %q is not bound to any configuration field", key)
		}
	}

	return nil
}
