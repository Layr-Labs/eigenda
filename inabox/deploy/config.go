package deploy

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
)

const (
	controllerGrpcPort = uint16(30000)
)

var logger = test.GetLogger()

func (env *Config) GetDeployer(name string) (*ContractDeployer, bool) {
	for _, deployer := range env.Deployers {
		if deployer.Name == name {
			return deployer, true
		}
	}
	return nil, false
}

// Constructs a mapping between service names/deployer names (e.g., 'dis0', 'opr1') and private keys
func (env *Config) loadPrivateKeys() error {
	logger.Info("Loading private keys using testbed")

	// Use testbed's LoadPrivateKeys function
	testbedKeys, err := testbed.LoadPrivateKeys(testbed.LoadPrivateKeysInput{
		NumOperators: env.Services.Counts.NumOpr,
		NumRelays:    env.Services.Counts.NumRelays,
	})
	if err != nil {
		return fmt.Errorf("failed to load private keys from testbed: %w", err)
	}

	// Convert testbed keys to our format
	if env.Pks == nil {
		env.Pks = &PkConfig{
			EcdsaMap: make(map[string]KeyInfo),
			BlsMap:   make(map[string]KeyInfo),
		}
	} else {
		// Initialize maps if they're nil
		if env.Pks.EcdsaMap == nil {
			env.Pks.EcdsaMap = make(map[string]KeyInfo)
		}
		if env.Pks.BlsMap == nil {
			env.Pks.BlsMap = make(map[string]KeyInfo)
		}
	}

	// Copy testbed keys to our structure
	for name, keyInfo := range testbedKeys.EcdsaMap {
		env.Pks.EcdsaMap[name] = KeyInfo{
			PrivateKey: keyInfo.PrivateKey,
			Password:   keyInfo.Password,
			KeyFile:    keyInfo.KeyFile,
		}
	}

	for name, keyInfo := range testbedKeys.BlsMap {
		env.Pks.BlsMap[name] = KeyInfo{
			PrivateKey: keyInfo.PrivateKey,
			Password:   keyInfo.Password,
			KeyFile:    keyInfo.KeyFile,
		}
	}

	// Add deployer keys if they don't exist (for backward compatibility)
	for _, d := range env.Deployers {
		if _, exists := env.Pks.EcdsaMap[d.Name]; !exists {
			// Use the same key as "deployer" if available
			if deployerKey, ok := env.Pks.EcdsaMap["deployer"]; ok {
				env.Pks.EcdsaMap[d.Name] = deployerKey
				env.Pks.BlsMap[d.Name] = env.Pks.BlsMap["deployer"]
			}
		}
	}

	logger.Info("Successfully loaded private keys", "ecdsaKeys", len(env.Pks.EcdsaMap), "blsKeys", len(env.Pks.BlsMap))

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
		CHURNER_EIGENDA_DIRECTORY:           env.EigenDA.EigenDADirectory,
		CHURNER_BLS_OPERATOR_STATE_RETRIVER: env.EigenDA.OperatorStateRetriever,
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

		DISPERSER_SERVER_EIGENDA_DIRECTORY:           env.EigenDA.EigenDADirectory,
		DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER: env.EigenDA.OperatorStateRetriever,
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

		DISPERSER_SERVER_EIGENDA_DIRECTORY:           env.EigenDA.EigenDADirectory,
		DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER: env.EigenDA.OperatorStateRetriever,
		DISPERSER_SERVER_EIGENDA_SERVICE_MANAGER:     env.EigenDA.ServiceManager,
		DISPERSER_SERVER_DISPERSER_VERSION:           "2",

		DISPERSER_SERVER_ENABLE_PAYMENT_METERER:  "true",
		DISPERSER_SERVER_RESERVED_ONLY:           "false",
		DISPERSER_SERVER_RESERVATIONS_TABLE_NAME: "e2e_v2_reservation",
		DISPERSER_SERVER_ON_DEMAND_TABLE_NAME:    "e2e_v2_ondemand",
		DISPERSER_SERVER_GLOBAL_RATE_TABLE_NAME:  "e2e_v2_global_reservation",
		DISPERSER_SERVER_CONTROLLER_ADDRESS:      fmt.Sprintf("localhost:%d", controllerGrpcPort),

		// DisperserV2 uses the V2 prover which always uses SRSOrder=2^28.
		// So it needs the trailing g2 points to generate correct length commitments.
		DISPERSER_SERVER_G2_TRAILING_PATH:               "../resources/srs/g2.trailing.point",
		DISPERSER_SERVER_ONCHAIN_STATE_REFRESH_INTERVAL: "1s",
	}

	if env.UseNewPayments {
		v.DISPERSER_SERVER_USE_CONTROLLER_MEDIATED_PAYMENTS = "true"
	} else {
		v.DISPERSER_SERVER_USE_CONTROLLER_MEDIATED_PAYMENTS = "false"
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
		BATCHER_OBJECT_STORAGE_BACKEND:        "s3",
		BATCHER_ENABLE_METRICS:                "true",
		BATCHER_METRICS_HTTP_PORT:             "9094",
		BATCHER_PULL_INTERVAL:                 "5s",
		BATCHER_EIGENDA_DIRECTORY:             env.EigenDA.EigenDADirectory,
		BATCHER_BLS_OPERATOR_STATE_RETRIVER:   env.EigenDA.OperatorStateRetriever,
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
		ENCODER_LOG_FORMAT:                               "text",
		ENCODER_AWS_REGION:                               "",
		ENCODER_AWS_ACCESS_KEY:                           "",
		ENCODER_AWS_SECRET_ACCESS_KEY:                    "",
		ENCODER_AWS_ENDPOINT_URL:                         "",
		ENCODER_GRPC_PORT:                                grpcPort,
		ENCODER_ENABLE_METRICS:                           "true",
		ENCODER_KZG_G1_PATH:                              "",
		ENCODER_KZG_G2_PATH:                              "",
		ENCODER_KZG_SRS_ORDER:                            "",
		ENCODER_KZG_SRS_NUMBER_TO_LOAD:                   "",
		ENCODER_KZG_CACHE_DIR:                            "",
		ENCODER_KZG_VERBOSE:                              "",
		ENCODER_KZG_NUM_WORKER:                           fmt.Sprint(runtime.GOMAXPROCS(0)),
		ENCODER_SERVER_MAX_CONCURRENT_REQUESTS_DANGEROUS: "16",
		ENCODER_SERVER_REQUEST_POOL_SIZE:                 "32",
		ENCODER_SERVER_REQUEST_QUEUE_SIZE:                "32",
		ENCODER_ENCODER_VERSION:                          "1",
	}

	env.applyDefaults(&v, "ENCODER", "enc", ind)

	return v
}

