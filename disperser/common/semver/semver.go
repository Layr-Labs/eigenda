package semver

import (
	"context"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// query operator host info endpoint if available
func GetSemverInfo(ctx context.Context, socket string, operatorId string, logger logging.Logger) string {
	conn, err := grpc.Dial(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "unreachable"
	}
	defer conn.Close()
	client := node.NewDispersalClient(conn)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(3))
	defer cancel()
	reply, err := client.NodeInfo(ctxWithTimeout, &node.NodeInfoRequest{})
	if err != nil {
		var semver string
		if strings.Contains(err.Error(), "Unimplemented") {
			semver = "<0.8.0"
		} else if strings.Contains(err.Error(), "DeadlineExceeded") {
			semver = "timeout"
		} else if strings.Contains(err.Error(), "Unavailable") {
			semver = "refused"
		} else {
			semver = "error"
		}

		logger.Warn("NodeInfo", "operatorId", operatorId, "semver", semver, "error", err)
		return semver
	}

	// local node source compiles without semver
	if reply.Semver == "" {
		reply.Semver = "src-compile"
	}

	logger.Info("NodeInfo", "operatorId", operatorId, "socker", socket, "semver", reply.Semver, "os", reply.Os, "arch", reply.Arch, "numCpu", reply.NumCpu, "memBytes", reply.MemBytes)
	return reply.Semver
}
