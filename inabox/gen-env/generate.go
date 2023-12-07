package genenv

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/inabox/config"
	"github.com/Layr-Labs/eigenda/inabox/utils"
	"gopkg.in/yaml.v2"
)

func GenerateConfigLock(rootPath, testName string) *config.ConfigLock {
	lock := &config.ConfigLock{}
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		log.Panicf("Error %s:", err.Error())
	}

	testPath := filepath.Join(rootPath, "inabox/testdata", testName)

	configPath := testPath + "/config.yaml"
	data := utils.MustReadFile(configPath)
	err = yaml.Unmarshal(data, &lock.Config)
	if err != nil {
		log.Panicf("Error %s:", err.Error())
	}

	lock.Pks = config.LoadPrivateKeys(&lock.Config, rootPath)

	if err != nil {
		log.Panicf("could not load private keys: %v", err)
	}

	lock.TestName = testName
	lock.Path = testPath
	lock.RootPath = rootPath

	port := lock.Config.Services.BasePort
	lock.Envs.Churner = GenerateChurnerVars(lock, 0, port)
	port += 1
	lock.Envs.Encoder = []config.EncoderVars{GenerateEncoderVars(lock, 0, 34000)}
	lock.Envs.Dispersers = make([]config.DisperserVars, lock.Config.Services.Counts.NumDis)
	for i := 0; i < lock.Config.Services.Counts.NumDis; i++ {
		lock.Envs.Dispersers[i] = GenerateDisperserVars(lock, i, port)
		port += 1
	}

	lock.Envs.Operators = make([]config.OperatorVars, lock.Config.Services.Counts.NumOpr)
	for i := 0; i < lock.Config.Services.Counts.NumOpr; i++ {
		dispersalPort := port
		retrievalPort := port + 1
		metricsPort := port + 2
		nodeApiPort := port + 3
		port += 4
		lock.Envs.Operators[i] = GenerateOperatorVars(lock, i, dispersalPort, retrievalPort, metricsPort, nodeApiPort)
	}
	lock.Envs.Batcher = GenerateBatcherVars(lock, 0)
	lock.Envs.Retriever = GenerateRetrieverVars(lock, 0, port)

	utils.MustWriteObjectToFile(filepath.Join(rootPath, "inabox", "testdata", testName, "config.lock"), lock)
	return lock
}

// GenerateDockerCompose all of the config for the test environment.
// Returns an object that corresponds to the participants of the
// current experiment.
func GenerateDockerCompose(lock *config.ConfigLock) {
	deployer, found := lock.Config.GetDeployer(lock.Config.EigenDA.Deployer)
	if !found {
		log.Panicf("Could not find deployer configuration for configured deployer name")
	}

	compose := NewDockerCompose()

	builder := DockerComposeBuilder{
		lock:    lock,
		compose: compose,
	}

	builder.AddAWSInitScriptService()
	if !deployer.LocalAnvil {
		builder.AddAnvilService(8545, "")
	}
	builder.AddEthInitScriptService(deployer.LocalAnvil)
	builder.AddGraphInitScriptService()

	builder.AddChurnerService(lock.Envs.Churner)
	for _, vars := range lock.Envs.Encoder {
		builder.AddEncoderService(vars)
	}

	for _, vars := range lock.Envs.Dispersers {
		builder.AddDisperserService(vars)
	}

	for _, vars := range lock.Envs.Operators {
		builder.AddOperatorService(vars)
	}

	builder.AddBatcherService(lock.Envs.Batcher)
	builder.AddRetrieverService(lock.Envs.Retriever)

	if lock.Config.Telemetry.IsNeeded {
		// run grafana agent
		// sd is required for accessing docker daemon
		// agent.yaml configures the grafana agent
		agentVolumes := append(
			lock.Config.Telemetry.DockerSd,
			lock.Config.Telemetry.ConfigPath+":/etc/agent/agent.yaml",
		)
		builder.AddTelemetryService("grafana-agent", "grafana/agent", agentVolumes)

		// run node exporter
		builder.AddTelemetryService("node-exporter", "prom/node-exporter", nil)
	}

	// Write to compose file
	composeYaml, err := yaml.Marshal(&compose)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	composeFile := lock.Path + "/docker-compose.gen.yml"
	utils.MustWriteFile(composeFile, composeYaml)
}

func CompileDockerCompose(rootPath, testName string) {
	filename := filepath.Join(rootPath, "inabox/testdata", testName, "docker-compose.yml")

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	// Ensure the file is closed at the end
	defer file.Close()

	buffer, err := utils.RunCommandAndCaptureOutput("docker-compose",
		"-f", filepath.Join(rootPath, "inabox/docker-compose.localstack.yml"),
		"-f", filepath.Join(rootPath, "inabox/docker-compose.thegraph.yml"),
		"-f", filepath.Join(rootPath, "inabox/testdata", testName, "docker-compose.gen.yml"),
		"config")
	if err != nil {
		log.Panicf("Could not compile docker-compose.yaml: %v", err)
	}

	// Write buffer to file
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