func (env *Config) generateEncoderV2Vars(ind int, grpcPort string) EncoderVars {
	v := EncoderVars{
		ENCODER_LOG_FORMAT:                               "text",
		ENCODER_AWS_REGION:                               "",
		ENCODER_AWS_ACCESS_KEY:                           "",
		ENCODER_AWS_SECRET_ACCESS_KEY:                    "",
		ENCODER_AWS_ENDPOINT_URL:                         "",
		ENCODER_GRPC_PORT:                                grpcPort,
		ENCODER_ENABLE_METRICS:                           "true",
		ENCODER_KZG_G1_PATH:                              "",
		ENCODER_KZG_G2_PATH:                              "",
		ENCODER_KZG_SRS_ORDER:                            "",
		ENCODER_KZG_SRS_NUMBER_TO_LOAD:                   "",
		ENCODER_KZG_CACHE_DIR:                            "",
		ENCODER_KZG_VERBOSE:                              "",
		ENCODER_KZG_NUM_WORKER:                           fmt.Sprint(runtime.GOMAXPROCS(0)),
		ENCODER_SERVER_MAX_CONCURRENT_REQUESTS_DANGEROUS: "16",
		ENCODER_SERVER_REQUEST_POOL_SIZE:                 "32",
		ENCODER_ENCODER_VERSION:                          "2",
		ENCODER_BLOB_STORE_BUCKET_NAME:                   "test-eigenda-blobstore",
		ENCODER_CHUNK_STORE_BUCKET_NAME:                  "test-eigenda-blobstore",
		ENCODER_SERVER_REQUEST_QUEUE_SIZE:                "32",
	}

	env.applyDefaults(&v, "ENCODER", "enc", ind)

	return v
}

