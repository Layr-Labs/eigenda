ifeq ($(wildcard ../.git/*),)
$(warning semver disabled - building from release zip)
GITCOMMIT := ""
GITDATE := ""
SEMVER := $(shell basename $(CURDIR))
else
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITDATE := $(shell git log -1 --format=%cd --date=unix)
SEMVER := $(shell docker run --rm --volume "$(PWD)/../:/repo" gittools/gitversion:5.12.0 /repo -output json -showvariable SemVer)
ifeq ($(SEMVER), )
$(warning semver disabled - docker not installed)
SEMVER := "0.0.0"
endif
endif

RELEASE_TAG := $(or $(RELEASE_TAG),latest)

build: clean
	go mod tidy
	go build -o ./bin/node ./cmd

build-plugin: clean
	go mod tidy
	go build -o ./bin/node_plugin ./plugin/cmd

proto:
	cd .. && make protoc

clean:
	rm -rf ./bin

lint:
	golangci-lint run

test:
	go test -short ./...

docker: docker-node docker-plugin

docker-node:
	cd ../ && docker build --build-arg SEMVER=${SEMVER} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITDATE=${GITDATE} . -t opr-node:${SEMVER} -t opr-node:${RELEASE_TAG} -f node/cmd/Dockerfile

docker-plugin:
	cd ../ && docker build --build-arg SEMVER=${SEMVER} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITDATE=${GITDATE} . -t opr-nodeplugin:${SEMVER} -t opr-nodeplugin:${RELEASE_TAG} -f node/plugin/cmd/Dockerfile

docker-node-group:
	cd ../ && GIT_SHORT_SHA=${GITCOMMIT} \
	docker buildx bake node-group

semver:
	echo "${SEMVER}"

run: build
	set -a && \
	source .env && \
	NODE_LOG_PATH=$${NODE_LOG_PATH_HOST} \
	NODE_G1_PATH=$${NODE_G1_PATH_HOST} \
	NODE_G2_POWER_OF_2_PATH=$${NODE_G2_PATH_HOST} \
	NODE_DB_PATH=$${NODE_DB_PATH_HOST} \
	NODE_CACHE_PATH=$${NODE_CACHE_PATH_HOST} \
	NODE_ECDSA_KEY_FILE=$${NODE_ECDSA_KEY_FILE_HOST} \
	NODE_BLS_KEY_FILE=$${NODE_BLS_KEY_FILE_HOST} \
	./bin/node

run-update-socket: build-plugin
	set -a && \
	source .env && \
	NODE_LOG_PATH=$${NODE_LOG_PATH_HOST} \
	NODE_G1_PATH=$${NODE_G1_PATH_HOST} \
	NODE_G2_POWER_OF_2_PATH=$${NODE_G2_PATH_HOST} \
	NODE_DB_PATH=$${NODE_DB_PATH_HOST} \
	NODE_CACHE_PATH=$${NODE_CACHE_PATH_HOST} \
	NODE_ECDSA_KEY_FILE=$${NODE_ECDSA_KEY_FILE_HOST} \
	NODE_BLS_KEY_FILE=$${NODE_BLS_KEY_FILE_HOST} \
	NODE_SOCKET="$${NODE_HOSTNAME}:$${NODE_DISPERSAL_PORT};$${NODE_RETRIEVAL_PORT};$${NODE_V2_DISPERSAL_PORT};$${NODE_V2_RETRIEVAL_PORT}" \
	./bin/node_plugin --operation=update-socket

run-update-socket-v1: build-plugin
	set -a && \
	source .env && \
	NODE_LOG_PATH=$${NODE_LOG_PATH_HOST} \
	NODE_G1_PATH=$${NODE_G1_PATH_HOST} \
	NODE_G2_POWER_OF_2_PATH=$${NODE_G2_PATH_HOST} \
	NODE_DB_PATH=$${NODE_DB_PATH_HOST} \
	NODE_CACHE_PATH=$${NODE_CACHE_PATH_HOST} \
	NODE_ECDSA_KEY_FILE=$${NODE_ECDSA_KEY_FILE_HOST} \
	NODE_BLS_KEY_FILE=$${NODE_BLS_KEY_FILE_HOST} \
	NODE_SOCKET="$${NODE_HOSTNAME}:$${NODE_DISPERSAL_PORT};$${NODE_RETRIEVAL_PORT}" \
	./bin/node_plugin --operation=update-socket


run-list-quorums: build-plugin
	set -a && \
	source .env && \
	NODE_LOG_PATH=$${NODE_LOG_PATH_HOST} \
	NODE_G1_PATH=$${NODE_G1_PATH_HOST} \
	NODE_G2_POWER_OF_2_PATH=$${NODE_G2_PATH_HOST} \
	NODE_DB_PATH=$${NODE_DB_PATH_HOST} \
	NODE_CACHE_PATH=$${NODE_CACHE_PATH_HOST} \
	NODE_ECDSA_KEY_FILE=$${NODE_ECDSA_KEY_FILE_HOST} \
	NODE_BLS_KEY_FILE=$${NODE_BLS_KEY_FILE_HOST} \
	NODE_SOCKET="$${NODE_HOSTNAME}:$${NODE_DISPERSAL_PORT};$${NODE_RETRIEVAL_PORT};$${NODE_V2_DISPERSAL_PORT};$${NODE_V2_RETRIEVAL_PORT}" \
	./bin/node_plugin --operation=list-quorums

run-opt-out: build-plugin
	set -a && \
	source .env && \
	NODE_LOG_PATH=$${NODE_LOG_PATH_HOST} \
	NODE_G1_PATH=$${NODE_G1_PATH_HOST} \
	NODE_G2_POWER_OF_2_PATH=$${NODE_G2_PATH_HOST} \
	NODE_DB_PATH=$${NODE_DB_PATH_HOST} \
	NODE_CACHE_PATH=$${NODE_CACHE_PATH_HOST} \
	NODE_ECDSA_KEY_FILE=$${NODE_ECDSA_KEY_FILE_HOST} \
	NODE_BLS_KEY_FILE=$${NODE_BLS_KEY_FILE_HOST} \
	NODE_SOCKET="$${NODE_HOSTNAME}:$${NODE_DISPERSAL_PORT};$${NODE_RETRIEVAL_PORT};$${NODE_V2_DISPERSAL_PORT};$${NODE_V2_RETRIEVAL_PORT}" \
	./bin/node_plugin --operation=opt-out

run-opt-in: build-plugin
	set -a && \
	source .env && \
	NODE_LOG_PATH=$${NODE_LOG_PATH_HOST} \
	NODE_G1_PATH=$${NODE_G1_PATH_HOST} \
	NODE_G2_POWER_OF_2_PATH=$${NODE_G2_PATH_HOST} \
	NODE_DB_PATH=$${NODE_DB_PATH_HOST} \
	NODE_CACHE_PATH=$${NODE_CACHE_PATH_HOST} \
	NODE_ECDSA_KEY_FILE=$${NODE_ECDSA_KEY_FILE_HOST} \
	NODE_BLS_KEY_FILE=$${NODE_BLS_KEY_FILE_HOST} \
	NODE_SOCKET="$${NODE_HOSTNAME}:$${NODE_DISPERSAL_PORT};$${NODE_RETRIEVAL_PORT};$${NODE_V2_DISPERSAL_PORT};$${NODE_V2_RETRIEVAL_PORT}" \
	./bin/node_plugin --operation=opt-in
