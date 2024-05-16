package plasma

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

	"github.com/Layr-Labs/op-plasma-eigenda/metrics"
	"github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

type PlasmaStore interface {
	// Get retrieves the given key if it's present in the key-value data store.
	Get(ctx context.Context, key []byte) ([]byte, error)
	// Put inserts the given value into the key-value data store.
	PutWithComm(ctx context.Context, key []byte, value []byte) error
	// Put inserts the given value into the key-value data store.
	PutWithoutComm(ctx context.Context, value []byte) (key []byte, err error)
}

type DAServer struct {
	log        log.Logger
	endpoint   string
	store      PlasmaStore
	m          metrics.Metricer
	tls        *rpc.ServerTLSConfig
	httpServer *http.Server
	listener   net.Listener
}

func NewDAServer(host string, port int, store PlasmaStore, log log.Logger, m metrics.Metricer) *DAServer {
	endpoint := net.JoinHostPort(host, strconv.Itoa(port))
	return &DAServer{
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

func (d *DAServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/get/", d.HandleGet)
	mux.HandleFunc("/put/", d.HandlePut)
	mux.HandleFunc("/health", d.Health)

	d.httpServer.Handler = mux

	listener, err := net.Listen("tcp", d.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	d.listener = listener

	d.endpoint = listener.Addr().String()

	d.log.Info("Starting DA server on", d.endpoint)
	errCh := make(chan error, 1)
	go func() {
		if d.tls != nil {
			if err := d.httpServer.ServeTLS(d.listener, "", ""); err != nil {
				errCh <- err
			}
		} else {
			if err := d.httpServer.Serve(d.listener); err != nil {
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

func (d *DAServer) Health(w http.ResponseWriter, r *http.Request) {
	d.log.Info("GET", "url", r.URL)
	recordDur := d.m.RecordRPCServerRequest("health")
	defer recordDur()

	w.WriteHeader(http.StatusOK)
}

func (d *DAServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	d.log.Info("GET", "url", r.URL)
	recordDur := d.m.RecordRPCServerRequest("put")
	defer recordDur()

	route := path.Dir(r.URL.Path)
	if route != "/get" {
		d.log.Info("invalid route", "route", route)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := path.Base(r.URL.Path)
	comm, err := hexutil.Decode(key)
	if err != nil {
		d.log.Info("failed to decode commitment bytes from hex", "err", err, "key", key)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decodedComm, err := DecodeEigenDACommitment(comm)
	if err != nil {
		d.log.Info("failed to decode commitment", "err", err, "key", key)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input, err := d.store.Get(r.Context(), decodedComm)
	if err != nil && errors.Is(err, ErrNotFound) {
		d.log.Info("key not found", "key", key)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		d.log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(input); err != nil {
		d.log.Error("failed to write response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *DAServer) HandlePut(w http.ResponseWriter, r *http.Request) {
	d.log.Info("PUT", "url", r.URL)
	recordDur := d.m.RecordRPCServerRequest("put")
	defer recordDur()

	route := path.Dir(r.URL.Path)
	if route != "/put" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var comm []byte
	if r.URL.Path == "/put" || r.URL.Path == "/put/" { // without commitment
		if comm, err = d.store.PutWithoutComm(r.Context(), input); err != nil {
			d.log.Error("Failed to store commitment to the DA server", "err", err, "comm", comm)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else { // with commitment (might be worth deleting if we never expect a commitment to be passed in the URL for this server type)
		key := path.Base(r.URL.Path)
		comm, err = hexutil.Decode(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := d.store.PutWithComm(r.Context(), comm, input); err != nil {
			d.log.Error("Failed to store commitment to the DA server", "err", err, "key", key)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// write out encoded commitment
	if _, err := w.Write(EigenDACommitment.Encode(comm)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (b *DAServer) Endpoint() string {
	return b.listener.Addr().String()
}

func (b *DAServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := b.httpServer.Shutdown(ctx); err != nil {
		b.log.Error("Failed to shutdown DA server", "err", err)
		return err
	}
	return nil
}
