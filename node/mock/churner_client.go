package mock

import (
	"context"

	churnerpb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/mock"
)

type ChurnerClient struct {
	mock.Mock
}

var _ node.ChurnerClient = (*ChurnerClient)(nil)

func (c *ChurnerClient) Churn(ctx context.Context, operatorAddress string, keyPair *core.KeyPair, quorumIDs []core.QuorumID) (*churnerpb.ChurnReply, error) {
	args := c.Called()
	var reply *churnerpb.ChurnReply
	if args.Get(0) != nil {
		reply = (args.Get(0)).(*churnerpb.ChurnReply)
	}

	var err error
	if args.Get(1) != nil {
		err = (args.Get(1)).(error)
	}
	return reply, err
}
