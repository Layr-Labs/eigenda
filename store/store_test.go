package store_test

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/stretchr/testify/require"
)

func validCfg() *store.Config {
	return &store.Config{
		RedisConfig: redis.Config{
			Endpoint: "localhost:6379",
			Password: "password",
			DB:       0,
			Eviction: 10 * time.Minute,
		},
		S3Config: s3.Config{
			Bucket:          "test-bucket",
			Path:            "",
			Endpoint:        "http://localhost:9000",
			EnableTLS:       false,
			AccessKeyID:     "access-key-id",
			AccessKeySecret: "access-key-secret",
		},
	}
}

func TestConfigVerification(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := validCfg()

		err := cfg.Check()
		require.NoError(t, err)
	})

	t.Run("MissingS3AccessKeys", func(t *testing.T) {
		cfg := validCfg()

		cfg.S3Config.CredentialType = s3.CredentialTypeStatic
		cfg.S3Config.Endpoint = "http://localhost:9000"
		cfg.S3Config.AccessKeyID = ""

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("MissingS3Credential", func(t *testing.T) {
		cfg := validCfg()

		cfg.S3Config.CredentialType = s3.CredentialTypeUnknown

		err := cfg.Check()
		require.Error(t, err)
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

	t.Run("BadRedisConfiguration", func(t *testing.T) {
		cfg := validCfg()
		cfg.RedisConfig.Endpoint = ""

		err := cfg.Check()
		require.Error(t, err)
	})
}
