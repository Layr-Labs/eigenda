package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	invalidDomain         = "invalid domain type"
	invalidCommitmentMode = "invalid commitment mode"

	GetRoute = "/get/"
	PutRoute = "/put/"

	DomainFilterKey   = "domain"
	CommitmentModeKey = "commitment_mode"
)

type Server struct {
	log        log.Logger
	endpoint   string
	router      *store.Router
	m          metrics.Metricer
	tls        *rpc.ServerTLSConfig
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(host string, port int, router *store.Router, log log.Logger, m metrics.Metricer) *Server {
	endpoint := net.JoinHostPort(host, strconv.Itoa(port))
	return &Server{
		m:        m,
		log:      log,
		endpoint: endpoint,
		router:    router,
		httpServer: &http.Server{
			Addr:              endpoint,
			ReadHeaderTimeout: 10 * time.Second,
			// aligned with existing blob finalization times
			WriteTimeout: 40 * time.Minute,
		},
	}
}

// WithMetrics is a middleware that records metrics for the route path.
func WithMetrics(handleFn func(http.ResponseWriter, *http.Request) error, m metrics.Metricer) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		recordDur := m.RecordRPCServerRequest(r.URL.Path)
		defer recordDur()

		return handleFn(w, r)
	}
}

// WithLogging is a middleware that logs the request method and URL.
func WithLogging(handleFn func(http.ResponseWriter, *http.Request) error, log log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("request", "method", r.Method, "url", r.URL)
		err := handleFn(w, r)
		if err != nil { // #nosec G104
			w.Write([]byte(err.Error()))
			log.Error(err.Error())
		}
	}
}

func (svr *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc(GetRoute, WithLogging(WithMetrics(svr.HandleGet, svr.m), svr.log))
	mux.HandleFunc(PutRoute, WithLogging(WithMetrics(svr.HandlePut, svr.m), svr.log))
	mux.HandleFunc("/health", WithLogging(svr.Health, svr.log))

	svr.httpServer.Handler = mux

	listener, err := net.Listen("tcp", svr.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	svr.listener = listener

	svr.endpoint = listener.Addr().String()

	svr.log.Info("Starting DA server", "endpoint", svr.endpoint)
	errCh := make(chan error, 1)
	go func() {
		if svr.tls != nil {
			if err := svr.httpServer.ServeTLS(svr.listener, "", ""); err != nil {
				errCh <- err
			}
		} else {
			if err := svr.httpServer.Serve(svr.listener); err != nil {
				errCh <- err
			}
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
func (svr *Server) Health(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func (svr *Server) HandleGet(w http.ResponseWriter, r *http.Request) error {

	ct, err := ReadCommitmentMode(r)
	if err != nil {
		svr.WriteBadRequest(w, invalidCommitmentMode)
		return err
	}
	key := path.Base(r.URL.Path)
	comm, err := commitments.StringToDecodedCommitment(key, ct)
	if err != nil {
		svr.log.Info("failed to decode commitment", "err", err, "commitment", comm)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	input, err := svr.router.Get(r.Context(), comm, ct)
	if err != nil && errors.Is(err, ErrNotFound) {
		svr.WriteNotFound(w, err.Error())
		return err
	}

	if err != nil {
		svr.WriteInternalError(w, err)
		return err
	}

	svr.WriteResponse(w, input)
	return nil
}

func (svr *Server) HandlePut(w http.ResponseWriter, r *http.Request) error {
	ct, err := ReadCommitmentMode(r)
	if err != nil {
		svr.WriteBadRequest(w, invalidCommitmentMode)
		return err
	}

	input, err := io.ReadAll(r.Body)
	if err != nil {
		svr.log.Error("Failed to read request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	key := path.Base(r.URL.Path)
	var commitment []byte

	if len(key) > 0 && key != "put" { // commitment key already provided (keccak256)
		comm, err := commitments.StringToDecodedCommitment(key, ct)
		if err != nil {
			svr.log.Info("failed to decode commitment", "err", err, "key", key)
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		commitment, err = svr.router.PutWithKey(context.Background(), comm, input)
		if err != nil {
			svr.log.Info("failed to put with key", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

	} else { // without
		commitment, err = svr.router.PutWithoutKey(context.Background(), input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}

	comm, err := commitments.EncodeCommitment(commitment, ct)
	if err != nil {
		svr.log.Info("failed to encode commitment", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	svr.log.Info(fmt.Sprintf("write commitment: %x\n", comm))
	// write out encoded commitment
	svr.WriteResponse(w, comm)
	return nil
}

func (svr *Server) WriteResponse(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
		svr.WriteInternalError(w, err)
	}
}

func (svr *Server) WriteInternalError(w http.ResponseWriter, err error) {
	svr.log.Error("internal server error", "err", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (svr *Server) WriteNotFound(w http.ResponseWriter, msg string) {
	svr.log.Info("not found", "msg", msg)
	w.WriteHeader(http.StatusNotFound)
}

func (svr *Server) WriteBadRequest(w http.ResponseWriter, msg string) {
	svr.log.Info("bad request", "msg", msg)
	w.WriteHeader(http.StatusBadRequest)
}

func (svr *Server) Port() int {
	// read from listener
	_, portStr, _ := net.SplitHostPort(svr.listener.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func ReadCommitmentMode(r *http.Request) (commitments.CommitmentMode, error) {
	query := r.URL.Query()
	key := query.Get(CommitmentModeKey)
	if key == "" { // default
		commit := path.Base(r.URL.Path)
		if len(commit) > 0 && commit != "put" { // provided commitment in request params
			decodedCommit, err := hexutil.Decode(commit)
			if err != nil {
				return commitments.SimpleCommitmentMode, err
			}

			switch decodedCommit[0] {
			case byte(commitments.GenericCommitmentType):
				return commitments.OptimismAltDA, nil

			case byte(commitments.Keccak256CommitmentType):
				return commitments.OptimismGeneric, nil

			default:
				return commitments.SimpleCommitmentMode, fmt.Errorf("unknown commit byte prefix")
	
			}
		}

		return commitments.OptimismAltDA, nil
	}

	return commitments.StringToCommitmentMode(key)
}

func (svr *Server) GetMemStats() *store.Stats {
	return svr.router.GetMemStore().Stats()
}

func (svr *Server) GetS3Stats() *store.Stats {
	return svr.router.GetS3Store().Stats()
}