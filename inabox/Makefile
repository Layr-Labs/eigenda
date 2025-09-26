dt := $(shell date '+%YY-%mM-%dD-%HH-%MM-%SS')
.PHONY: new-anvil deploy-all stop-infra run-e2e clean

new-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil.yaml "testdata/$(dt)/config.yaml"

# Start infra will start anvil, a graph node, and aws localstack services (s3 and dynamodb).
# It will also deploy the contracts and the subgraphs onto the graph node.
# After this you can run `make start-services`.
start-infra: new-anvil
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
start-services:
	./bin.sh start-detached

stop-services:
	./bin.sh stop-detached

run-e2e:
	go test ./tests -v -config=../templates/testconfig-anvil.yaml

clean:
	rm -rf testdata/*
