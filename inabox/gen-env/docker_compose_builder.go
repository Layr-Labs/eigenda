package genenv

import (
	"fmt"
	"log"

	"github.com/Layr-Labs/eigenda/inabox/config"
)

type DockerComposeBuilder struct {
	lock    *config.ConfigLock
	compose *DockerCompose
}

func (e *DockerComposeBuilder) AddTelemetryService(name, image string, volumes []string) {
	e.compose.Services[name] = map[string]interface{}{
		"image":  image,
		"volume": volumes,
	}
}

// Used to generate a docker compose file corresponding to the test environment
func (e *DockerComposeBuilder) GenService(name, image, dockerfilePath string, env map[string]string, ports []string) {

	for i, port := range ports {
		ports[i] = port + ":" + port
	}

	e.compose.Services[name] = map[string]interface{}{
		// "image":    image,
		"environment": env,
		"ports":       ports,
		"volumes": []string{
			e.lock.Path + ":/data",
			e.lock.RootPath + "/inabox/secrets:/secrets",
			e.lock.RootPath + "/inabox/resources:/resources",
		},
		"build": Build{
			Context:    dockerBuildContext,
			Dockerfile: dockerfilePath,
		},
		"extra_hosts": []string{
			"host.docker.internal:host-gateway",
		},
		"networks": []string{
			"eigenda-demo",
		},
		"depends_on": map[string]interface{}{
			"eth-init-script": map[string]interface{}{
				"condition": "service_healthy",
			},
			"graph-init-script": map[string]interface{}{
				"condition": "service_healthy",
			},
			"aws-init-script": map[string]interface{}{
				"condition": "service_healthy",
			},
		},
	}
}

func (e *DockerComposeBuilder) AddEthInitScriptService(localAnvil bool) {
	// Generate eth-init-script
	e.compose.Services["eth-init-script"] = map[string]interface{}{
		"volumes": []string{
			e.lock.Path + ":/data",
			e.lock.RootPath + "/contracts:/contracts",
		},
		"build": Build{
			Context:    dockerBuildContext,
			Dockerfile: ethInitDockerfile,
		},
		"extra_hosts": []string{
			"host.docker.internal:host-gateway",
		},
		"networks": []string{
			"eigenda-demo",
		},
		"healthcheck": map[string]interface{}{
			"test": []string{
				"CMD", "nc", "-z", "localhost", "8080",
			},
			"interval":     "1s",
			"timeout":      "5s",
			"retries":      60 * 20, // 20 minutes
			"start_period": "1s",
		},
		"ports": []string{"8080:8080"},
	}
	if !localAnvil {
		e.compose.Services["eth-init-script"]["depends_on"] = map[string]interface{}{
			"anvil": map[string]string{
				"condition": "service_healthy",
			},
		}
	} else {
		portsInterface := e.compose.Services["eth-init-script"]["ports"]
		ports, ok := portsInterface.([]string)
		if !ok {
			// Handle the error if the type assertion fails
			log.Fatal("Type assertion failed: 'ports' is not of type []string")
		}
		ports = append(ports, "8545:8545")
		e.compose.Services["eth-init-script"]["ports"] = ports
	}
}

func (e *DockerComposeBuilder) AddGraphInitScriptService() {
	// Generate graph-init-script
	e.compose.Services["graph-init-script"] = map[string]interface{}{
		"build": Build{
			Context:    dockerBuildContext,
			Dockerfile: graphInitDockerfile,
		},
		"volumes": []string{
			e.lock.Path + ":/data",
		},
		"extra_hosts": []string{
			"host.docker.internal:host-gateway",
		},
		"networks": []string{
			"eigenda-demo",
		},
		"healthcheck": map[string]interface{}{
			"test": []string{
				"CMD", "curl", "-f", "http://localhost:8080/health",
			},
			"interval":     "1s",
			"timeout":      "5s",
			"retries":      60 * 3, // 3 minutes
			"start_period": "5s",
		},
		"depends_on": map[string]interface{}{
			"graph-node": map[string]string{
				"condition": "service_healthy",
			},
		},
	}
}

func (e *DockerComposeBuilder) AddChurnerService(vars config.ChurnerVars) {
	e.GenService(vars.CHURNER_NAME, churnerImage, churnerDockerfile, vars.GetEnvMap(), []string{fmt.Sprint(vars.CHURNER_GRPC_PORT)})
	// return fmt.Sprintf("%s:%s", vars.CHURNER_HOSTNAME, vars.CHURNER_GRPC_PORT)
}

func (e *DockerComposeBuilder) AddDisperserService(vars config.DisperserVars) {
	// name := fmt.Sprintf("dis%v", index)
	// disperserConfig := e.GenerateDisperserVars(index, fmt.Sprint(port))
	e.GenService(vars.DISPERSER_NAME, disImage, disDockerfile, vars.GetEnvMap(), []string{fmt.Sprint(vars.DISPERSER_SERVER_GRPC_PORT)})
}

