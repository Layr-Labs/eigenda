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
	--retriever.grpc-port 32099 \
	--retriever.timeout 10s \
	--retriever.bls-operator-state-retriever 0x5f3f1dBD7B74C6B46e8c44f98792A1dAf8d69154 \
	--retriever.eigenda-service-manager 0x851356ae760d987E095750cCeb3bC6014560891C \
	--kzg.g1-path ../inabox/resources/kzg/g1.point \
	--kzg.g2-path ../inabox/resources/kzg/g2.point \
	--kzg.cache-path ../inabox/resources/kzg/SRSTables \
	--kzg.srs-load 3000 \
	--kzg.srs-order 3000 \
	--chain.rpc http://localhost:8545 \
	--chain.private-key="" \
	--thegraph.endpoint localhost:8000