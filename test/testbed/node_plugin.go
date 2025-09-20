package testbed

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// OperatorConfig holds configuration for a node plugin operator
type OperatorConfig struct {
	NODE_HOSTNAME                    string
	NODE_DISPERSAL_PORT              string
	NODE_RETRIEVAL_PORT              string
	NODE_V2_DISPERSAL_PORT           string
	NODE_V2_RETRIEVAL_PORT           string
	NODE_ECDSA_KEY_FILE              string
	NODE_BLS_KEY_FILE                string
	NODE_ECDSA_KEY_PASSWORD          string
	NODE_BLS_KEY_PASSWORD            string
	NODE_QUORUM_ID_LIST              string
	NODE_CHAIN_RPC                   string
	NODE_EIGENDA_DIRECTORY           string
	NODE_BLS_OPERATOR_STATE_RETRIVER string
	NODE_EIGENDA_SERVICE_MANAGER     string
	NODE_CHURNER_URL                 string
}

// RunNodePlugin runs the node plugin directly using go run
func RunNodePlugin(ctx context.Context, operation string, operator OperatorConfig, logger logging.Logger) error {
	socket := string(core.MakeOperatorSocket(
		operator.NODE_HOSTNAME,
		operator.NODE_DISPERSAL_PORT,
		operator.NODE_RETRIEVAL_PORT,
		operator.NODE_V2_DISPERSAL_PORT,
		operator.NODE_V2_RETRIEVAL_PORT,
	))

	// Get the path to the node plugin cmd directory relative to this file
	_, filename, _, _ := runtime.Caller(0)
	testbedDir := filepath.Dir(filename)
	rootDir := filepath.Join(testbedDir, "..", "..")
	pluginCmdPath := filepath.Join(rootDir, "node", "plugin", "cmd")

	// Run the plugin directly with go run
	cmd := exec.CommandContext(ctx, "go", "run", pluginCmdPath,
		"--operation", operation,
		"--ecdsa-key-file", operator.NODE_ECDSA_KEY_FILE,
		"--bls-key-file", operator.NODE_BLS_KEY_FILE,
		"--ecdsa-key-password", operator.NODE_ECDSA_KEY_PASSWORD,
		"--bls-key-password", operator.NODE_BLS_KEY_PASSWORD,
		"--socket", socket,
		"--quorum-id-list", operator.NODE_QUORUM_ID_LIST,
		"--chain-rpc", operator.NODE_CHAIN_RPC,
		"--eigenda-directory", operator.NODE_EIGENDA_DIRECTORY,
		"--bls-operator-state-retriever", operator.NODE_BLS_OPERATOR_STATE_RETRIVER,
		"--eigenda-service-manager", operator.NODE_EIGENDA_SERVICE_MANAGER,
		"--churner-url", operator.NODE_CHURNER_URL,
		"--num-confirmations", "0",
	)

	logger.Info("Running node plugin", "operation", operation)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run node plugin: %w, output: %s", err, string(output))
	}

	logger.Info("Node plugin executed successfully", "output", string(output))
	return nil
}
