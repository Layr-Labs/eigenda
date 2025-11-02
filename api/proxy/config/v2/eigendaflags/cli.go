package eigendaflags

import (
	"fmt"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	clients_v2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config/eigendaflags"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/urfave/cli/v2"
)

var (
	DisperserFlagName               = withFlagPrefix("disperser-rpc")
	DisableTLSFlagName              = withFlagPrefix("disable-tls")
	BlobStatusPollIntervalFlagName  = withFlagPrefix("blob-status-poll-interval")
	PointEvaluationDisabledFlagName = withFlagPrefix("disable-point-evaluation")

	PutRetriesFlagName                                = withFlagPrefix("put-retries")
	SignerPaymentKeyHexFlagName                       = withFlagPrefix("signer-payment-key-hex")
	DisperseBlobTimeoutFlagName                       = withFlagPrefix("disperse-blob-timeout")
	BlobCertifiedTimeoutFlagName                      = withFlagPrefix("blob-certified-timeout")
	CertVerifierRouterOrImmutableVerifierAddrFlagName = withFlagPrefix(
		"cert-verifier-router-or-immutable-verifier-addr",
	)
	EigenDADirectoryFlagName        = withFlagPrefix("eigenda-directory")
	RelayTimeoutFlagName            = withFlagPrefix("relay-timeout")
	ValidatorTimeoutFlagName        = withFlagPrefix("validator-timeout")
	ContractCallTimeoutFlagName     = withFlagPrefix("contract-call-timeout")
	BlobParamsVersionFlagName       = withFlagPrefix("blob-version")
	EthRPCURLFlagName               = withFlagPrefix("eth-rpc")
	MaxBlobLengthFlagName           = withFlagPrefix("max-blob-length")
	NetworkFlagName                 = withFlagPrefix("network")
	RBNRecencyWindowSizeFlagName    = withFlagPrefix("rbn-recency-window-size")
	RelayConnectionPoolSizeFlagName = withFlagPrefix("relay-connection-pool-size")

	ClientLedgerModeFlagName            = withFlagPrefix("client-ledger-mode")
	PaymentVaultMonitorIntervalFlagName = withFlagPrefix("payment-vault-monitor-interval")
)

func withFlagPrefix(s string) string {
	return "eigenda.v2." + s
}

func withEnvPrefix(envPrefix, s string) string {
	return envPrefix + "_EIGENDA_V2_" + s
}

