package genenv

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
	"gopkg.in/yaml.v2"
)

const GraphUrl = "http://graph-node:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"

// GenerateDockerCompose all of the config for the test environment.
// Returns an object that corresponds to the participants of the
// current experiment.
func GenerateDockerCompose(lock *config.ConfigLock) {
	compose := NewDockerCompose()

	builder := DockerComposeBuilder{
		lock:    lock,
		compose: compose,
	}

	port := lock.Config.Services.BasePort

	builder.AddAWSInitScriptService()
	// TODO: Generate anvil-state.json using local anvil instance and forge script --broadcast to speed up local deploys.
	builder.AddAnvilService(8545, "")
	builder.AddEthInitScriptService()

	builder.AddGraphInitScriptService()

	// Generate churner service
	churnerURL := builder.AddChurnerService(port)
	port += 1

	// Generate disperser services
	for i := 0; i < lock.Config.Services.Counts.NumDis; i++ {
		builder.AddDisperserService(i, port)
		port += 2
	}

	// Generate operator services
	for i := 0; i < lock.Config.Services.Counts.NumOpr; i++ {
		builder.AddOperatorService(i, port, port+1, port+2, port+3, churnerURL)
		port += 5
	}

	builder.AddBatcherService()

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
	writeFile(composeFile, composeYaml)
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

type DockerComposeBuilder struct {
	lock    *config.ConfigLock
	compose *DockerCompose
}

func (e *DockerComposeBuilder) applyOverrides(c any, prefix, stub string, ind int) {

	pv := reflect.ValueOf(c)
	v := pv.Elem()

	prefix += "_"

	for key, value := range e.lock.Config.Services.Variables["globals"] {
		field := v.FieldByName(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

	for key, value := range e.lock.Config.Services.Variables[stub] {
		field := v.FieldByName(prefix + key)
		fmt.Println(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

	for key, value := range e.lock.Config.Services.Variables[fmt.Sprintf("%v%v", stub, ind)] {
		field := v.FieldByName(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

}

// Generates churner .env
func (e *DockerComposeBuilder) GenerateChurnerVars(ind int, graphUrl, grpcPort string) ChurnerVars {
	v := ChurnerVars{
		CHURNER_HOSTNAME:                    "",
		CHURNER_GRPC_PORT:                   grpcPort,
		CHURNER_BLS_OPERATOR_STATE_RETRIVER: e.lock.Config.EigenDA.OperatorStateRetreiver,
		CHURNER_EIGENDA_SERVICE_MANAGER:     e.lock.Config.EigenDA.ServiceManager,

		CHURNER_CHAIN_RPC:   "",
		CHURNER_PRIVATE_KEY: strings.TrimPrefix(e.lock.Pks.EcdsaMap[e.lock.Config.EigenDA.Deployer].PrivateKey, "0x"),

		CHURNER_GRAPH_URL:             graphUrl,
		CHURNER_INDEXER_PULL_INTERVAL: "1s",

		CHURNER_STD_LOG_LEVEL:  "debug",
		CHURNER_FILE_LOG_LEVEL: "trace",
		CHURNER_ENABLE_METRICS: "true",
	}

	e.applyOverrides(&v, "CHURNER", "churner", ind)
	v.CHURNER_HOSTNAME = "churner"

	return v
}

// Generates disperser .env
func (e *DockerComposeBuilder) GenerateDisperserVars(ind int, grpcPort string) DisperserVars {
	v := DisperserVars{
		DISPERSER_SERVER_S3_BUCKET_NAME:         "test-eigenda-blobstore",
		DISPERSER_SERVER_DYNAMODB_TABLE_NAME:    "test-BlobMetadata",
		DISPERSER_SERVER_RATE_BUCKET_TABLE_NAME: "test-BucketStore",
		DISPERSER_SERVER_GRPC_PORT:              grpcPort,
		DISPERSER_SERVER_ENABLE_METRICS:         "true",
		DISPERSER_SERVER_METRICS_HTTP_PORT:      "9093",
		DISPERSER_SERVER_CHAIN_RPC:              "",
		DISPERSER_SERVER_PRIVATE_KEY:            "123",

		DISPERSER_SERVER_REGISTERED_QUORUM_ID:       "0",
		DISPERSER_SERVER_TOTAL_UNAUTH_THROUGHPUT:    "10000000",
		DISPERSER_SERVER_PER_USER_UNAUTH_THROUGHPUT: "32000",
		DISPERSER_SERVER_ENABLE_RATELIMITER:         "true",
		DISPERSER_SERVER_BUCKET_STORE_SIZE:          "1000",

		DISPERSER_SERVER_BUCKET_SIZES:       "5s",
		DISPERSER_SERVER_BUCKET_MULTIPLIERS: "1",
		DISPERSER_SERVER_COUNT_FAILED:       "true",

		DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER: e.lock.Config.EigenDA.OperatorStateRetreiver,
		DISPERSER_SERVER_EIGENDA_SERVICE_MANAGER:     e.lock.Config.EigenDA.ServiceManager,
	}

	e.applyOverrides(&v, "DISPERSER_SERVER", "dis", ind)

	return v

}

// Generates batcher .env
func (e *DockerComposeBuilder) GenerateBatcherVars(ind int, key, graphUrl string) BatcherVars {
	v := BatcherVars{
		BATCHER_S3_BUCKET_NAME:              "test-eigenda-blobstore",
		BATCHER_DYNAMODB_TABLE_NAME:         "test-BlobMetadata",
		BATCHER_ENABLE_METRICS:              "true",
		BATCHER_METRICS_HTTP_PORT:           "9094",
		BATCHER_PULL_INTERVAL:               "5s",
		BATCHER_BLS_OPERATOR_STATE_RETRIVER: e.lock.Config.EigenDA.OperatorStateRetreiver,
		BATCHER_EIGENDA_SERVICE_MANAGER:     e.lock.Config.EigenDA.ServiceManager,
		BATCHER_SRS_ORDER:                   "300000",
		BATCHER_CHAIN_RPC:                   "",
		BATCHER_PRIVATE_KEY:                 key[2:],
		BATCHER_STD_LOG_LEVEL:               "debug",
		BATCHER_FILE_LOG_LEVEL:              "trace",
		BATCHER_GRAPH_URL:                   graphUrl,
		BATCHER_USE_GRAPH:                   "true",
		BATCHER_BATCH_SIZE_LIMIT:            "10240", // 10 GiB
		BATCHER_INDEXER_PULL_INTERVAL:       "1s",
		BATCHER_AWS_REGION:                  "",
		BATCHER_AWS_ACCESS_KEY_ID:           "",
		BATCHER_AWS_SECRET_ACCESS_KEY:       "",
		BATCHER_AWS_ENDPOINT_URL:            "",
		BATCHER_FINALIZER_INTERVAL:          "6m",
		BATCHER_ENCODING_REQUEST_QUEUE_SIZE: "500",
	}

	e.applyOverrides(&v, "BATCHER", "batcher", ind)

	return v
}

func (env *DockerComposeBuilder) GenerateEncoderVars(ind int, grpcPort string) EncoderVars {
	v := EncoderVars{
		DISPERSER_ENCODER_GRPC_PORT:               grpcPort,
		DISPERSER_ENCODER_ENABLE_METRICS:          "true",
		DISPERSER_ENCODER_G1_PATH:                 "",
		DISPERSER_ENCODER_G2_PATH:                 "",
		DISPERSER_ENCODER_SRS_ORDER:               "",
		DISPERSER_ENCODER_CACHE_PATH:              "",
		DISPERSER_ENCODER_VERBOSE:                 "",
		DISPERSER_ENCODER_NUM_WORKERS:             fmt.Sprint(runtime.GOMAXPROCS(0)),
		DISPERSER_ENCODER_MAX_CONCURRENT_REQUESTS: "16",
		DISPERSER_ENCODER_REQUEST_POOL_SIZE:       "32",
	}

	env.applyOverrides(&v, "DISPERSER_ENCODER", "enc", ind)

	return v
}

// Generates DA node .env
func (e *DockerComposeBuilder) GenerateOperatorVars(ind int, name, churnerUrl, dispersalPort, retrievalPort, metricsPort, nodeApiPort string) OperatorVars {
	blsKey := e.lock.Pks.BlsMap[name].PrivateKey
	ecdsaKey := e.lock.Pks.EcdsaMap[name].PrivateKey

	v := OperatorVars{
		NODE_HOSTNAME:                    "",
		NODE_DISPERSAL_PORT:              dispersalPort,
		NODE_RETRIEVAL_PORT:              retrievalPort,
		NODE_ENABLE_METRICS:              "true",
		NODE_METRICS_PORT:                metricsPort,
		NODE_ENABLE_NODE_API:             "true",
		NODE_API_PORT:                    nodeApiPort,
		NODE_TIMEOUT:                     "10",
		NODE_QUORUM_ID_LIST:              "0",
		NODE_DB_PATH:                     "/db",
		NODE_ENABLE_TEST_MODE:            "true",
		NODE_TEST_PRIVATE_BLS:            blsKey,
		NODE_PRIVATE_KEY:                 ecdsaKey,
		NODE_BLS_OPERATOR_STATE_RETRIVER: e.lock.Config.EigenDA.OperatorStateRetreiver,
		NODE_EIGENDA_SERVICE_MANAGER:     e.lock.Config.EigenDA.ServiceManager,
		NODE_REGISTER_AT_NODE_START:      "true",
		NODE_CHURNER_URL:                 churnerUrl,
		NODE_EXPIRATION_POLL_INTERVAL:    "10",
		NODE_G1_PATH:                     "",
		NODE_G2_PATH:                     "",
		NODE_CACHE_PATH:                  "",
		NODE_SRS_ORDER:                   "",
		NODE_NUM_WORKERS:                 fmt.Sprint(runtime.GOMAXPROCS(0)),
		NODE_VERBOSE:                     "true",
		NODE_CHAIN_RPC:                   "",
		NODE_STD_LOG_LEVEL:               "debug",
		NODE_FILE_LOG_LEVEL:              "trace",
		NODE_NUM_BATCH_VALIDATORS:        "128",
		NODE_PUBLIC_IP_PROVIDER:          "mockip",
		NODE_PUBLIC_IP_CHECK_INTERVAL:    "10s",
	}

	e.applyOverrides(&v, "NODE", "opr", ind)

	return v

}

// Generates retriever .env
func (e *DockerComposeBuilder) GenerateRetrieverVars(ind int, graphUrl, grpcPort string) RetrieverVars {
	v := RetrieverVars{
		RETRIEVER_HOSTNAME:                    "",
		RETRIEVER_GRPC_PORT:                   grpcPort,
		RETRIEVER_TIMEOUT:                     "10s",
		RETRIEVER_BLS_OPERATOR_STATE_RETRIVER: e.lock.Config.EigenDA.OperatorStateRetreiver,
		RETRIEVER_EIGENDA_SERVICE_MANAGER:     e.lock.Config.EigenDA.ServiceManager,
		RETRIEVER_NUM_CONNECTIONS:             "10",

		RETRIEVER_G1_PATH:             "",
		RETRIEVER_G2_PATH:             "",
		RETRIEVER_CACHE_PATH:          "",
		RETRIEVER_SRS_ORDER:           "",
		RETRIEVER_NUM_WORKERS:         fmt.Sprint(runtime.GOMAXPROCS(0)),
		RETRIEVER_VERBOSE:             "true",
		RETRIEVER_CACHE_ENCODED_BLOBS: "false",

		RETRIEVER_INDEXER_PULL_INTERVAL: "1s",

		RETRIEVER_STD_LOG_LEVEL:  "debug",
		RETRIEVER_FILE_LOG_LEVEL: "trace",
	}

	e.applyOverrides(&v, "RETRIEVER", "retriever", ind)

	return v
}

func (e *DockerComposeBuilder) AddEthInitScriptService() {
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
			"interval":     "5s",
			"timeout":      "5s",
			"retries":      12 * 20, // 20 minutes
			"start_period": "1s",
		},
		"ports": []string{"8080:8080"},
		"depends_on": map[string]interface{}{
			"anvil": map[string]string{
				"condition": "service_healthy",
			},
		},
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
			"interval":     "5s",
			"timeout":      "5s",
			"retries":      12 * 3, // 3 minutes
			"start_period": "5s",
		},
		"depends_on": map[string]interface{}{
			"eth-init-script": map[string]string{
				"condition": "service_healthy",
			},
		},
	}
}

func (e *DockerComposeBuilder) AddChurnerService(port int) (churnerURL string) {
	churnerConfig := e.GenerateChurnerVars(0, GraphUrl, fmt.Sprint(port))
	e.GenService("churner", churnerImage, churnerDockerfile, churnerConfig.getEnvMap(), []string{fmt.Sprint(port)})
	return fmt.Sprintf("%s:%s", churnerConfig.CHURNER_HOSTNAME, churnerConfig.CHURNER_GRPC_PORT)
}

func (e *DockerComposeBuilder) AddDisperserService(index, port int) {
	name := fmt.Sprintf("dis%v", index)
	disperserConfig := e.GenerateDisperserVars(index, fmt.Sprint(port))
	e.GenService(name, disImage, disDockerfile, disperserConfig.getEnvMap(), []string{fmt.Sprint(port)})
}

func (e *DockerComposeBuilder) AddOperatorService(index, dispersalPort, retrievalPort, metricsPort, nodeAPIPort int, churnerURL string) {
	name := fmt.Sprintf("opr%v", index)

	operatorConfig := e.GenerateOperatorVars(
		index,
		name,
		churnerURL,
		fmt.Sprint(dispersalPort),
		fmt.Sprint(retrievalPort),
		fmt.Sprint(metricsPort),
		fmt.Sprint(nodeAPIPort),
	)

	ports := []string{
		fmt.Sprintf("%v:%v", dispersalPort, dispersalPort),
		fmt.Sprintf("%v:%v", retrievalPort, retrievalPort),
		fmt.Sprintf("%v:%v", metricsPort, metricsPort),
		fmt.Sprintf("%v:%v", nodeAPIPort, nodeAPIPort),
	}

	e.compose.Services[name] = map[string]interface{}{
		// "image":    image,
		"environment": operatorConfig.getEnvMap(),
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

	nginxService := name + "_nginx"
	e.compose.Services[nginxService] = map[string]interface{}{
		"image": "nginx:latest",
		"environment": []string{
			"REQUEST_LIMIT=1r/s",
			"NODE_HOST=" + name,
		},
		"ports": ports,
		"volumes": []string{
			e.lock.RootPath + "/node/cmd/resources/rate-limit-nginx.conf:/etc/nginx/templates/default.conf.template:ro",
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

func (e *DockerComposeBuilder) AddBatcherService() {
	name := "batcher0"
	key, _ := e.lock.GetKey(name)
	batcherConfig := e.GenerateBatcherVars(0, key, GraphUrl)
	e.GenService(name, batcherImage, batcherDockerfile, batcherConfig.getEnvMap(), []string{})
}

func (e *DockerComposeBuilder) AddEncoderService() {
	// TODO: Add more encoders
	encoderConfig := e.GenerateEncoderVars(0, "34000")
	e.GenService("encoder0", encoderImage, encoderDockerfile, encoderConfig.getEnvMap(), []string{"34000"})
}

func (e *DockerComposeBuilder) AddRetrieverService(port int) {
	retrieverConfig := e.GenerateRetrieverVars(0, GraphUrl, fmt.Sprint(port))
	e.GenService("retriever0", retrieverImage, retrieverDockerfile, retrieverConfig.getEnvMap(), []string{fmt.Sprint(port)})
}

func (e *DockerComposeBuilder) AddAnvilService(port int, stateFile string) {
	cmd := fmt.Sprintf("anvil --host 0.0.0.0 --port %d", port)
	if stateFile != "" {
		cmd = fmt.Sprintf("%v --load-state %v", cmd, stateFile)
	}
	e.compose.Services["anvil"] = map[string]interface{}{
		"image":   "ghcr.io/foundry-rs/foundry:latest",
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
