package rest

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

// CompatibilityConfig ... CompatibilityConfig stores values useful to external services for checking compatibility
// with the proxy instance, such as version, chainID, and recency window size. These values are returned by the rest
// servers /config endpoint.
type CompatibilityConfig struct {
	// Current proxy version in the format {MAJOR}.{MINOR}.{PATCH}-{META} e.g: 2.4.0-43-g3b4f9f40. The version
	// is injected at build using `git describe --tags --always --dirty`. This allows a service to perform a
	// minimum version supported check.
	Version string `json:"version"`
	// The ChainID of the connected ethClient. This allows a service to check which chain the proxy is connected
	// to. If the proxy has memstore enabled, a ChainID of "memstore" will be set.
	ChainID string `json:"chain_id"`
	// The EigenDA directory address. This allows a service to verify which contracts are being used by the proxy.
	DirectoryAddress string `json:"directory_address"`
	// The cert verifier router or immutable contract address. This allows a service to verify the cert verifier being
	// used by the proxy.
	CertVerifierAddress string `json:"cert_verifier_address"`
	// The max supported payload size in bytes supported by the proxy instance. Calculated from `MaxBlobSizeBytes`.
	MaxPayloadSizeBytes uint32 `json:"max_payload_size_bytes"`
	// The recency window size. This allows a service (e.g batch poster) to check alignment with the proxy instance.
	RecencyWindowSize uint64 `json:"recency_window_size"`
	// The APIs currently enabled on the rest server
	APIsEnabled []string `json:"apis_enabled"`
	// Whether the proxy is in read-only mode (no signer payment key)
	ReadOnlyMode bool `json:"read_only_mode"`
}

// Config ... Config for the proxy HTTP server
type Config struct {
	Host             string
	Port             int
	APIsEnabled      *enablement.RestApisEnabled
	CompatibilityCfg CompatibilityConfig
}

func (c *Config) BuildCompatibilityConfig(
	version string,
	chainID string,
	clientConfigV2 common.ClientConfigV2,
	readOnly bool,
) error {
	maxPayloadSize, err := codec.BlobSymbolsToMaxPayloadSize(
		uint32(clientConfigV2.MaxBlobSizeBytes / encoding.BYTES_PER_SYMBOL))
	if err != nil {
		return fmt.Errorf("calculate max payload size: %w", err)
	}

	// Remove 'v' prefix from version string if present for compatibility with eigenda/common/version helper funcs
	if len(version) > 0 {
		versionRunes := []rune(version)
		if versionRunes[0] == 'v' || versionRunes[0] == 'V' {
			version = string(versionRunes[1:])
		}
	}

	c.CompatibilityCfg = CompatibilityConfig{
		Version:             version,
		ChainID:             chainID,
		DirectoryAddress:    clientConfigV2.EigenDADirectory,
		CertVerifierAddress: clientConfigV2.EigenDACertVerifierOrRouterAddress,
		MaxPayloadSizeBytes: maxPayloadSize,
		RecencyWindowSize:   clientConfigV2.RBNRecencyWindowSize,
		APIsEnabled:         c.APIsEnabled.ToStringSlice(),
		ReadOnlyMode:        readOnly,
	}
	return nil
}

type Server struct {
	log        logging.Logger
	endpoint   string
	certMgr    store.IEigenDAManager
	keccakMgr  store.IKeccakManager
	m          metrics.Metricer
	httpServer *http.Server
	listener   net.Listener
	config     Config
}

func NewServer(
	cfg Config,
	certMgr store.IEigenDAManager,
	keccakMgr store.IKeccakManager,
	log logging.Logger,
	m metrics.Metricer,
) *Server {
	endpoint := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	return &Server{
		m:         m,
		log:       log,
		endpoint:  endpoint,
		certMgr:   certMgr,
		keccakMgr: keccakMgr,
		config:    cfg,
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

	svr.log.Info("Starting REST ALT DA server", "endpoint", svr.endpoint)
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
	svr.certMgr.SetDispersalBackend(backend)
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
