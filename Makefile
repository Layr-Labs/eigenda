GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d--%H:%M:%S')
GIT_TAG := $(shell git describe --tags --always --dirty)

LDFLAGSSTRING +=-X main.Commit=$(GIT_COMMIT)
LDFLAGSSTRING +=-X main.Date=$(BUILD_TIME)
LDFLAGSSTRING +=-X main.Version=$(GIT_TAG)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

.PHONY: eigenda-proxy
eigenda-proxy:
	env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/eigenda-proxy ./cmd/server

.PHONY: docker-build
docker-build:
	# we only use this to build the docker image locally, so we give it the dev tag as a reminder
	@docker build -t ghcr.io/layr-labs/eigenda-proxy:dev .

run-memstore-server:
	./bin/eigenda-proxy --memstore.enabled --metrics.enabled

disperse-test-blob:
	curl -X POST -d my-blob-content http://127.0.0.1:3100/put/

clean:
	rm bin/eigenda-proxy

test-unit:
	go test ./... -parallel 4

# Local V1/V2 E2E tests, leveraging op-e2e framework. Also tests the standard client against the proxy.
test-e2e-local:
	INTEGRATION=true go test -timeout 1m ./e2e -parallel 4
	INTEGRATION_V2=true go test -timeout 1m ./e2e -parallel 4

# E2E tests against holesky testnet
# Holesky is currently broken after recent pectra hardfork.
# This test is thus flaky depending on whether the testnet producing blocks or not
# at the time it is run...
# In good cases it runs in ~20 mins, so we set a timeout of 30 mins.
# The test failing in CI is currently expected however, so expect to have to re-run it.
# See https://dora.holesky.ethpandaops.io/epochs for block production status.
test-e2e-holesky:
	# Add the -v flag to be able to observe logs as the run is happening on CI
	# given that this test takes >20 mins to run. Good to have early feedback when needed.
	TESTNET=true go test -v -timeout 30m ./e2e  -parallel 4

# E2E test which fuzzes the proxy client server integration and op client keccak256 with malformed inputs
test-e2e-fuzz:
	FUZZ=true go test ./e2e -fuzz -v -fuzztime=5m

.PHONY: lint
lint:
	@if ! command -v golangci-lint  &> /dev/null; \
	then \
    	echo "golangci-lint command could not be found...."; \
		echo "You can install via 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'"; \
		echo "or visit https://golangci-lint.run/welcome/install/ for other installation methods."; \
    	exit 1; \
	fi
	@golangci-lint run

.PHONY: format
format:
	# We also format line lengths. The length here should match that in the lll linter in .golangci.yml
	@if ! command -v golines  &> /dev/null; \
	then \
    	echo "golines command could not be found...."; \
		echo "You can install via 'go install github.com/segmentio/golines@latest'"; \
		echo "or visit https://github.com/segmentio/golines for other installation methods."; \
    	exit 1; \
	fi
	@go fmt ./...
	@golines --write-output --shorten-comments --max-len 120 .

go-gen-mocks:
	@echo "generating go mocks..."
	@GO111MODULE=on go generate --run "mockgen*" ./...

op-devnet-allocs:
	@echo "Generating devnet allocs..."
	@./scripts/op-devnet-allocs.sh

benchmark:
	go test -benchmem -run=^$ -bench . ./e2e -test.parallel 4

.PHONY: \
	clean \
	test
