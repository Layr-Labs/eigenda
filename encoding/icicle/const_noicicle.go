//go:build !icicle

package icicle

// IsAvailable indicates whether the icicle library is available,
// which is the case when the binary was compiled with the icicle build tag.
// Note that this does not guarantee that the GPU device is available at runtime.
const IsAvailable = false
