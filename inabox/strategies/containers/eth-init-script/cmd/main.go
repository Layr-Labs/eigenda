package main

import (
	"log"
	"net"
)

func main() {
	// cfg := config.OpenCwdConfigLock()

	// log.Print("Deploying experiment...")

	// // TODO: Check if contracts have already been deployed
	// log.Print("Deploying EigenDA")
	// // get deployer
	// deployer, ok := cfg.Config.GetDeployer(cfg.Config.EigenDA.Deployer)
	// if !ok {
	// 	log.Panicf("Deployer improperly configured")
	// }
	// deployerKey, _ := cfg.GetKey(deployer.Name)
	// eigendaDeployConfig := ethinitscript.GenerateEigenDADeployConfig(cfg)
	// serviceManagerAddr := cfg.Config.EigenDA.ServiceManager
	// ethinitscript.DeployEigenDAContracts(deployerKey, deployer, eigendaDeployConfig, serviceManagerAddr)

	// log.Print("Test environment has succesfully deployed!")

	// Indicates that the script is done and other docker compose services can start
	startHealthcheckServer()
}

func startHealthcheckServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()
	log.Println("TCP server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		conn.Close() // Close the connection immediately for the health check purpose
	}
}