func (e *DockerComposeBuilder) AddOperatorService(vars config.OperatorVars) {
	// index, dispersalPort, retrievalPort, metricsPort, nodeAPIPort int, churnerURL string
	// name := fmt.Sprintf("opr%v", index)

	// operatorConfig := e.GenerateOperatorVars(
	// 	index,
	// 	name,
	// 	churnerURL,
	// 	fmt.Sprint(dispersalPort),
	// 	fmt.Sprint(retrievalPort),
	// 	fmt.Sprint(metricsPort),
	// 	fmt.Sprint(nodeAPIPort),
	// )

	// ports := []string{
	// 	fmt.Sprintf("%v:%v", dispersalPort, dispersalPort),
	// 	fmt.Sprintf("%v:%v", retrievalPort, retrievalPort),
	// 	fmt.Sprintf("%v:%v", metricsPort, metricsPort),
	// 	fmt.Sprintf("%v:%v", nodeAPIPort, nodeAPIPort),
	// }

	ports := []string{
		fmt.Sprintf("%v:%v", vars.NODE_DISPERSAL_PORT, vars.NODE_DISPERSAL_PORT),
		fmt.Sprintf("%v:%v", vars.NODE_RETRIEVAL_PORT, vars.NODE_RETRIEVAL_PORT),
		fmt.Sprintf("%v:%v", vars.NODE_METRICS_PORT, vars.NODE_METRICS_PORT),
		fmt.Sprintf("%v:%v", vars.NODE_API_PORT, vars.NODE_API_PORT),
	}

	e.compose.Services[vars.NODE_NAME] = map[string]interface{}{
		// "image":    image,
		"environment": vars.GetEnvMap(),
		"volumes": []string{
			e.lock.Path + ":/data",
			e.lock.RootPath + "/inabox/secrets:/secrets",
			e.lock.RootPath + "/inabox/resources:/resources",
		},
		"build": Build{
			Context:    dockerBuildContext,
			Dockerfile: nodeDockerfile,
		},
		"extra_hosts": []string{
			"host.docker.internal:host-gateway",
		},
		"networks": []string{
			"eigenda-demo",
		},
		"ports": ports,
		"depends_on": map[string]interface{}{
			"eth-init-script": map[string]interface{}{
				"condition": "service_healthy",
			},
			"graph-init-script": map[string]interface{}{
				"condition": "service_healthy",
			},
			"aws-init-script": map[string]interface{}{
				"condition": "service_healthy",
			},
		},
	}
}

func (e *DockerComposeBuilder) AddBatcherService(vars config.BatcherVars) {
	e.GenService(vars.BATCHER_NAME, batcherImage, batcherDockerfile, vars.GetEnvMap(), []string{})
}

func (e *DockerComposeBuilder) AddEncoderService(vars config.EncoderVars) {
	e.GenService(vars.DISPERSER_ENCODER_NAME, encoderImage, encoderDockerfile, vars.GetEnvMap(), []string{vars.DISPERSER_ENCODER_GRPC_PORT})
}

func (e *DockerComposeBuilder) AddRetrieverService(vars config.RetrieverVars) {
	// retrieverConfig := e.GenerateRetrieverVars(0, fmt.Sprint(port))
	e.GenService(vars.RETRIEVER_NAME, retrieverImage, retrieverDockerfile, vars.GetEnvMap(), []string{fmt.Sprint(vars.RETRIEVER_GRPC_PORT)})
}

func (e *DockerComposeBuilder) AddAnvilService(port int, stateFile string) {
	cmd := fmt.Sprintf("anvil --host 0.0.0.0 --port %d", port)
	if stateFile != "" {
		cmd = fmt.Sprintf("%v --load-state %v", cmd, stateFile)
	}
	e.compose.Services["anvil"] = map[string]interface{}{
		"image":   "ghcr.io/foundry-rs/foundry:nightly-f5c91995f80b5bf3b4c29c934d414cc198c9e7a8",
		"command": []string{cmd},
		"ports":   []string{fmt.Sprintf("%v:%v", port, port)},
		"volumes": []string{
			e.lock.Path + ":/data",
			e.lock.RootPath + "/inabox/secrets:/secrets",
			e.lock.RootPath + "/inabox/resources:/resources",
		},
		"networks": []string{
			"eigenda-demo",
		},
		"healthcheck": map[string]interface{}{
			"test":         []string{"CMD", "nc", "-z", "localhost", "8545"},
			"interval":     "2s",
			"timeout":      "5s",
			"retries":      10,
			"start_period": "10s",
		},
	}
}

func (e *DockerComposeBuilder) AddAWSInitScriptService() {
	// Generate eth-init-script
	e.compose.Services["aws-init-script"] = map[string]interface{}{
		"build": Build{
			Context:    dockerBuildContext,
			Dockerfile: awsInitDockerfile,
		},
		"networks": []string{
			"eigenda-demo",
		},
		"healthcheck": map[string]interface{}{
			"test": []string{
				"CMD", "nc", "-z", "localhost", "8080",
			},
			"interval":     "5s",
			"timeout":      "5s",
			"retries":      12 * 1, // 1 minute
			"start_period": "1s",
		},
		"depends_on": map[string]interface{}{
			"localstack": map[string]string{
				"condition": "service_healthy",
			},
		},
	}
}
