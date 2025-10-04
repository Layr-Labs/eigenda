package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func validCfg() *Config {
	return &Config{}
}

func TestConfigVerification(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := validCfg()

		err := cfg.Check()
		require.NoError(t, err)
	})

	t.Run("InvalidFallbackTarget", func(t *testing.T) {
		cfg := validCfg()
		cfg.FallbackTargets = []string{"postgres"}

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("InvalidCacheTarget", func(t *testing.T) {
		cfg := validCfg()
		cfg.CacheTargets = []string{"postgres"}

		err := cfg.Check()
		require.Error(t, err)
	})
	t.Run("InvalidCacheTarget", func(t *testing.T) {
		cfg := validCfg()
		cfg.CacheTargets = []string{"postgres"}

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("DuplicateCacheTargets", func(t *testing.T) {
		cfg := validCfg()
		cfg.CacheTargets = []string{"s3", "s3"}

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("DuplicateFallbackTargets", func(t *testing.T) {
		cfg := validCfg()
		cfg.FallbackTargets = []string{"s3", "s3"}

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("OverlappingCacheFallbackTargets", func(t *testing.T) {
		cfg := validCfg()
		cfg.FallbackTargets = []string{"s3"}
		cfg.CacheTargets = []string{"s3"}

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("ErrorOnSecondaryInsertFailure: flag OFF, async OFF", func(t *testing.T) {
		cfg := validCfg()
		cfg.AsyncPutWorkers = 0
		cfg.ErrorOnSecondaryInsertFailure = false

		err := cfg.Check()
		require.NoError(t, err)
	})

	t.Run("ErrorOnSecondaryInsertFailure: flag OFF, async ON", func(t *testing.T) {
		cfg := validCfg()
		cfg.AsyncPutWorkers = 5
		cfg.ErrorOnSecondaryInsertFailure = false

		err := cfg.Check()
		require.NoError(t, err)
	})

	t.Run("ErrorOnSecondaryInsertFailure: flag ON, async OFF", func(t *testing.T) {
		cfg := validCfg()
		cfg.AsyncPutWorkers = 0
		cfg.ErrorOnSecondaryInsertFailure = true

		err := cfg.Check()
		require.NoError(t, err)
	})

	t.Run("ErrorOnSecondaryInsertFailure: flag ON, async ON (invalid)", func(t *testing.T) {
		cfg := validCfg()
		cfg.AsyncPutWorkers = 5
		cfg.ErrorOnSecondaryInsertFailure = true

		err := cfg.Check()
		require.Error(t, err)
		require.Contains(t, err.Error(), "requires synchronous writes")
	})
}
