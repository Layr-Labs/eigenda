package node

import (
	"errors"
	"fmt"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/node/flags"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/urfave/cli"
)

const (
	// Min number of seconds for the ExpirationPollIntervalSecFlag.
	minExpirationPollIntervalSec   = 3
	minReachabilityPollIntervalSec = 10
	AppName                        = "da-node"
)

var (
	// QuorumNames maps quorum IDs to their names.
	// this is used for eigen metrics
	QuorumNames = map[core.QuorumID]string{
		0: "eth_quorum",
		1: "eignen_quorum",
	}
	SemVer    = "0.0.0"
	GitCommit = ""
	GitDate   = ""
)

// Config contains all of the configuration information for a DA node.
type Config struct {
	Hostname                       string
	RetrievalPort                  string
	DispersalPort                  string
	InternalRetrievalPort          string
	InternalDispersalPort          string
	EnableNodeApi                  bool
	NodeApiPort                    string
	EnableMetrics                  bool
	MetricsPort                    string
	OnchainMetricsInterval         int64
	Timeout                        time.Duration
	RegisterNodeAtStart            bool
	ExpirationPollIntervalSec      uint64
	EnableTestMode                 bool
	OverrideBlockStaleMeasure      int64
	OverrideStoreDurationBlocks    int64
	QuorumIDList                   []core.QuorumID
	DbPath                         string
	LogPath                        string
	PrivateBls                     string
	ID                             core.OperatorID
	BLSOperatorStateRetrieverAddr  string
	EigenDAServiceManagerAddr      string
	PubIPProvider                  string
	PubIPCheckInterval             time.Duration
	ChurnerUrl                     string
	DataApiUrl                     string
	NumBatchValidators             int
	NumBatchDeserializationWorkers int
	EnableGnarkBundleEncoding      bool
	ClientIPHeader                 string
	UseSecureGrpc                  bool
	ReachabilityPollIntervalSec    uint64
	DisableNodeInfoResources       bool

	BLSRemoteSignerEnabled   bool
	BLSRemoteSignerUrl       string
	BLSPublicKeyHex          string
	BLSKeyPassword           string
	BLSSignerTLSCertFilePath string

	EthClientConfig geth.EthClientConfig
	LoggerConfig    common.LoggerConfig
	EncoderConfig   kzg.KzgConfig
}

