package eigendaflags

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common/consts"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/urfave/cli/v2"
)

// TODO: we should eventually move all of these flags into the eigenda repo

var (
	DisperserRPCFlagName                 = withFlagPrefix("disperser-rpc")
	ResponseTimeoutFlagName              = withFlagPrefix("response-timeout")
	ConfirmationTimeoutFlagName          = withFlagPrefix("confirmation-timeout")
	StatusQueryRetryIntervalFlagName     = withFlagPrefix("status-query-retry-interval")
	StatusQueryTimeoutFlagName           = withFlagPrefix("status-query-timeout")
	DisableTLSFlagName                   = withFlagPrefix("disable-tls")
	CustomQuorumIDsFlagName              = withFlagPrefix("custom-quorum-ids")
	SignerPrivateKeyHexFlagName          = withFlagPrefix("signer-private-key-hex")
	PutBlobEncodingVersionFlagName       = withFlagPrefix("put-blob-encoding-version")
	DisablePointVerificationModeFlagName = withFlagPrefix("disable-point-verification-mode")
	WaitForFinalizationFlagName          = withFlagPrefix("wait-for-finalization")
	ConfirmationDepthFlagName            = withFlagPrefix("confirmation-depth")
	EthRPCURLFlagName                    = withFlagPrefix("eth-rpc")
	SvcManagerAddrFlagName               = withFlagPrefix("svc-manager-addr")
	// Flags that are proxy specific, and not used by the eigenda-client
	PutRetriesFlagName = withFlagPrefix("put-retries")
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
			Name:     ResponseTimeoutFlagName,
			Usage:    "Flag used to configure the underlying disperser-client. Total time to wait for the disperseBlob call to return or disperseAuthenticatedBlob stream to finish and close.",
			Value:    60 * time.Second,
			EnvVars:  []string{withEnvPrefix(envPrefix, "RESPONSE_TIMEOUT")},
			Category: category,
		},
		&cli.DurationFlag{
			Name: ConfirmationTimeoutFlagName,
			Usage: `The total amount of time that the client will spend waiting for EigenDA
			to "confirm" (include onchain) a blob after it has been dispersed. Note that
			we stick to "confirm" here but this really means InclusionTimeout,
			not confirmation in the sense of confirmation depth.
			
			If ConfirmationTimeout time passes and the blob is not yet confirmed,
			the client will return an api.ErrorFailover to let the caller failover to EthDA.`,
			Value:    15 * time.Minute,
			EnvVars:  []string{withEnvPrefix(envPrefix, "CONFIRMATION_TIMEOUT")},
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
		&cli.UintSliceFlag{
			Name:     CustomQuorumIDsFlagName,
			Usage:    "Custom quorum IDs for writing blobs. Should not include default quorums 0 or 1.",
			Value:    cli.NewUintSlice(),
			EnvVars:  []string{withEnvPrefix(envPrefix, "CUSTOM_QUORUM_IDS")},
			Category: category,
		},
		&cli.StringFlag{
			Name:     SignerPrivateKeyHexFlagName,
			Usage:    "Hex-encoded signer private key. Used for authn/authz and rate limits on EigenDA disperser. Should not be associated with an Ethereum address holding any funds.",
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
			// This flag is DEPRECATED. Use ConfirmationDepthFlagName, which accept "finalization" or a number <64.
			Name:     WaitForFinalizationFlagName,
			Usage:    "Wait for blob finalization before returning from PutBlob.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "WAIT_FOR_FINALIZATION")},
			Value:    false,
			Category: category,
			Hidden:   true,
			Action: func(_ *cli.Context, _ bool) error {
				return fmt.Errorf("flag --%s is deprecated, instead use --%s finalized", WaitForFinalizationFlagName, ConfirmationDepthFlagName)
			},
		},
		&cli.StringFlag{
			Name: ConfirmationDepthFlagName,
			Usage: "Number of Ethereum blocks to wait after the blob's batch has been included on-chain, " +
				"before returning from PutBlob calls. Can either be a number or 'finalized'.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CONFIRMATION_DEPTH")},
			Value:    "0",
			Category: category,
			Action: func(_ *cli.Context, val string) error {
				return validateConfirmationFlag(val)
			},
		},
		&cli.StringFlag{
			Name:     EthRPCURLFlagName,
			Usage:    "URL of the Ethereum RPC endpoint. Needed to confirm blobs landed onchain.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "ETH_RPC")},
			Category: category,
			Required: false,
		},
		&cli.StringFlag{
			Name:     SvcManagerAddrFlagName,
			Usage:    "Address of the EigenDAServiceManager contract. Required to confirm blobs landed onchain. See https://github.com/Layr-Labs/eigenlayer-middleware/?tab=readme-ov-file#current-mainnet-deployment",
			EnvVars:  []string{withEnvPrefix(envPrefix, "SERVICE_MANAGER_ADDR")},
			Category: category,
			Required: false,
		},
		// Flags that are proxy specific, and not used by the eigenda-client
		// TODO: should we move this to a more specific category, like EIGENDA_STORE?
		&cli.UintFlag{
			Name:     PutRetriesFlagName,
			Usage:    "Number of times to retry blob dispersals.",
			Value:    3,
			EnvVars:  []string{withEnvPrefix(envPrefix, "PUT_RETRIES")},
			Category: category,
		},
	}
}

