package eigendaflags

import (
	"fmt"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/config/eigendaflags"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	clients_v2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/urfave/cli/v2"
)

var (
	DisperserFlagName               = withFlagPrefix("disperser-rpc")
	DisableTLSFlagName              = withFlagPrefix("disable-tls")
	BlobStatusPollIntervalFlagName  = withFlagPrefix("blob-status-poll-interval")
	PointEvaluationDisabledFlagName = withFlagPrefix("disable-point-evaluation")

	PutRetriesFlagName           = withFlagPrefix("put-retries")
	SignerPaymentKeyHexFlagName  = withFlagPrefix("signer-payment-key-hex")
	DisperseBlobTimeoutFlagName  = withFlagPrefix("disperse-blob-timeout")
	BlobCertifiedTimeoutFlagName = withFlagPrefix("blob-certified-timeout")
	CertVerifierAddrFlagName     = withFlagPrefix("cert-verifier-addr")
	RelayTimeoutFlagName         = withFlagPrefix("relay-timeout")
	ContractCallTimeoutFlagName  = withFlagPrefix("contract-call-timeout")
	BlobParamsVersionFlagName    = withFlagPrefix("blob-version")
	EthRPCURLFlagName            = withFlagPrefix("eth-rpc")
	MaxBlobLengthFlagName        = withFlagPrefix("max-blob-length")
)

func withFlagPrefix(s string) string {
	return "eigenda.v2." + s
}

func withEnvPrefix(envPrefix, s string) string {
	return envPrefix + "_EIGENDA_V2_" + s
}
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     DisperserFlagName,
			Usage:    "RPC endpoint of the EigenDA disperser.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "DISPERSER_RPC")},
			Category: category,
		},
		&cli.BoolFlag{
			Name:     DisableTLSFlagName,
			Usage:    "Disable TLS for gRPC communication with the EigenDA disperser and retrieval subnet.",
			Value:    false,
			EnvVars:  []string{withEnvPrefix(envPrefix, "GRPC_DISABLE_TLS")},
			Category: category,
		},
		&cli.StringFlag{
			Name:     SignerPaymentKeyHexFlagName,
			Usage:    "Hex-encoded signer private key. Used for authorizing payments with EigenDA disperser. Should not be associated with an Ethereum address holding any funds.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "SIGNER_PRIVATE_KEY_HEX")},
			Category: category,
		},
		&cli.BoolFlag{
			Name: PointEvaluationDisabledFlagName,
			Usage: "Disables IFFT transformation done during payload encoding. " +
				"Using this mode results in blobs that can't be proven.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "DISABLE_POINT_EVALUATION")},
			Value:    false,
			Category: category,
		},
		&cli.StringFlag{
			Name:     EthRPCURLFlagName,
			Usage:    "URL of the Ethereum RPC endpoint.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "ETH_RPC")},
			Category: category,
			Required: false,
		},
		&cli.UintFlag{
			Name:     PutRetriesFlagName,
			Usage:    "Number of times to retry blob dispersals before serving an error response.",
			Value:    3,
			EnvVars:  []string{withEnvPrefix(envPrefix, "PUT_RETRIES")},
			Category: category,
		},
		&cli.DurationFlag{
			Name:     DisperseBlobTimeoutFlagName,
			Usage:    "Maximum amount of time to wait for a blob to disperse against v2 protocol.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "DISPERSE_BLOB_TIMEOUT")},
			Category: category,
			Required: false,
			Value:    time.Minute * 2,
		},
		&cli.DurationFlag{
			Name:     BlobCertifiedTimeoutFlagName,
			Usage:    "Maximum amount of time to wait for blob certification against the on-chain EigenDACertVerifier.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CERTIFY_BLOB_TIMEOUT")},
			Category: category,
			Required: false,
			Value:    time.Second * 30,
		},
		&cli.StringFlag{
			Name: CertVerifierAddrFlagName,
			Usage: "Address of the EigenDACertVerifier contract. " +
				"Required for performing eth_calls to verify EigenDA certificates.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CERT_VERIFIER_ADDR")},
			Category: category,
			Required: false,
		},
		&cli.DurationFlag{
			Name:     ContractCallTimeoutFlagName,
			Usage:    "Timeout used when performing smart contract call operation (i.e, eth_call).",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CONTRACT_CALL_TIMEOUT")},
			Category: category,
			Value:    10 * time.Second,
			Required: false,
		},
		&cli.DurationFlag{
			Name:     RelayTimeoutFlagName,
			Usage:    "Timeout used when querying an individual relay for blob contents.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "RELAY_TIMEOUT")},
			Category: category,
			Value:    10 * time.Second,
			Required: false,
		},
		&cli.DurationFlag{
			Name:     BlobStatusPollIntervalFlagName,
			Usage:    "Duration to query for blob status updates during dispersal.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "BLOB_STATUS_POLL_INTERVAL")},
			Category: category,
			Value:    1 * time.Second,
			Required: false,
		},
		&cli.UintFlag{
			Name: BlobParamsVersionFlagName,
			Usage: `Blob params version used when dispersing. 
					   This refers to a global version maintained by EigenDA governance 
					   and is injected in the BlobHeader before dispersing.
					   Currently only supports (0).`,
			EnvVars:  []string{withEnvPrefix(envPrefix, "BLOB_PARAMS_VERSION")},
			Category: category,
			Value:    uint(0),
			Required: false,
		},
		&cli.StringFlag{
			Name: MaxBlobLengthFlagName,
			Usage: `Maximum blob length to be written or read from EigenDA. Determines the number of SRS points
						loaded into memory for KZG commitments. Example units: '30MiB', '4Kb', '30MB'. Maximum size
						slightly exceeds 1GB.`,
			EnvVars:  []string{withEnvPrefix(envPrefix, "MAX_BLOB_LENGTH")},
			Value:    "16MiB",
			Category: category,
		},
	}
}

