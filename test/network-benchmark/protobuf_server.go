package network_benchmark

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/grpc/relay"
)

var _ TestServer = &protobufServer{}
var _ relay.ThroughputTestServer = &protobufServer{}

type protobufServer struct {
	relay.UnimplementedThroughputTestServer
	randomData *reusableRandomness
}

func NewProtobufServer() relay.ThroughputTestServer {
	return &protobufServer{}
}

func (s *protobufServer) SetRandomData(randomData *reusableRandomness) {
	s.randomData = randomData
}

func (s *protobufServer) GetData(
	ctx context.Context,
	request *relay.ThroughputTestRequest) (*relay.ThroughputTestResponse, error) {

	data := s.randomData.getData(request.Size, request.Seed)
	response := &relay.ThroughputTestResponse{
		Data: data,
	}

	return response, nil
}
