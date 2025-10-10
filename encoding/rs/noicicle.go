//go:build !icicle

package rs

import (
	"errors"
)

func createIcicleBackend(enableGPU bool) (EncoderDevice, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}
