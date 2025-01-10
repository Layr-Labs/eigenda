package deploy

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

func (env *Config) GetDeployer(name string) (*ContractDeployer, bool) {

	for _, deployer := range env.Deployers {
		if deployer.Name == name {
			return deployer, true
		}
	}
	return nil, false
}

// Constructs a mapping between service names/deployer names (e.g., 'dis0', 'opr1') and private keys. Order of priority: Map, List, File
func (env *Config) loadPrivateKeys() error {

	// construct full list of names
	// nTotal := env.Services.Counts.NumDis + env.Services.Counts.NumOpr + env.Services.Counts.NumRet + env.Services.Counts.NumSeq + env.Services.Counts.NumCha
	// names := make([]string, len(env.Deployers)+nTotal)
	names := make([]string, 0)
	for _, d := range env.Deployers {
		names = append(names, d.Name)
	}
	addNames := func(prefix string, num int) {
		for i := 0; i < num; i++ {
			names = append(names, fmt.Sprintf("%v%v", prefix, i))
		}
	}
	addNames("dis", 2)
	addNames("opr", env.Services.Counts.NumOpr)
	addNames("staker", env.Services.Counts.NumOpr)
	addNames("retriever", 1)
	addNames("relay", env.Services.Counts.NumRelays)

	log.Println("service names:", names)

	// Collect private keys from file
	keyPath := "secrets"

	// Read ECDSA private keys
	fileData := readFile(filepath.Join(keyPath, "ecdsa_keys/private_key_hex.txt"))
	ecdsaPks := strings.Split(string(fileData), "\n")
	// Read ECDSA passwords
	fileData = readFile(filepath.Join(keyPath, "ecdsa_keys/password.txt"))
	ecdsaPwds := strings.Split(string(fileData), "\n")
	// Read BLS private keys
	fileData = readFile(filepath.Join(keyPath, "bls_keys/private_key_hex.txt"))
	blsPks := strings.Split(string(fileData), "\n")
	// Read BLS passwords
	fileData = readFile(filepath.Join(keyPath, "bls_keys/password.txt"))
	blsPwds := strings.Split(string(fileData), "\n")

	if len(ecdsaPks) != len(blsPks) || len(blsPks) != len(ecdsaPwds) || len(ecdsaPwds) != len(blsPwds) {
		return errors.New("the number of keys and passwords for ECDSA and BLS must be the same")
	}

	// Add missing items to map
	if env.Pks.EcdsaMap == nil {
		env.Pks.EcdsaMap = make(map[string]KeyInfo)
	}
	if env.Pks.BlsMap == nil {
		env.Pks.BlsMap = make(map[string]KeyInfo)
	}

	ind := 0
	for _, name := range names {
		_, exists := env.Pks.EcdsaMap[name]
		if !exists {

			if ind >= len(ecdsaPks) {
				return errors.New("not enough pks")
			}

			env.Pks.EcdsaMap[name] = KeyInfo{
				PrivateKey: ecdsaPks[ind],
				Password:   ecdsaPwds[ind],
				KeyFile:    fmt.Sprintf("%s/ecdsa_keys/keys/%v.ecdsa.key.json", keyPath, ind+1),
			}
			env.Pks.BlsMap[name] = KeyInfo{
				PrivateKey: blsPks[ind],
				Password:   blsPwds[ind],
				KeyFile:    fmt.Sprintf("%s/bls_keys/keys/%v.bls.key.json", keyPath, ind+1),
			}

			ind++
		}
	}

	return nil
}

