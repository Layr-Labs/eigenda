package secret

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

// DecodeHook is a mapstructure decode hook that handles Secret types.
// It converts string inputs from config files or environment variables into Secret instances.
//
// Usage:
//
//	decoderConfig := &mapstructure.DecoderConfig{
//	    DecodeHook: mapstructure.ComposeDecodeHookFunc(
//	        secret.DecodeHook,
//	        // other hooks...
//	    ),
//	}
var DecodeHook mapstructure.DecodeHookFunc = func(from reflect.Type, to reflect.Type, data any) (any, error) {
	// Check if source is a string or []byte
	if from.Kind() != reflect.String && !(from.Kind() == reflect.Slice && from.Elem().Kind() == reflect.Uint8) {
		return data, nil
	}

	// Check if target type is a slice of pointers to Secret
	if to.Kind() == reflect.Slice {
		// Check if the slice element is a pointer to Secret
		if to.Elem().Kind() == reflect.Ptr && to.Elem().Elem() == reflect.TypeOf((*Secret)(nil)).Elem() {
			// Get the source data as a string
			var sourceStr string
			switch v := data.(type) {
			case string:
				sourceStr = v
			case []byte:
				sourceStr = string(v)
			default:
				// If it's not a string or []byte then we can't handle it here
				return nil, fmt.Errorf("cannot convert %v to slice of Secrets", from)
			}

			// If the source string is empty, return an empty slice
			if len(sourceStr) == 0 {
				return []*Secret{}, nil
			}

			// Split the string by commas and create a slice of secrets
			parts := strings.Split(sourceStr, ",")
			secrets := make([]*Secret, len(parts))
			for i, part := range parts {
				secrets[i] = NewSecret(strings.TrimSpace(part))
			}
			return secrets, nil
		}
		return data, nil
	}

	// Check if target type is a pointer to Secret
	if to.Kind() != reflect.Ptr {
		return data, nil
	}

	elem := to.Elem()
	// Check if this is a Secret type
	if elem != reflect.TypeOf((*Secret)(nil)).Elem() {
		return data, nil
	}

	// Get the source data as a string
	var sourceStr string
	switch v := data.(type) {
	case string:
		sourceStr = v
	case []byte:
		sourceStr = string(v)
	default:
		// If it's not a string or []byte, let mapstructure handle it normally
		return data, nil
	}

	return NewSecret(sourceStr), nil
}
