package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	pbv2 "github.com/Layr-Labs/eigenda/api/grpc/retriever/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/retriever"
	retrivereth "github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigenda/retriever/flags"
	retrieverv2 "github.com/Layr-Labs/eigenda/retriever/v2"
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
	app.Name = "retriever"
	app.Usage = "EigenDA Retriever"
	app.Description = "Service for collecting coded chunks and decode the original data"
	app.Flags = flags.Flags
	app.Action = RetrieverMain
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RetrieverMain(ctx *cli.Context) error {
	log.Println("Initializing Retriever")
	hostname := ctx.String(flags.HostnameFlag.Name)
	port := ctx.String(flags.GrpcPortFlag.Name)
	addr := fmt.Sprintf("%s:%s", hostname, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300)
	gs := grpc.NewServer(
		opt,
		grpc.ChainUnaryInterceptor(
		// TODO(ian-shim): Add interceptors
		// correlation.UnaryServerInterceptor(),
		// logger.UnaryServerInterceptor(*s.logger.Logger),
		),
	)

	config, err := retriever.NewConfig(ctx)
	if err != nil {
		log.Fatalf("failed to parse the command line flags: %v", err)
	}
	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	nodeClient := clients.NewNodeClient(config.Timeout)

	config.EncoderConfig.LoadG2Points = true
	v, err := verifier.NewVerifier(&config.EncoderConfig, nil)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}
	gethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	tx, err := eth.NewReader(logger, gethClient, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}
	cs, err := eth.NewChainState(tx, gethClient, logger)
	if err != nil {
		log.Fatalln("could not start chain state", err)
	}

	if config.EigenDAVersion == 1 {

		agn := &core.StdAssignmentCoordinator{}
		retrievalClient, err := clients.NewRetrievalClient(logger, cs, agn, nodeClient, v, config.NumConnections)
		if err != nil {
			log.Fatalln("could not start tcp listener", err)
		}

		chainClient := retrivereth.NewChainClient(gethClient, logger)
		retrieverServiceServer := retriever.NewServer(config, logger, retrievalClient, cs, chainClient)
		// This only start the metrics server; consider unwrapping the function
		retrieverServiceServer.Start(context.Background())

		// Register reflection service on gRPC server
		// This makes "grpcurl -plaintext localhost:9000 list" command work
		reflection.Register(gs)

		pb.RegisterRetrieverServer(gs, retrieverServiceServer)

		// Register Server for Health Checks
		name := pb.Retriever_ServiceDesc.ServiceName
		healthcheck.RegisterHealthServer(name, gs)

		log.Printf("server listening at %s", addr)
		return gs.Serve(listener)
	}

	if config.EigenDAVersion == 2 {
		retrievalClient := clientsv2.NewRetrievalClient(logger, tx, cs, v, config.NumConnections)
		retrieverServiceServer := retrieverv2.NewServer(config, logger, retrievalClient, cs)
		retrieverServiceServer.Start(context.Background())

		// Register reflection service on gRPC server
		// This makes "grpcurl -plaintext localhost:9000 list" command work
		reflection.Register(gs)

		pbv2.RegisterRetrieverServer(gs, retrieverServiceServer)

		// Register Server for Health Checks
		name := pb.Retriever_ServiceDesc.ServiceName
		healthcheck.RegisterHealthServer(name, gs)

		log.Printf("server listening at %s", addr)
		return gs.Serve(listener)
	}

	return errors.New("invalid EigenDA version")
}