func ReadConfig(ctx *cli.Context) clients.EigenDAClientConfig {
	waitForFinalization, confirmationDepth := parseConfirmationFlag(ctx.String(ConfirmationDepthFlagName))
	return clients.EigenDAClientConfig{
		RPC:                          ctx.String(DisperserRPCFlagName),
		ResponseTimeout:              ctx.Duration(ResponseTimeoutFlagName),
		ConfirmationTimeout:          ctx.Duration(ConfirmationTimeoutFlagName),
		StatusQueryRetryInterval:     ctx.Duration(StatusQueryRetryIntervalFlagName),
		StatusQueryTimeout:           ctx.Duration(StatusQueryTimeoutFlagName),
		DisableTLS:                   ctx.Bool(DisableTLSFlagName),
		CustomQuorumIDs:              ctx.UintSlice(CustomQuorumIDsFlagName),
		SignerPrivateKeyHex:          ctx.String(SignerPrivateKeyHexFlagName),
		PutBlobEncodingVersion:       codecs.BlobEncodingVersion(ctx.Uint(PutBlobEncodingVersionFlagName)),
		DisablePointVerificationMode: ctx.Bool(DisablePointVerificationModeFlagName),
		WaitForFinalization:          waitForFinalization,
		WaitForConfirmationDepth:     confirmationDepth,
		EthRpcUrl:                    ctx.String(EthRPCURLFlagName),
		SvcManagerAddr:               ctx.String(SvcManagerAddrFlagName),
	}
}

// parse the val (either "finalized" or a number) into waitForFinalization (bool) and confirmationDepth (uint64).
func parseConfirmationFlag(val string) (bool, uint64) {
	if val == "finalized" {
		return true, 0
	}

	depth, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic("this should never happen, as the flag is validated before this point")
	}

	return false, depth
}

func validateConfirmationFlag(val string) error {
	if val == "finalized" {
		return nil
	}

	depth, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return fmt.Errorf("confirmation-depth must be either 'finalized' or a number, got: %s", val)
	}

	if depth >= uint64(consts.EthHappyPathFinalizationDepthBlocks) {
		// We keep this low (<128) to avoid requiring an archive node (see how this is used in CertVerifier).
		// Note: assuming here that no sane person would ever need to set this to a number >64.
		// But perhaps someone testing crazy reorg scenarios where finalization takes >2 epochs might want to set this to a higher number.
		// Do keep in mind if you ever change this that it might affect a LOT of validators on your rollup who would now need an archival node.
		return fmt.Errorf("confirmation depth set to %d, which is > 2 epochs (64). Use 'finalized' instead", depth)
	}

	return nil
}
