package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda-proxy/config"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

type Server struct {
	log        logging.Logger
	endpoint   string
	sm         store.IManager
	m          metrics.Metricer
	httpServer *http.Server
	listener   net.Listener
	config     config.ServerConfig
}

func NewServer(
	cfg config.ServerConfig,
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

// BuildAndStartProxyServer constructs a new proxy server, and starts it
func BuildAndStartProxyServer(
	ctx context.Context,
	logger logging.Logger,
	metrics metrics.Metricer,
	appConfig config.AppConfig,
) (*Server, error) {
	storageManager, err := store.NewStorageManagerBuilder(
		ctx,
		logger,
		metrics,
		appConfig.EigenDAConfig.StorageConfig,
		appConfig.EigenDAConfig.MemstoreConfig,
		appConfig.EigenDAConfig.MemstoreEnabled,
		appConfig.EigenDAConfig.KzgConfig,
		appConfig.EigenDAConfig.ClientConfigV1,
		appConfig.EigenDAConfig.VerifierConfigV1,
		appConfig.EigenDAConfig.ClientConfigV2,
		appConfig.SecretConfig,
	).Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	proxyServer := NewServer(appConfig.EigenDAConfig.ServerConfig, storageManager, logger, metrics)
	router := mux.NewRouter()
	proxyServer.RegisterRoutes(router)
	if appConfig.EigenDAConfig.MemstoreEnabled {
		memconfig.NewHandlerHTTP(logger, appConfig.EigenDAConfig.MemstoreConfig).RegisterMemstoreConfigHandlers(router)
	}

	if err := proxyServer.Start(router); err != nil {
		return nil, fmt.Errorf("failed to start the DA server: %w", err)
	}

	logger.Info("Started EigenDA proxy server")

	return proxyServer, nil
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

func (svr *Server) writeResponse(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
	}
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
