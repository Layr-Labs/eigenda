package network_benchmark

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/grpc/relay"
)

type protobufClient struct {
	client relay.ThroughputTestClient
}

func newProtobufClient(client relay.ThroughputTestClient) *protobufClient {
	return &protobufClient{
		client: client,
	}
}

// getData retrieves data from the server with the specified size and seed
func (c *protobufClient) getData(size int64, seed int64) ([]byte, error) {
	request := &relay.ThroughputTestRequest{
		Size: size,
		Seed: seed,
	}
	response, err := c.client.GetData(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
