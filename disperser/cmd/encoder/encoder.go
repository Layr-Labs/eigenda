package main

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type EncoderGRPCServer struct {
	Server *encoder.Server
}

func NewEncoderGRPCServer(config Config, _logger logging.Logger) (*EncoderGRPCServer, error) {
	logger := _logger.With("component", "EncoderGRPCServer")

	// TODO: If encoder is v1, we need to load g2 points, otherwise we don't
	p, err := prover.NewProver(&config.EncoderConfig, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}

	metrics := encoder.NewMetrics(config.MetricsConfig.HTTPPort, logger)
	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Encoder", "socket", httpSocket)
	}

	server := encoder.NewServer(*config.ServerConfig, logger, p, metrics)

	return &EncoderGRPCServer{
		Server: server,
	}, nil
}

func (d *EncoderGRPCServer) Start(ctx context.Context) error {
	// TODO: Start Metrics
	return d.Server.Start()
}

func (d *EncoderGRPCServer) Close() {
	d.Server.Close()
}
