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

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum/go-ethereum/log"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	invalidDomain = "invalid domain type"
)

const (
	GetRoute = "/get/"
	PutRoute = "/put/"

	DomainFilterKey = "domain"
)

type Server struct {
	log        log.Logger
	endpoint   string
	store      store.Store
	m          metrics.Metricer
	tls        *rpc.ServerTLSConfig
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(host string, port int, store store.Store, log log.Logger, m metrics.Metricer) *Server {
	endpoint := net.JoinHostPort(host, strconv.Itoa(port))
	return &Server{
		m:        m,
		log:      log,
		endpoint: endpoint,
		store:    store,
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

func WithLogging(handleFn func(http.ResponseWriter, *http.Request) error, log log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("request", "method", r.Method, "url", r.URL)
		err := handleFn(w, r)
		if err != nil {
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
	domain, err := ReadDomainFilter(r)
	if err != nil {
		svr.WriteBadRequest(w, invalidDomain)
		return err
	}

	key := path.Base(r.URL.Path)
	comm, err := eigenda.StringToCommit(key)
	if err != nil {
		svr.log.Info("failed to decode commitment", "err", err, "key", key)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	input, err := svr.store.Get(r.Context(), comm, domain)
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
	input, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var comm []byte
	if comm, err = svr.store.Put(r.Context(), input); err != nil {
		svr.WriteInternalError(w, err)
		return err
	}

	// write out encoded commitment
	svr.WriteResponse(w, eigenda.Commitment.Encode(comm))
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

func ReadDomainFilter(r *http.Request) (common.DomainType, error) {
	query := r.URL.Query()
	key := query.Get(DomainFilterKey)
	if key == "" { // default
		return common.BinaryDomain, nil
	}
	dt := common.StrToDomainType(key)
	if dt == common.UnknownDomain {
		return common.UnknownDomain, common.ErrInvalidDomainType
	}

	return dt, nil
}
