LINTER_VERSION = v1.52.1
LINTER_URL = https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
GET_LINT_CMD = "curl -sSfL $(LINTER_URL) | sh -s -- -b $(go env GOPATH)/bin $(LINTER_VERSION)"

GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d--%H:%M:%S')
GIT_TAG := $(shell git describe --tags --always --dirty)

LDFLAGSSTRING +=-X main.Commit=$(GIT_COMMIT)
LDFLAGSSTRING +=-X main.Date=$(BUILD_TIME)
LDFLAGSSTRING +=-X main.Version=$(GIT_TAG)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

E2ETEST = INTEGRATION=true go test -timeout 1m -v ./e2e -parallel 4 -deploy-config ../.devnet/devnetL1.json
HOLESKYTEST = TESTNET=true go test -timeout 50m -v ./e2e  -parallel 4 -deploy-config ../.devnet/devnetL1.json

.PHONY: eigenda-proxy
eigenda-proxy:
	env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/eigenda-proxy ./cmd/server

.PHONY: docker-build
docker-build:
	# we only use this to build the docker image locally, so we give it the dev tag as a reminder
	@docker build -t ghcr.io/layr-labs/eigenda-proxy:dev .

run-minio:
	docker run -p 4566:9000 -d -e "MINIO_ROOT_USER=minioadmin" -e "MINIO_ROOT_PASSWORD=minioadmin" --name minio minio/minio server /data

run-redis:
	docker run -p 9001:6379 -d --name redis redis

stop-minio:
	@if [ -n "$$(docker ps -q -f name=minio)" ]; then \
		docker stop minio && docker rm minio; \
	fi

stop-redis:
	@if [ -n "$$(docker ps -q -f name=redis)" ]; then \
		docker stop redis && docker rm redis; \
	fi

run-memstore-server:
	./bin/eigenda-proxy --memstore.enabled

disperse-test-blob:
	curl -X POST -d my-blob-content http://127.0.0.1:3100/put/

clean:
	rm bin/eigenda-proxy

test:
	go test -v ./... -parallel 4 

e2e-test: stop-minio stop-redis run-minio run-redis
	$(E2ETEST) && \
	make stop-minio && \
	make stop-redis

holesky-test: stop-minio stop-redis run-minio run-redis
	$(HOLESKYTEST) && \
	make stop-minio && \
	make stop-redis

.PHONY: lint
lint:
	@if ! test -f  &> /dev/null; \
	then \
    	echo "golangci-lint command could not be found...."; \
		echo "\nTo install, please run $(GET_LINT_CMD)"; \
		echo "\nBuild instructions can be found at: https://golangci-lint.run/usage/install/."; \
    	exit 1; \
	fi

	@golangci-lint run

go-gen-mocks:
	@echo "generating go mocks..."
	@GO111MODULE=on go generate --run "mockgen*" ./...

install-lint:
	@echo "Installing golangci-lint..."
	@sh -c $(GET_LINT_CMD)

gosec:
	@echo "Running security scan with gosec..."
	gosec ./...

submodules:
	git submodule update --init --recursive

op-devnet-allocs:
	@echo "Generating devnet allocs..."
	@./scripts/op-devnet-allocs.sh

.PHONY: \
	clean \
	test
