package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	coreindexer "github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/retriever"
	retrivereth "github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigenda/retriever/flags"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/shurcooL/graphql"
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
	v, err := verifier.NewVerifier(&config.EncoderConfig, false)
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

	tx, err := eth.NewTransactor(logger, gethClient, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}
	cs := eth.NewChainState(tx, gethClient)
	rpcClient, err := rpc.Dial(config.EthClientConfig.RPCURL)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	var ics core.IndexedChainState
	if config.UseGraph {
		logger.Info("Using graph node")
		querier := graphql.NewClient(config.GraphUrl, nil)
		logger.Info("Connecting to subgraph", "url", config.GraphUrl)
		ics = thegraph.NewIndexedChainState(cs, querier, logger)
	} else {
		logger.Info("Using built-in indexer")

		indexer, err := coreindexer.CreateNewIndexer(
			&config.IndexerConfig,
			gethClient,
			rpcClient,
			config.EigenDAServiceManagerAddr,
			logger,
		)
		if err != nil {
			return err
		}
		ics, err = coreindexer.NewIndexedChainState(cs, indexer)
		if err != nil {
			return err
		}
	}

	agn := &core.StdAssignmentCoordinator{}
	retrievalClient, err := clients.NewRetrievalClient(logger, ics, agn, nodeClient, v, config.NumConnections)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}

	chainClient := retrivereth.NewChainClient(gethClient, logger)
	retrieverServiceServer := retriever.NewServer(config, logger, retrievalClient, v, ics, chainClient)
	if err = retrieverServiceServer.Start(context.Background()); err != nil {
		log.Fatalln("failed to start retriever service server", err)
	}

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
