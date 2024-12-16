package common

import (
	"bytes"
	"crypto/sha256"
	"time"
	"unsafe"

	"github.com/fxamacker/cbor/v2"
)

// PrefixEnvVar returns the environment variable name with the given prefix and suffix
func PrefixEnvVar(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	if suffix == "" {
		return prefix
	}
	return prefix + "_" + suffix
}

// PrefixFlag returns the flag name with the given prefix and suffix
func PrefixFlag(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	if suffix == "" {
		return prefix
	}
	return prefix + "." + suffix
}

// Hash returns the sha256 hash of the given value
func Hash[T any](t T) ([]byte, error) {
	bytes, err := EncodeToBytes(t)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	hasher.Write(bytes)
	return hasher.Sum(nil), nil
}

// EncodeToBytes encodes the given value to bytes
func EncodeToBytes[T any](t T) ([]byte, error) {
	size := int(unsafe.Sizeof(t))
	buffer := bytes.NewBuffer(make([]byte, 0, size))
	err := cbor.NewEncoder(buffer).Encode(t)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeFromBytes decodes the given bytes to the given value
func DecodeFromBytes[T any](b []byte) (T, error) {
	var t T
	buffer := bytes.NewBuffer(b)
	err := cbor.NewDecoder(buffer).Decode(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}

// ToMilliseconds converts the given duration to milliseconds. Unlike duration.Milliseconds(), this function returns
// a float64 with nanosecond precision (at least, as much precision as floating point numbers allow).
func ToMilliseconds(duration time.Duration) float64 {
	return float64(duration.Nanoseconds()) / float64(time.Millisecond)
}
