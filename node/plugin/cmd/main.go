package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/plugin"
	blssigner "github.com/Layr-Labs/eigensdk-go/signer/bls"
	blssignerTypes "github.com/Layr-Labs/eigensdk-go/signer/bls/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		plugin.OperationFlag,
		plugin.EcdsaKeyFileFlag,
		plugin.BlsKeyFileFlag,
		plugin.EcdsaKeyPasswordFlag,
		plugin.BlsKeyPasswordFlag,
		plugin.SocketFlag,
		plugin.QuorumIDListFlag,
		plugin.ChainRpcUrlFlag,
		plugin.EigenDADirectoryFlag,
		plugin.BlsOperatorStateRetrieverFlag,
		plugin.EigenDAServiceManagerFlag,
		plugin.ChurnerUrlFlag,
		plugin.NumConfirmationsFlag,
		plugin.PubIPProviderFlag,
		plugin.BLSRemoteSignerUrlFlag,
		plugin.BLSPublicKeyHexFlag,
		plugin.BLSSignerCertFileFlag,
		plugin.BLSSignerAPIKeyFlag,
	}
	app.Name = "eigenda-node-plugin"
	app.Usage = "EigenDA Node Plugin"
	app.Version = fmt.Sprintf("%s %s %s", node.SemVer, node.GitCommit, node.GitDate)
	app.Description = "Run one time operations like avs opt-in/opt-out for EigenDA Node"
	app.Action = pluginOps
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Application failed.", "Message:", err)
	}
}

