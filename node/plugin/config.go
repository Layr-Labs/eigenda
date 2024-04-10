package plugin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/node/flags"
	"github.com/urfave/cli"
)

var (
	/* Required Flags */

	PubIPProviderFlag = cli.StringFlag{
		Name:     "public-ip-provider",
		Usage:    "The ip provider service used to obtain a operator's public IP [seeip (default), ipify)",
		Required: true,
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "PUBLIC_IP_PROVIDER"),
	}

	// The operation to run.
	OperationFlag = cli.StringFlag{
		Name:     "operation",
		Required: true,
		Usage:    "Supported operations: opt-in, opt-out, list-quorums",
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "OPERATION"),
	}

	// The files for encrypted private keys.
	EcdsaKeyFileFlag = cli.StringFlag{
		Name:     "ecdsa-key-file",
		Required: true,
		Usage:    "Path to the encrypted ecdsa key",
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "ECDSA_KEY_FILE"),
	}
	BlsKeyFileFlag = cli.StringFlag{
		Name:     "bls-key-file",
		Required: true,
		Usage:    "Path to the encrypted bls key",
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "BLS_KEY_FILE"),
	}

	// The passwords to decrypt the private keys.
	EcdsaKeyPasswordFlag = cli.StringFlag{
		Name:     "ecdsa-key-password",
		Required: true,
		Usage:    "Password to decrypt the ecdsa key",
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "ECDSA_KEY_PASSWORD"),
	}
	BlsKeyPasswordFlag = cli.StringFlag{
		Name:     "bls-key-password",
		Required: true,
		Usage:    "Password to decrypt the bls key",
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "BLS_KEY_PASSWORD"),
	}

	// The socket and the quorums to register.
	SocketFlag = cli.StringFlag{
		Name:     "socket",
		Required: true,
		Usage:    "The socket of the EigenDA Node for serving dispersal and retrieval",
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "SOCKET"),
	}
	QuorumIDListFlag = cli.StringFlag{
		Name:     "quorum-id-list",
		Usage:    "Comma separated list of quorum IDs that the node will opt-in or opt-out, depending on the OperationFlag. If OperationFlag is opt-in, all quorums should not have been registered already; if it's opt-out, all quorums should have been registered already",
		Required: true,
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "QUORUM_ID_LIST"),
	}

	// The chain and contract addresses to register with.
	ChainRpcUrlFlag = cli.StringFlag{
		Name:     "chain-rpc",
		Usage:    "Chain rpc url",
		Required: true,
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "CHAIN_RPC"),
	}
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     "bls-operator-state-retriever",
		Usage:    "Address of the BLS Operator State Retriever",
		Required: true,
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "BLS_OPERATOR_STATE_RETRIEVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     "eigenda-service-manager",
		Usage:    "Address of the EigenDA Service Manager",
		Required: true,
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
	ChurnerUrlFlag = cli.StringFlag{
		Name:     "churner-url",
		Usage:    "URL of the Churner",
		Required: true,
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "CHURNER_URL"),
	}
	NumConfirmationsFlag = cli.IntFlag{
		Name:     "num-confirmations",
		Usage:    "Number of confirmations to wait for",
		Required: false,
		Value:    3,
		EnvVar:   common.PrefixEnvVar(flags.EnvVarPrefix, "NUM_CONFIRMATIONS"),
	}
)

type Config struct {
	PubIPProvider                 string
	Operation                     string
	EcdsaKeyFile                  string
	BlsKeyFile                    string
	EcdsaKeyPassword              string
	BlsKeyPassword                string
	Socket                        string
	QuorumIDList                  []core.QuorumID
	ChainRpcUrl                   string
	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
	ChurnerUrl                    string
	NumConfirmations              int
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	idsStr := strings.Split(ctx.GlobalString(QuorumIDListFlag.Name), ",")
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

	op := ctx.GlobalString(OperationFlag.Name)
	if op != "opt-in" && op != "opt-out" && op != "list-quorums" {
		return nil, errors.New("unsupported operation type")
	}

	return &Config{
		PubIPProvider:                 ctx.GlobalString(PubIPProviderFlag.Name),
		Operation:                     op,
		EcdsaKeyPassword:              ctx.GlobalString(EcdsaKeyPasswordFlag.Name),
		BlsKeyPassword:                ctx.GlobalString(BlsKeyPasswordFlag.Name),
		EcdsaKeyFile:                  ctx.GlobalString(EcdsaKeyFileFlag.Name),
		BlsKeyFile:                    ctx.GlobalString(BlsKeyFileFlag.Name),
		Socket:                        ctx.GlobalString(SocketFlag.Name),
		QuorumIDList:                  ids,
		ChainRpcUrl:                   ctx.GlobalString(ChainRpcUrlFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(EigenDAServiceManagerFlag.Name),
		ChurnerUrl:                    ctx.GlobalString(ChurnerUrlFlag.Name),
		NumConfirmations:              ctx.GlobalInt(NumConfirmationsFlag.Name),
	}, nil
}
