package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGeneratorConfigVerify(t *testing.T) {
	t.Run("default config should pass validation", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		err := cfg.Verify()
		assert.NoError(t, err)
	})

	t.Run("invalid MbPerSecond", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.MbPerSecond = 0
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MbPerSecond")
	})

	t.Run("invalid BlobSizeMb", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.BlobSizeMb = -1
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "BlobSizeMb")
	})

	t.Run("invalid ValidatorVerificationFraction - too high", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.ValidatorVerificationFraction = 1.5
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValidatorVerificationFraction")
	})

	t.Run("invalid ValidatorVerificationFraction - negative", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.ValidatorVerificationFraction = -0.1
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ValidatorVerificationFraction")
	})

	t.Run("invalid SubmissionParallelism", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.SubmissionParallelism = 0
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SubmissionParallelism")
	})

	t.Run("invalid PprofHttpPort when enabled", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.EnablePprof = true
		cfg.PprofHttpPort = 70000
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PprofHttpPort")
	})

	t.Run("invalid PprofHttpPort ignored when disabled", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.EnablePprof = false
		cfg.PprofHttpPort = 70000
		err := cfg.Verify()
		assert.NoError(t, err)
	})

	t.Run("negative FrequencyAcceleration", func(t *testing.T) {
		cfg := DefaultLoadGeneratorConfig()
		cfg.FrequencyAcceleration = -1.0
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "FrequencyAcceleration")
	})
}

func TestTrafficGeneratorConfigVerify(t *testing.T) {
	t.Run("invalid load config propagates error", func(t *testing.T) {
		cfg := DefaultTrafficGeneratorConfig()
		cfg.Load.MbPerSecond = 0
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "load generator config verification failed")
	})

	t.Run("invalid environment config propagates error", func(t *testing.T) {
		cfg := DefaultTrafficGeneratorConfig()
		// Default environment config has required fields that aren't set
		err := cfg.Verify()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "environment config verification failed")
	})
}
