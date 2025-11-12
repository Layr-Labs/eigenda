package arbitrum_altda

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// determine which upstream example DA Server config fields are absolutely necessary.
// the example DA Provider server provided by OCL uses:
//   - String JWT auth
//
// JWT authorization is supported by the AnyTrust and example Custom DA servers
// and toggled in core nitro as a client config field
// TODO: Add support for JWT authentication
//
// The Custom DA server implementation is a thin wrapper over the existing proxy
// storage abstractions with lightweight translation from the existing critical
// REST status code signals (i.e, "drop cert", "failover") into arbitrum specific
// errors
type Config struct {
	Host string
	Port int
}

type Server struct {
	cfg      *Config
	svr      *http.Server
	listener net.Listener
}

// NewServer constructs the RPC server
func NewServer(ctx context.Context, cfg *Config, h IHandlers) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on tcp: %w", err)
	}

	rpcServer := rpc.NewServer()
	if err := rpcServer.RegisterName("daprovider", h); err != nil {
		return nil, fmt.Errorf("failed to register daprovider: %w", err)
	}

	rpcServer.SetHTTPBodyLimit(int(common.MaxServerPOSTRequestBodySize))

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return nil, errors.New("failed getting provider server address from listener")
	}

	svr := &http.Server{
		Addr:    "http://" + addr.String(),
		Handler: rpcServer,
	}

	return &Server{
		cfg:      cfg,
		svr:      svr,
		listener: listener,
	}, nil

}

// Port returns the port that the server is listening on.
// Useful in case Config.Port was set to 0 to let the OS assign a random port.
func (svr *Server) Port() int {
	// read from listener
	_, portStr, _ := net.SplitHostPort(svr.listener.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func (s *Server) Addr() string {
	return s.svr.Addr
}

// Start serves a tcp listener on an independent go routine
func (s *Server) Start() error {
	go func() {
		if err := s.svr.Serve(s.listener); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			println(fmt.Sprintf("provider server's Serve method returned a non http.ErrServerClosed error: %s", err.Error()))
		}
	}()

	return nil
}

// Stop is a shutdown function
func (s *Server) Stop() error {
	if err := s.svr.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
