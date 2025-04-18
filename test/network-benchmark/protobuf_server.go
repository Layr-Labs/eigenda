package network_benchmark

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/api/grpc/relay"
)

var _ TestServer = &protobufServer{}
var _ relay.ThroughputTestServer = &protobufServer{}

type protobufServer struct {
	relay.UnimplementedThroughputTestServer
	randomData     *reusableRandomness
	requestsServed atomic.Uint64
}

func NewProtobufServer() TestServer {
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

	requestsServed := s.requestsServed.Add(1)
	if requestsServed%10000 == 0 {
		fmt.Printf("Requests served: %d\r", requestsServed)
	}

	return response, nil
}
