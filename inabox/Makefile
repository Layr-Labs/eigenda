dt := $(shell date '+%YY-%mM-%dD-%HH-%MM-%SS')
.PHONY: run-e2e-test start-inabox stop-inabox new-testdata-dir clean start-infra stop-infra start-services stop-services

# Starts a short-lived inabox devnet, and runs integration/e2e tests against it.
run-e2e-test:
	go test ./tests -v -config=../templates/testconfig-anvil.yaml

# Starts a long-lived inabox local devnet.
start-inabox: start-infra start-services
	@echo "If everything went well, inabox should now be running."
	@echo
	@echo "export these variables:"
	@echo "ETH_RPC_URL=http://localhost:8545"
	@echo "EIGENDA_DIRECTORY_ADDR=$(shell cat ../contracts/script/output/eigenda_deploy_output.json | jq -r .eigenDADirectory)"
	@echo "EIGENDA_DISPERSER_V1_URL=localhost:32004"
	@echo "EIGENDA_DISPERSER_V2_URL=localhost:32005"
	@echo
	@echo "You can query other contract addresses from the directory:"
	@echo 'cast call $$EIGENDA_DIRECTORY_ADDR "getAddress(string)(address)" "CERT_VERIFIER_ROUTER"'
	@echo "You can query the disperser v2 by using:"
	@echo 'grpcurl -plaintext $$EIGENDA_DISPERSER_V2_URL list'
	@echo
	@echo "Infra components (anvil, graph, aws localstack) are managed by docker."
	@echo "Run 'docker ps' to see and manage them."
	@echo
	@echo "EigenDA services (disperser, validators, etc) are ran as local processes."
	@echo "Their config is available under `pwd`/testdata/_latest/envs"
	@echo "Their logs are available under `pwd`/testdata/_latest/logs"

stop-inabox: stop-services stop-infra

############################################################################
# Below section are lower level commands. Most people will not need these. #
############################################################################

# Every inabox run (whether as a local devnet or as part of integration tests)
# uses a directory under inabox/testdata/ to store its configs, logs, db state,
# pid file, etc.
new-testdata-dir:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil.yaml "testdata/$(dt)/config.yaml"
# We use _latest so that it appears before the other directories under testdata/
# because there are some hardcoded assumptions in the deploy process that read the "last"
# directory alphabetically.
# TODO(samlaf): we should probably move to using the _latest directory instead.
	ln -sfn $(shell pwd)/testdata/$(dt) testdata/_latest

clean:
	rm -rf testdata/*

# Start infra will start anvil, a graph node, and aws localstack services (s3 and dynamodb)
# docker containers. It will also deploy the contracts and the subgraphs onto the graph node.
# After this you can run `make start-services`.
start-infra:
	go run ./deploy/cmd -localstack-port 4570

# Using filter based on ancestor doesn't seem to work with grep expressions,
# so we need to match the exact version that is spun up in the golang code.
# If we ever change the version and forget to update here we'll leave some dangling containers.
# TODO(samlaf): we prob should start all containers with a inabox specific label, so that we
# can instead filter and kill all docker containers that contain a specific label.
stop-infra:
# Stop anvil
	docker ps -q --filter "ancestor=ghcr.io/foundry-rs/foundry" | xargs -r docker stop 2>/dev/null || true
# Stop localstack based on container name
	docker ps -q --filter "ancestor=localstack/localstack:4.7.0" | xargs -r docker stop 2>/dev/null || true
# Stop graph node stack based on container names
	docker ps -q --filter "ancestor=graphprotocol/graph-node:v0.35.0" | xargs -r docker stop 2>/dev/null || true
	docker ps -q --filter "ancestor=ipfs/kubo:v0.24.0" | xargs -r docker stop 2>/dev/null || true
	docker ps -q --filter "ancestor=postgres:13" | xargs -r docker stop 2>/dev/null || true

# Starts all the services: disperser, validators, churner, etc.
# These are run as local processes, not containerized. Their pids are stored inside the testdata dir
# and they can be killed using `make stop-services`. Their configs and logs are also available
# inside the testdata dir.
start-services:
	./bin.sh start-detached

stop-services:
	./bin.sh stop-detached
