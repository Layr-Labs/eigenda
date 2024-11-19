package flags

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli"
)

const (
	FlagPrefix   = "node"
	EnvVarPrefix = "NODE"
)

var (
	/* Required Flags */

	HostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "hostname"),
		Usage:    "Hostname at which node is available",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "HOSTNAME"),
	}
	DispersalPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dispersal-port"),
		Usage:    "Port at which node registers to listen for dispersal calls",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "DISPERSAL_PORT"),
	}
	RetrievalPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "retrieval-port"),
		Usage:    "Port at which node registers to listen for retrieval calls",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "RETRIEVAL_PORT"),
	}
	InternalDispersalPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "internal-dispersal-port"),
		Usage:    "Port at which node listens for dispersal calls (used when node is behind NGINX)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "INTERNAL_DISPERSAL_PORT"),
	}
	InternalRetrievalPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "internal-retrieval-port"),
		Usage:    "Port at which node listens for retrieval calls (used when node is behind NGINX)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "INTERNAL_RETRIEVAL_PORT"),
	}
	EnableNodeApiFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-node-api"),
		Usage:    "enable node-api to serve eigenlayer-cli node-api calls",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENABLE_NODE_API"),
	}
	NodeApiPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "node-api-port"),
		Usage:    "Port at which node listens for eigenlayer-cli node-api calls",
		Required: false,
		Value:    "9091",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "API_PORT"),
	}
	EnableMetricsFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "enable prometheus to serve metrics collection",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENABLE_METRICS"),
	}
	MetricsPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-port"),
		Usage:    "Port at which node listens for metrics calls",
		Required: false,
		Value:    "9091",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "METRICS_PORT"),
	}
	OnchainMetricsIntervalFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "onchain-metrics-interval"),
		Usage:    "The interval in seconds at which the node polls the onchain state of the operator and update metrics. <=0 means no poll",
		Required: false,
		Value:    "180",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ONCHAIN_METRICS_INTERVAL"),
	}
	TimeoutFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "Amount of time to wait for GPRC",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "TIMEOUT"),
	}
	QuorumIDListFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "quorum-id-list"),
		Usage:    "Comma separated list of quorum IDs that the node will participate in. There should be at least one quorum ID. This list must not contain quorums node is already registered with.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "QUORUM_ID_LIST"),
	}
	DbPathFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "db-path"),
		Usage:    "Path for level db",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "DB_PATH"),
	}
	// The files for encrypted private keys.
	BlsKeyFileFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-key-file"),
		Required: false,
		Usage:    "Path to the encrypted bls private key",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BLS_KEY_FILE"),
	}
	EcdsaKeyFileFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "ecdsa-key-file"),
		Required: false,
		Usage:    "Path to the encrypted ecdsa private key",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ECDSA_KEY_FILE"),
	}
	// Passwords to decrypt the private keys.
	BlsKeyPasswordFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-key-password"),
		Required: false,
		Usage:    "Password to decrypt bls private key",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BLS_KEY_PASSWORD"),
	}
	EcdsaKeyPasswordFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "ecdsa-key-password"),
		Required: false,
		Usage:    "Password to decrypt ecdsa private key",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ECDSA_KEY_PASSWORD"),
	}
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "Address of the BLS Operator State Retriever",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "Address of the EigenDA Service Manager",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
	ChurnerUrlFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "churner-url"),
		Usage:    "URL of the Churner",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "CHURNER_URL"),
	}
	ChurnerUseSecureGRPC = cli.BoolTFlag{
		Name:     common.PrefixFlag(FlagPrefix, "churner-use-secure-grpc"),
		Usage:    "Whether to use secure GRPC connection to Churner",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "CHURNER_USE_SECURE_GRPC"),
	}
	PubIPProviderFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "public-ip-provider"),
		Usage:    "The ip provider service used to obtain a node's public IP [seeip (default), ipify)",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "PUBLIC_IP_PROVIDER"),
	}
	PubIPCheckIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "public-ip-check-interval"),
		Usage:    "Interval at which to check for changes in the node's public IP (Ex: 10s). If set to 0, the check will be disabled.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "PUBLIC_IP_CHECK_INTERVAL"),
	}

	/* Optional Flags */

	// This flag is used to control if the DA Node registers itself when it starts.
	// This is useful for testing and for hosted node where we don't want to have
	// mannual operation with CLI to register.
	// By default, it will not register itself at start.
	RegisterAtNodeStartFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "register-at-node-start"),
		Usage:    "Whether to register the node for EigenDA when it starts",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "REGISTER_AT_NODE_START"),
	}
	ExpirationPollIntervalSecFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "expiration-poll-interval"),
		Usage:    "How often (in second) to poll status and expire outdated blobs",
		Required: false,
		Value:    "180",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "EXPIRATION_POLL_INTERVAL"),
	}
	ReachabilityPollIntervalSecFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "reachability-poll-interval"),
		Usage:    "How often (in second) to check if node is reachabile from Disperser",
		Required: false,
		Value:    "60",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "REACHABILITY_POLL_INTERVAL"),
	}
	// Optional DataAPI URL. If not set, reachability checks are disabled
	DataApiUrlFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dataapi-url"),
		Usage:    "URL of the DataAPI",
		Required: false,
		Value:    "",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "DATAAPI_URL"),
	}
	// NumBatchValidators is the maximum number of parallel workers used to
	// validate a batch (defaults to 128).
	NumBatchValidatorsFlag = cli.IntFlag{
		Name:     "num-batch-validators",
		Usage:    "maximum number of parallel workers used to validate a batch (defaults to 128)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "NUM_BATCH_VALIDATORS"),
		Value:    128,
	}
	NumBatchDeserializationWorkersFlag = cli.IntFlag{
		Name:     "num-batch-deserialization-workers",
		Usage:    "maximum number of parallel workers used to deserialize a batch (defaults to 128)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "NUM_BATCH_DESERIALIZATION_WORKERS"),
		Value:    128,
	}
	EnableGnarkBundleEncodingFlag = cli.BoolFlag{
		Name:     "enable-gnark-bundle-encoding",
		Usage:    "Enable Gnark bundle encoding for chunks",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENABLE_GNARK_BUNDLE_ENCODING"),
	}

	// Test only, DO NOT USE the following flags in production

	// This flag controls whether other test flags can take effect.
	// By default, it is not test mode.
	EnableTestModeFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-test-mode"),
		Usage:    "Whether to run as test mode. This flag needs to be enabled for other test flags to take effect",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENABLE_TEST_MODE"),
	}

	// Corresponding to the BLOCK_STALE_MEASURE defined onchain in
	// contracts/src/core/EigenDAServiceManagerStorage.sol
	// This flag is used to override the value from the chain. The target use case is testing.
	OverrideBlockStaleMeasureFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "override-block-stale-measure"),
		Usage:    "The maximum amount of blocks in the past that the service will consider stake amounts to still be valid. This is used to override the value set on chain. <=0 means no override",
		Required: false,
		Value:    "-1",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "OVERRIDE_BLOCK_STALE_MEASURE"),
	}
	// Corresponding to the STORE_DURATION_BLOCKS defined onchain in
	// contracts/src/core/EigenDAServiceManagerStorage.sol
	// This flag is used to override the value from the chain. The target use case is testing.
	OverrideStoreDurationBlocksFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "override-store-duration-blocks"),
		Usage:    "Unit of measure (in blocks) for which data will be stored for after confirmation. This is used to override the value set on chain. <=0 means no override",
		Required: false,
		Value:    "-1",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "OVERRIDE_STORE_DURATION_BLOCKS"),
	}
	// DO NOT set plain private key in flag in production.
	// When test mode is enabled, the DA Node will take private BLS key from this flag.
	TestPrivateBlsFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "test-private-bls"),
		Usage:    "Test BLS private key for node operator",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "TEST_PRIVATE_BLS"),
	}
	ClientIPHeaderFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "client-ip-header"),
		Usage:    "The name of the header used to get the client IP address. If set to empty string, the IP address will be taken from the connection. The rightmost value of the header will be used.",
		Required: false,
		Value:    "",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "CLIENT_IP_HEADER"),
	}

	DisableNodeInfoResourcesFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disable-node-info-resources"),
		Usage:    "Disable system resource information (OS, architecture, CPU, memory) on the NodeInfo API",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "DISABLE_NODE_INFO_RESOURCES"),
	}

	BLSRemoteSignerEnabledFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-remote-signer-enabled"),
		Usage:    "Set to true to enable the BLS remote signer",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BLS_REMOTE_SIGNER_ENABLED"),
	}

	BLSRemoteSignerUrlFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-remote-signer-url"),
		Usage:    "The URL of the BLS remote signer",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BLS_REMOTE_SIGNER_URL"),
	}

	BLSPublicKeyHexFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-public-key-hex"),
		Usage:    "The hex-encoded public key of the BLS signer",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BLS_PUBLIC_KEY_HEX"),
	}

	BLSSignerCertFileFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-signer-cert-file"),
		Usage:    "The path to the BLS signer certificate file",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BLS_SIGNER_CERT_FILE"),
	}
)

