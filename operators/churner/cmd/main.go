package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core/eth"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/Layr-Labs/eigenda/operators/churner/flags"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "churner"
	app.Usage = "EigenDA Churner"
	app.Description = "Service manages contract registrations, facilitates operator removal, and gathers deregistration information from operators."
	app.Flags = flags.Flags
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func run(ctx *cli.Context) error {
	log.Println("Initializing churner")
	hostname := "0.0.0.0"
	port := ctx.String(flags.GrpcPortFlag.Name)
	addr := fmt.Sprintf("%s:%s", hostname, port)
	log.Println("Starting churner server at", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300)
	gs := grpc.NewServer(
		opt,
		grpc.ChainUnaryInterceptor(),
	)

	config, err := churner.NewConfig(ctx)
	if err != nil {
		log.Fatalf("failed to parse the command line flags: %v", err)
	}
	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	log.Println("Starting geth client")
	gethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	tx, err := eth.NewWriter(logger, gethClient, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Fatalln("could not create new transactor", err)
	}

	cs := coreeth.NewChainState(tx, gethClient)

	logger.Info("Using graph node")

	logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)
	indexer := thegraph.MakeIndexedChainState(config.ChainStateConfig, cs, logger)

	metrics := churner.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	cn, err := churner.NewChurner(config, indexer, tx, logger, metrics)
	if err != nil {
		log.Fatalln("cannot create churner", err)
	}

	churnerServer := churner.NewServer(config, cn, logger, metrics)
	if err = churnerServer.Start(config.MetricsConfig); err != nil {
		log.Fatalln("failed to start churner server", err)
	}

	// Register reflection service on gRPC server
	// This makes "grpcurl -plaintext localhost:9000 list" command work
	reflection.Register(gs)

	pb.RegisterChurnerServer(gs, churnerServer)

	// Register Server for Health Checks
	name := pb.Churner_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, gs)

	log.Printf("churner server listening at %s", addr)
	return gs.Serve(listener)
}
