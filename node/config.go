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
	minExpirationPollIntervalSec = 3
	AppName                      = "da-node"
	SemVer                       = "0.5.0"
	GitCommit                    = ""
	GitDate                      = ""
)

var (
	// QuorumNames maps quorum IDs to their names.
	// this is used for eigen metrics
	QuorumNames = map[core.QuorumID]string{
		0: "eth_quorum",
		1: "permissioned_quorum",
	}
)

// Config contains all of the configuration information for a DA node.
type Config struct {
	Hostname                      string
	RetrievalPort                 string
	DispersalPort                 string
	InternalRetrievalPort         string
	InternalDispersalPort         string
	EnableNodeApi                 bool
	NodeApiPort                   string
	EnableMetrics                 bool
	MetricsPort                   string
	Timeout                       time.Duration
	RegisterNodeAtStart           bool
	ExpirationPollIntervalSec     uint64
	EnableTestMode                bool
	OverrideBlockStaleMeasure     int64
	OverrideStoreDurationBlocks   int64
	QuorumIDList                  []core.QuorumID
	DbPath                        string
	LogPath                       string
	PrivateBls                    string
	ID                            core.OperatorID
	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
	PubIPProvider                 string
	PubIPCheckInterval            time.Duration
	ChurnerUrl                    string
	NumBatchValidators            int
	ClientIPHeader                string
	UseSecureGrpc                 bool

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
	if expirationPollIntervalSec <= minExpirationPollIntervalSec {
		return nil, errors.New("the expiration-poll-interval flag must be greater than 3 seconds")
	}

	testMode := ctx.GlobalBool(flags.EnableTestModeFlag.Name)

	// Decrypt ECDSA key
	var ethClientConfig geth.EthClientConfig
	if !testMode {
		keyContents, err := os.ReadFile(ctx.GlobalString(flags.EcdsaKeyFileFlag.Name))
		if err != nil {
			return nil, fmt.Errorf("could not read ECDSA key file: %v", err)
		}
		sk, err := keystore.DecryptKey(keyContents, ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name))
		if err != nil {
			return nil, fmt.Errorf("could not decrypt the ECDSA file: %s", ctx.GlobalString(flags.EcdsaKeyFileFlag.Name))
		}
		ethClientConfig = geth.ReadEthClientConfigRPCOnly(ctx)
		ethClientConfig.PrivateKeyString = fmt.Sprintf("%x", crypto.FromECDSA(sk.PrivateKey))
	} else {
		ethClientConfig = geth.ReadEthClientConfig(ctx)
	}

	// Decrypt BLS key
	var privateBls string
	if !testMode {
		kp, err := bls.ReadPrivateKeyFromFile(ctx.GlobalString(flags.BlsKeyFileFlag.Name), ctx.GlobalString(flags.BlsKeyPasswordFlag.Name))
		if err != nil {
			return nil, fmt.Errorf("could not read or decrypt the BLS private key: %v", err)
		}
		privateBls = kp.PrivKey.String()
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
		Hostname:                      ctx.GlobalString(flags.HostnameFlag.Name),
		DispersalPort:                 ctx.GlobalString(flags.DispersalPortFlag.Name),
		RetrievalPort:                 ctx.GlobalString(flags.RetrievalPortFlag.Name),
		InternalDispersalPort:         internalDispersalFlag,
		InternalRetrievalPort:         internalRetrievalFlag,
		EnableNodeApi:                 ctx.GlobalBool(flags.EnableNodeApiFlag.Name),
		NodeApiPort:                   ctx.GlobalString(flags.NodeApiPortFlag.Name),
		EnableMetrics:                 ctx.GlobalBool(flags.EnableMetricsFlag.Name),
		MetricsPort:                   ctx.GlobalString(flags.MetricsPortFlag.Name),
		Timeout:                       timeout,
		RegisterNodeAtStart:           ctx.GlobalBool(flags.RegisterAtNodeStartFlag.Name),
		ExpirationPollIntervalSec:     expirationPollIntervalSec,
		EnableTestMode:                testMode,
		OverrideBlockStaleMeasure:     ctx.GlobalInt64(flags.OverrideBlockStaleMeasureFlag.Name),
		OverrideStoreDurationBlocks:   ctx.GlobalInt64(flags.OverrideStoreDurationBlocksFlag.Name),
		QuorumIDList:                  ids,
		DbPath:                        ctx.GlobalString(flags.DbPathFlag.Name),
		PrivateBls:                    privateBls,
		EthClientConfig:               ethClientConfig,
		EncoderConfig:                 kzg.ReadCLIConfig(ctx),
		LoggerConfig:                  *loggerConfig,
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		PubIPProvider:                 ctx.GlobalString(flags.PubIPProviderFlag.Name),
		PubIPCheckInterval:            ctx.GlobalDuration(flags.PubIPCheckIntervalFlag.Name),
		ChurnerUrl:                    ctx.GlobalString(flags.ChurnerUrlFlag.Name),
		NumBatchValidators:            ctx.GlobalInt(flags.NumBatchValidatorsFlag.Name),
		ClientIPHeader:                ctx.GlobalString(flags.ClientIPHeaderFlag.Name),
		UseSecureGrpc:                 ctx.GlobalBoolT(flags.ChurnerUseSecureGRPC.Name),
	}, nil
}
