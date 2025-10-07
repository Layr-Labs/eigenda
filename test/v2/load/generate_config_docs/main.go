package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/test/v2/load"
)

func main() {
	outputPath := filepath.Join("docs", "config", "load_generator_config.md")
	err := generateTrafficGeneratorConfigDocs(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating config docs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully generated config documentation at docs/config/load_generator_config.md")
}

// Generate config documentation for TrafficGeneratorConfig and save it to the specified destination.
func generateTrafficGeneratorConfigDocs(destination string) error {
	err := config.DocumentConfig(
		"Traffic Generator (v2)",
		load.DefaultTrafficGeneratorConfig,
		"TRAFFIC_GENERATOR",
		[]string{
			"github.com/Layr-Labs/eigenda/test/v2/client",
			"github.com/Layr-Labs/eigenda/test/v2/load",
		},
		destination,
		true,
	)

	if err != nil {
		return fmt.Errorf("failed to generate traffic generator config docs: %w", err)
	}

	return nil
}
