package server

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/stretchr/testify/require"
)

func validCfg() *Config {
	return &Config{
		RedisCfg: store.RedisConfig{
			Endpoint: "localhost:6379",
			Password: "password",
			DB:       0,
			Eviction: 10 * time.Minute,
		},
		S3Config: store.S3Config{
			Bucket:          "test-bucket",
			Path:            "",
			Endpoint:        "http://localhost:9000",
			AccessKeyID:     "access-key-id",
			AccessKeySecret: "access-key-secret",
		},
		ClientConfig: clients.EigenDAClientConfig{
			RPC:                          "http://localhost:8545",
			StatusQueryRetryInterval:     5 * time.Second,
			StatusQueryTimeout:           30 * time.Minute,
			DisableTLS:                   true,
			ResponseTimeout:              10 * time.Second,
			CustomQuorumIDs:              []uint{1, 2, 3},
			SignerPrivateKeyHex:          "private-key-hex",
			PutBlobEncodingVersion:       0,
			DisablePointVerificationMode: false,
		},
		G1Path:                  "path/to/g1",
		G2PowerOfTauPath:        "path/to/g2",
		CacheDir:                "path/to/cache",
		MaxBlobLength:           "2MiB",
		CertVerificationEnabled: false,
		SvcManagerAddr:          "0x1234567890abcdef",
		EthRPC:                  "http://localhost:8545",
		EthConfirmationDepth:    12,
		MemstoreEnabled:         true,
		MemstoreBlobExpiration:  25 * time.Minute,
	}
}

func TestConfigVerification(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := validCfg()

		err := cfg.Check()
		require.NoError(t, err)
	})

	t.Run("InvalidMaxBlobLength", func(t *testing.T) {
		cfg := validCfg()
		cfg.MaxBlobLength = "0kzg"

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("0MaxBlobLength", func(t *testing.T) {
		cfg := validCfg()
		cfg.MaxBlobLength = "0kib"

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("CertVerificationEnabled", func(t *testing.T) {
		// when eigenDABackend is enabled (memstore.enabled = false),
		// some extra fields are required.
		t.Run("MissingSvcManagerAddr", func(t *testing.T) {
			cfg := validCfg()
			// cert verification only makes sense when memstore is disabled (we use eigenda as backend)
			cfg.MemstoreEnabled = false
			cfg.CertVerificationEnabled = true
			cfg.SvcManagerAddr = ""

			err := cfg.Check()
			require.Error(t, err)
		})

		t.Run("MissingEthRPC", func(t *testing.T) {
			cfg := validCfg()
			// cert verification only makes sense when memstore is disabled (we use eigenda as backend)
			cfg.MemstoreEnabled = false
			cfg.CertVerificationEnabled = true
			cfg.EthRPC = ""

			err := cfg.Check()
			require.Error(t, err)
		})

		t.Run("CantDoCertVerificationWhenMemstoreEnabled", func(t *testing.T) {
			cfg := validCfg()
			cfg.MemstoreEnabled = true
			cfg.CertVerificationEnabled = true

			err := cfg.Check()
			require.Error(t, err)
		})
	})

	t.Run("MissingS3AccessKeys", func(t *testing.T) {
		cfg := validCfg()

		cfg.S3Config.S3CredentialType = store.S3CredentialStatic
		cfg.S3Config.Endpoint = "http://localhost:9000"
		cfg.S3Config.AccessKeyID = ""

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("MissingS3Credential", func(t *testing.T) {
		cfg := validCfg()

		cfg.S3Config.S3CredentialType = store.S3CredentialUnknown

		err := cfg.Check()
		require.Error(t, err)
	})

	t.Run("MissingEigenDADisperserRPC", func(t *testing.T) {
		cfg := validCfg()
		cfg.ClientConfig.RPC = ""
		cfg.MemstoreEnabled = false

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
		cfg.RedisCfg.Endpoint = ""

		err := cfg.Check()
		require.Error(t, err)
	})
}
