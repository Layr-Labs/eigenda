# Make sure the help command stays first, so that it's printed by default when `make` is called without arguments
.PHONY: help compile-el compile-dl clean protoc lint build unit-tests integration-tests-churner integration-tests-indexer integration-tests-inabox integration-tests-inabox-nochurner integration-tests-graph-indexer

ifeq ($(wildcard .git/*),)
$(warning semver disabled - building from release zip)
GITCOMMIT := ""
GITDATE := ""
SEMVER := $(shell basename $(CURDIR))
else
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITDATE := $(shell git log -1 --format=%cd --date=unix)
SEMVER := $(shell docker run --rm --volume "$(PWD):/repo" gittools/gitversion:5.12.0 /repo -output json -showvariable SemVer)
ifeq ($(SEMVER), )
$(warning semver disabled - docker not installed)
SEMVER := "0.0.0"
endif
endif

RELEASE_TAG := $(or $(RELEASE_TAG),latest)

PROTOS := ./api/proto
PROTOS_DISPERSER := ./disperser/api/proto
PROTO_GEN := ./api/grpc
PROTO_GEN_DISPERSER_PATH = ./disperser/api/grpc

help: ## prints this help message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

compile-el: ## compiles Ethereum contracts for the el environment
	cd contracts && ./compile.sh compile-el

compile-dl: ## compiles distributed ledger contracts
	cd contracts && ./compile.sh compile-dl

clean: ## cleans up generated files
	find $(PROTO_GEN) -name "*.pb.go" -type f | xargs rm -rf
	mkdir -p $(PROTO_GEN)
	find $(PROTO_GEN_DISPERSER_PATH) -name "*.pb.go" -type f | xargs rm -rf
	mkdir -p $(PROTO_GEN_DISPERSER_PATH)

protoc: clean ## generates protobuf files
	protoc -I $(PROTOS) \
	--go_out=$(PROTO_GEN) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_GEN) \
	--go-grpc_opt=paths=source_relative \
	$(PROTOS)/**/*.proto
	# Generate Protobuf for sub directories of ./api/proto/disperser
	protoc -I $(PROTOS_DISPERSER) -I $(PROTOS) \
	--go_out=$(PROTO_GEN_DISPERSER_PATH) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_GEN_DISPERSER_PATH) \
	--go-grpc_opt=paths=source_relative \
	$(PROTOS_DISPERSER)/**/*.proto

lint: ## runs all linters
	golint -set_exit_status ./...
	go tool fix ./..
	golangci-lint run

build: ## builds all components
	cd operators/churner && make build
	cd disperser && make build
	cd node && make build
	cd retriever && make build
	cd tools/traffic && make build
	cd tools/kzgpad && make build

dataapi-build: ## builds dataapi cli
	cd disperser && go build -o ./bin/dataapi ./cmd/dataapi

unit-tests: ## runs unit tests
	./test.sh

integration-tests-churner: ## runs integration tests for churner
	go test -v ./churner/tests

integration-tests-indexer: ## runs integration tests for indexer
	go test -v ./core/indexer

integration-tests-node-plugin: ## runs integration tests for the node plugin
	go test -v ./node/plugin/tests

integration-tests-inabox: ## runs all integration tests in a boxed environment
	make build
	cd inabox && make run-e2e

integration-tests-inabox-nochurner: ## runs all integration tests in a boxed environment without churner
	make build
	cd inabox && make run-e2e-nochurner

integration-tests-graph-indexer: ## runs integration tests for the graph indexer
	make build
	go test -v ./core/thegraph

integration-tests-dataapi: ## runs integration tests for the dataapi
	make dataapi-build
	go test -v ./disperser/dataapi

docker-release-build: ## builds docker images for release
	RELEASE_TAG=${SEMVER} docker compose -f docker-compose-release.yaml build --build-arg SEMVER=${SEMVER} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITDATE=${GITDATE} ${PUSH_FLAG}

semver: ## displays the current semantic version
	echo "${SEMVER}"
