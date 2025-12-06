build:
# We build the apiserver individually to change its name to "server" instead of "apiserver"
	go build -o ./bin/server ./cmd/apiserver
# All the other binaries (dataapi, encoder, batcher, etc) are then built together.
	go build -o ./bin ./...

clean:
	rm -rf ./bin

# Below are example run commands. They are not maintained so likely to be out of date.

run_batcher: build
	./bin/batcher \
	--batcher.pull-interval 10s \
	--batcher.bls-operator-state-retriever 0x9d4454B023096f34B160D6B654540c56A1F81688 \
	--batcher.eigenda-service-manager 0x67d269191c92Caf3cD7723F116c85e6E9bf55933 \
	--chain.rpc http://localhost:8545 \
	--chain.private-key ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
	--batcher.aws.region us-east-1 \
	--batcher.aws.access-key-id xyz \
	--batcher.aws.secret-access-key hello \
	--batcher.aws.endpoint-url http://0.0.0.0:4566 \
	--batcher.s3-bucket-name test-eigenda-blobstore \
	--batcher.dynamodb-table-name test-BlobMetadata \
	--encoder-socket 34000 \
	--batcher.enable-metrics \
	--batcher.graph-url false \
	--batcher.batch-size-limit 10000 \
	--batcher.use-graph false \
	--batcher.srs-order 3000 \
	--encoding-timeout 10s \
	--attestation-timeout 11s \
	--chain-read-timeout 12s \
	--chain-write-timeout 13s

run_server: build
	./bin/server \
	--grpc-port 51001 \
	--aws.region us-east-1 \
	--aws.access-key-id xyz \
	--aws.secret-access-key hello \
	--aws.endpoint-url http://0.0.0.0:4566

run_encoder: build
	export ENCODER_VERSION=1 && \
	export ENCODER_ENABLE_METRICS=true && \
	export ENCODER_METRICS_PORT=9109 && \
	export ENCODER_LOG_LEVEL=debug && \
	export ENCODER_LOG_FORMAT=json && \
	export ENCODER_LOG_COLOR=false && \
	export ENCODER_GRPC_PORT=34000 && \
	export ENCODER_KZG_G1_PATH=../resources/srs/g1.point && \
	export ENCODER_KZG_G2_PATH=../resources/srs/g2.point && \
	export ENCODER_KZG_CACHE_DIR=../resources/srs/SRSTables && \
	export ENCODER_KZG_SRS_ORDER=3000 && \
	export ENCODER_KZG_SRS_NUMBER_TO_LOAD=3000 && \
	export ENCODER_KZG_NUM_WORKER=12 && \
	./bin/encoder

# You can override these defaults via CLI or environment variables
run_blobapi: build
	./bin/blobapi \
	--disperser-server.grpc-port 51002 \
	--disperser-server.enable-metrics=false \
	--auth.registered-quorum 1 \
	--auth.total-unauth-byte-rate 1000000 \
	--auth.per-user-unauth-byte-rate 10000 \
	--auth.total-unauth-blob-rate 100 \
	--auth.per-user-unauth-blob-rate 10 \
	--auth.retrieval-blob-rate 100 \
	--auth.retrieval-throughput 100000 \
	--relay.grpc-port 52002 \
	--relay.relay-keys 1 \
	--relay.enable-metrics=false \
