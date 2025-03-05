package config

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda-proxy/verify/v1"
	"github.com/Layr-Labs/eigenda/api/clients"
	v2_clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/stretchr/testify/require"
)

func validCfg() ProxyConfig {
	maxBlobLengthBytes, err := common.ParseBytesAmount("2MiB")
	if err != nil {
		panic(err)
	}
	proxyCfg := ProxyConfig{
		EdaClientConfigV1: clients.EigenDAClientConfig{
			RPC:                          "http://localhost:8545",
			StatusQueryRetryInterval:     5 * time.Second,
			StatusQueryTimeout:           30 * time.Minute,
			DisableTLS:                   true,
			ResponseTimeout:              10 * time.Second,
			CustomQuorumIDs:              []uint{1, 2, 3},
			SignerPrivateKeyHex:          "private-key-hex",
			PutBlobEncodingVersion:       0,
			DisablePointVerificationMode: false,
			SvcManagerAddr:               "0x00000000069",
			EthRpcUrl:                    "http://localhosts",
		},
		EdaVerifierConfigV1: verify.Config{
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
		MemstoreEnabled: false,
		MemstoreConfig: memconfig.NewSafeConfig(
			memconfig.Config{
				BlobExpiration: 25 * time.Minute,
			}),
		EdaClientConfigV2: common.ClientConfigV2{
			Enabled: true,
			DisperserClientCfg: v2_clients.DisperserClientConfig{
				Hostname:          "http://localhost",
				Port:              "9999",
				UseSecureGrpcFlag: true,
			},
			ServiceManagerAddress:      "0x1234567890abcdef",
			EigenDACertVerifierAddress: "0x0000000000032443134",
		},
		EigenDAV2Enabled: true,
	}

	return proxyCfg
}

func TestConfigVerification(t *testing.T) {
	t.Run(
		"ValidConfig", func(t *testing.T) {
			cfg := validCfg()

			err := cfg.Check()
			require.NoError(t, err)
		})

	t.Run(
		"CertVerificationEnabled", func(t *testing.T) {
			// when eigenDABackend is enabled (memstore.enabled = false),
			// some extra fields are required.
			t.Run(
				"MissingSvcManagerAddr", func(t *testing.T) {
					cfg := validCfg()
					// cert verification only makes sense when memstore is disabled (we use eigenda as backend)
					cfg.MemstoreEnabled = false
					cfg.EdaVerifierConfigV1.VerifyCerts = true
					cfg.EdaVerifierConfigV1.SvcManagerAddr = ""

					err := cfg.Check()
					require.Error(t, err)
				})

			t.Run(
				"MissingEthRPC", func(t *testing.T) {
					cfg := validCfg()
					// cert verification only makes sense when memstore is disabled (we use eigenda as backend)
					cfg.MemstoreEnabled = false
					cfg.EdaVerifierConfigV1.VerifyCerts = true
					cfg.EdaVerifierConfigV1.RPCURL = ""

					err := cfg.Check()
					require.Error(t, err)
				})

			t.Run(
				"CantDoCertVerificationWhenMemstoreEnabled", func(t *testing.T) {
					cfg := validCfg()
					cfg.MemstoreEnabled = true
					cfg.EdaVerifierConfigV1.VerifyCerts = true

					err := cfg.Check()
					require.Error(t, err)
				})

			t.Run(
				"EigenDAClientFieldsAreDefaultSetWhenMemStoreEnabled", func(t *testing.T) {
					cfg := validCfg()
					cfg.MemstoreEnabled = true
					cfg.EdaVerifierConfigV1.VerifyCerts = false
					cfg.EdaVerifierConfigV1.RPCURL = ""
					cfg.EdaVerifierConfigV1.SvcManagerAddr = ""

					err := cfg.Check()
					require.NoError(t, err)
					require.True(t, len(cfg.EdaClientConfigV1.EthRpcUrl) > 1)
					require.True(t, len(cfg.EdaClientConfigV1.SvcManagerAddr) > 1)
				})

			t.Run(
				"FailWhenEigenDAClientFieldsAreUnsetAndMemStoreDisabled", func(t *testing.T) {
					cfg := validCfg()
					cfg.MemstoreEnabled = false
					cfg.EdaClientConfigV1.EthRpcUrl = ""
					cfg.EdaClientConfigV1.SvcManagerAddr = ""

					err := cfg.Check()
					require.Error(t, err)
				})
			t.Run(
				"FailWhenRequiredEigenDAV2FieldsAreUnset", func(t *testing.T) {
					cfg := validCfg()
					cfg.EdaClientConfigV2.DisperserClientCfg.Hostname = ""
					require.Error(t, cfg.Check())
				})
		})

}
