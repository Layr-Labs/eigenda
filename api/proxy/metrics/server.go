package metrics

import (
	"net"
	"strconv"

	ophttp "github.com/ethereum-optimism/optimism/op-service/httputil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config ... Metrics server configuration
type Config struct {
	Host    string
	Port    int
	Enabled bool
}

func NewServer(registry *prometheus.Registry, cfg Config) *ophttp.HTTPServer {
	address := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	h := promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)

	return ophttp.NewHTTPServer(address, h)
}
