build:
# We build the apiserver individually to change its name to "server" instead of "apiserver"
	go build -o ./bin/server ./cmd/apiserver
# All the other binaries (dataapi, encoder, etc) are then built together.
	go build -o ./bin ./...

clean:
	rm -rf ./bin

# Below are example run commands. They are not maintained so likely to be out of date.

run_server: build
	./bin/server \
	--grpc-port 51001 \
	--aws.region us-east-1 \
	--aws.access-key-id xyz \
	--aws.secret-access-key hello \
	--aws.endpoint-url http://0.0.0.0:4566

run_encoder: build
	./bin/encoder \
	--disperser-encoder.grpc-port 34000 \
  --disperser-encoder.metrics-http-port 9109 \
  --kzg.g1-path ../resources/srs/g1.point \
  --kzg.g2-path ../resources/srs/g2.point \
  --kzg.cache-path ../resources/srs/SRSTables \
  --kzg.srs-order 3000 \
  --kzg.num-workers 12 \
  --disperser-encoder.log.level-std debug \
  --disperser-encoder.log.level-file debug

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