var requiredFlags = []cli.Flag{
	HostnameFlag,
	DispersalPortFlag,
	RetrievalPortFlag,
	EnableMetricsFlag,
	MetricsPortFlag,
	OnchainMetricsIntervalFlag,
	EnableNodeApiFlag,
	NodeApiPortFlag,
	TimeoutFlag,
	QuorumIDListFlag,
	DbPathFlag,
	BlsKeyFileFlag,
	BlsKeyPasswordFlag,
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
	PubIPProviderFlag,
	PubIPCheckIntervalFlag,
	ChurnerUrlFlag,
}

var optionalFlags = []cli.Flag{
	RegisterAtNodeStartFlag,
	ExpirationPollIntervalSecFlag,
	ReachabilityPollIntervalSecFlag,
	EnableTestModeFlag,
	OverrideBlockStaleMeasureFlag,
	OverrideStoreDurationBlocksFlag,
	TestPrivateBlsFlag,
	NumBatchValidatorsFlag,
	NumBatchDeserializationWorkersFlag,
	InternalDispersalPortFlag,
	InternalRetrievalPortFlag,
	ClientIPHeaderFlag,
	ChurnerUseSecureGRPC,
	EcdsaKeyFileFlag,
	EcdsaKeyPasswordFlag,
	DataApiUrlFlag,
	DisableNodeInfoResourcesFlag,
	EnableGnarkBundleEncodingFlag,
	BLSRemoteSignerEnabledFlag,
	BLSRemoteSignerUrlFlag,
	BLSPublicKeyHexFlag,
	BLSSignerCertFileFlag,
}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, kzg.CLIFlags(EnvVarPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(EnvVarPrefix)...)
	Flags = append(Flags, common.LoggerCLIFlags(EnvVarPrefix, FlagPrefix)...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag
