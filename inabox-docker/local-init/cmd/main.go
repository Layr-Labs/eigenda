package main

import (
	"log"
	"net"
	"os"

	localinit "github.com/Layr-Labs/eigenda/inabox-docker/local-init"
	"github.com/urfave/cli/v2"
)

var (
	localstackFlagName = "localstack-port"

	metadataTableName = "test-BlobMetadata"
	bucketTableName   = "test-BucketStore"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  localstackFlagName,
				Value: "4570",
				Usage: "path to the config file",
			},
		},
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {

	err := localinit.DeployResources(nil, ctx.String(localstackFlagName), metadataTableName, bucketTableName)
	if err != nil {
		return err
	}

	config := localinit.NewTestConfig()

	defer config.SaveTestConfig()

	log.Print("Deploying experiment...")

	// Create a new experiment and deploy the contracts
	err = config.LoadPrivateKeys()
	if err != nil {
		log.Panicf("could not load private keys: %v", err)
	}

	// if config.EigenDA.Deployer != "" && !config.IsEigenDADeployed() {
	log.Print("Deploying EigenDA")
	config.DeployEigenDAContracts()
	// }

	if deployer, ok := config.GetDeployer(config.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		startBlock := localinit.GetLatestBlockNumber(config.Deployers[0].RPC)
		config.DeploySubgraphs(startBlock)
	}

	log.Print("Test environment has succesfully deployed!")

	// Indicates that the script is done and other docker compose services can start
	startHealthcheckServer()

	return nil
}

func startHealthcheckServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()
	log.Println("TCP server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		conn.Close() // Close the connection immediately for the health check purpose
	}
}
