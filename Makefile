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

E2EFUZZTEST = FUZZ=true go test ./e2e -fuzz -v -fuzztime=15m

.PHONY: eigenda-proxy
eigenda-proxy:
	env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/eigenda-proxy ./cmd/server

.PHONY: docker-build
docker-build:
	# we only use this to build the docker image locally, so we give it the dev tag as a reminder
	@docker build -t ghcr.io/layr-labs/eigenda-proxy:dev .

run-memstore-server:
	./bin/eigenda-proxy --memstore.enabled --eigenda.cert-verification-disabled --eigenda.eth-rpc http://localhost:8545 --eigenda.svc-manager-addr 0x123 --metrics.enabled

disperse-test-blob:
	curl -X POST -d my-blob-content http://127.0.0.1:3100/put/

clean:
	rm bin/eigenda-proxy

# Unit tests
test:
	go test ./... -parallel 4 

# E2E tests, leveraging op-e2e
e2e-test:
	INTEGRATION=true go test -timeout 1m ./e2e -parallel 4

# E2E test which fuzzes the proxy client server integration and op client keccak256 with malformed inputs
e2e-fuzz-test:
	$(E2EFUZZTEST)

# E2E tests against holesky testnet
holesky-test:
	TESTNET=true go test -timeout 50m ./e2e  -parallel 4

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

.PHONY: format
format:
	@go fmt ./...

go-gen-mocks:
	@echo "generating go mocks..."
	@GO111MODULE=on go generate --run "mockgen*" ./...

install-lint:
	@echo "Installing golangci-lint..."
	@sh -c $(GET_LINT_CMD)

op-devnet-allocs:
	@echo "Generating devnet allocs..."
	@./scripts/op-devnet-allocs.sh

benchmark:
	go test -benchmem -run=^$ -bench . ./e2e -test.parallel 4

.PHONY: \
	clean \
	test
