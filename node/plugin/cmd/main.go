package main

import (
	"context"
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
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
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
		plugin.BlsOperatorStateRetrieverFlag,
		plugin.EigenDAServiceManagerFlag,
		plugin.ChurnerUrlFlag,
		plugin.NumConfirmationsFlag,
		plugin.PubIPProviderFlag,
	}
	app.Name = "eigenda-node-plugin"
	app.Usage = "EigenDA Node Plugin"
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

	kp, err := bls.ReadPrivateKeyFromFile(config.BlsKeyFile, config.BlsKeyPassword)
	if err != nil {
		log.Printf("Error: failed to read or decrypt the BLS private key: %v", err)
		return
	}
	g1point := &core.G1Point{
		G1Affine: kp.PubKey.G1Affine,
	}
	keyPair := &core.KeyPair{
		PrivKey: kp.PrivKey,
		PubKey:  g1point,
	}
	log.Printf("Info: Bls key read and decrypted from %s", config.BlsKeyFile)

	operatorID := keyPair.GetPubKeyG1().GetOperatorID()

	sk, privateKey, err := plugin.GetECDSAPrivateKey(config.EcdsaKeyFile, config.EcdsaKeyPassword)
	if err != nil {
		log.Printf("Error: failed to read or decrypt the ECDSA private key: %v", err)
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

	tx, err := eth.NewWriter(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Printf("Error: failed to create EigenDA transactor: %v", err)
		return
	}

	_, dispersalPort, retrievalPort, err := core.ParseOperatorSocket(config.Socket)
	if err != nil {
		log.Printf("Error: failed to parse operator socket: %v", err)
		return
	}

	socket := config.Socket
	if isLocalhost(socket) {
		pubIPProvider := pubip.ProviderOrDefault(config.PubIPProvider)
		socket, err = node.SocketAddress(context.Background(), pubIPProvider, dispersalPort, retrievalPort)
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
		KeyPair:             keyPair,
		OperatorId:          keyPair.GetPubKeyG1().GetOperatorID(),
		QuorumIDs:           config.QuorumIDList,
		RegisterNodeAtStart: false,
	}
	churnerClient := node.NewChurnerClient(config.ChurnerUrl, true, operator.Timeout, logger)
	if config.Operation == plugin.OperationOptIn {
		log.Printf("Info: Operator with Operator Address: %x is opting in to EigenDA", sk.Address)
		err = node.RegisterOperator(context.Background(), operator, tx, churnerClient, logger.With("component", "NodeOperator"))
		if err != nil {
			log.Printf("Error: failed to opt-in EigenDA Node Network for operator ID: %x, operator address: %x, error: %v", operatorID, sk.Address, err)
			return
		}
		log.Printf("Info: successfully opt-in the EigenDA, for operator ID: %x, operator address: %x, socket: %s, and quorums: %v", operatorID, sk.Address, config.Socket, config.QuorumIDList)
	} else if config.Operation == plugin.OperationOptOut {
		log.Printf("Info: Operator with Operator Address: %x and OperatorID: %x is opting out of EigenDA", sk.Address, operatorID)
		err = node.DeregisterOperator(context.Background(), operator, keyPair, tx)
		if err != nil {
			log.Printf("Error: failed to opt-out EigenDA Node Network for operator ID: %x, operator address: %x, quorums: %v, error: %v", operatorID, sk.Address, config.QuorumIDList, err)
			return
		}
		log.Printf("Info: successfully opt-out the EigenDA, for operator ID: %x, operator address: %x", operatorID, sk.Address)
	} else if config.Operation == plugin.OperationUpdateSocket {
		log.Printf("Info: Operator with Operator Address: %x is updating its socket: %s", sk.Address, config.Socket)
		err = node.UpdateOperatorSocket(context.Background(), tx, config.Socket)
		if err != nil {
			log.Printf("Error: failed to update socket for operator ID: %x, operator address: %x, socket: %s, error: %v", operatorID, sk.Address, config.Socket, err)
			return
		}
		log.Printf("Info: successfully updated socket, for operator ID: %x, operator address: %x, socket: %s", operatorID, sk.Address, config.Socket)
	} else if config.Operation == plugin.OperationListQuorums {
		quorumIds, err := tx.GetRegisteredQuorumIdsForOperator(context.Background(), operatorID)
		if err != nil {
			log.Printf("Error: failed to get quorum(s) for operatorID: %x, operator address: %x, error: %v", operatorID, sk.Address, err)
			return
		}
		log.Printf("Info: operator ID: %x, operator address: %x, current quorums: %v", operatorID, sk.Address, quorumIds)
	} else {
		log.Fatalf("Fatal: unsupported operation: %s", config.Operation)
	}
}

func isLocalhost(socket string) bool {
	return strings.Contains(socket, "localhost") || strings.Contains(socket, "127.0.0.1") || strings.Contains(socket, "0.0.0.0")
}
