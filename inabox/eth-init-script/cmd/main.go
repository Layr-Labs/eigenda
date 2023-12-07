package main

import (
	"log"
	"net"
	"os"

	"github.com/Layr-Labs/eigenda/inabox/config"
	ethinitscript "github.com/Layr-Labs/eigenda/inabox/eth-init-script"
	"github.com/urfave/cli"
)

const LocalAnvilFlagName = "local-anvil"

func main() {
	app := &cli.App{
		Flags:  []cli.Flag{},
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	cfg := config.OpenCwdConfigLock()

	_ = ethinitscript.SetupEigenDA(cfg, "/contracts")
	// we don't need the cmd object bec we won't ever use it to stop anvil

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