func (env *Config) generateControllerVars(
	ind int,
	graphUrl string) ControllerVars {

	v := ControllerVars{
		CONTROLLER_LOG_FORMAT:                         "text",
		CONTROLLER_DYNAMODB_TABLE_NAME:                "test-BlobMetadata-v2",
		CONTROLLER_EIGENDA_CONTRACT_DIRECTORY_ADDRESS: env.EigenDA.EigenDADirectory,
		CONTROLLER_USE_GRAPH:                          "true",
		CONTROLLER_GRAPH_URL:                          graphUrl,
		CONTROLLER_ENCODING_PULL_INTERVAL:             "1s",
		CONTROLLER_AVAILABLE_RELAYS:                   "0,1,2,3",
		CONTROLLER_DISPATCHER_PULL_INTERVAL:           "3s",
		CONTROLLER_ATTESTATION_TIMEOUT:                "5s",
		CONTROLLER_BATCH_ATTESTATION_TIMEOUT:          "6s",
		CONTROLLER_CHAIN_RPC:                          "",
		CONTROLLER_PRIVATE_KEY:                        "123",
		CONTROLLER_NUM_CONFIRMATIONS:                  "0",
		CONTROLLER_INDEXER_PULL_INTERVAL:              "1s",
		CONTROLLER_AWS_REGION:                         "",
		CONTROLLER_AWS_ACCESS_KEY_ID:                  "",
		CONTROLLER_AWS_SECRET_ACCESS_KEY:              "",
		CONTROLLER_AWS_ENDPOINT_URL:                   "",
		CONTROLLER_ENCODER_ADDRESS:                    "0.0.0.0:34001",
		CONTROLLER_BATCH_METADATA_UPDATE_PERIOD:       "100ms",
		// set to 5 to ensure payload disperser checkDACert calls pass in integration_v2 test since
		// disperser chooses rbn = latest_block_number - finalization_block_delay
		CONTROLLER_FINALIZATION_BLOCK_DELAY:                "5",
		CONTROLLER_DISPERSER_STORE_CHUNKS_SIGNING_DISABLED: "false",
		CONTROLLER_DISPERSER_KMS_KEY_ID:                    env.DisperserKMSKeyID,
		CONTROLLER_DISPERSER_ID:                            "0",
	}

	if env.UseNewPayments {
		v.CONTROLLER_GRPC_SERVER_ENABLE = "true"
		v.CONTROLLER_GRPC_PAYMENT_AUTHENTICATION = "true"
		v.CONTROLLER_GRPC_PORT = fmt.Sprintf("%d", controllerGrpcPort)
		v.CONTROLLER_ON_DEMAND_PAYMENTS_TABLE_NAME = "e2e_v2_ondemand"
		v.CONTROLLER_PAYMENT_VAULT_UPDATE_INTERVAL = "1s"
	} else {
		v.CONTROLLER_GRPC_SERVER_ENABLE = "false"
		v.CONTROLLER_GRPC_PAYMENT_AUTHENTICATION = "false"
	}
	env.applyDefaults(&v, "CONTROLLER", "controller", ind)

	return v
}

func (env *Config) generateProxyVars(ind int) ProxyVars {
	v := ProxyVars{
		EIGENDA_PROXY_APIS_TO_ENABLE:             "op-generic,standard,metrics",
		EIGENDA_PROXY_STORAGE_BACKENDS_TO_ENABLE: "V2", // we only enable V2
		EIGENDA_PROXY_STORAGE_DISPERSAL_BACKEND:  "V2",
		// V2 Variables
		// TODO(samlaf): this private key should be read from the output config file instead of hardcoded.
		EIGENDA_PROXY_EIGENDA_V2_SIGNER_PRIVATE_KEY_HEX: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded",
		// TODO(samlaf): this should not be hardcoded
		EIGENDA_PROXY_EIGENDA_V2_ETH_RPC:                                         "http://localhost:8545",
		EIGENDA_PROXY_EIGENDA_V2_MAX_BLOB_LENGTH:                                 "16MiB",
		EIGENDA_PROXY_EIGENDA_V2_CERT_VERIFIER_ROUTER_OR_IMMUTABLE_VERIFIER_ADDR: env.EigenDA.CertVerifierRouter,
		EIGENDA_PROXY_EIGENDA_V2_RBN_RECENCY_WINDOW_SIZE:                         "0",
		// TODO(samlaf): this should not be hardcoded
		EIGENDA_PROXY_EIGENDA_V2_DISPERSER_RPC:     "localhost:32005",
		EIGENDA_PROXY_EIGENDA_V2_EIGENDA_DIRECTORY: env.EigenDA.EigenDADirectory,
		EIGENDA_PROXY_EIGENDA_V2_GRPC_DISABLE_TLS:  "true",
		// SRS paths
		EIGENDA_PROXY_EIGENDA_TARGET_KZG_G1_PATH:          "../resources/srs/g1.point",
		EIGENDA_PROXY_EIGENDA_TARGET_KZG_G2_PATH:          "../resources/srs/g2.point",
		EIGENDA_PROXY_EIGENDA_TARGET_KZG_G2_TRAILING_PATH: "../resources/srs/g2.trailing.point",
	}
	env.applyDefaults(&v, "EIGENDA_PROXY", "proxy", ind)
	return v
}

