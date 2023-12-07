package config

import "reflect"

type DisperserVars struct {
	DISPERSER_NAME string

	DISPERSER_SERVER_S3_BUCKET_NAME string

	DISPERSER_SERVER_DYNAMODB_TABLE_NAME string

	DISPERSER_SERVER_GRPC_PORT string

	DISPERSER_SERVER_RATE_BUCKET_TABLE_NAME string

	DISPERSER_SERVER_BLS_OPERATOR_STATE_RETRIVER string

	DISPERSER_SERVER_EIGENDA_SERVICE_MANAGER string

	DISPERSER_SERVER_METRICS_HTTP_PORT string

	DISPERSER_SERVER_ENABLE_METRICS string

	DISPERSER_SERVER_ENABLE_RATELIMITER string

	DISPERSER_SERVER_CHAIN_RPC string

	DISPERSER_SERVER_PRIVATE_KEY string

	DISPERSER_SERVER_STD_LOG_LEVEL string

	DISPERSER_SERVER_FILE_LOG_LEVEL string

	DISPERSER_SERVER_BUCKET_SIZES string

	DISPERSER_SERVER_BUCKET_MULTIPLIERS string

	DISPERSER_SERVER_COUNT_FAILED string

	DISPERSER_SERVER_BUCKET_STORE_SIZE string

	DISPERSER_SERVER_AWS_REGION string

	DISPERSER_SERVER_AWS_ACCESS_KEY_ID string

	DISPERSER_SERVER_AWS_SECRET_ACCESS_KEY string

	DISPERSER_SERVER_AWS_ENDPOINT_URL string

	DISPERSER_SERVER_REGISTERED_QUORUM_ID string

	DISPERSER_SERVER_TOTAL_UNAUTH_THROUGHPUT string

	DISPERSER_SERVER_PER_USER_UNAUTH_THROUGHPUT string
}

func (vars DisperserVars) GetEnvMap() map[string]string {
	v := reflect.ValueOf(vars)
	envMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		envMap[v.Type().Field(i).Name] = v.Field(i).String()
	}
	return envMap
}

type BatcherVars struct {
	BATCHER_NAME string

	BATCHER_S3_BUCKET_NAME string

	BATCHER_DYNAMODB_TABLE_NAME string

	BATCHER_PULL_INTERVAL string

	BATCHER_BLS_OPERATOR_STATE_RETRIVER string

	BATCHER_EIGENDA_SERVICE_MANAGER string

	BATCHER_ENCODER_ADDRESS string

	BATCHER_ENABLE_METRICS string

	BATCHER_GRAPH_URL string

	BATCHER_BATCH_SIZE_LIMIT string

	BATCHER_SRS_ORDER string

	BATCHER_USE_GRAPH string

	BATCHER_METRICS_HTTP_PORT string

	BATCHER_INDEXER_DATA_DIR string

	BATCHER_TIMEOUT string

	BATCHER_NUM_CONNECTIONS string

	BATCHER_FINALIZER_INTERVAL string

	BATCHER_ENCODING_REQUEST_QUEUE_SIZE string

	BATCHER_CHAIN_RPC string

	BATCHER_PRIVATE_KEY string

	BATCHER_STD_LOG_LEVEL string

	BATCHER_FILE_LOG_LEVEL string

	BATCHER_INDEXER_PULL_INTERVAL string

	BATCHER_AWS_REGION string

	BATCHER_AWS_ACCESS_KEY_ID string

	BATCHER_AWS_SECRET_ACCESS_KEY string

	BATCHER_AWS_ENDPOINT_URL string
}

func (vars BatcherVars) GetEnvMap() map[string]string {
	v := reflect.ValueOf(vars)
	envMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		envMap[v.Type().Field(i).Name] = v.Field(i).String()
	}
	return envMap
}

type EncoderVars struct {
	DISPERSER_ENCODER_NAME string

	DISPERSER_ENCODER_GRPC_PORT string

	DISPERSER_ENCODER_METRICS_HTTP_PORT string

	DISPERSER_ENCODER_ENABLE_METRICS string

	DISPERSER_ENCODER_MAX_CONCURRENT_REQUESTS string

	DISPERSER_ENCODER_REQUEST_POOL_SIZE string

	DISPERSER_ENCODER_G1_PATH string

	DISPERSER_ENCODER_G2_PATH string

	DISPERSER_ENCODER_CACHE_PATH string

	DISPERSER_ENCODER_SRS_ORDER string

	DISPERSER_ENCODER_NUM_WORKERS string

	DISPERSER_ENCODER_VERBOSE string

	DISPERSER_ENCODER_CACHE_ENCODED_BLOBS string

	DISPERSER_ENCODER_STD_LOG_LEVEL string

	DISPERSER_ENCODER_FILE_LOG_LEVEL string
}

