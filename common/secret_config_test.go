package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func validSecretConfig() SecretConfigV2 {
	secretConfig := SecretConfigV2{
		SignerPaymentKey: "0x000000000000000",
		EthRPCURL:        "http://localhost:8545",
	}

	return secretConfig
}

func TestValidSecretConfig(t *testing.T) {
	cfg := validSecretConfig()

	err := cfg.Check()
	require.NoError(t, err)
}

func TestSignerPaymentKeyMissing(t *testing.T) {
	cfg := validSecretConfig()
	cfg.SignerPaymentKey = ""

	err := cfg.Check()
	require.Error(t, err)
}

func TestEthRPCMissing(t *testing.T) {
	cfg := validSecretConfig()
	cfg.EthRPCURL = ""

	err := cfg.Check()
	require.Error(t, err)
}
