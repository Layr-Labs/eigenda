package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
	ethinitscript "github.com/Layr-Labs/eigenda/inabox/strategies/containers/eth-init-script"
)

func main() {
	cfg := config.OpenCwdConfigLock()

	cmd := startAnvil()

	log.Print("Deploying contracts...")
	deployer, ok := cfg.Config.GetDeployer(cfg.Config.EigenDA.Deployer)
	if !ok {
		log.Panicf("Deployer improperly configured")
	}
	deployer.RPC = "http://localhost:8545"
	deployerKey, _ := cfg.GetKey(deployer.Name)
	eigendaDeployConfig := ethinitscript.GenerateEigenDADeployConfig(cfg)
	serviceManagerAddr := cfg.Config.EigenDA.ServiceManager
	ethinitscript.DeployEigenDAContracts(deployerKey, deployer, eigendaDeployConfig, serviceManagerAddr)
	log.Print("Test environment succesfully deployed.")

	stopAnvil(cmd)
}

func startAnvil() *exec.Cmd {
	cmd := exec.Command("anvil", "--dump-state", "/data/anvil-state.json")

	// Set the output to the corresponding os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		log.Panicf("Error starting anvil: %v", err)
	}

	waitForAnvilToStart()

	return cmd
}

func waitForAnvilToStart() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Timeout reached, anvil still not available.")
			return
		case <-ticker.C:
			conn, err := net.Dial("tcp", "localhost:8545")
			if err != nil {
				log.Println("Anvil not available yet, checking again...")
			} else {
				conn.Close()
				log.Println("Anvil up.")
				return
			}
		}
	}
}

func stopAnvil(cmd *exec.Cmd) {
	// Send SIGINT
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		log.Panicf("Error while sending SIGINT to anvil: %v", err)
	}

	// Wait for the command to exit
	if err := cmd.Wait(); err != nil {
		log.Panicf("Error waiting for anvil to exit: %v", err)
	}

	log.Print("Anvil exited.")
}
