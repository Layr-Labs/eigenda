package arbitrum_altda

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/node"

	"github.com/ethereum/go-ethereum/rpc"
)

// The ALT DA server implementation is a thin wrapper over the existing
// storage abstractions with lightweight translation from the existing critical
// REST status code signals (i.e, "drop cert", "failover") into arbitrum specific
// errors
type Config struct {
	Host      string
	Port      int
	JWTSecret string
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

	var handler http.Handler
	// go-ethereum puts specific constraints on JWT usage; ie:
	//     - HS256 is the only supported symmetric key schema
	//     - only signed claim for token payload is the IAT (issued at timestamp)
	//
	// see https://github.com/ethereum/go-ethereum/blob/v1.16.7/node/jwt_auth.go#L28-L45
	//
	// go-ethereum uses JWT for authenticated communication with consensus client where
	// the HS256 symmetric private key is copied between server domains. it's assumed
	// this is only used for local or enclosed service environments that aren't shared with open internet.
	//
	// for arbitrum, this is used for secure communication between rollup nodes and the
	// CustomDA server.
	if cfg.JWTSecret != "" {
		jwt, err := fetchJWTSecret(cfg.JWTSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch JWT secret: %w", err)
		}
		handler = node.NewHTTPHandlerStack(rpcServer, nil, nil, jwt)
	} else {
		handler = rpcServer
	}

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return nil, errors.New("failed getting provider server address from listener")
	}

	svr := &http.Server{
		Addr:    "http://" + addr.String(),
		Handler: handler,
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

// fetchJWTSecret processes a HS256 private key from a user provided text file
//
// this is a refactor of:
// https://github.com/OffchainLabs/nitro/blob/9eda1777a836c13916caac493ee1e2796c536afc/daprovider/server/provider_server.go#L76-L88
func fetchJWTSecret(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not read JWT Secret at file %s : %w", fileName, err)
	}

	jwtSecret := gethcommon.FromHex(strings.TrimSpace(string(data)))
	if length := len(jwtSecret); length != 32 {
		return nil, fmt.Errorf("invalid length detected for JWT token, expected 32 bytes but got %d", length)
	}

	return jwtSecret, nil
}
