package pprof

import (
	"fmt"
	"net/http"
	pprofhttp "net/http/pprof"
	"os"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type PprofProfiler struct {
	logger   logging.Logger
	httpPort string
}

func NewPprofProfiler(httpPort string, logger logging.Logger) *PprofProfiler {
	return &PprofProfiler{
		logger:   logger.With("component", "PprofProfiler"),
		httpPort: httpPort,
	}
}

// Start the pprof server
func (p *PprofProfiler) Start() {
	host := os.Getenv("EIGENDA_PPROF_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%s", host, p.httpPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprofhttp.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprofhttp.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprofhttp.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprofhttp.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprofhttp.Trace)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		p.logger.Error("pprof server failed", "error", err, "pprofAddr", addr)
	}
}
