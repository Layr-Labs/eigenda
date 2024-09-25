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
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	GetRoute = "/get/"
	PutRoute = "/put/"
	Put      = "put"

	CommitmentModeKey = "commitment_mode"
)

type Server struct {
	log        log.Logger
	endpoint   string
	router     store.IRouter
	m          metrics.Metricer
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(host string, port int, router store.IRouter, log log.Logger,
	m metrics.Metricer) *Server {
	endpoint := net.JoinHostPort(host, strconv.Itoa(port))
	return &Server{
		m:        m,
		log:      log,
		endpoint: endpoint,
		router:   router,
		httpServer: &http.Server{
			Addr:              endpoint,
			ReadHeaderTimeout: 10 * time.Second,
			// aligned with existing blob finalization times
			WriteTimeout: 40 * time.Minute,
		},
	}
}

// WithMetrics is a middleware that records metrics for the route path.
func WithMetrics(
	handleFn func(http.ResponseWriter, *http.Request) (commitments.CommitmentMeta, error),
	m metrics.Metricer,
) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		recordDur := m.RecordRPCServerRequest(r.Method)

		meta, err := handleFn(w, r)
		if err != nil {
			var metaErr MetaError
			if errors.As(err, &metaErr) {
				recordDur(w.Header().Get("status"), string(metaErr.Meta.Mode), string(metaErr.Meta.CertVersion))
			} else {
				recordDur(w.Header().Get("status"), string("NoCommitmentMode"), string("NoCertVersion"))
			}
			return err
		}
		// we assume that every route will set the status header
		recordDur(w.Header().Get("status"), string(meta.Mode), string(meta.CertVersion))
		return nil
	}
}