func (env *Config) applyDefaults(c any, prefix, stub string, ind int) {

	pv := reflect.ValueOf(c)
	v := pv.Elem()

	prefix += "_"

	for key, value := range env.Services.Variables["globals"] {
		field := v.FieldByName(prefix + key)
		if field.IsValid() && field.CanSet() && field.String() == "" {
			field.SetString(value)
		}
	}

	for key, value := range env.Services.Variables[stub] {
		field := v.FieldByName(prefix + key)
		fmt.Println(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

	for key, value := range env.Services.Variables[fmt.Sprintf("%v%v", stub, ind)] {
		field := v.FieldByName(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

}

// Generates churner .env
func (env *Config) generateChurnerVars(ind int, graphUrl, logPath, grpcPort string) ChurnerVars {
	v := ChurnerVars{
		CHURNER_LOG_FORMAT:                  "text",
		CHURNER_HOSTNAME:                    "",
		CHURNER_GRPC_PORT:                   grpcPort,
		CHURNER_BLS_OPERATOR_STATE_RETRIVER: env.EigenDA.OperatorStateRetreiver,
		CHURNER_EIGENDA_SERVICE_MANAGER:     env.EigenDA.ServiceManager,

		CHURNER_CHAIN_RPC:   "",
		CHURNER_PRIVATE_KEY: strings.TrimPrefix(env.Pks.EcdsaMap[env.EigenDA.Deployer].PrivateKey, "0x"),

		CHURNER_GRAPH_URL:             graphUrl,
		CHURNER_INDEXER_PULL_INTERVAL: "1s",

		CHURNER_ENABLE_METRICS:          "true",
		CHURNER_METRICS_HTTP_PORT:       "9095",
		CHURNER_CHURN_APPROVAL_INTERVAL: "900s",
	}

	env.applyDefaults(&v, "CHURNER", "churner", ind)

	return v
}

// Generates disperser .env
func (env *Config) generateDisperserVars(ind int, logPath, dbPath, grpcPort string) DisperserVars {
	v := DisperserVars{
		DISPERSER_SERVER_LOG_FORMAT:             "text",
		DISPERSER_SERVER_S3_BUCKET_NAME:         "test-eigenda-blobstore",
		DISPERSER_SERVER_DYNAMODB_TABLE_NAME:    "test-BlobMetadata",
		DISPERSER_SERVER_RATE_BUCKET_TABLE_NAME: "",
		DISPERSER_SERVER_RATE_BUCKET_STORE_SIZE: "100000",
		DISPERSER_SERVER_GRPC_PORT:              grpcPort,
		DISPERSER_SERVER_ENABLE_METRICS:         "true",
		DISPERSER_SERVER_METRICS_HTTP_PORT:      "9093",
		DISPERSER_SERVER_CHAIN_RPC:              "",
		DISPERSER_SERVER_PRIVATE_KEY:            "123",
		DISPERSER_SERVER_NUM_CONFIRMATIONS:      "0",

		DISPERSER_SERVER_REGISTERED_QUORUM_ID:      "0,1",
		DISPERSER_SERVER_TOTAL_UNAUTH_BYTE_RATE:    "10000000,10000000",
		DISPERSER_SERVER_PER_USER_UNAUTH_BYTE_RATE: "32000,32000",
		DISPERSER_SERVER_TOTAL_UNAUTH_BLOB_RATE:    "10,10",
		DISPERSER_SERVER_PER_USER_UNAUTH_BLOB_RATE: "2,2",
		DISPERSER_SERVER_ENABLE_RATELIMITER:        "true",

		DISPERSER_SERVER_RETRIEVAL_BLOB_RATE: "4",
		DISPERSER_SERVER_RETRIEVAL_BYTE_RATE: "10000000",

		DISPERSER_SERVER_BUCKET_SIZES:       "5s",
		DISPERSER_SERVER_BUCKET_MULTIPLIERS: "1",
		DISPERSER_SERVER_COUNT_FAILED:       "true",

		DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER: env.EigenDA.OperatorStateRetreiver,
		DISPERSER_SERVER_EIGENDA_SERVICE_MANAGER:     env.EigenDA.ServiceManager,
	}

	env.applyDefaults(&v, "DISPERSER_SERVER", "dis", ind)

	return v

}

func (env *Config) generateDisperserV2Vars(ind int, logPath, dbPath, grpcPort string) DisperserVars {
	v := DisperserVars{
		DISPERSER_SERVER_LOG_FORMAT:             "text",
		DISPERSER_SERVER_S3_BUCKET_NAME:         "test-eigenda-blobstore",
		DISPERSER_SERVER_DYNAMODB_TABLE_NAME:    "test-BlobMetadata-v2",
		DISPERSER_SERVER_RATE_BUCKET_TABLE_NAME: "",
		DISPERSER_SERVER_RATE_BUCKET_STORE_SIZE: "100000",
		DISPERSER_SERVER_GRPC_PORT:              grpcPort,
		DISPERSER_SERVER_ENABLE_METRICS:         "true",
		DISPERSER_SERVER_METRICS_HTTP_PORT:      "9093",
		DISPERSER_SERVER_CHAIN_RPC:              "",
		DISPERSER_SERVER_PRIVATE_KEY:            "123",
		DISPERSER_SERVER_NUM_CONFIRMATIONS:      "0",

		DISPERSER_SERVER_REGISTERED_QUORUM_ID:      "0,1",
		DISPERSER_SERVER_TOTAL_UNAUTH_BYTE_RATE:    "10000000,10000000",
		DISPERSER_SERVER_PER_USER_UNAUTH_BYTE_RATE: "32000,32000",
		DISPERSER_SERVER_TOTAL_UNAUTH_BLOB_RATE:    "10,10",
		DISPERSER_SERVER_PER_USER_UNAUTH_BLOB_RATE: "2,2",
		DISPERSER_SERVER_ENABLE_RATELIMITER:        "true",

		DISPERSER_SERVER_RETRIEVAL_BLOB_RATE: "4",
		DISPERSER_SERVER_RETRIEVAL_BYTE_RATE: "10000000",

		DISPERSER_SERVER_BUCKET_SIZES:       "5s",
		DISPERSER_SERVER_BUCKET_MULTIPLIERS: "1",
		DISPERSER_SERVER_COUNT_FAILED:       "true",

		DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER: env.EigenDA.OperatorStateRetreiver,
		DISPERSER_SERVER_EIGENDA_SERVICE_MANAGER:     env.EigenDA.ServiceManager,
		DISPERSER_SERVER_DISPERSER_VERSION:           "2",

		DISPERSER_SERVER_ENABLE_PAYMENT_METERER:  "true",
		DISPERSER_SERVER_RESERVATIONS_TABLE_NAME: "e2e_v2_reservation",
		DISPERSER_SERVER_ON_DEMAND_TABLE_NAME:    "e2e_v2_ondemand",
		DISPERSER_SERVER_GLOBAL_RATE_TABLE_NAME:  "e2e_v2_global_reservation",
	}

	env.applyDefaults(&v, "DISPERSER_SERVER", "dis", ind)

	return v
}

// Generates batcher .env
func (env *Config) generateBatcherVars(ind int, key, graphUrl, logPath string) BatcherVars {
	v := BatcherVars{
		BATCHER_LOG_FORMAT:                    "text",
		BATCHER_S3_BUCKET_NAME:                "test-eigenda-blobstore",
		BATCHER_DYNAMODB_TABLE_NAME:           "test-BlobMetadata",
		BATCHER_ENABLE_METRICS:                "true",
		BATCHER_METRICS_HTTP_PORT:             "9094",
		BATCHER_PULL_INTERVAL:                 "5s",
		BATCHER_BLS_OPERATOR_STATE_RETRIVER:   env.EigenDA.OperatorStateRetreiver,
		BATCHER_EIGENDA_SERVICE_MANAGER:       env.EigenDA.ServiceManager,
		BATCHER_SRS_ORDER:                     "300000",
		BATCHER_CHAIN_RPC:                     "",
		BATCHER_PRIVATE_KEY:                   key[2:],
		BATCHER_GRAPH_URL:                     graphUrl,
		BATCHER_USE_GRAPH:                     "true",
		BATCHER_BATCH_SIZE_LIMIT:              "10240", // 10 GiB
		BATCHER_INDEXER_PULL_INTERVAL:         "1s",
		BATCHER_AWS_REGION:                    "",
		BATCHER_AWS_ACCESS_KEY_ID:             "",
		BATCHER_AWS_SECRET_ACCESS_KEY:         "",
		BATCHER_AWS_ENDPOINT_URL:              "",
		BATCHER_FINALIZER_INTERVAL:            "6m",
		BATCHER_ENCODING_REQUEST_QUEUE_SIZE:   "500",
		BATCHER_NUM_CONFIRMATIONS:             "0",
		BATCHER_MAX_BLOBS_TO_FETCH_FROM_STORE: "100",
		BATCHER_FINALIZATION_BLOCK_DELAY:      "0",
		BATCHER_KMS_KEY_DISABLE:               "true",
	}

	env.applyDefaults(&v, "BATCHER", "batcher", ind)

	return v
}

func (env *Config) generateEncoderVars(ind int, grpcPort string) EncoderVars {
	v := EncoderVars{
		DISPERSER_ENCODER_LOG_FORMAT:              "text",
		DISPERSER_ENCODER_AWS_REGION:              "",
		DISPERSER_ENCODER_AWS_ACCESS_KEY_ID:       "",
		DISPERSER_ENCODER_AWS_SECRET_ACCESS_KEY:   "",
		DISPERSER_ENCODER_AWS_ENDPOINT_URL:        "",
		DISPERSER_ENCODER_GRPC_PORT:               grpcPort,
		DISPERSER_ENCODER_ENABLE_METRICS:          "true",
		DISPERSER_ENCODER_G1_PATH:                 "",
		DISPERSER_ENCODER_G2_PATH:                 "",
		DISPERSER_ENCODER_SRS_ORDER:               "",
		DISPERSER_ENCODER_SRS_LOAD:                "",
		DISPERSER_ENCODER_CACHE_PATH:              "",
		DISPERSER_ENCODER_VERBOSE:                 "",
		DISPERSER_ENCODER_NUM_WORKERS:             fmt.Sprint(runtime.GOMAXPROCS(0)),
		DISPERSER_ENCODER_MAX_CONCURRENT_REQUESTS: "16",
		DISPERSER_ENCODER_REQUEST_POOL_SIZE:       "32",
		DISPERSER_ENCODER_REQUEST_QUEUE_SIZE:      "32",
	}

	env.applyDefaults(&v, "DISPERSER_ENCODER", "enc", ind)

	return v
}

func (env *Config) generateEncoderV2Vars(ind int, grpcPort string) EncoderVars {
	v := EncoderVars{
		DISPERSER_ENCODER_LOG_FORMAT:              "text",
		DISPERSER_ENCODER_AWS_REGION:              "",
		DISPERSER_ENCODER_AWS_ACCESS_KEY_ID:       "",
		DISPERSER_ENCODER_AWS_SECRET_ACCESS_KEY:   "",
		DISPERSER_ENCODER_AWS_ENDPOINT_URL:        "",
		DISPERSER_ENCODER_GRPC_PORT:               grpcPort,
		DISPERSER_ENCODER_ENABLE_METRICS:          "true",
		DISPERSER_ENCODER_G1_PATH:                 "",
		DISPERSER_ENCODER_G2_PATH:                 "",
		DISPERSER_ENCODER_SRS_ORDER:               "",
		DISPERSER_ENCODER_SRS_LOAD:                "",
		DISPERSER_ENCODER_CACHE_PATH:              "",
		DISPERSER_ENCODER_VERBOSE:                 "",
		DISPERSER_ENCODER_NUM_WORKERS:             fmt.Sprint(runtime.GOMAXPROCS(0)),
		DISPERSER_ENCODER_MAX_CONCURRENT_REQUESTS: "16",
		DISPERSER_ENCODER_REQUEST_POOL_SIZE:       "32",
		DISPERSER_ENCODER_ENCODER_VERSION:         "2",
		DISPERSER_ENCODER_S3_BUCKET_NAME:          "test-eigenda-blobstore",
		DISPERSER_ENCODER_REQUEST_QUEUE_SIZE:      "32",
	}

	env.applyDefaults(&v, "DISPERSER_ENCODER", "enc", ind)

	return v
}

func (env *Config) generateControllerVars(ind int, graphUrl string) ControllerVars {
	v := ControllerVars{
		CONTROLLER_LOG_FORMAT:                              "text",
		CONTROLLER_DYNAMODB_TABLE_NAME:                     "test-BlobMetadata-v2",
		CONTROLLER_BLS_OPERATOR_STATE_RETRIVER:             env.EigenDA.OperatorStateRetreiver,
		CONTROLLER_EIGENDA_SERVICE_MANAGER:                 env.EigenDA.ServiceManager,
		CONTROLLER_USE_GRAPH:                               "true",
		CONTROLLER_GRAPH_URL:                               graphUrl,
		CONTROLLER_ENCODING_PULL_INTERVAL:                  "1s",
		CONTROLLER_AVAILABLE_RELAYS:                        "0,1,2,3",
		CONTROLLER_DISPATCHER_PULL_INTERVAL:                "3s",
		CONTROLLER_NODE_REQUEST_TIMEOUT:                    "5s",
		CONTROLLER_NUM_CONNECTIONS_TO_NODES:                "10",
		CONTROLLER_CHAIN_RPC:                               "",
		CONTROLLER_PRIVATE_KEY:                             "123",
		CONTROLLER_NUM_CONFIRMATIONS:                       "0",
		CONTROLLER_INDEXER_PULL_INTERVAL:                   "1s",
		CONTROLLER_AWS_REGION:                              "",
		CONTROLLER_AWS_ACCESS_KEY_ID:                       "",
		CONTROLLER_AWS_SECRET_ACCESS_KEY:                   "",
		CONTROLLER_AWS_ENDPOINT_URL:                        "",
		CONTROLLER_ENCODER_ADDRESS:                         "0.0.0.0:34001",
		CONTROLLER_FINALIZATION_BLOCK_DELAY:                "0",
		CONTROLLER_DISPERSER_STORE_CHUNKS_SIGNING_DISABLED: "true",
	}
	env.applyDefaults(&v, "CONTROLLER", "controller", ind)

	return v
}

func (env *Config) generateRelayVars(ind int, graphUrl, grpcPort string) RelayVars {
	v := RelayVars{
		RELAY_LOG_FORMAT:                            "text",
		RELAY_GRPC_PORT:                             grpcPort,
		RELAY_BUCKET_NAME:                           "test-eigenda-blobstore",
		RELAY_METADATA_TABLE_NAME:                   "test-BlobMetadata-v2",
		RELAY_RELAY_IDS:                             fmt.Sprint(ind),
		RELAY_BLS_OPERATOR_STATE_RETRIEVER_ADDR:     env.EigenDA.OperatorStateRetreiver,
		RELAY_EIGEN_DA_SERVICE_MANAGER_ADDR:         env.EigenDA.ServiceManager,
		RELAY_PRIVATE_KEY:                           "123",
		RELAY_GRAPH_URL:                             graphUrl,
		RELAY_ONCHAIN_STATE_REFRESH_INTERVAL:        "1s",
		RELAY_MAX_CONCURRENT_GET_CHUNK_OPS_CLIENT:   "10",
		RELAY_MAX_GET_CHUNK_BYTES_PER_SECOND_CLIENT: "100000000",
		RELAY_AUTHENTICATION_DISABLED:               "false",
	}
	env.applyDefaults(&v, "RELAY", "relay", ind)

	return v
}

// Generates DA node .env
func (env *Config) generateOperatorVars(ind int, name, key, churnerUrl, logPath, dbPath, dispersalPort, retrievalPort, metricsPort, nodeApiPort string) OperatorVars {

	max, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	// max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("Could not generate key: %v", err)
	}

	//String representation of n in base 32
	blsKey := n.Text(10)

	blsKeyFile := env.Pks.BlsMap[name].KeyFile
	blsPassword := env.Pks.BlsMap[name].Password
	ecdsaKeyFile := env.Pks.EcdsaMap[name].KeyFile
	ecdsaPassword := env.Pks.EcdsaMap[name].Password

	v := OperatorVars{
		NODE_LOG_FORMAT:                       "text",
		NODE_HOSTNAME:                         "",
		NODE_DISPERSAL_PORT:                   dispersalPort,
		NODE_RETRIEVAL_PORT:                   retrievalPort,
		NODE_INTERNAL_DISPERSAL_PORT:          dispersalPort,
		NODE_INTERNAL_RETRIEVAL_PORT:          retrievalPort,
		NODE_ENABLE_METRICS:                   "true",
		NODE_METRICS_PORT:                     metricsPort,
		NODE_ENABLE_NODE_API:                  "true",
		NODE_API_PORT:                         nodeApiPort,
		NODE_TIMEOUT:                          "10s",
		NODE_QUORUM_ID_LIST:                   "0,1",
		NODE_DB_PATH:                          dbPath,
		NODE_ENABLE_TEST_MODE:                 "false", // using encrypted key in inabox
		NODE_TEST_PRIVATE_BLS:                 blsKey,
		NODE_BLS_KEY_FILE:                     blsKeyFile,
		NODE_ECDSA_KEY_FILE:                   ecdsaKeyFile,
		NODE_BLS_KEY_PASSWORD:                 blsPassword,
		NODE_ECDSA_KEY_PASSWORD:               ecdsaPassword,
		NODE_BLS_OPERATOR_STATE_RETRIVER:      env.EigenDA.OperatorStateRetreiver,
		NODE_EIGENDA_SERVICE_MANAGER:          env.EigenDA.ServiceManager,
		NODE_REGISTER_AT_NODE_START:           "true",
		NODE_CHURNER_URL:                      churnerUrl,
		NODE_CHURNER_USE_SECURE_GRPC:          "false",
		NODE_EXPIRATION_POLL_INTERVAL:         "10",
		NODE_G1_PATH:                          "",
		NODE_G2_PATH:                          "",
		NODE_G2_POWER_OF_2_PATH:               "",
		NODE_CACHE_PATH:                       "",
		NODE_SRS_ORDER:                        "",
		NODE_SRS_LOAD:                         "",
		NODE_NUM_WORKERS:                      fmt.Sprint(runtime.GOMAXPROCS(0)),
		NODE_VERBOSE:                          "true",
		NODE_CHAIN_RPC:                        "",
		NODE_PRIVATE_KEY:                      key[2:],
		NODE_NUM_BATCH_VALIDATORS:             "128",
		NODE_PUBLIC_IP_PROVIDER:               "mockip",
		NODE_PUBLIC_IP_CHECK_INTERVAL:         "10s",
		NODE_NUM_CONFIRMATIONS:                "0",
		NODE_ONCHAIN_METRICS_INTERVAL:         "-1",
		NODE_ENABLE_V2:                        "true",
		NODE_DISABLE_DISPERSAL_AUTHENTICATION: "true",
	}

	env.applyDefaults(&v, "NODE", "opr", ind)
	v.NODE_G2_PATH = ""
	return v

}

// Generates retriever .env
func (env *Config) generateRetrieverVars(ind int, key string, graphUrl, logPath, grpcPort string) RetrieverVars {
	v := RetrieverVars{
		RETRIEVER_LOG_FORMAT:                  "text",
		RETRIEVER_HOSTNAME:                    "",
		RETRIEVER_GRPC_PORT:                   grpcPort,
		RETRIEVER_TIMEOUT:                     "10s",
		RETRIEVER_BLS_OPERATOR_STATE_RETRIVER: env.EigenDA.OperatorStateRetreiver,
		RETRIEVER_EIGENDA_SERVICE_MANAGER:     env.EigenDA.ServiceManager,
		RETRIEVER_NUM_CONNECTIONS:             "10",

		RETRIEVER_CHAIN_RPC:   "",
		RETRIEVER_PRIVATE_KEY: key[2:],

		RETRIEVER_G1_PATH:             "",
		RETRIEVER_G2_PATH:             "",
		RETRIEVER_G2_POWER_OF_2_PATH:  "",
		RETRIEVER_CACHE_PATH:          "",
		RETRIEVER_SRS_ORDER:           "",
		RETRIEVER_SRS_LOAD:            "",
		RETRIEVER_NUM_WORKERS:         fmt.Sprint(runtime.GOMAXPROCS(0)),
		RETRIEVER_VERBOSE:             "true",
		RETRIEVER_CACHE_ENCODED_BLOBS: "false",
		RETRIEVER_GRAPH_URL:           graphUrl,
		RETRIEVER_GRAPH_BACKOFF:       "1s",
		RETRIEVER_GRAPH_MAX_RETRIES:   "3",
	}

	v.RETRIEVER_G2_PATH = ""

	env.applyDefaults(&v, "RETRIEVER", "retriever", ind)

	return v
}

// Used to generate a docker compose file corresponding to the test environment
func (env *Config) genService(compose testbed, name, image, envFile string, ports []string) {

	for i, port := range ports {
		ports[i] = port + ":" + port
	}

	compose.Services[name] = map[string]interface{}{
		"image":    image,
		"env_file": []string{envFile},
		"ports":    ports,
		"volumes": []string{
			env.Path + ":/data",
			env.rootPath + "/inabox/secrets:/secrets",
			env.rootPath + "/inabox/resources:/resources",
		},
		"extra_hosts": []string{
			"host.docker.internal:host-gateway",
		},
	}
}

// Used to generate a docker compose file corresponding to the test environment
func (env *Config) genNodeService(compose testbed, name, image, envFile string, ports []string) {

	for i, port := range ports {
		ports[i] = port + ":" + port
	}

	compose.Services[name] = map[string]interface{}{
		"image":    image,
		"env_file": []string{envFile},
		"volumes": []string{
			env.Path + ":/data",
			env.rootPath + "/inabox/secrets:/secrets",
			env.rootPath + "/inabox/resources:/resources",
		},
		"extra_hosts": []string{
			"host.docker.internal:host-gateway",
		},
		// "environment": []string{
		// 	"NODE_HOSTNAME=" + name,
		// },
	}

	nginxService := name + "_nginx"
	compose.Services[nginxService] = map[string]interface{}{
		"image":    "nginx:latest",
		"env_file": []string{envFile},
		"environment": []string{
			"REQUEST_LIMIT=1r/s",
			"BURST_LIMIT=2",
			"NODE_HOST=" + name,
		},
		"depends_on": []string{name},
		"ports":      ports,
		"volumes": []string{
			env.rootPath + "/node/cmd/resources/nginx-local.conf:/etc/nginx/templates/default.conf.template:ro",
		},
	}
}

func genTelemetryServices(compose testbed, name, image string, volumes []string) {
	compose.Services[name] = map[string]interface{}{
		"image":  image,
		"volume": volumes,
	}
}

func (env *Config) getPaths(name string) (logPath, dbPath, envFilename, envFile string) {
	if env.Environment.IsLocal() {
		logPath = ""
		dbPath = "testdata/" + env.TestName + "/db/" + name
	} else {
		logPath = "/data/logs/" + name
		dbPath = "/data/db/" + name
	}

	envFilename = "envs/" + name + ".env"
	envFile = "testdata/" + env.TestName + "/" + envFilename
	return
}

func (env *Config) getKey(name string) (key, address string) {
	key = env.Pks.EcdsaMap[name].PrivateKey
	log.Printf("name: %s, key: %v", name, key)
	address = GetAddress(key)
	return
}

// GenerateAllVariables all of the config for the test environment.
// Returns an object that corresponds to the participants of the
// current experiment.
func (env *Config) GenerateAllVariables() {
	// hardcode graphurl for now
	graphUrl := "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"

	// Create envs directory
	createDirectory(env.Path + "/envs")
	changeDirectory(env.rootPath + "/inabox")

	// Gather keys
	// keyData := readFile(gethPrivateKeys)
	// keys := strings.Split(string(keyData), "\n")
	// id := 1

	// Create compose file
	composeFile := env.Path + "/docker-compose.yml"
	servicesMap := make(map[string]map[string]interface{})
	compose := testbed{
		Services: servicesMap,
	}

	// Create participants
	port := env.Services.BasePort

	// Generate churners
	name := "churner"
	port += 2
	logPath, _, filename, envFile := env.getPaths(name)
	churnerConfig := env.generateChurnerVars(0, graphUrl, logPath, fmt.Sprint(port))
	writeEnv(churnerConfig.getEnvMap(), envFile)
	env.Churner = churnerConfig
	env.genService(
		compose, name, churnerImage,
		filename, []string{fmt.Sprint(port)})

	churnerUrl := fmt.Sprintf("%s:%s", churnerConfig.CHURNER_HOSTNAME, churnerConfig.CHURNER_GRPC_PORT)

	// Generate disperser nodes

	grpcPort := fmt.Sprint(port + 1)
	port += 2

	name = "dis0"
	logPath, dbPath, filename, envFile := env.getPaths(name)
	disperserConfig := env.generateDisperserVars(0, logPath, dbPath, grpcPort)
	writeEnv(disperserConfig.getEnvMap(), envFile)
	env.Dispersers = append(env.Dispersers, disperserConfig)
	env.genService(
		compose, name, disImage,
		filename, []string{grpcPort})

	// v2 disperser
	grpcPort = fmt.Sprint(port + 1)
	port += 2

	name = "dis1"
	logPath, dbPath, filename, envFile = env.getPaths(name)

	// Convert key to address
	disperserConfig = env.generateDisperserV2Vars(0, logPath, dbPath, grpcPort)
	writeEnv(disperserConfig.getEnvMap(), envFile)
	env.Dispersers = append(env.Dispersers, disperserConfig)

	env.genService(
		compose, name, disImage,
		filename, []string{grpcPort})

	for i := 0; i < env.Services.Counts.NumOpr; i++ {
		metricsPort := fmt.Sprint(port + 1) // port
		dispersalPort := fmt.Sprint(port + 2)
		retrievalPort := fmt.Sprint(port + 3)
		nodeApiPort := fmt.Sprint(port + 4)
		port += 5

		name := fmt.Sprintf("opr%v", i)
		logPath, dbPath, filename, envFile := env.getPaths(name)
		key, _ := env.getKey(name)

		// Convert key to address

		operatorConfig := env.generateOperatorVars(i, name, key, churnerUrl, logPath, dbPath, dispersalPort, retrievalPort, fmt.Sprint(metricsPort), nodeApiPort)
		writeEnv(operatorConfig.getEnvMap(), envFile)
		env.Operators = append(env.Operators, operatorConfig)

		env.genNodeService(
			compose, name, nodeImage,
			filename, []string{dispersalPort, retrievalPort})
	}

	// Batcher
	name = "batcher0"
	logPath, _, filename, envFile = env.getPaths(name)
	key, _ := env.getKey(name)
	batcherConfig := env.generateBatcherVars(0, key, graphUrl, logPath)
	writeEnv(batcherConfig.getEnvMap(), envFile)
	env.Batcher = append(env.Batcher, batcherConfig)
	env.genService(
		compose, name, batcherImage,
		filename, []string{})

	// Encoders
	// TODO: Add more encoders
	name = "enc0"
	_, _, filename, envFile = env.getPaths(name)
	encoderConfig := env.generateEncoderVars(0, "34000")
	writeEnv(encoderConfig.getEnvMap(), envFile)
	env.Encoder = append(env.Encoder, encoderConfig)
	env.genService(
		compose, name, encoderImage,
		filename, []string{"34000"})

	// v2 encoder
	name = "enc1"
	_, _, filename, envFile = env.getPaths(name)
	encoderConfig = env.generateEncoderV2Vars(0, "34001")
	writeEnv(encoderConfig.getEnvMap(), envFile)
	env.Encoder = append(env.Encoder, encoderConfig)
	env.genService(
		compose, name, encoderImage,
		filename, []string{"34001"})

	// Stakers
	for i := 0; i < env.Services.Counts.NumOpr; i++ {

		name := fmt.Sprintf("staker%v", i)
		key, address := env.getKey(name)

		// Create staker paritipants
		participant := Staker{
			Address:    address,
			PrivateKey: key[2:],
		}
		env.Stakers = append(env.Stakers, participant)
	}

	// Relays
	for i := 0; i < env.Services.Counts.NumRelays; i++ {
		name := fmt.Sprintf("relay%v", i)
		grpcPort := fmt.Sprint(port + 1)
		port += 2
		_, _, filename, envFile := env.getPaths(name)
		relayConfig := env.generateRelayVars(i, graphUrl, grpcPort)
		writeEnv(relayConfig.getEnvMap(), envFile)
		env.Relays = append(env.Relays, relayConfig)
		env.genService(
			compose, name, relayImage,
			filename, []string{grpcPort})
	}

	name = "retriever0"
	key, _ = env.getKey(name)
	logPath, _, _, envFile = env.getPaths(name)
	retrieverConfig := env.generateRetrieverVars(0, key, graphUrl, logPath, fmt.Sprint(port+1))
	writeEnv(retrieverConfig.getEnvMap(), envFile)
	env.Retriever = retrieverConfig

	// Controller
	name = "controller0"
	_, _, _, envFile = env.getPaths(name)
	controllerConfig := env.generateControllerVars(0, graphUrl)
	writeEnv(controllerConfig.getEnvMap(), envFile)
	env.Controller = controllerConfig

	if env.Environment.IsLocal() {

		if env.Telemetry.IsNeeded {
			// sd is required for accessing docker daemon
			// agent.yaml configures the grafana agent
			agentVolumes := append(
				env.Telemetry.DockerSd,
				env.Telemetry.ConfigPath+":/etc/agent/agent.yaml",
			)

			// run grafana agent
			genTelemetryServices(compose, "grafana-agent", "grafana/agent", agentVolumes)
			// run node exporter
			genTelemetryServices(compose, "node-exporter", "prom/node-exporter", nil)
		}

		// Write to compose file
		composeYaml, err := yaml.Marshal(&compose)
		if err != nil {
			log.Panicf("Error: %s", err.Error())
		}
		writeFile(composeFile, composeYaml)
	}
}
