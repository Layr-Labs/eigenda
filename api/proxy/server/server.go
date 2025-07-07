package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

// Config ... Config for the proxy HTTP server
type Config struct {
	Host string
	Port int
	// EnabledAPIs contains the list of API types that are enabled.
	// When empty (default), no special API endpoints are registered.
	// Example: If it contains "admin", administrative endpoints like
	// /admin/eigenda-dispersal-backend will be available.
	EnabledAPIs []string
}

// IsAPIEnabled checks if a specific API type is enabled
func (c *Config) IsAPIEnabled(apiType string) bool {
	return slices.Contains(c.EnabledAPIs, apiType)
}

type Server struct {
	log        logging.Logger
	endpoint   string
	sm         store.IManager
	m          metrics.Metricer
	httpServer *http.Server
	listener   net.Listener
	config     Config
}

func NewServer(
	cfg Config,
	sm store.IManager,
	log logging.Logger,
	m metrics.Metricer,
) *Server {
	endpoint := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	return &Server{
		m:        m,
		log:      log,
		endpoint: endpoint,
		sm:       sm,
		config:   cfg,
		httpServer: &http.Server{
			Addr:              endpoint,
			ReadHeaderTimeout: 10 * time.Second,
			// aligned with existing blob finalization times
			WriteTimeout: 40 * time.Minute,
		},
	}
}

func (svr *Server) Start(r *mux.Router) error {
	svr.httpServer.Handler = r

	listener, err := net.Listen("tcp", svr.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	svr.listener = listener

	svr.endpoint = listener.Addr().String()

	svr.log.Info("Starting DA server", "endpoint", svr.endpoint)
	errCh := make(chan error, 1)
	go func() {
		if err := svr.httpServer.Serve(svr.listener); err != nil {
			errCh <- err
		}
	}()

	// verify that the server comes up
	tick := time.NewTimer(10 * time.Millisecond)
	defer tick.Stop()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server failed: %w", err)
	case <-tick.C:
		return nil
	}
}

func (svr *Server) Endpoint() string {
	return svr.listener.Addr().String()
}

func (svr *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.httpServer.Shutdown(ctx); err != nil {
		svr.log.Error("Failed to shutdown proxy server", "err", err)
		return err
	}
	return nil
}

// SetDispersalBackend configures which version of eigenDA the server disperses to
func (svr *Server) SetDispersalBackend(backend common.EigenDABackend) {
	svr.sm.SetDispersalBackend(backend)
}

func (svr *Server) Port() int {
	// read from listener
	_, portStr, _ := net.SplitHostPort(svr.listener.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func parseCertVersion(w http.ResponseWriter, r *http.Request) (certs.VersionByte, error) {
	vars := mux.Vars(r)
	// only GET routes use gorilla parsed vars to separate header bytes from the raw commitment bytes.
	// POST routes parse them by hand because they neeed to send the entire
	// request (including the type/version header bytes) to the server.
	// TODO: perhaps for consistency we should also use gorilla vars for POST routes,
	// and then just reconstruct the full commitment in the handlers?
	versionByteHex, isGETRoute := vars[routingVarNameVersionByteHex]
	if !isGETRoute {
		// TODO: this seems like a bug... used in metrics for POST route, so we'll just always return v0??
		return certs.V0VersionByte, nil
	}
	versionByte, err := hex.DecodeString(versionByteHex)
	if err != nil {
		return 0, fmt.Errorf("decode version byte %s: %w", versionByteHex, err)
	}
	if len(versionByte) != 1 {
		return 0, fmt.Errorf("version byte is not a single byte: %s", versionByteHex)
	}
	certVersion, err := certs.ByteToVersion(versionByte[0])
	if err != nil {
		errWithHexContext := fmt.Errorf("unsupported version byte %x: %w", versionByte, err)
		http.Error(w, errWithHexContext.Error(), http.StatusBadRequest)
		return 0, errWithHexContext
	}
	return certVersion, nil
}
