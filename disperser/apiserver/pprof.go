package apiserver

import (
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type PprofConfig struct {
	HTTPPort    string
	EnablePprof bool
}

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
func (p *PprofProfiler) Start(port string, logger logging.Logger) {
	pprofAddr := fmt.Sprintf("%s:%s", disperser.Localhost, port)
	mux := http.NewServeMux()

	if err := http.ListenAndServe(pprofAddr, mux); err != nil {
		p.logger.Error("pprof server failed", "error", err)
	}
}