func ReadClientConfigV2(ctx *cli.Context) (common.ClientConfigV2, error) {
	disperserConfig, err := readDisperserCfg(ctx)
	if err != nil {
		return common.ClientConfigV2{}, fmt.Errorf("read disperser config: %w", err)
	}

	maxBlobLengthFlagContents := ctx.String(MaxBlobLengthFlagName)
	maxBlobLengthBytes, err := eigendaflags.ParseMaxBlobLength(maxBlobLengthFlagContents)
	if err != nil {
		return common.ClientConfigV2{}, fmt.Errorf(
			"parse max blob length flag \"%v\": %w", maxBlobLengthFlagContents, err)
	}

	return common.ClientConfigV2{
		DisperserClientCfg:         disperserConfig,
		PayloadDisperserCfg:        readPayloadDisperserCfg(ctx),
		RelayPayloadRetrieverCfg:   readRetrievalConfig(ctx),
		PutRetries:                 ctx.Uint(PutRetriesFlagName),
		MaxBlobSizeBytes:           maxBlobLengthBytes,
		EigenDACertVerifierAddress: ctx.String(CertVerifierAddrFlagName),
	}, nil
}

func ReadSecretConfigV2(ctx *cli.Context) common.SecretConfigV2 {
	return common.SecretConfigV2{
		SignerPaymentKey: ctx.String(SignerPaymentKeyHexFlagName),
		EthRPCURL:        ctx.String(EthRPCURLFlagName),
	}
}

func readPayloadClientConfig(ctx *cli.Context) clients_v2.PayloadClientConfig {
	polyForm := codecs.PolynomialFormEval

	// if point evaluation mode is disabled then blob is treated as coefficients and
	// not iFFT'd before dispersal and FFT'd on retrieval
	if ctx.Bool(PointEvaluationDisabledFlagName) {
		polyForm = codecs.PolynomialFormCoeff
	}

	return clients_v2.PayloadClientConfig{
		PayloadPolynomialForm: polyForm,
		// #nosec G115 - only overflow on incorrect user input
		BlobVersion: uint16(ctx.Int(BlobParamsVersionFlagName)),
	}
}

func readPayloadDisperserCfg(ctx *cli.Context) payloaddispersal.PayloadDisperserConfig {
	payCfg := readPayloadClientConfig(ctx)

	return payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig:    payCfg,
		DisperseBlobTimeout:    ctx.Duration(DisperseBlobTimeoutFlagName),
		BlobCertifiedTimeout:   ctx.Duration(BlobCertifiedTimeoutFlagName),
		BlobStatusPollInterval: ctx.Duration(BlobStatusPollIntervalFlagName),
		ContractCallTimeout:    ctx.Duration(ContractCallTimeoutFlagName),
	}
}

func readDisperserCfg(ctx *cli.Context) (clients_v2.DisperserClientConfig, error) {
	disperserAddressString := ctx.String(DisperserFlagName)

	if disperserAddressString == "" {
		return clients_v2.DisperserClientConfig{}, nil
	}

	hostStr, portStr, err := net.SplitHostPort(disperserAddressString)
	if err != nil {
		return clients_v2.DisperserClientConfig{}, fmt.Errorf("split host port '%s': %w", disperserAddressString, err)
	}

	return clients_v2.DisperserClientConfig{
		Hostname:          hostStr,
		Port:              portStr,
		UseSecureGrpcFlag: !ctx.Bool(DisableTLSFlagName),
	}, nil
}

func readRetrievalConfig(ctx *cli.Context) payloadretrieval.RelayPayloadRetrieverConfig {
	return payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: readPayloadClientConfig(ctx),
		RelayTimeout:        ctx.Duration(RelayTimeoutFlagName),
	}
}
