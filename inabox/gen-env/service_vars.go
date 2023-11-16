package genenv

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/Layr-Labs/eigenda/inabox/config"
)

// Generates churner .env
func GenerateChurnerVars(lock *config.ConfigLock, ind int, grpcPort int) config.ChurnerVars {
	v := config.ChurnerVars{
		CHURNER_NAME:                        fmt.Sprintf("churner%d", ind),
		CHURNER_HOSTNAME:                    "",
		CHURNER_GRPC_PORT:                   fmt.Sprint(grpcPort),
		CHURNER_BLS_OPERATOR_STATE_RETRIVER: lock.Config.EigenDA.OperatorStateRetreiver,
		CHURNER_EIGENDA_SERVICE_MANAGER:     lock.Config.EigenDA.ServiceManager,

		CHURNER_CHAIN_RPC:   "",
		CHURNER_PRIVATE_KEY: strings.TrimPrefix(lock.Pks.EcdsaMap[lock.Config.EigenDA.Deployer].PrivateKey, "0x"),

		CHURNER_GRAPH_URL:             "",
		CHURNER_INDEXER_PULL_INTERVAL: "1s",

		CHURNER_STD_LOG_LEVEL:  "debug",
		CHURNER_FILE_LOG_LEVEL: "trace",
		CHURNER_ENABLE_METRICS: "true",
	}

	applyOverrides(lock, &v, "CHURNER", "churner", ind)
	v.CHURNER_HOSTNAME = "churner"

	return v
}

// Generates disperser .env
func GenerateDisperserVars(lock *config.ConfigLock, ind int, grpcPort int) config.DisperserVars {
	v := config.DisperserVars{
		DISPERSER_NAME:                          fmt.Sprintf("dis%d", ind),
		DISPERSER_SERVER_S3_BUCKET_NAME:         "test-eigenda-blobstore",
		DISPERSER_SERVER_DYNAMODB_TABLE_NAME:    "test-BlobMetadata",
		DISPERSER_SERVER_RATE_BUCKET_TABLE_NAME: "test-BucketStore",
		DISPERSER_SERVER_GRPC_PORT:              fmt.Sprint(grpcPort),
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

		DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER: lock.Config.EigenDA.OperatorStateRetreiver,
		DISPERSER_SERVER_EIGENDA_SERVICE_MANAGER:     lock.Config.EigenDA.ServiceManager,
	}

	applyOverrides(lock, &v, "DISPERSER_SERVER", "dis", ind)

	return v

}

// Generates batcher .env
func GenerateBatcherVars(lock *config.ConfigLock, ind int) config.BatcherVars {
	name := fmt.Sprintf("batcher%d", ind)
	key, _ := lock.GetKey(name)
	v := config.BatcherVars{
		BATCHER_NAME:                        name,
		BATCHER_S3_BUCKET_NAME:              "test-eigenda-blobstore",
		BATCHER_DYNAMODB_TABLE_NAME:         "test-BlobMetadata",
		BATCHER_ENABLE_METRICS:              "true",
		BATCHER_METRICS_HTTP_PORT:           "9094",
		BATCHER_PULL_INTERVAL:               "5s",
		BATCHER_BLS_OPERATOR_STATE_RETRIVER: lock.Config.EigenDA.OperatorStateRetreiver,
		BATCHER_EIGENDA_SERVICE_MANAGER:     lock.Config.EigenDA.ServiceManager,
		BATCHER_SRS_ORDER:                   "300000",
		BATCHER_CHAIN_RPC:                   "",
		BATCHER_PRIVATE_KEY:                 key[2:],
		BATCHER_STD_LOG_LEVEL:               "debug",
		BATCHER_FILE_LOG_LEVEL:              "trace",
		BATCHER_GRAPH_URL:                   "",
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

	applyOverrides(lock, &v, "BATCHER", "batcher", ind)

	return v
}

func GenerateEncoderVars(lock *config.ConfigLock, ind int, grpcPort int) config.EncoderVars {
	v := config.EncoderVars{
		DISPERSER_ENCODER_NAME:                    fmt.Sprintf("encoder%d", ind),
		DISPERSER_ENCODER_GRPC_PORT:               fmt.Sprint(grpcPort),
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

	applyOverrides(lock, &v, "DISPERSER_ENCODER", "enc", ind)

	return v
}

// Generates DA node .env
func GenerateOperatorVars(lock *config.ConfigLock, ind int, dispersalPort, retrievalPort, metricsPort, nodeApiPort int) config.OperatorVars {
	name := fmt.Sprintf("opr%d", ind)
	blsKey := lock.Pks.BlsMap[name].PrivateKey
	ecdsaKey := lock.Pks.EcdsaMap[name].PrivateKey

	v := config.OperatorVars{
		NODE_NAME:                        name,
		NODE_HOSTNAME:                    "",
		NODE_DISPERSAL_PORT:              fmt.Sprint(dispersalPort),
		NODE_RETRIEVAL_PORT:              fmt.Sprint(retrievalPort),
		NODE_ENABLE_METRICS:              "true",
		NODE_METRICS_PORT:                fmt.Sprint(metricsPort),
		NODE_ENABLE_NODE_API:             "true",
		NODE_API_PORT:                    fmt.Sprint(nodeApiPort),
		NODE_TIMEOUT:                     "10",
		NODE_QUORUM_ID_LIST:              "0",
		NODE_DB_PATH:                     "/db",
		NODE_ENABLE_TEST_MODE:            "true",
		NODE_TEST_PRIVATE_BLS:            blsKey,
		NODE_PRIVATE_KEY:                 ecdsaKey,
		NODE_BLS_OPERATOR_STATE_RETRIVER: lock.Config.EigenDA.OperatorStateRetreiver,
		NODE_EIGENDA_SERVICE_MANAGER:     lock.Config.EigenDA.ServiceManager,
		NODE_REGISTER_AT_NODE_START:      "true",
		NODE_CHURNER_URL:                 fmt.Sprintf("%s:%s", lock.Envs.Churner.CHURNER_HOSTNAME, lock.Envs.Churner.CHURNER_GRPC_PORT),
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

	applyOverrides(lock, &v, "NODE", "opr", ind)

	return v

}

// Generates retriever .env
func GenerateRetrieverVars(lock *config.ConfigLock, ind int, grpcPort int) config.RetrieverVars {
	v := config.RetrieverVars{
		RETRIEVER_NAME:                        fmt.Sprintf("retriever%d", ind),
		RETRIEVER_HOSTNAME:                    "",
		RETRIEVER_GRPC_PORT:                   fmt.Sprint(grpcPort),
		RETRIEVER_TIMEOUT:                     "10s",
		RETRIEVER_BLS_OPERATOR_STATE_RETRIVER: lock.Config.EigenDA.OperatorStateRetreiver,
		RETRIEVER_EIGENDA_SERVICE_MANAGER:     lock.Config.EigenDA.ServiceManager,
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

	applyOverrides(lock, &v, "RETRIEVER", "retriever", ind)

	return v
}

func applyOverrides(lock *config.ConfigLock, c any, prefix, stub string, ind int) {

	pv := reflect.ValueOf(c)
	v := pv.Elem()

	prefix += "_"

	for key, value := range lock.Config.Services.Variables["globals"] {
		field := v.FieldByName(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

	for key, value := range lock.Config.Services.Variables[stub] {
		field := v.FieldByName(prefix + key)
		fmt.Println(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

	for key, value := range lock.Config.Services.Variables[fmt.Sprintf("%v%v", stub, ind)] {
		field := v.FieldByName(prefix + key)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

}
