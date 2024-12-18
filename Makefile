.PHONY: compile-el compile-dl clean protoc lint build unit-tests integration-tests-churner integration-tests-indexer integration-tests-inabox integration-tests-inabox-nochurner integration-tests-graph-indexer

ifeq ($(wildcard .git/*),)
$(warning semver disabled - building from release zip)
GITCOMMIT := ""
GITSHA := ""
GITDATE := ""
BRANCH := ""
SEMVER := $(shell basename $(CURDIR))
else
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITDATE := $(shell git log -1 --format=%cd --date=unix)
GITSHA := $(shell git rev-parse HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD | sed 's/[^[:alnum:]\.\_\-]/-/g')
SEMVER := $(shell docker run --rm --volume "$(PWD):/repo" gittools/gitversion:5.12.0 /repo -output json -showvariable SemVer)
ifeq ($(SEMVER), )
$(warning semver disabled - docker not installed)
SEMVER := "0.0.0"
endif
endif

RELEASE_TAG := $(or $(RELEASE_TAG),latest)

compile-el:
	cd contracts && ./compile.sh compile-el

compile-dl:
	cd contracts && ./compile.sh compile-dl

clean:
	./api/builder/clean.sh

# Builds the protobuf files inside a docker container.
protoc: clean
	./api/builder/protoc-docker.sh
	./api/builder/generate-docs.sh

# Builds the protobuf files locally (i.e. without docker).
protoc-local: clean
	./api/builder/protoc.sh

lint:
	golint -set_exit_status ./...
	go tool fix ./..
	golangci-lint run

build:
	cd operators/churner && make build
	cd disperser && make build
	cd node && make build
	cd retriever && make build
	cd tools/traffic && make build
	cd tools/kzgpad && make build
	cd relay && make build

dataapi-build:
	cd disperser && go build -o ./bin/dataapi ./cmd/dataapi

unit-tests:
	./test.sh

fuzz-tests:
	go test --fuzz=FuzzParseSignatureKMS -fuzztime=5m ./common

integration-tests-churner:
	go test -v ./churner/tests

integration-tests-indexer:
	go test -v ./core/indexer

integration-tests-node-plugin:
	go test -v ./node/plugin/tests

integration-tests-inabox:
	make build
	cd inabox && make run-e2e

integration-tests-inabox-nochurner:
	make build
	cd inabox && make run-e2e-nochurner

integration-tests-graph-indexer:
	make build
	go test -v ./core/thegraph

integration-tests-dataapi:
	make dataapi-build
	go test -v ./disperser/dataapi

docker-release-build:
	BUILD_TAG=${SEMVER} SEMVER=${SEMVER} GITDATE=${GITDATE} GIT_SHA=${GITSHA} GIT_SHORT_SHA=${GITCOMMIT} \
	docker buildx bake node-group-release ${PUSH_FLAG}

semver:
	echo "${SEMVER}"
