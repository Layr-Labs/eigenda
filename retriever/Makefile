clean:
	rm -rf ./bin

build: clean
	# cd .. && make protoc
	go mod tidy
	go build -o ./bin/server ./cmd

run: build
	DA_RETRIEVER_HOSTNAME=localhost \
	DA_RETRIEVER_GRPC_PORT=50051 \
	DA_RETRIEVER_TIMEOUT=10s \
	./bin/server \
	--retriever.hostname localhost \
	--retriever.grpc-port 32011 \
	--retriever.timeout 10s \
	--retriever.bls-operator-state-retriever 0x9d4454B023096f34B160D6B654540c56A1F81688 \
	--retriever.eigenda-service-manager 0x67d269191c92Caf3cD7723F116c85e6E9bf55933 \
	--kzg.g1-path ../inabox/resources/kzg/g1.point \
	--kzg.g2-path ../inabox/resources/kzg/g2.point \
	--kzg.cache-path ../inabox/resources/kzg/SRSTables \
	--kzg.srs-order 3000 \
	--chain.rpc http://localhost:8545 \
	--chain.private-key=""