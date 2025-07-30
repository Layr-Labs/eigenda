package client

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	proxycommon "github.com/Layr-Labs/eigenda/api/proxy/common"
	proxyconfig "github.com/Layr-Labs/eigenda/api/proxy/config"
	proxymetrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/server"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

// ProxyWrapper starts an instance of the proxy in background goroutines, and then facilitates communication with it.
// This is intended to be used as a lightweight test utility, not as something that should be deployed outside of
// test settings.
type ProxyWrapper struct {
	proxyServer *server.Server
	client      *standard_client.Client
}

// Start a proxy in the background of this process (as opposed to the "normal" pattern of running a proxy in a
// separate process), and return a handle for communicating with the proxy.
func NewProxyWrapper(
	ctx context.Context,
	logger logging.Logger,
	proxyConfig *proxyconfig.AppConfig) (*ProxyWrapper, error) {

	err := proxyConfig.Check()
	if err != nil {
		return nil, fmt.Errorf("check proxy config: %w", err)
	}

	proxyMetrics := proxymetrics.NewMetrics("default")

	storeManager, err := builder.BuildStoreManager(
		ctx,
		logger,
		proxyMetrics,
		proxyConfig.StoreBuilderConfig,
		proxyConfig.SecretConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("build store manager: %w", err)
	}

	proxyServer := server.NewServer(proxyConfig.ServerConfig, storeManager, logger, proxyMetrics)

	router := mux.NewRouter()
	proxyServer.RegisterRoutes(router)
	proxyServer.SetDispersalBackend(proxycommon.V2EigenDABackend)
	err = proxyServer.Start(router)
	if err != nil {
		return nil, fmt.Errorf("start proxy server: %w", err)
	}

	client := standard_client.New(
		&standard_client.Config{
			URL: fmt.Sprintf("http://localhost:%d", proxyConfig.ServerConfig.Port),
		})

	return &ProxyWrapper{
		proxyServer: proxyServer,
		client:      client,
	}, nil
}

// Stop the proxy server gracefully.
func (w *ProxyWrapper) Stop() error {
	err := w.proxyServer.Stop()
	if err != nil {
		return fmt.Errorf("stop proxy server: %w", err)
	}

	return nil
}

// Disperse a payload to EigenDA. Returns a byte array representing the blob cert.
func (w *ProxyWrapper) SendPayload(ctx context.Context, payload []byte) ([]byte, error) {
	header, err := w.client.SetData(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("set data: %w", err)
	}
	return header, nil
}

// Fetch and verify a payload from EigenDA using the blob cert.
func (w *ProxyWrapper) GetPayload(ctx context.Context, cert []byte) ([]byte, error) {
	data, err := w.client.GetData(ctx, cert)
	if err != nil {
		return nil, fmt.Errorf("get data: %w", err)
	}
	return data, nil
}