func (env *Config) generateRelayVars(ind int, graphUrl, grpcPort string) RelayVars {
	v := RelayVars{
		RELAY_LOG_FORMAT:                            "text",
		RELAY_GRPC_PORT:                             grpcPort,
		RELAY_BUCKET_NAME:                           "test-eigenda-blobstore",
		RELAY_METADATA_TABLE_NAME:                   "test-BlobMetadata-v2",
		RELAY_RELAY_KEYS:                            fmt.Sprint(ind),
		RELAY_EIGENDA_DIRECTORY:                     env.EigenDA.EigenDADirectory,
		RELAY_BLS_OPERATOR_STATE_RETRIEVER_ADDR:     env.EigenDA.OperatorStateRetriever,
		RELAY_EIGEN_DA_SERVICE_MANAGER_ADDR:         env.EigenDA.ServiceManager,
		RELAY_PRIVATE_KEY:                           "123",
		RELAY_GRAPH_URL:                             graphUrl,
		RELAY_ONCHAIN_STATE_REFRESH_INTERVAL:        "1s",
		RELAY_MAX_CONCURRENT_GET_CHUNK_OPS_CLIENT:   "10",
		RELAY_MAX_GET_CHUNK_BYTES_PER_SECOND_CLIENT: "100000000",
		RELAY_AUTHENTICATION_DISABLED:               "false",
		RELAY_ENABLE_METRICS:                        "true",
	}
	env.applyDefaults(&v, "RELAY", "relay", ind)

	return v
}

// Generates DA node .env
func (env *Config) generateOperatorVars(ind int, name, key, churnerUrl, logPath, dbPath, dispersalPort, retrievalPort, v2DispersalPort, v2RetrievalPort, metricsPort, nodeApiPort string) OperatorVars {

	max, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	// max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.Fatal("Could not generate key", "error", err)
	}

	//String representation of n in base 32
	blsKey := n.Text(10)

	blsKeyFile := env.Pks.BlsMap[name].KeyFile
	blsPassword := env.Pks.BlsMap[name].Password
	ecdsaKeyFile := env.Pks.EcdsaMap[name].KeyFile
	ecdsaPassword := env.Pks.EcdsaMap[name].Password

	v := OperatorVars{
		NODE_LOG_FORMAT:               "text",
		NODE_HOSTNAME:                 "",
		NODE_DISPERSAL_PORT:           dispersalPort,
		NODE_RETRIEVAL_PORT:           retrievalPort,
		NODE_INTERNAL_DISPERSAL_PORT:  dispersalPort,
		NODE_INTERNAL_RETRIEVAL_PORT:  retrievalPort,
		NODE_V2_DISPERSAL_PORT:        v2DispersalPort,
		NODE_V2_RETRIEVAL_PORT:        v2RetrievalPort,
		NODE_ENABLE_METRICS:           "true",
		NODE_METRICS_PORT:             metricsPort,
		NODE_ENABLE_NODE_API:          "true",
		NODE_API_PORT:                 nodeApiPort,
		NODE_TIMEOUT:                  "10s",
		NODE_QUORUM_ID_LIST:           "0,1",
		NODE_DB_PATH:                  dbPath,
		NODE_LITT_DB_STORAGE_PATHS:    dbPath,
		NODE_ENABLE_TEST_MODE:         "false", // using encrypted key in inabox
		NODE_TEST_PRIVATE_BLS:         blsKey,
		NODE_BLS_KEY_FILE:             blsKeyFile,
		NODE_ECDSA_KEY_FILE:           ecdsaKeyFile,
		NODE_BLS_KEY_PASSWORD:         blsPassword,
		NODE_ECDSA_KEY_PASSWORD:       ecdsaPassword,
		NODE_EIGENDA_DIRECTORY:        env.EigenDA.EigenDADirectory,
		NODE_REGISTER_AT_NODE_START:   "true",
		NODE_CHURNER_URL:              churnerUrl,
		NODE_CHURNER_USE_SECURE_GRPC:  "false",
		NODE_RELAY_USE_SECURE_GRPC:    "false",
		NODE_EXPIRATION_POLL_INTERVAL: "10",
		NODE_G1_PATH:                  "",
		NODE_G2_PATH:                  "",
		NODE_G2_POWER_OF_2_PATH:       "",
		NODE_CACHE_PATH:               "",
		NODE_SRS_ORDER:                "",
		NODE_SRS_LOAD:                 "",
		NODE_NUM_WORKERS:              fmt.Sprint(runtime.GOMAXPROCS(0)),
		NODE_VERBOSE:                  "true",
		NODE_CHAIN_RPC:                "",
		NODE_PRIVATE_KEY:              key[2:],
		NODE_NUM_BATCH_VALIDATORS:     "128",
		NODE_PUBLIC_IP_PROVIDER:       "mockip",
		NODE_PUBLIC_IP_CHECK_INTERVAL: "10s",
		NODE_NUM_CONFIRMATIONS:        "0",
		NODE_ONCHAIN_METRICS_INTERVAL: "-1",
		NODE_RUNTIME_MODE:             "v1-and-v2",
	}

	env.applyDefaults(&v, "NODE", "opr", ind)
	v.NODE_G2_PATH = ""
	return v

}

