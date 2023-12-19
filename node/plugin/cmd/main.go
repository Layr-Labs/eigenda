package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/plugin"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

// Returns the decrypted ECDSA private key from the given file.
func getECDSAPrivateKey(keyFile string, password string) (*keystore.Key, *string, error) {
	keyContents, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	sk, err := keystore.DecryptKey(keyContents, password)
	if err != nil {
		return nil, nil, err
	}
	privateKey := fmt.Sprintf("%x", crypto.FromECDSA(sk.PrivateKey))
	return sk, &privateKey, nil
}

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

	sk, privateKey, err := getECDSAPrivateKey(config.EcdsaKeyFile, config.EcdsaKeyPassword)
	if err != nil {
		log.Printf("Error: failed to read or decrypt the ECDSA private key: %v", err)
		return
	}
	log.Printf("Info: ECDSA key read and decrypted from %s", config.EcdsaKeyFile)

	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	if err != nil {
		log.Printf("Error: failed to create a EigenDA logger: %v", err)
		return
	}

	ethConfig := geth.EthClientConfig{
		RPCURL:           config.ChainRpcUrl,
		PrivateKeyString: *privateKey,
		NumConfirmations: config.NumConfirmations,
	}
	client, err := geth.NewClient(ethConfig, logger)
	if err != nil {
		log.Printf("Error: failed to create eth client: %v", err)
		return
	}
	log.Printf("Info: ethclient created for url: %s", config.ChainRpcUrl)

	tx, err := eth.NewTransactor(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Printf("Error: failed to create EigenDA transactor: %v", err)
		return
	}

	s := strings.Split(config.Socket, ":")
	if len(s) != 2 {
		log.Printf("Error: invalid socket address format: %s", config.Socket)
		return
	}
	hostname := s[0]
	dispersalPort := s[1]

	s = strings.Split(config.Socket, ";")
	if len(s) != 2 {
		log.Printf("Error: invalid socket address format, missing retrieval port: %s", config.Socket)
		return
	}
	retrievalPort := s[1]

	socket := string(core.MakeOperatorSocket(hostname, dispersalPort, retrievalPort))
	if isLocalhost(socket) {
		pubIPProvider := pubip.ProviderOrDefault(config.PubIPProvider)
		socket, err = node.SocketAddress(context.Background(), pubIPProvider, dispersalPort, retrievalPort)
		if err != nil {
			log.Printf("Error: failed to get socket address from ip provider: %v", err)
			return
		}
	}

	operator := &node.Operator{
		Socket:     socket,
		Timeout:    10 * time.Second,
		KeyPair:    keyPair,
		OperatorId: keyPair.GetPubKeyG1().GetOperatorID(),
		QuorumIDs:  config.QuorumIDList,
	}
	if config.Operation == "opt-in" {
		log.Printf("Info: Operator with Operator Address: %x is opting in to EigenDA", sk.Address)
		err = node.RegisterOperator(context.Background(), operator, tx, config.ChurnerUrl, true, logger)
		if err != nil {
			log.Printf("Error: failed to opt-in EigenDA Node Network for operator ID: %x, operator address: %x, error: %v", operatorID, sk.Address, err)
			return
		}
		log.Printf("Info: successfully opt-in the EigenDA, for operator ID: %x, operator address: %x, socket: %s, and quorums: %v", operatorID, sk.Address, config.Socket, config.QuorumIDList)
	} else if config.Operation == "opt-out" {
		log.Printf("Info: Operator with Operator Address: %x and OpearatorID: %x is opting out of EigenDA", sk.Address, operatorID)
		err = node.DeregisterOperator(context.Background(), keyPair, tx)
		if err != nil {
			log.Printf("Error: failed to opt-out EigenDA Node Network for operator ID: %x, operator address: %x, error: %v", operatorID, sk.Address, err)
			return
		}
		log.Printf("Info: successfully opt-out the EigenDA, for operator ID: %x, operator address: %x", operatorID, sk.Address)
	} else {
		log.Fatalf("Fatal: unsupported operation: %s", config.Operation)
	}
}

func isLocalhost(socket string) bool {
	return strings.Contains(socket, "localhost") || strings.Contains(socket, "127.0.0.1") || strings.Contains(socket, "0.0.0.0")
}
