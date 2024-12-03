package server

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/stretchr/testify/require"
)

func validCfg() *Config {
	maxBlobLengthBytes, err := common.ParseBytesAmount("2MiB")
	if err != nil {
		panic(err)
	}
	return &Config{
		EdaClientConfig: clients.EigenDAClientConfig{
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
		VerifierConfig: verify.Config{
			KzgConfig: &kzg.KzgConfig{
				G1Path:         "path/to/g1",
				G2PowerOf2Path: "path/to/g2",
				CacheDir:       "path/to/cache",
				SRSOrder:       maxBlobLengthBytes / 32,
			},
			VerifyCerts:          false,
			SvcManagerAddr:       "0x1234567890abcdef",
			RPCURL:               "http://localhost:8545",
			EthConfirmationDepth: 12,
		},
		MemstoreEnabled: true,
		MemstoreConfig: memstore.Config{
			BlobExpiration: 25 * time.Minute,
		},
	}
}

func TestConfigVerification(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := validCfg()

		err := cfg.Check()
		require.NoError(t, err)
	})

	t.Run("CertVerificationEnabled", func(t *testing.T) {
		// when eigenDABackend is enabled (memstore.enabled = false),
		// some extra fields are required.
		t.Run("MissingSvcManagerAddr", func(t *testing.T) {
			cfg := validCfg()
			// cert verification only makes sense when memstore is disabled (we use eigenda as backend)
			cfg.MemstoreEnabled = false
			cfg.VerifierConfig.VerifyCerts = true
			cfg.VerifierConfig.SvcManagerAddr = ""

			err := cfg.Check()
			require.Error(t, err)
		})

		t.Run("MissingEthRPC", func(t *testing.T) {
			cfg := validCfg()
			// cert verification only makes sense when memstore is disabled (we use eigenda as backend)
			cfg.MemstoreEnabled = false
			cfg.VerifierConfig.VerifyCerts = true
			cfg.VerifierConfig.RPCURL = ""

			err := cfg.Check()
			require.Error(t, err)
		})

		t.Run("CantDoCertVerificationWhenMemstoreEnabled", func(t *testing.T) {
			cfg := validCfg()
			cfg.MemstoreEnabled = true
			cfg.VerifierConfig.VerifyCerts = true

			err := cfg.Check()
			require.Error(t, err)
		})

		t.Run("EigenDAClientFieldsAreDefaultSetWhenMemStoreEnabled", func(t *testing.T) {
			cfg := validCfg()
			cfg.MemstoreEnabled = true
			cfg.VerifierConfig.VerifyCerts = false
			cfg.VerifierConfig.RPCURL = ""
			cfg.VerifierConfig.SvcManagerAddr = ""

			err := cfg.Check()
			require.NoError(t, err)
			require.True(t, len(cfg.EdaClientConfig.EthRpcUrl) > 1)
			require.True(t, len(cfg.EdaClientConfig.SvcManagerAddr) > 1)
		})

		t.Run("FailWhenEigenDAClientFieldsAreUnsetAndMemStoreDisabled", func(t *testing.T) {
			cfg := validCfg()
			cfg.MemstoreEnabled = false
			cfg.VerifierConfig.RPCURL = ""
			cfg.VerifierConfig.SvcManagerAddr = ""

			err := cfg.Check()
			require.Error(t, err)
		})
	})

}
