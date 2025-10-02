build:
	go build -o ./bin/server ./cmd

clean:
	rm -rf ./bin


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
	--kzg.g1-path ../resources/srs/g1.point \
	--kzg.g2-path ../resources/srs/g2.point \
	--kzg.cache-path ../resources/srs/SRSTables \
	--kzg.srs-order 3000 \
	--chain.rpc http://localhost:8545 \
	--chain.private-key=""