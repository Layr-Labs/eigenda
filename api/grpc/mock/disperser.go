package mock

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/api/grpc/disperser"

	"google.golang.org/grpc"
)

func MakeStreamMock(ctx context.Context) *StreamMock {
	return &StreamMock{
		ctx:            ctx,
		recvToServer:   make(chan *disperser.AuthenticatedRequest, 10),
		sentFromServer: make(chan *disperser.AuthenticatedReply, 10),
		closed:         false,
	}
}

type StreamMock struct {
	grpc.ServerStream
	ctx            context.Context
	recvToServer   chan *disperser.AuthenticatedRequest
	sentFromServer chan *disperser.AuthenticatedReply
	closed         bool
}

func (m *StreamMock) Context() context.Context {
	return m.ctx
}

func (m *StreamMock) Send(resp *disperser.AuthenticatedReply) error {
	m.sentFromServer <- resp
	return nil
}

func (m *StreamMock) Recv() (*disperser.AuthenticatedRequest, error) {
	req, more := <-m.recvToServer
	if !more {
		return nil, errors.New("empty")
	}
	return req, nil
}

func (m *StreamMock) SendFromClient(req *disperser.AuthenticatedRequest) error {
	if m.closed {
		return errors.New("closed")
	}
	m.recvToServer <- req
	return nil
}

func (m *StreamMock) RecvToClient() (*disperser.AuthenticatedReply, error) {
	response, more := <-m.sentFromServer
	if !more {
		return nil, errors.New("empty")
	}
	return response, nil
}

func (m *StreamMock) Close() {
	close(m.recvToServer)
	close(m.sentFromServer)
	m.closed = true
}
