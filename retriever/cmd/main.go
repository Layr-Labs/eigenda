package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/indexer/inmem"
	"github.com/Layr-Labs/eigenda/retriever"
	retrivereth "github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigenda/retriever/flags"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
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
	port := ctx.String(flags.GrpcPortFlag.Name)
	addr := fmt.Sprintf("0.0.0.0:%s", port)
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

	config := retriever.NewConfig(ctx)
	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	nodeClient := clients.NewNodeClient(config.Timeout)
	encoder, err := encoding.NewEncoder(config.EncoderConfig)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}
	gethClient, err := geth.NewClient(config.EthClientConfig, logger)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	// TODO(ian-shim): uncomment when https://github.com/Layr-Labs/eigenda-internal/issues/77 is done
	// store, err := leveldb.NewHeaderStore(config.IndexerDataDir)
	// if err != nil {
	// 	return err
	// }
	store := inmem.NewHeaderStore()

	tx, err := eth.NewTransactor(logger, gethClient, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}
	cs := eth.NewChainState(tx, gethClient)
	rpcClient, err := rpc.Dial(config.EthClientConfig.RPCURL)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}
	indexedState, err := indexer.NewIndexedChainState(&config.IndexerConfig, common.HexToAddress(config.EigenDAServiceManagerAddr), cs, store, gethClient, rpcClient, logger)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	agn := &core.StdAssignmentCoordinator{}
	retrievalClient := clients.NewRetrievalClient(logger, indexedState, agn, nodeClient, encoder, config.NumConnections)

	chainClient := retrivereth.NewChainClient(gethClient, logger)
	retrieverServiceServer := retriever.NewServer(config, logger, retrievalClient, encoder, indexedState, chainClient)
	if err = retrieverServiceServer.Start(context.Background()); err != nil {
		log.Fatalln("failed to start retriever service server", err)
	}

	// Register reflection service on gRPC server
	// This makes "grpcurl -plaintext localhost:9000 list" command work
	reflection.Register(gs)

	pb.RegisterRetrieverServer(gs, retrieverServiceServer)

	// Register Server for Health Checks
	healthcheck.RegisterHealthServer(gs)

	log.Printf("server listening at %s", addr)
	return gs.Serve(listener)
}