func (vars EncoderVars) GetEnvMap() map[string]string {
	v := reflect.ValueOf(vars)
	envMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		envMap[v.Type().Field(i).Name] = v.Field(i).String()
	}
	return envMap
}

type OperatorVars struct {
	NODE_NAME string

	NODE_HOSTNAME string

	NODE_DISPERSAL_PORT string

	NODE_RETRIEVAL_PORT string

	NODE_ENABLE_METRICS string

	NODE_METRICS_PORT string

	NODE_ENABLE_NODE_API string

	NODE_API_PORT string

	NODE_TIMEOUT string

	NODE_QUORUM_ID_LIST string

	NODE_DB_PATH string

	NODE_BLS_KEY_FILE string

	NODE_ECDSA_KEY_FILE string

	NODE_BLS_KEY_PASSWORD string

	NODE_ECDSA_KEY_PASSWORD string

	NODE_BLS_OPERATOR_STATE_RETRIVER string

	NODE_EIGENDA_SERVICE_MANAGER string

	NODE_CHURNER_URL string

	NODE_REGISTER_AT_NODE_START string

	NODE_EXPIRATION_POLL_INTERVAL string

	NODE_ENABLE_TEST_MODE string

	NODE_OVERRIDE_BLOCK_STALE_MEASURE string

	NODE_OVERRIDE_STORE_DURATION_BLOCKS string

	NODE_TEST_PRIVATE_BLS string

	NODE_NUM_BATCH_VALIDATORS string

	NODE_G1_PATH string

	NODE_G2_PATH string

	NODE_CACHE_PATH string

	NODE_SRS_ORDER string

	NODE_NUM_WORKERS string

	NODE_VERBOSE string

	NODE_CACHE_ENCODED_BLOBS string

	NODE_CHAIN_RPC string

	NODE_PRIVATE_KEY string

	NODE_STD_LOG_LEVEL string

	NODE_FILE_LOG_LEVEL string

	NODE_PUBLIC_IP_PROVIDER string

	NODE_PUBLIC_IP_CHECK_INTERVAL string
}

func (vars OperatorVars) GetEnvMap() map[string]string {
	v := reflect.ValueOf(vars)
	envMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		envMap[v.Type().Field(i).Name] = v.Field(i).String()
	}
	return envMap
}

type RetrieverVars struct {
	RETRIEVER_NAME string

	RETRIEVER_HOSTNAME string

	RETRIEVER_GRPC_PORT string

	RETRIEVER_TIMEOUT string

	RETRIEVER_BLS_OPERATOR_STATE_RETRIVER string

	RETRIEVER_EIGENDA_SERVICE_MANAGER string

	RETRIEVER_NUM_CONNECTIONS string

	RETRIEVER_DATA_DIR string

	RETRIEVER_METRICS_HTTP_PORT string

	RETRIEVER_G1_PATH string

	RETRIEVER_G2_PATH string

	RETRIEVER_CACHE_PATH string

	RETRIEVER_SRS_ORDER string

	RETRIEVER_NUM_WORKERS string

	RETRIEVER_VERBOSE string

	RETRIEVER_CACHE_ENCODED_BLOBS string

	RETRIEVER_CHAIN_RPC string

	RETRIEVER_PRIVATE_KEY string

	RETRIEVER_STD_LOG_LEVEL string

	RETRIEVER_FILE_LOG_LEVEL string

	RETRIEVER_INDEXER_PULL_INTERVAL string
}

func (vars RetrieverVars) GetEnvMap() map[string]string {
	v := reflect.ValueOf(vars)
	envMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		envMap[v.Type().Field(i).Name] = v.Field(i).String()
	}
	return envMap
}

type ChurnerVars struct {
	CHURNER_NAME string

	CHURNER_HOSTNAME string

	CHURNER_GRPC_PORT string

	CHURNER_GRAPH_URL string

	CHURNER_BLS_OPERATOR_STATE_RETRIVER string

	CHURNER_EIGENDA_SERVICE_MANAGER string

	CHURNER_CHAIN_RPC string

	CHURNER_PRIVATE_KEY string

	CHURNER_STD_LOG_LEVEL string

	CHURNER_FILE_LOG_LEVEL string

	CHURNER_INDEXER_PULL_INTERVAL string

	CHURNER_ENABLE_METRICS string
}

func (vars ChurnerVars) GetEnvMap() map[string]string {
	v := reflect.ValueOf(vars)
	envMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		envMap[v.Type().Field(i).Name] = v.Field(i).String()
	}
	return envMap
}
