GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d--%H:%M:%S')
GIT_TAG := $(shell git describe --tags --always --dirty)

LDFLAGSSTRING +=-X main.Commit=$(GIT_COMMIT)
LDFLAGSSTRING +=-X main.Date=$(BUILD_TIME)
LDFLAGSSTRING +=-X main.Version=$(GIT_TAG)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

build:
	env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/eigenda-proxy ./cmd/server

clean:
	rm bin/eigenda-proxy

docker-build:
	# we only use this to build the docker image locally, so we give it the dev tag as a reminder
	@docker build -t ghcr.io/layr-labs/eigenda-proxy:dev .

run-memstore-server: build
	./bin/eigenda-proxy --memstore.enabled --metrics.enabled

disperse-test-blob:
	curl -X POST -d my-blob-content http://127.0.0.1:3100/put/

# Runs all tests, excluding e2e
test-unit:
	gotestsum --format pkgname-and-test-fails -- `go list ./... | grep -v ./e2e` -parallel 4

# E2E tests using local memstore, leveraging op-e2e framework. Also tests the standard client against the proxy.
test-e2e-local:
	BACKEND=memstore gotestsum --format testname -- -v -timeout 10m ./e2e -parallel 8

# E2E tests using holesky testnet backend, leveraging op-e2e framework. Also tests the standard client against the proxy.
# If holesky tests are failing, consider checking https://dora.holesky.ethpandaops.io/epochs for block production status.
test-e2e-testnet:
	BACKEND=testnet gotestsum --format testname -- -v -timeout 20m ./e2e -parallel 32

## Equivalent to `test-e2e-testnet`, but against preprod instead of testnet
test-e2e-preprod:
	BACKEND=preprod gotestsum --format testname -- -v -timeout 20m ./e2e -parallel 32

# Very simple fuzzer which generates random bytes arrays and sends them to the proxy using the standard client.
# To clean the cached corpus, run `go clean -fuzzcache` before running this.
test-fuzz:
	go test ./fuzz -fuzz=FuzzProxyClientServerV1 -fuzztime=1m
	go test ./fuzz -fuzz=FuzzProxyClientServerV2 -fuzztime=1m

.PHONY: lint
lint:
	golangci-lint run

.PHONY: format
format:
	# We also format line lengths. The length here should match that in the lll linter in .golangci.yml
	go fmt ./...
	golines --write-output --shorten-comments --max-len 120 .

## calls --help on binary and routes output to file while ignoring dynamic fields specific
## to indivdual builds (e.g, version)
gen-static-help-output: build
	@echo "Storing binary output to docs/help_out.txt"
	@./bin/eigenda-proxy --help | sed '/^VERSION:/ {N;d;}' > docs/help_out.txt

mocks:
	@echo "generating go mocks..."
	@GO111MODULE=on go generate --run "mockgen*" ./...

op-devnet-allocs:
	@echo "Generating devnet allocs..."
	@./scripts/op-devnet-allocs.sh

benchmark:
	go test -benchmem -run=^$ -bench . ./benchmark -test.parallel 4

deps:
	mise install

.PHONY: build clean docker-build test lint format benchmark deps mocks
