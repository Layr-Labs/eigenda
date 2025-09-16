package arbitrum_altda

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

// determine which upstream example DA Server config fields are absolutely necessary.
// the example DA Provider server provided by OCL uses:
//   - String JWT auth
//   - Write toggling
//   - Server request body limits
//
// JWT authorization is supported by the AnyTrust and example Custom DA servers
// and toggled in core nitro as a client config field
// TODO: Add support for JWT authentication
//
// Write toggling allows a user to use the ALT DA server in "read only" mode
// which is typical behavior for all rollup operator nodes outside of the batch poster
// TODO: Add support for "read only" mode
//
// Server request body limits are used to impose a maximum on the allowed request body size
// TODO: Understand if this is actually used/respected in DA client anywhere
// TODO: Add env ingestion for these values
// TODO: Determine a proper default given EigenDA's large batch uniqueness. The defaultBodyLimit
// //    is currently set to 5 mib which is too low for integration EigenDA
//
// The ALT DA server implementation should be a thin wrapper over the existing
// storage abstractions with lightweight translation from the existing critical
// REST status code signals (i.e, "drop cert", "failover") into arbitrum specific
// errors
type Config struct {
	Enable bool
	Host   string
	Port   int
}

type Server struct {
	cfg      *Config
	svr      *http.Server
	listener net.Listener
}

// NewServer constructs the RPC server
func NewServer(ctx context.Context, cfg *Config, h *Handlers) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on tcp: %w", err)
	}

	rpcServer := rpc.NewServer()
	if err := rpcServer.RegisterName("daprovider", h); err != nil {
		return nil, fmt.Errorf("failed to register daprovider: %w", err)
	}

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