// nolint: funlen
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
			Name: SignerPaymentKeyHexFlagName,
			Usage: "Optional hex-encoded signer private key. Used for authorizing payments with EigenDA disperser in PUT routes. " +
				"If not provided, proxy will be started in read-only mode, and will not be able to submit blobs to EigenDA. " +
				"Should not be associated with an Ethereum address holding any funds.",
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
		&cli.IntFlag{
			Name: PutRetriesFlagName,
			Usage: "Total number of times to try blob dispersals before serving an error response." +
				">0 = try dispersal that many times. <0 = retry indefinitely. 0 is not permitted (causes startup error).",
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
			Name: CertVerifierRouterOrImmutableVerifierAddrFlagName,
			Usage: "Address of either the EigenDACertVerifierRouter or immutable EigenDACertVerifier (V3 or above) contract. " +
				"Required for performing eth_calls to verify EigenDA certificates, as well as fetching " +
				"required_quorums and signature_thresholds needed when creating new EigenDA certificates during dispersals (POST routes).",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CERT_VERIFIER_ROUTER_OR_IMMUTABLE_VERIFIER_ADDR")},
			Category: category,
			Required: false,
		},
		&cli.StringFlag{
			Name:     EigenDADirectoryFlagName,
			Usage:    "Address of the EigenDA directory contract, which points to all other EigenDA contract addresses. This is the only contract entrypoint needed offchain..",
			EnvVars:  []string{withEnvPrefix(envPrefix, "EIGENDA_DIRECTORY")},
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
			Name: ValidatorTimeoutFlagName,
			Usage: "Timeout used when retrieving chunks directly from EigenDA validators. " +
				"This is a secondary retrieval method, in case retrieval from the relay network fails.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "VALIDATOR_TIMEOUT")},
			Category: category,
			Value:    2 * time.Minute,
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
			Usage: `Blob params version used when dispersing. This refers to a global version maintained by EigenDA
governance and is injected in the BlobHeader before dispersing. Currently only supports (0).`,
			EnvVars:  []string{withEnvPrefix(envPrefix, "BLOB_PARAMS_VERSION")},
			Category: category,
			Value:    uint(0),
			Required: false,
		},
		&cli.StringFlag{
			Name: MaxBlobLengthFlagName,
			Usage: `Maximum blob length (base 2) to be written or read from EigenDA. Determines the number of SRS points
loaded into memory for KZG commitments. Example units: '15MiB', '4Kib'.`,
			EnvVars:  []string{withEnvPrefix(envPrefix, "MAX_BLOB_LENGTH")},
			Value:    "16MiB",
			Category: category,
		},
		&cli.StringFlag{
			Name: NetworkFlagName,
			Usage: fmt.Sprintf(`The EigenDA network that is being used. This is an optional flag, 
to configure default values for different EigenDA contracts and disperser URL. 
See https://github.com/Layr-Labs/eigenda/blob/master/api/proxy/common/eigenda_network.go
for the exact values getting set by this flag. All of those values can also be manually
set via their respective flags, and take precedence over the default values set by the network flag.
If all of those other flags are manually configured, the network flag may be omitted. 
Permitted EigenDANetwork values include %s, %s, & %s.`,
				common.MainnetEigenDANetwork,
				common.HoodiTestnetEigenDANetwork,
				common.SepoliaTestnetEigenDANetwork,
			),
			EnvVars:  []string{withEnvPrefix(envPrefix, "NETWORK")},
			Category: category,
		},
		&cli.Uint64Flag{
			Name: RBNRecencyWindowSizeFlagName,
			Usage: `Allowed distance (in L1 blocks) between the eigenDA cert's reference 
block number (RBN) and the L1 block number at which the cert was included 
in the rollup's batch inbox. A cert is valid when cert.RBN < certL1InclusionBlock <= cert.RBN + rbnRecencyWindowSize, 
and otherwise is considered stale and verification will fail, and a 418 HTTP error will be returned.
This check is optional and will be skipped when set to 0.`,
			Value:    0,
			EnvVars:  []string{withEnvPrefix(envPrefix, "RBN_RECENCY_WINDOW_SIZE")},
			Category: category,
		},
		&cli.Uint64Flag{
			Name:     RelayConnectionPoolSizeFlagName,
			Usage:    "Number of gRPC connections to maintain to each relay.",
			Value:    1,
			EnvVars:  []string{withEnvPrefix(envPrefix, "RELAY_CONNECTION_POOL_SIZE")},
			Category: category,
			Required: false,
		},
		&cli.StringFlag{
			Name: ClientLedgerModeFlagName,
			Usage: "Payment mode for the client. Options: 'legacy', 'reservation-only', 'on-demand-only', " +
				"'reservation-and-on-demand'. The current default is 'legacy', which means that payments will be tracked " +
				"via the bin-based model, which is in the process of being deprecated. Eventually, the 'legacy' option " +
				"will be removed, once the migration to the new leaky bucket payment model is complete.",
			Value:    "legacy",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CLIENT_LEDGER_MODE")},
			Category: category,
			Required: false,
		},
		&cli.DurationFlag{
			Name: PaymentVaultMonitorIntervalFlagName,
			Usage: "Interval at which clients poll to check for changes to the PaymentVault contract (relevant " +
				"updates include changes to reservation parameters, and new on-demand payment deposits)",
			Value:    30 * time.Second,
			EnvVars:  []string{withEnvPrefix(envPrefix, "PAYMENT_VAULT_MONITOR_INTERVAL")},
			Category: category,
			Required: false,
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

	var eigenDANetwork common.EigenDANetwork
	networkString := ctx.String(NetworkFlagName)
	if networkString != "" {
		eigenDANetwork, err = common.EigenDANetworkFromString(networkString)
		if err != nil {
			return common.ClientConfigV2{}, fmt.Errorf("parse eigenDANetwork: %w", err)
		}
	}

	eigenDADirectory := ctx.String(EigenDADirectoryFlagName)
	if eigenDADirectory == "" {
		if networkString == "" {
			return common.ClientConfigV2{},
				fmt.Errorf("either EigenDA Directory contract address or EigenDANetwork enum must be specified")
		}
		eigenDANetwork, err := common.EigenDANetworkFromString(networkString)
		if err != nil {
			return common.ClientConfigV2{}, fmt.Errorf("parse eigenDANetwork: %w", err)
		}
		eigenDADirectory = eigenDANetwork.GetEigenDADirectory()
	}

	return common.ClientConfigV2{
		DisperserClientCfg:           disperserConfig,
		PayloadDisperserCfg:          readPayloadDisperserCfg(ctx),
		RelayPayloadRetrieverCfg:     readRelayRetrievalConfig(ctx),
		ValidatorPayloadRetrieverCfg: readValidatorRetrievalConfig(ctx),
		PutTries:                     ctx.Int(PutRetriesFlagName),
		MaxBlobSizeBytes:             maxBlobLengthBytes,
		// we don't expose this configuration to users, as all production use cases should have
		// both retrieval methods enabled. This could be exposed in the future, if necessary.
		// Note the order of these retrievers, which is significant: the relay retriever will be
		// tried first, and the validator retriever will only be tried if the relay retriever fails
		RetrieversToEnable: []common.RetrieverType{
			common.RelayRetrieverType,
			common.ValidatorRetrieverType,
		},
		EigenDACertVerifierOrRouterAddress: ctx.String(CertVerifierRouterOrImmutableVerifierAddrFlagName),
		EigenDADirectory:                   eigenDADirectory,
		RBNRecencyWindowSize:               ctx.Uint64(RBNRecencyWindowSizeFlagName),
		EigenDANetwork:                     eigenDANetwork,
		RelayConnectionPoolSize:            ctx.Uint(RelayConnectionPoolSizeFlagName),
		ClientLedgerMode:                   clientledger.ParseClientLedgerMode(ctx.String(ClientLedgerModeFlagName)),
		VaultMonitorInterval:               ctx.Duration(PaymentVaultMonitorIntervalFlagName),
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
		BlobCompleteTimeout:    ctx.Duration(BlobCertifiedTimeoutFlagName),
		BlobStatusPollInterval: ctx.Duration(BlobStatusPollIntervalFlagName),
		ContractCallTimeout:    ctx.Duration(ContractCallTimeoutFlagName),
	}
}

