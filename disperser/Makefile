clean:
	rm -rf ./bin

build: build_server build_batcher build_encoder build_dataapi build_controller build_blobapi

build_batcher:
	go build -o ./bin/batcher ./cmd/batcher

build_server:
	go build -o ./bin/server ./cmd/apiserver

build_encoder:
	go build -o ./bin/encoder ./cmd/encoder

build_dataapi:
	go build -o ./bin/dataapi ./cmd/dataapi

build_controller:
	go build -o ./bin/controller ./cmd/controller

build_blobapi:
	go build -o ./bin/blobapi ./cmd/blobapi

run_batcher: build_batcher
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

run_server: build_server
	./bin/server \
	--grpc-port 51001 \
	--aws.region us-east-1 \
	--aws.access-key-id xyz \
	--aws.secret-access-key hello \
	--aws.endpoint-url http://0.0.0.0:4566

run_encoder: build_encoder
	./bin/encoder \
	--disperser-encoder.grpc-port 34000 \
  --disperser-encoder.metrics-http-port 9109 \
  --kzg.g1-path ../inabox/resources/kzg/g1.point \
  --kzg.g2-path ../inabox/resources/kzg/g2.point \
  --kzg.cache-path ../inabox/resources/kzg/SRSTables \
  --kzg.srs-order 3000 \
  --kzg.num-workers 12 \
  --disperser-encoder.log.level-std debug \
  --disperser-encoder.log.level-file debug

run_dataapi_tests:
	go test -v disperser/dataapi

# You can override these defaults via CLI or environment variables
run_blobapi: build_blobapi
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
