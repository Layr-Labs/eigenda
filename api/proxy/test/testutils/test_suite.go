package testutils

import (
	"context"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/api/proxy/config"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/server"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

// TestSuite contains necessary objects, to be able to execute a proxy test
type TestSuite struct {
	Ctx     context.Context
	Log     logging.Logger
	Metrics *proxy_metrics.EmulatedMetricer
	Server  *server.Server
}

// TestSuiteWithLogger returns a function which overrides the logger for a TestSuite
func TestSuiteWithLogger(log logging.Logger) func(*TestSuite) {
	return func(ts *TestSuite) {
		ts.Log = log
	}
}

// CreateTestSuite creates a test suite.
//
// It accepts parameters indicating which type of Backend to use, and a test config.
// It also accepts a variadic options parameter, which contains functions that operate on a TestSuite object.
// These options allow for configuration control over the TestSuite.
func CreateTestSuite(
	appConfig config.AppConfig,
	options ...func(*TestSuite),
) (TestSuite, func()) {
	ts := &TestSuite{
		Ctx:     context.Background(),
		Log:     logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{}),
		Metrics: proxy_metrics.NewEmulatedMetricer(),
	}
	// Override the defaults with the provided options, if present.
	for _, option := range options {
		option(ts)
	}

	ctx, logger, metrics := ts.Ctx, ts.Log, ts.Metrics

	if err := appConfig.Check(); err != nil {
		panic(err)
	}
	// Commenting out because it clutters the log outputs in CI too much.
	// We should prob take in a *testing.T and use t.Logf instead, so that logs
	// only appear when the test fails.
	// configString, err := appConfig.StoreBuilderConfig.ToString()
	// if err != nil {
	// 	panic(fmt.Sprintf("convert config json to string: %v", err))
	// }
	//
	// logger.Infof(
	// 	"Creating EigenDA proxy server for testSuite with config (\"*****\" fields are hidden): %v",
	// 	configString,
	// )

	storeManager, err := builder.BuildStoreManager(
		ctx,
		logger,
		metrics,
		appConfig.StoreBuilderConfig,
		appConfig.SecretConfig,
	)
	if err != nil {
		panic(fmt.Sprintf("build storage manager: %v", err.Error()))
	}

	proxyServer := server.NewServer(appConfig.ServerConfig, storeManager, logger, metrics)
	router := mux.NewRouter()
	proxyServer.RegisterRoutes(router)
	if appConfig.StoreBuilderConfig.MemstoreEnabled {
		memconfig.NewHandlerHTTP(logger, appConfig.StoreBuilderConfig.MemstoreConfig).
			RegisterMemstoreConfigHandlers(router)
	}

	if err := proxyServer.Start(router); err != nil {
		panic(fmt.Sprintf("start proxy server: %v", err.Error()))
	}

	kill := func() {
		if err := proxyServer.Stop(); err != nil {
			logger.Error("failed to stop proxy server", "err", err)
		}
	}

	return TestSuite{
		Ctx:     ctx,
		Log:     logger,
		Metrics: metrics,
		Server:  proxyServer,
	}, kill
}

func (ts *TestSuite) Address() string {
	// read port from listener
	port := ts.Server.Port()

	return fmt.Sprintf("%s://%s:%d", transport, host, port)
}