func readDisperserCfg(ctx *cli.Context) (clients_v2.DisperserClientConfig, error) {
	disperserAddressString := ctx.String(DisperserFlagName)
	if disperserAddressString == "" {
		networkString := ctx.String(NetworkFlagName)
		if networkString == "" {
			return clients_v2.DisperserClientConfig{},
				fmt.Errorf("either disperser address or EigenDANetwork must be specified")
		}

		eigenDANetwork, err := common.EigenDANetworkFromString(networkString)
		if err != nil {
			return clients_v2.DisperserClientConfig{}, fmt.Errorf("parse eigenDANetwork: %w", err)
		}

		disperserAddressString = eigenDANetwork.GetDisperserAddress()
	}

	hostStr, portStr, err := net.SplitHostPort(disperserAddressString)
	if err != nil {
		return clients_v2.DisperserClientConfig{},
			fmt.Errorf("split host port '%s': %w", disperserAddressString, err)
	}

	return clients_v2.DisperserClientConfig{
		Hostname:          hostStr,
		Port:              portStr,
		UseSecureGrpcFlag: !ctx.Bool(DisableTLSFlagName),
	}, nil
}

func readRelayRetrievalConfig(ctx *cli.Context) payloadretrieval.RelayPayloadRetrieverConfig {
	return payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: readPayloadClientConfig(ctx),
		RelayTimeout:        ctx.Duration(RelayTimeoutFlagName),
	}
}

func readValidatorRetrievalConfig(ctx *cli.Context) payloadretrieval.ValidatorPayloadRetrieverConfig {
	return payloadretrieval.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: readPayloadClientConfig(ctx),
		RetrievalTimeout:    ctx.Duration(ValidatorTimeoutFlagName),
	}
}
