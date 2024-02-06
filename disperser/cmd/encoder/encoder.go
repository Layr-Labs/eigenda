package main

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
)

type EncoderGRPCServer struct {
	Server *encoder.Server
}

func NewEncoderGRPCServer(config Config, logger common.Logger) (*EncoderGRPCServer, error) {

	coreEncoder, err := encoding.NewEncoder(config.EncoderConfig, true)
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

	server := encoder.NewServer(*config.ServerConfig, logger, coreEncoder, metrics)

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