// Generates retriever .env
func (env *Config) generateRetrieverVars(ind int, key string, graphUrl, logPath, grpcPort string) RetrieverVars {
	v := RetrieverVars{
		RETRIEVER_LOG_FORMAT:              "text",
		RETRIEVER_HOSTNAME:                "",
		RETRIEVER_GRPC_PORT:               grpcPort,
		RETRIEVER_TIMEOUT:                 "10s",
		RETRIEVER_EIGENDA_DIRECTORY:       env.EigenDA.EigenDADirectory,
		RETRIEVER_EIGENDA_SERVICE_MANAGER: env.EigenDA.ServiceManager,
		RETRIEVER_NUM_CONNECTIONS:         "10",

		RETRIEVER_CHAIN_RPC:   "",
		RETRIEVER_PRIVATE_KEY: key[2:],

		RETRIEVER_G1_PATH:             "",
		RETRIEVER_G2_PATH:             "",
		RETRIEVER_CACHE_PATH:          "",
		RETRIEVER_SRS_ORDER:           "",
		RETRIEVER_SRS_LOAD:            "",
		RETRIEVER_NUM_WORKERS:         fmt.Sprint(runtime.GOMAXPROCS(0)),
		RETRIEVER_VERBOSE:             "true",
		RETRIEVER_CACHE_ENCODED_BLOBS: "false",
	}

	v.RETRIEVER_G2_PATH = ""

	env.applyDefaults(&v, "RETRIEVER", "retriever", ind)

	return v
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

func (env *Config) getKey(name string) (key, address string, err error) {
	key = env.Pks.EcdsaMap[name].PrivateKey
	logger.Debug("Getting key", "name", name, "key", key)
	address, err = GetAddress(key)
	if err != nil {
		logger.Error("Failed to get address", "error", err)
		return "", "", fmt.Errorf("failed to get address: %w", err)
	}

	return key, address, nil
}

// GenerateAllVariables all of the config for the test environment.
// Returns an object that corresponds to the participants of the
// current experiment.
func (env *Config) GenerateAllVariables() error {
	// hardcode graphurl for now
	graphUrl := "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"

	env.localstackEndpoint = "http://localhost:4570"
	env.localstackRegion = "us-east-1"

	// Create envs directory
	if err := createDirectory(env.Path + "/envs"); err != nil {
		return fmt.Errorf("failed to create envs directory: %w", err)
	}

	logger.Info("Changing directories", "path", env.rootPath+"/inabox")
	if err := changeDirectory(env.rootPath + "/inabox"); err != nil {
		return fmt.Errorf("failed to change directories: %w", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	// Create participants
	port := env.Services.BasePort

	// Generate churners
	name := "churner"
	port += 2
	logPath, _, _, envFile := env.getPaths(name)
	churnerConfig := env.generateChurnerVars(0, graphUrl, logPath, fmt.Sprint(port))
	if err := writeEnv(churnerConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Churner = churnerConfig
	churnerUrl := fmt.Sprintf("%s:%s", churnerConfig.CHURNER_HOSTNAME, churnerConfig.CHURNER_GRPC_PORT)

	// Generate disperser nodes

	grpcPort := fmt.Sprint(port + 1)
	port += 2

	name = "dis0"
	logPath, dbPath, _, envFile := env.getPaths(name)
	disperserConfig := env.generateDisperserVars(0, logPath, dbPath, grpcPort)
	if err := writeEnv(disperserConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Dispersers = append(env.Dispersers, disperserConfig)

	// v2 disperser
	grpcPort = fmt.Sprint(port + 1)
	port += 2

	name = "dis1"
	logPath, dbPath, _, envFile = env.getPaths(name)

	// Convert key to address
	disperserConfig = env.generateDisperserV2Vars(0, logPath, dbPath, grpcPort)
	if err := writeEnv(disperserConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Dispersers = append(env.Dispersers, disperserConfig)

	for i := 0; i < env.Services.Counts.NumOpr; i++ {
		metricsPort := fmt.Sprint(port + 1) // port
		dispersalPort := fmt.Sprint(port + 2)
		retrievalPort := fmt.Sprint(port + 3)
		v2DispersalPort := fmt.Sprint(port + 4)
		v2RetrievalPort := fmt.Sprint(port + 5)
		nodeApiPort := fmt.Sprint(port + 6)
		port += 7

		name := fmt.Sprintf("opr%v", i)
		logPath, dbPath, _, envFile := env.getPaths(name)
		key, _, err := env.getKey(name)
		if err != nil {
			return fmt.Errorf("failed to get key for %s: %w", name, err)
		}

		// Convert key to address

		operatorConfig := env.generateOperatorVars(i, name, key, churnerUrl, logPath, dbPath, dispersalPort, retrievalPort, v2DispersalPort, v2RetrievalPort, fmt.Sprint(metricsPort), nodeApiPort)
		if err := writeEnv(operatorConfig.getEnvMap(), envFile); err != nil {
			return fmt.Errorf("failed to write env file: %w", err)
		}
		env.Operators = append(env.Operators, operatorConfig)
	}

	// Batcher
	name = "batcher0"
	logPath, _, _, envFile = env.getPaths(name)
	key, _, err := env.getKey(name)
	if err != nil {
		return fmt.Errorf("failed to get key for %s: %w", name, err)
	}

	batcherConfig := env.generateBatcherVars(0, key, graphUrl, logPath)
	if err := writeEnv(batcherConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Batcher = append(env.Batcher, batcherConfig)

	// Encoders
	// TODO: Add more encoders
	name = "enc0"
	_, _, _, envFile = env.getPaths(name)
	encoderConfig := env.generateEncoderVars(0, "34000")
	if err := writeEnv(encoderConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Encoder = append(env.Encoder, encoderConfig)

	// v2 encoder
	name = "enc1"
	_, _, _, envFile = env.getPaths(name)
	encoderConfig = env.generateEncoderV2Vars(0, "34001")
	if err := writeEnv(encoderConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Encoder = append(env.Encoder, encoderConfig)

	// Stakers
	for i := 0; i < env.Services.Counts.NumOpr; i++ {

		name := fmt.Sprintf("staker%v", i)
		key, address, err := env.getKey(name)
		if err != nil {
			return fmt.Errorf("failed to get key for %s: %w", name, err)
		}

		// Create staker participants
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
		_, _, _, envFile := env.getPaths(name)
		relayConfig := env.generateRelayVars(i, graphUrl, grpcPort)
		if err := writeEnv(relayConfig.getEnvMap(), envFile); err != nil {
			return fmt.Errorf("failed to write env file: %w", err)
		}
		env.Relays = append(env.Relays, relayConfig)
	}

	name = "retriever0"
	key, _, err = env.getKey(name)
	if err != nil {
		return fmt.Errorf("failed to get key for %s: %w", name, err)
	}

	logPath, _, _, envFile = env.getPaths(name)
	retrieverConfig := env.generateRetrieverVars(0, key, graphUrl, logPath, fmt.Sprint(port+1))
	if err := writeEnv(retrieverConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Retriever = retrieverConfig

	// Controller
	name = "controller0"
	_, _, _, envFile = env.getPaths(name)
	controllerConfig := env.generateControllerVars(0, graphUrl)
	if err := writeEnv(controllerConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Controller = controllerConfig

	// Proxy
	name = "proxy0"
	_, _, _, envFile = env.getPaths(name)
	proxyConfig := env.generateProxyVars(0)
	if err := writeEnv(proxyConfig.getEnvMap(), envFile); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}
	env.Proxy = proxyConfig

	return nil
}