// NewConfig parses the Config from the provided flags or environment variables and
// returns a Config.
func NewConfig(ctx *cli.Context) (*Config, error) {
	timeout, err := time.ParseDuration(ctx.GlobalString(flags.TimeoutFlag.Name))
	if err != nil {
		return &Config{}, err
	}

	idsStr := strings.Split(ctx.GlobalString(flags.QuorumIDListFlag.Name), ",")
	ids := make([]core.QuorumID, 0)
	for _, id := range idsStr {
		val, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, core.QuorumID(val))
	}
	if len(ids) == 0 {
		return nil, errors.New("no quorum ids provided")
	}

	expirationPollIntervalSec := ctx.GlobalUint64(flags.ExpirationPollIntervalSecFlag.Name)
	if expirationPollIntervalSec < minExpirationPollIntervalSec {
		return nil, fmt.Errorf("the expiration-poll-interval flag must be >= %d seconds", minExpirationPollIntervalSec)
	}

	reachabilityPollIntervalSec := ctx.GlobalUint64(flags.ReachabilityPollIntervalSecFlag.Name)
	if reachabilityPollIntervalSec != 0 && reachabilityPollIntervalSec < minReachabilityPollIntervalSec {
		return nil, fmt.Errorf("the reachability-poll-interval flag must be >= %d seconds or 0 to disable", minReachabilityPollIntervalSec)
	}

	testMode := ctx.GlobalBool(flags.EnableTestModeFlag.Name)

	// Configuration options that require the Node Operator ECDSA key at runtime
	registerNodeAtStart := ctx.GlobalBool(flags.RegisterAtNodeStartFlag.Name)
	pubIPCheckInterval := ctx.GlobalDuration(flags.PubIPCheckIntervalFlag.Name)
	needECDSAKey := registerNodeAtStart || pubIPCheckInterval > 0
	if registerNodeAtStart && (ctx.GlobalString(flags.EcdsaKeyFileFlag.Name) == "" || ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name) == "") {
		return nil, fmt.Errorf("%s and %s are required if %s is enabled", flags.EcdsaKeyFileFlag.Name, flags.EcdsaKeyPasswordFlag.Name, flags.RegisterAtNodeStartFlag.Name)
	}
	if pubIPCheckInterval > 0 && (ctx.GlobalString(flags.EcdsaKeyFileFlag.Name) == "" || ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name) == "") {
		return nil, fmt.Errorf("%s and %s are required if %s is > 0", flags.EcdsaKeyFileFlag.Name, flags.EcdsaKeyPasswordFlag.Name, flags.PubIPCheckIntervalFlag.Name)
	}

	var ethClientConfig geth.EthClientConfig
	if !testMode {
		ethClientConfig = geth.ReadEthClientConfigRPCOnly(ctx)
		if needECDSAKey {
			// Decrypt ECDSA key
			keyContents, err := os.ReadFile(ctx.GlobalString(flags.EcdsaKeyFileFlag.Name))
			if err != nil {
				return nil, fmt.Errorf("could not read ECDSA key file: %v", err)
			}
			sk, err := keystore.DecryptKey(keyContents, ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name))
			if err != nil {
				return nil, fmt.Errorf("could not decrypt the ECDSA file: %s", ctx.GlobalString(flags.EcdsaKeyFileFlag.Name))
			}
			ethClientConfig.PrivateKeyString = fmt.Sprintf("%x", crypto.FromECDSA(sk.PrivateKey))
		}
	} else {
		ethClientConfig = geth.ReadEthClientConfig(ctx)
	}

	// check if BLS remote signer configuration is provided
	blsRemoteSignerEnabled := ctx.GlobalBool(flags.BLSRemoteSignerEnabledFlag.Name)
	if blsRemoteSignerEnabled && (ctx.GlobalString(flags.BLSRemoteSignerUrlFlag.Name) == "" || ctx.GlobalString(flags.BLSPublicKeyHexFlag.Name) == "") {
		return nil, fmt.Errorf("BLS remote signer URL and Public Key Hex is required if BLS remote signer is enabled")
	}
	if !blsRemoteSignerEnabled && (ctx.GlobalString(flags.BlsKeyFileFlag.Name) == "" || ctx.GlobalString(flags.BlsKeyPasswordFlag.Name) == "") {
		return nil, fmt.Errorf("BLS key file and password is required if BLS remote signer is disabled")
	}

	// Decrypt BLS key
	var privateBls string
	if !testMode {
		// If remote signer fields are empty then try to read the BLS key from the file
		if !blsRemoteSignerEnabled {
			kp, err := bls.ReadPrivateKeyFromFile(ctx.GlobalString(flags.BlsKeyFileFlag.Name), ctx.GlobalString(flags.BlsKeyPasswordFlag.Name))
			if err != nil {
				return nil, fmt.Errorf("could not read or decrypt the BLS private key: %v", err)
			}
			privateBls = kp.PrivKey.String()
		}
	} else {
		privateBls = ctx.GlobalString(flags.TestPrivateBlsFlag.Name)
	}

	internalDispersalFlag := ctx.GlobalString(flags.InternalDispersalPortFlag.Name)
	internalRetrievalFlag := ctx.GlobalString(flags.InternalRetrievalPortFlag.Name)
	if internalDispersalFlag == "" {
		internalDispersalFlag = ctx.GlobalString(flags.DispersalPortFlag.Name)
	}
	if internalRetrievalFlag == "" {
		internalRetrievalFlag = ctx.GlobalString(flags.RetrievalPortFlag.Name)
	}

	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}

	return &Config{
		Hostname:                       ctx.GlobalString(flags.HostnameFlag.Name),
		DispersalPort:                  ctx.GlobalString(flags.DispersalPortFlag.Name),
		RetrievalPort:                  ctx.GlobalString(flags.RetrievalPortFlag.Name),
		InternalDispersalPort:          internalDispersalFlag,
		InternalRetrievalPort:          internalRetrievalFlag,
		EnableNodeApi:                  ctx.GlobalBool(flags.EnableNodeApiFlag.Name),
		NodeApiPort:                    ctx.GlobalString(flags.NodeApiPortFlag.Name),
		EnableMetrics:                  ctx.GlobalBool(flags.EnableMetricsFlag.Name),
		MetricsPort:                    ctx.GlobalString(flags.MetricsPortFlag.Name),
		OnchainMetricsInterval:         ctx.GlobalInt64(flags.OnchainMetricsIntervalFlag.Name),
		Timeout:                        timeout,
		RegisterNodeAtStart:            registerNodeAtStart,
		ExpirationPollIntervalSec:      expirationPollIntervalSec,
		ReachabilityPollIntervalSec:    reachabilityPollIntervalSec,
		EnableTestMode:                 testMode,
		OverrideBlockStaleMeasure:      ctx.GlobalInt64(flags.OverrideBlockStaleMeasureFlag.Name),
		OverrideStoreDurationBlocks:    ctx.GlobalInt64(flags.OverrideStoreDurationBlocksFlag.Name),
		QuorumIDList:                   ids,
		DbPath:                         ctx.GlobalString(flags.DbPathFlag.Name),
		PrivateBls:                     privateBls,
		EthClientConfig:                ethClientConfig,
		EncoderConfig:                  kzg.ReadCLIConfig(ctx),
		LoggerConfig:                   *loggerConfig,
		BLSOperatorStateRetrieverAddr:  ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:      ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		PubIPProvider:                  ctx.GlobalString(flags.PubIPProviderFlag.Name),
		PubIPCheckInterval:             pubIPCheckInterval,
		ChurnerUrl:                     ctx.GlobalString(flags.ChurnerUrlFlag.Name),
		DataApiUrl:                     ctx.GlobalString(flags.DataApiUrlFlag.Name),
		NumBatchValidators:             ctx.GlobalInt(flags.NumBatchValidatorsFlag.Name),
		NumBatchDeserializationWorkers: ctx.GlobalInt(flags.NumBatchDeserializationWorkersFlag.Name),
		EnableGnarkBundleEncoding:      ctx.Bool(flags.EnableGnarkBundleEncodingFlag.Name),
		ClientIPHeader:                 ctx.GlobalString(flags.ClientIPHeaderFlag.Name),
		UseSecureGrpc:                  ctx.GlobalBoolT(flags.ChurnerUseSecureGRPC.Name),
		DisableNodeInfoResources:       ctx.GlobalBool(flags.DisableNodeInfoResourcesFlag.Name),
		BLSRemoteSignerUrl:             ctx.GlobalString(flags.BLSRemoteSignerUrlFlag.Name),
		BLSPublicKeyHex:                ctx.GlobalString(flags.BLSPublicKeyHexFlag.Name),
		BLSKeyPassword:                 ctx.GlobalString(flags.BlsKeyPasswordFlag.Name),
		BLSSignerTLSCertFilePath:       ctx.GlobalString(flags.BLSSignerCertFileFlag.Name),
		BLSRemoteSignerEnabled:         blsRemoteSignerEnabled,
	}, nil
}
