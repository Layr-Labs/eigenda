package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gorilla/mux"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	Put = "put"

	CommitmentModeKey = "commitment_mode"
)

type Server struct {
	log        log.Logger
	endpoint   string
	sm         store.IManager
	m          metrics.Metricer
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(host string, port int, sm store.IManager, log log.Logger,
	m metrics.Metricer) *Server {
	endpoint := net.JoinHostPort(host, strconv.Itoa(port))
	return &Server{
		m:        m,
		log:      log,
		endpoint: endpoint,
		sm:       sm,
		httpServer: &http.Server{
			Addr:              endpoint,
			ReadHeaderTimeout: 10 * time.Second,
			// aligned with existing blob finalization times
			WriteTimeout: 40 * time.Minute,
		},
	}
}

func (svr *Server) Start() error {
	r := mux.NewRouter()
	svr.registerRoutes(r)
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

func parseVersionByte(w http.ResponseWriter, r *http.Request) (byte, error) {
	vars := mux.Vars(r)
	// only GET routes use gorilla parsed vars to separate header bytes from the raw commitment bytes.
	// POST routes parse them by hand because they neeed to send the entire
	// request (including the type/version header bytes) to the server.
	// TODO: perhaps for consistency we should also use gorilla vars for POST routes,
	// and then just reconstruct the full commitment in the handlers?
	versionByteHex, isGETRoute := vars[routingVarNameVersionByteHex]
	if !isGETRoute {
		// v0 is hardcoded in POST routes for now (see handlers.go that also have this hardcoded)
		// TODO: change this once we introduce v1/v2 certs
		return byte(commitments.CertV0), nil
	}
	versionByte, err := hex.DecodeString(versionByteHex)
	if err != nil {
		return 0, fmt.Errorf("decode version byte %s: %w", versionByteHex, err)
	}
	if len(versionByte) != 1 {
		return 0, fmt.Errorf("version byte is not a single byte: %s", versionByteHex)
	}
	switch versionByte[0] {
	case byte(commitments.CertV0):
		return versionByte[0], nil
	default:
		http.Error(w, fmt.Sprintf("unsupported version byte %x", versionByte), http.StatusBadRequest)
		return 0, fmt.Errorf("unsupported version byte %x", versionByte)
	}
}
