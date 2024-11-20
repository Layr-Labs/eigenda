package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/flags"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
)

func StartProxySvr(cliCtx *cli.Context) error {
	log := oplog.NewLogger(oplog.AppOut(cliCtx), oplog.ReadCLIConfig(cliCtx)).New("role", "eigenda_proxy")
	oplog.SetGlobalLogHandler(log.Handler())
	log.Info("Starting EigenDA Proxy Server", "version", Version, "date", Date, "commit", Commit)

	cfg := server.ReadCLIConfig(cliCtx)
	if err := cfg.Check(); err != nil {
		return err
	}
	err := prettyPrintConfig(cliCtx, log)
	if err != nil {
		return fmt.Errorf("failed to pretty print config: %w", err)
	}

	m := metrics.NewMetrics("default")

	ctx, ctxCancel := context.WithCancel(cliCtx.Context)
	defer ctxCancel()

	sm, err := server.LoadStoreManager(ctx, cfg, log, m)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	server := server.NewServer(cliCtx.String(flags.ListenAddrFlagName), cliCtx.Int(flags.PortFlagName), sm, log, m)

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start the DA server: %w", err)
	}

	log.Info("Started EigenDA proxy server")

	defer func() {
		if err := server.Stop(); err != nil {
			log.Error("failed to stop DA server", "err", err)
		}

		log.Info("successfully shutdown API server")
	}()

	if cfg.MetricsCfg.Enabled {
		log.Debug("starting metrics server", "addr", cfg.MetricsCfg.ListenAddr, "port", cfg.MetricsCfg.ListenPort)
		svr, err := m.StartServer(cfg.MetricsCfg.ListenAddr, cfg.MetricsCfg.ListenPort)
		if err != nil {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		defer func() {
			if err := svr.Stop(context.Background()); err != nil {
				log.Error("failed to stop metrics server", "err", err)
			}
		}()
		log.Info("started metrics server", "addr", svr.Addr())
		m.RecordUp()
	}

	return ctxinterrupt.Wait(cliCtx.Context)
}

// TODO: we should probably just change EdaClientConfig struct definition in eigenda-client
func prettyPrintConfig(cliCtx *cli.Context, log log.Logger) error {
	// we read a new config which we modify to hide private info in order to log the rest
	cfg := server.ReadCLIConfig(cliCtx)
	if cfg.EigenDAConfig.EdaClientConfig.SignerPrivateKeyHex != "" {
		cfg.EigenDAConfig.EdaClientConfig.SignerPrivateKeyHex = "*****" // marshaling defined in client config
	}
	if cfg.EigenDAConfig.EdaClientConfig.EthRpcUrl != "" {
		cfg.EigenDAConfig.EdaClientConfig.EthRpcUrl = "*****" // hiding as RPC providers typically use sensitive API keys within
	}

	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	log.Info(fmt.Sprintf("Initializing EigenDA proxy server with config (\"*****\" fields are hidden): %v", string(configJSON)))
	return nil
}