// WithLogging is a middleware that logs the request method and URL.
func WithLogging(
	handleFn func(http.ResponseWriter, *http.Request) error,
	log log.Logger,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("request", "method", r.Method, "url", r.URL)
		err := handleFn(w, r)
		if err != nil { // #nosec G104
			w.Write([]byte(err.Error())) //nolint:errcheck // ignore error
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
func (svr *Server) Health(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

// HandleGet handles the GET request for commitments.
// Note: even when an error is returned, the commitment meta is still returned,
// because it is needed for metrics (see the WithMetrics middleware).
// TODO: we should change this behavior and instead use a custom error that contains the commitment meta.
func (svr *Server) HandleGet(w http.ResponseWriter, r *http.Request) (commitments.CommitmentMeta, error) {
	meta, err := ReadCommitmentMeta(r)
	if err != nil {
		err = fmt.Errorf("invalid commitment mode: %w", err)
		svr.WriteBadRequest(w, err)
		return commitments.CommitmentMeta{}, err
	}
	key := path.Base(r.URL.Path)
	comm, err := commitments.StringToDecodedCommitment(key, meta.Mode)
	if err != nil {
		err = fmt.Errorf("failed to decode commitment from key %v (commitment mode %v): %w", key, meta.Mode, err)
		svr.WriteBadRequest(w, err)
		return commitments.CommitmentMeta{}, MetaError{
			Err:  err,
			Meta: meta,
		}
	}

	input, err := svr.router.Get(r.Context(), comm, meta.Mode)
	if err != nil {
		err = fmt.Errorf("get request failed with commitment %v (commitment mode %v): %w", comm, meta.Mode, err)
		if errors.Is(err, ErrNotFound) {
			svr.WriteNotFound(w, err)
		} else {
			svr.WriteInternalError(w, err)
		}
		return commitments.CommitmentMeta{}, MetaError{
			Err:  err,
			Meta: meta,
		}
	}

	svr.WriteResponse(w, input)
	return meta, nil
}

// HandlePut handles the PUT request for commitments.
// Note: even when an error is returned, the commitment meta is still returned,
// because it is needed for metrics (see the WithMetrics middleware).
// TODO: we should change this behavior and instead use a custom error that contains the commitment meta.
func (svr *Server) HandlePut(w http.ResponseWriter, r *http.Request) (commitments.CommitmentMeta, error) {
	meta, err := ReadCommitmentMeta(r)
	if err != nil {
		err = fmt.Errorf("invalid commitment mode: %w", err)
		svr.WriteBadRequest(w, err)
		return commitments.CommitmentMeta{}, err
	}
	// ReadCommitmentMeta function invoked inside HandlePut will not return a valid certVersion
	// Current simple fix is using the hardcoded default value of 0 (also the only supported value)
	//TODO: smarter decode needed when there's more than one version
	meta.CertVersion = byte(commitments.CertV0)

	input, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)
		svr.WriteBadRequest(w, err)
		return commitments.CommitmentMeta{}, MetaError{
			Err:  err,
			Meta: meta,
		}
	}

	key := path.Base(r.URL.Path)
	var comm []byte

	if len(key) > 0 && key != Put { // commitment key already provided (keccak256)
		comm, err = commitments.StringToDecodedCommitment(key, meta.Mode)
		if err != nil {
			err = fmt.Errorf("failed to decode commitment from key %v (commitment mode %v): %w", key, meta.Mode, err)
			svr.WriteBadRequest(w, err)
			return commitments.CommitmentMeta{}, MetaError{
				Err:  err,
				Meta: meta,
			}
		}
	}

	commitment, err := svr.router.Put(r.Context(), meta.Mode, comm, input)
	if err != nil {
		err = fmt.Errorf("put request failed with commitment %v (commitment mode %v): %w", comm, meta.Mode, err)

		if errors.Is(err, store.ErrEigenDAOversizedBlob) || errors.Is(err, store.ErrProxyOversizedBlob) {
			// we add here any error that should be returned as a 400 instead of a 500.
			// currently only includes oversized blob requests
			svr.WriteBadRequest(w, err)
			return meta, err
		}

		svr.WriteInternalError(w, err)
		return commitments.CommitmentMeta{}, MetaError{
			Err:  err,
			Meta: meta,
		}
	}

	responseCommit, err := commitments.EncodeCommitment(commitment, meta.Mode)
	if err != nil {
		err = fmt.Errorf("failed to encode commitment %v (commitment mode %v): %w", commitment, meta.Mode, err)
		svr.WriteInternalError(w, err)
		return commitments.CommitmentMeta{}, MetaError{
			Err:  err,
			Meta: meta,
		}
	}

	svr.log.Info(fmt.Sprintf("response commitment: %x\n", responseCommit))
	// write out encoded commitment
	svr.WriteResponse(w, responseCommit)
	return meta, nil
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

func (svr *Server) WriteNotFound(w http.ResponseWriter, err error) {
	svr.log.Info("not found", "err", err)
	w.WriteHeader(http.StatusNotFound)
}

func (svr *Server) WriteBadRequest(w http.ResponseWriter, err error) {
	svr.log.Info("bad request", "err", err)
	w.WriteHeader(http.StatusBadRequest)
}

func (svr *Server) Port() int {
	// read from listener
	_, portStr, _ := net.SplitHostPort(svr.listener.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

// Read both commitment mode and version
func ReadCommitmentMeta(r *http.Request) (commitments.CommitmentMeta, error) {
	// label requests with commitment mode and version
	ct, err := ReadCommitmentMode(r)
	if err != nil {
		return commitments.CommitmentMeta{}, err
	}
	if ct == "" {
		return commitments.CommitmentMeta{}, fmt.Errorf("commitment mode is empty")
	}
	cv, err := ReadCommitmentVersion(r, ct)
	if err != nil {
		// default to version 0
		return commitments.CommitmentMeta{Mode: ct, CertVersion: cv}, err
	}
	return commitments.CommitmentMeta{Mode: ct, CertVersion: cv}, nil
}

func ReadCommitmentMode(r *http.Request) (commitments.CommitmentMode, error) {
	query := r.URL.Query()
	key := query.Get(CommitmentModeKey)
	if key != "" {
		return commitments.StringToCommitmentMode(key)
	}

	commit := path.Base(r.URL.Path)
	if len(commit) > 0 && commit != Put { // provided commitment in request params (op keccak256)
		if !strings.HasPrefix(commit, "0x") {
			commit = "0x" + commit
		}

		decodedCommit, err := hexutil.Decode(commit)
		if err != nil {
			return "", err
		}

		if len(decodedCommit) < 3 {
			return "", fmt.Errorf("commitment is too short")
		}

		switch decodedCommit[0] {
		case byte(commitments.GenericCommitmentType):
			return commitments.OptimismGeneric, nil

		case byte(commitments.Keccak256CommitmentType):
			return commitments.OptimismKeccak, nil

		default:
			return commitments.SimpleCommitmentMode, fmt.Errorf("unknown commit byte prefix")
		}
	}
	return commitments.OptimismGeneric, nil
}

func ReadCommitmentVersion(r *http.Request, mode commitments.CommitmentMode) (byte, error) {
	commit := path.Base(r.URL.Path)
	if len(commit) > 0 && commit != Put { // provided commitment in request params (op keccak256)
		if !strings.HasPrefix(commit, "0x") {
			commit = "0x" + commit
		}

		decodedCommit, err := hexutil.Decode(commit)
		if err != nil {
			return 0, err
		}

		if len(decodedCommit) < 3 {
			return 0, fmt.Errorf("commitment is too short")
		}

		if mode == commitments.OptimismGeneric || mode == commitments.SimpleCommitmentMode {
			return decodedCommit[2], nil
		}

		return decodedCommit[0], nil
	}
	return 0, nil
}

func (svr *Server) GetEigenDAStats() *store.Stats {
	return svr.router.GetEigenDAStore().Stats()
}

func (svr *Server) GetS3Stats() *store.Stats {
	return svr.router.GetS3Store().Stats()
}

func (svr *Server) GetStoreStats(bt store.BackendType) (*store.Stats, error) {
	// first check if the store is a cache
	for _, cache := range svr.router.Caches() {
		if cache.BackendType() == bt {
			return cache.Stats(), nil
		}
	}

	// then check if the store is a fallback
	for _, fallback := range svr.router.Fallbacks() {
		if fallback.BackendType() == bt {
			return fallback.Stats(), nil
		}
	}

	return nil, fmt.Errorf("store not found")
}
