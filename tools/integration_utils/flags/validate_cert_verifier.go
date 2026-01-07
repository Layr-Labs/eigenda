package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	ValidateCertVerifierEnvPrefix = "VALIDATE_CERT_VERIFIER"
)

var (
	ValidateCertVerifierJsonRPCURLFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "json-rpc-url"),
		Usage:    "JSON RPC URL for Ethereum client",
		EnvVar:   common.PrefixEnvVar(ValidateCertVerifierEnvPrefix, "JSON_RPC_URL"),
		Required: true,
	}
	ValidateCertVerifierSignerAuthKeyFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signer-auth-key"),
		Usage:    "Private key for signing dispersal requests (hex format, without 0x prefix)",
		EnvVar:   common.PrefixEnvVar(ValidateCertVerifierEnvPrefix, "SIGNER_AUTH_KEY"),
		Required: true,
	}
	ValidateCertVerifierCertVerifierAddrFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "cert-verifier-address"),
		Usage:    "Address of the EigenDACertVerifier contract (optional, defaults to network value)",
		EnvVar:   common.PrefixEnvVar(ValidateCertVerifierEnvPrefix, "CERT_VERIFIER_ADDRESS"),
		Required: false,
	}
	ValidateCertVerifierSrsPathFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "srs-path"),
		Usage:    "Path to SRS files directory",
		EnvVar:   common.PrefixEnvVar(ValidateCertVerifierEnvPrefix, "SRS_PATH"),
		Value:    "resources/srs",
		Required: false,
	}
)

var ValidateCertVerifierFlags = []cli.Flag{
	NetworkFlag,
	ValidateCertVerifierJsonRPCURLFlag,
	ValidateCertVerifierSignerAuthKeyFlag,
	ValidateCertVerifierCertVerifierAddrFlag,
	ValidateCertVerifierSrsPathFlag,
}
