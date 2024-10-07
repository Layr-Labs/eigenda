package eigendaflags

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

// All of these flags are deprecated and will be removed in release v2.0.0
// we leave them here with actions that crash the program to ensure they are not used,
// and to make it easier for users to find the new flags (instead of silently crashing late during execution
// because some flag's env var was changed but the user forgot to update it)
var (
	DeprecatedDisperserRPCFlagName                 = withDeprecatedFlagPrefix("disperser-rpc")
	DeprecatedStatusQueryRetryIntervalFlagName     = withDeprecatedFlagPrefix("status-query-retry-interval")
	DeprecatedStatusQueryTimeoutFlagName           = withDeprecatedFlagPrefix("status-query-timeout")
	DeprecatedDisableTLSFlagName                   = withDeprecatedFlagPrefix("disable-tls")
	DeprecatedResponseTimeoutFlagName              = withDeprecatedFlagPrefix("response-timeout")
	DeprecatedCustomQuorumIDsFlagName              = withDeprecatedFlagPrefix("custom-quorum-ids")
	DeprecatedSignerPrivateKeyHexFlagName          = withDeprecatedFlagPrefix("signer-private-key-hex")
	DeprecatedPutBlobEncodingVersionFlagName       = withDeprecatedFlagPrefix("put-blob-encoding-version")
	DeprecatedDisablePointVerificationModeFlagName = withDeprecatedFlagPrefix("disable-point-verification-mode")
	DeprecatedWaitForFinalizationFlagName          = withDeprecatedFlagPrefix("wait-for-finalization")
)

func withDeprecatedFlagPrefix(s string) string {
	return "eigenda-" + s
}

func withDeprecatedEnvPrefix(envPrefix, s string) string {
	return envPrefix + "_" + s
}

// CLIFlags ... used for EigenDA client configuration
func DeprecatedCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     DeprecatedDisperserRPCFlagName,
			Usage:    "RPC endpoint of the EigenDA disperser.",
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "DISPERSER_RPC")},
			Category: category,
			Action: func(*cli.Context, string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedDisperserRPCFlagName, withDeprecatedEnvPrefix(envPrefix, "DISPERSER_RPC"),
					DisperserRPCFlagName, withEnvPrefix(envPrefix, "DISPERSER_RPC"))
			},
		},
		&cli.DurationFlag{
			Name:     DeprecatedStatusQueryTimeoutFlagName,
			Usage:    "Duration to wait for a blob to finalize after being sent for dispersal. Default is 30 minutes.",
			Value:    30 * time.Minute,
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "STATUS_QUERY_TIMEOUT")},
			Category: category,
			Action: func(*cli.Context, time.Duration) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedStatusQueryTimeoutFlagName, withDeprecatedEnvPrefix(envPrefix, "STATUS_QUERY_TIMEOUT"),
					StatusQueryTimeoutFlagName, withEnvPrefix(envPrefix, "STATUS_QUERY_TIMEOUT"))
			},
		},
		&cli.DurationFlag{
			Name:     DeprecatedStatusQueryRetryIntervalFlagName,
			Usage:    "Interval between retries when awaiting network blob finalization. Default is 5 seconds.",
			Value:    5 * time.Second,
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "STATUS_QUERY_INTERVAL")},
			Category: category,
			Action: func(*cli.Context, time.Duration) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedStatusQueryRetryIntervalFlagName, withDeprecatedEnvPrefix(envPrefix, "STATUS_QUERY_INTERVAL"),
					StatusQueryRetryIntervalFlagName, withEnvPrefix(envPrefix, "STATUS_QUERY_INTERVAL"))
			},
		},
		&cli.BoolFlag{
			Name:     DeprecatedDisableTLSFlagName,
			Usage:    "Disable TLS for gRPC communication with the EigenDA disperser. Default is false.",
			Value:    false,
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "GRPC_DISABLE_TLS")},
			Category: category,
			Action: func(*cli.Context, bool) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedDisableTLSFlagName, withDeprecatedEnvPrefix(envPrefix, "GRPC_DISABLE_TLS"),
					DisableTLSFlagName, withEnvPrefix(envPrefix, "GRPC_DISABLE_TLS"))
			},
		},
		&cli.DurationFlag{
			Name:     DeprecatedResponseTimeoutFlagName,
			Usage:    "Total time to wait for a response from the EigenDA disperser. Default is 60 seconds.",
			Value:    60 * time.Second,
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "RESPONSE_TIMEOUT")},
			Category: category,
			Action: func(*cli.Context, time.Duration) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedResponseTimeoutFlagName, withDeprecatedEnvPrefix(envPrefix, "RESPONSE_TIMEOUT"),
					ResponseTimeoutFlagName, withEnvPrefix(envPrefix, "RESPONSE_TIMEOUT"))
			},
		},
		&cli.UintSliceFlag{
			Name:     DeprecatedCustomQuorumIDsFlagName,
			Usage:    "Custom quorum IDs for writing blobs. Should not include default quorums 0 or 1.",
			Value:    cli.NewUintSlice(),
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "CUSTOM_QUORUM_IDS")},
			Category: category,
			Action: func(*cli.Context, []uint) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedCustomQuorumIDsFlagName, withDeprecatedEnvPrefix(envPrefix, "CUSTOM_QUORUM_IDS"),
					CustomQuorumIDsFlagName, withEnvPrefix(envPrefix, "CUSTOM_QUORUM_IDS"))
			},
		},
		&cli.StringFlag{
			Name:     DeprecatedSignerPrivateKeyHexFlagName,
			Usage:    "Hex-encoded signer private key. This key should not be associated with an Ethereum address holding any funds.",
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "SIGNER_PRIVATE_KEY_HEX")},
			Category: category,
			Action: func(*cli.Context, string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedSignerPrivateKeyHexFlagName, withDeprecatedEnvPrefix(envPrefix, "SIGNER_PRIVATE_KEY_HEX"),
					SignerPrivateKeyHexFlagName, withEnvPrefix(envPrefix, "SIGNER_PRIVATE_KEY_HEX"))
			},
		},
		&cli.UintFlag{
			Name:     DeprecatedPutBlobEncodingVersionFlagName,
			Usage:    "Blob encoding version to use when writing blobs from the high-level interface.",
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "PUT_BLOB_ENCODING_VERSION")},
			Value:    0,
			Category: category,
			Action: func(*cli.Context, uint) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedPutBlobEncodingVersionFlagName, withDeprecatedEnvPrefix(envPrefix, "PUT_BLOB_ENCODING_VERSION"),
					PutBlobEncodingVersionFlagName, withEnvPrefix(envPrefix, "PUT_BLOB_ENCODING_VERSION"))
			},
		},
		&cli.BoolFlag{
			Name:     DeprecatedDisablePointVerificationModeFlagName,
			Usage:    "Disable point verification mode. This mode performs IFFT on data before writing and FFT on data after reading. Disabling requires supplying the entire blob for verification against the KZG commitment.",
			EnvVars:  []string{withDeprecatedEnvPrefix(envPrefix, "DISABLE_POINT_VERIFICATION_MODE")},
			Value:    false,
			Category: category,
			Action: func(*cli.Context, bool) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedDisablePointVerificationModeFlagName, withDeprecatedEnvPrefix(envPrefix, "DISABLE_POINT_VERIFICATION_MODE"),
					DisablePointVerificationModeFlagName, withEnvPrefix(envPrefix, "DISABLE_POINT_VERIFICATION_MODE"))
			},
		},
	}
}