func pluginOps(ctx *cli.Context) {
	config, err := plugin.NewConfig(ctx)
	if err != nil {
		log.Printf("Error: failed to parse the command line flags: %v", err)
		return
	}
	log.Printf("Info: plugin configs and flags parsed")

	signerCfg := blssignerTypes.SignerConfig{
		PublicKeyHex:     config.BLSPublicKeyHex,
		CerberusUrl:      config.BLSRemoteSignerUrl,
		CerberusPassword: config.BlsKeyPassword,
		TLSCertFilePath:  config.BLSSignerCertFile,
		Path:             config.BlsKeyFile,
		Password:         config.BlsKeyPassword,
		CerberusAPIKey:   config.BLSSignerAPIKey,
	}
	if config.BLSRemoteSignerUrl != "" {
		signerCfg.SignerType = blssignerTypes.Cerberus
	} else {
		signerCfg.SignerType = blssignerTypes.Local
	}
	signer, err := blssigner.NewSigner(signerCfg)
	if err != nil {
		log.Printf("Error: failed to create BLS signer: %v", err)
		return
	}

	opID, err := signer.GetOperatorId()
	if err != nil {
		log.Printf("Error: failed to get operator ID: %v", err)
		return
	}
	operatorID, err := core.OperatorIDFromHex(opID)
	if err != nil {
		log.Printf("Error: failed to convert operator ID: %v", err)
		return
	}
	pubKeyG1Hex := signer.GetPublicKeyG1()
	pubKeyG1, err := hex.DecodeString(pubKeyG1Hex)
	if err != nil {
		log.Printf("Error: failed to decode public key G1: %v", err)
		return
	}
	pubKeyG1Point := new(core.G1Point)
	pubKeyG1Point, err = pubKeyG1Point.Deserialize(pubKeyG1)
	if err != nil {
		log.Printf("Error: failed to deserialize public key G1: %v", err)
		return
	}

	sk, privateKey, err := plugin.GetECDSAPrivateKey(config.EcdsaKeyFile, config.EcdsaKeyPassword)
	if err != nil {
		log.Printf("Error: failed to read or decrypt the ECDSA from file (%s) for private key: %v", config.EcdsaKeyFile, err)
		return
	}
	log.Printf("Info: ECDSA key read and decrypted from %s", config.EcdsaKeyFile)

	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		log.Printf("Error: failed to create logger: %v", err)
		return
	}

	ethConfig := geth.EthClientConfig{
		RPCURLs:          []string{config.ChainRpcUrl},
		PrivateKeyString: *privateKey,
		NumConfirmations: config.NumConfirmations,
	}
	client, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	if err != nil {
		log.Printf("Error: failed to create eth client: %v", err)
		return
	}
	log.Printf("Info: ethclient created for url: %s", config.ChainRpcUrl)

	tx, err := eth.NewWriter(logger, client, config.EigenDADirectory, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Printf("Error: failed to create EigenDA transactor: %v", err)
		return
	}

	_, dispersalPort, retrievalPort, v2DispersalPort, v2RetrievalPort, err := core.ParseOperatorSocket(config.Socket)
	if err != nil {
		log.Printf("Error: failed to parse operator socket: %v", err)
		return
	}

	socket := config.Socket
	if isLocalhost(socket) {
		pubIPProvider := pubip.ProviderOrDefault(logger, config.PubIPProvider)
		socket, err = node.SocketAddress(context.Background(), pubIPProvider, dispersalPort, retrievalPort, v2DispersalPort, v2RetrievalPort)
		if err != nil {
			log.Printf("Error: failed to get socket address from ip provider: %v", err)
			return
		}
	}

	operator := &node.Operator{
		Address:             sk.Address.Hex(),
		Socket:              socket,
		Timeout:             10 * time.Second,
		PrivKey:             sk.PrivateKey,
		Signer:              signer,
		OperatorId:          operatorID,
		QuorumIDs:           config.QuorumIDList,
		RegisterNodeAtStart: false,
	}
	churnerClient := node.NewChurnerClient(config.ChurnerUrl, true, operator.Timeout, logger)
	switch config.Operation {
	case plugin.OperationOptIn:
		log.Printf("Info: Operator with Operator Address: %x is opting in to EigenDA", sk.Address)
		err = node.RegisterOperator(context.Background(), operator, tx, churnerClient, logger.With("component", "NodeOperator"))
		if err != nil {
			log.Printf("Error: failed to opt-in EigenDA Node Network for operator ID: %x, operator address: %x, error: %v", operatorID, sk.Address, err)
			return
		}
		log.Printf("Info: successfully opt-in the EigenDA, for operator ID: %x, operator address: %x, socket: %s, and quorums: %v", operatorID, sk.Address, config.Socket, config.QuorumIDList)
	case plugin.OperationOptOut:
		log.Printf("Info: Operator with Operator Address: %x and OperatorID: %x is opting out of EigenDA", sk.Address, operatorID)
		err = node.DeregisterOperator(context.Background(), operator, pubKeyG1Point, tx)
		if err != nil {
			log.Printf("Error: failed to opt-out EigenDA Node Network for operator ID: %x, operator address: %x, quorums: %v, error: %v", operatorID, sk.Address, config.QuorumIDList, err)
			return
		}
		log.Printf("Info: successfully opt-out the EigenDA, for operator ID: %x, operator address: %x", operatorID, sk.Address)
	case plugin.OperationUpdateSocket:
		log.Printf("Info: Operator with Operator Address: %x is updating its socket: %s", sk.Address, config.Socket)
		err = node.UpdateOperatorSocket(context.Background(), tx, config.Socket)
		if err != nil {
			log.Printf("Error: failed to update socket for operator ID: %x, operator address: %x, socket: %s, error: %v", operatorID, sk.Address, config.Socket, err)
			return
		}
		log.Printf("Info: successfully updated socket, for operator ID: %x, operator address: %x, socket: %s", operatorID, sk.Address, config.Socket)
	case plugin.OperationListQuorums:
		quorumIds, err := tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
		if err != nil {
			log.Printf("Error: failed to get quorum(s) for operatorID: %x, operator address: %x, error: %v", operatorID, sk.Address, err)
			return
		}
		log.Printf("Info: operator ID: %x, operator address: %x, current quorums: %v", operatorID, sk.Address, quorumIds)
	default:
		log.Fatalf("Fatal: unsupported operation: %s", config.Operation)
	}
}

func isLocalhost(socket string) bool {
	return strings.Contains(socket, "localhost") || strings.Contains(socket, "127.0.0.1") || strings.Contains(socket, "0.0.0.0")
}
