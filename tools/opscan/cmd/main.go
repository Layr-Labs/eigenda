package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigenda/tools/opscan"
	"github.com/Layr-Labs/eigenda/tools/opscan/flags"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	version   = ""
	gitCommit = ""
	gitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s,%s,%s", version, gitCommit, gitDate)
	app.Name = "opscan"
	app.Description = "operator network scanner"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = RunScan
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunScan(ctx *cli.Context) error {
	config, err := opscan.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	subgraphApi := subgraph.NewApi(config.SubgraphEndpoint, config.SubgraphEndpoint)
	subgraphClient := dataapi.NewSubgraphClient(subgraphApi, logger)

	if config.OperatorId != "" {
		operatorInfo, err := subgraphClient.QueryOperatorInfoByOperatorId(context.Background(), config.OperatorId)
		if err != nil {
			logger.Warn("failed to fetch operator info", "operatorId", config.OperatorId, "error", err)
			return errors.New("operator info not found")
		}

		operatorSocket := core.OperatorSocket(operatorInfo.Socket)
		retrievalSocket := operatorSocket.GetRetrievalSocket()
		getNodeInfo(context.Background(), retrievalSocket, config.OperatorId, logger)

	}
	return nil
}

func getNodeInfo(ctx context.Context, socket string, operatorId string, logger logging.Logger) {
	conn, err := grpc.Dial(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to dial grpc operator socket", "operatorId", operatorId, "socket", socket, "error", err)
		return
	}
	defer conn.Close()
	client := node.NewRetrievalClient(conn)
	reply, err := client.NodeInfo(ctx, &node.NodeInfoRequest{})
	if err != nil {
		logger.Info("NodeInfo", "operatorId", operatorId, "semver", "unknown")
		return
	}

	logger.Info("NodeInfo", "operatorId", operatorId, "semver", reply.Semver, "os", reply.Os, "arch", reply.Arch, "numCpu", reply.NumCpu, "memBytes", reply.MemBytes)
}
