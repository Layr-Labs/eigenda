package eigendaflags

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/urfave/cli/v2"
)

// TODO: we should eventually move all of these flags into the eigenda repo

var (
	DisperserRPCFlagName                 = withFlagPrefix("disperser-rpc")
	StatusQueryRetryIntervalFlagName     = withFlagPrefix("status-query-retry-interval")
	StatusQueryTimeoutFlagName           = withFlagPrefix("status-query-timeout")
	DisableTLSFlagName                   = withFlagPrefix("disable-tls")
	ResponseTimeoutFlagName              = withFlagPrefix("response-timeout")
	CustomQuorumIDsFlagName              = withFlagPrefix("custom-quorum-ids")
	SignerPrivateKeyHexFlagName          = withFlagPrefix("signer-private-key-hex")
	PutBlobEncodingVersionFlagName       = withFlagPrefix("put-blob-encoding-version")
	DisablePointVerificationModeFlagName = withFlagPrefix("disable-point-verification-mode")
	WaitForFinalizationFlagName          = withFlagPrefix("wait-for-finalization")
)

func withFlagPrefix(s string) string {
	return "eigenda." + s
}

func withEnvPrefix(envPrefix, s string) string {
	return envPrefix + "_EIGENDA_" + s
}

// CLIFlags ... used for EigenDA client configuration
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     DisperserRPCFlagName,
			Usage:    "RPC endpoint of the EigenDA disperser.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "DISPERSER_RPC")},
			Category: category,
		},
		&cli.DurationFlag{
			Name:     StatusQueryTimeoutFlagName,
			Usage:    "Duration to wait for a blob to finalize after being sent for dispersal. Default is 30 minutes.",
			Value:    30 * time.Minute,
			EnvVars:  []string{withEnvPrefix(envPrefix, "STATUS_QUERY_TIMEOUT")},
			Category: category,
		},
		&cli.DurationFlag{
			Name:     StatusQueryRetryIntervalFlagName,
			Usage:    "Interval between retries when awaiting network blob finalization. Default is 5 seconds.",
			Value:    5 * time.Second,
			EnvVars:  []string{withEnvPrefix(envPrefix, "STATUS_QUERY_INTERVAL")},
			Category: category,
		},
		&cli.BoolFlag{
			Name:     DisableTLSFlagName,
			Usage:    "Disable TLS for gRPC communication with the EigenDA disperser. Default is false.",
			Value:    false,
			EnvVars:  []string{withEnvPrefix(envPrefix, "GRPC_DISABLE_TLS")},
			Category: category,
		},
		&cli.DurationFlag{
			Name:     ResponseTimeoutFlagName,
			Usage:    "Total time to wait for a response from the EigenDA disperser. Default is 60 seconds.",
			Value:    60 * time.Second,
			EnvVars:  []string{withEnvPrefix(envPrefix, "RESPONSE_TIMEOUT")},
			Category: category,
		},
		&cli.UintSliceFlag{
			Name:     CustomQuorumIDsFlagName,
			Usage:    "Custom quorum IDs for writing blobs. Should not include default quorums 0 or 1.",
			Value:    cli.NewUintSlice(),
			EnvVars:  []string{withEnvPrefix(envPrefix, "CUSTOM_QUORUM_IDS")},
			Category: category,
		},
		&cli.StringFlag{
			Name:     SignerPrivateKeyHexFlagName,
			Usage:    "Hex-encoded signer private key. This key should not be associated with an Ethereum address holding any funds.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "SIGNER_PRIVATE_KEY_HEX")},
			Category: category,
		},
		&cli.UintFlag{
			Name:     PutBlobEncodingVersionFlagName,
			Usage:    "Blob encoding version to use when writing blobs from the high-level interface.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "PUT_BLOB_ENCODING_VERSION")},
			Value:    0,
			Category: category,
		},
		&cli.BoolFlag{
			Name:     DisablePointVerificationModeFlagName,
			Usage:    "Disable point verification mode. This mode performs IFFT on data before writing and FFT on data after reading. Disabling requires supplying the entire blob for verification against the KZG commitment.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "DISABLE_POINT_VERIFICATION_MODE")},
			Value:    false,
			Category: category,
		},
		&cli.BoolFlag{
			Name:     WaitForFinalizationFlagName,
			Usage:    "Wait for blob finalization before returning from PutBlob.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "WAIT_FOR_FINALIZATION")},
			Value:    false,
			Category: category,
		},
	}
}

func ReadConfig(ctx *cli.Context) clients.EigenDAClientConfig {
	return clients.EigenDAClientConfig{
		RPC:                          ctx.String(DisperserRPCFlagName),
		StatusQueryRetryInterval:     ctx.Duration(StatusQueryRetryIntervalFlagName),
		StatusQueryTimeout:           ctx.Duration(StatusQueryTimeoutFlagName),
		DisableTLS:                   ctx.Bool(DisableTLSFlagName),
		ResponseTimeout:              ctx.Duration(ResponseTimeoutFlagName),
		CustomQuorumIDs:              ctx.UintSlice(CustomQuorumIDsFlagName),
		SignerPrivateKeyHex:          ctx.String(SignerPrivateKeyHexFlagName),
		PutBlobEncodingVersion:       codecs.BlobEncodingVersion(ctx.Uint(PutBlobEncodingVersionFlagName)),
		DisablePointVerificationMode: ctx.Bool(DisablePointVerificationModeFlagName),
		WaitForFinalization:          ctx.Bool(WaitForFinalizationFlagName),
	}
}
