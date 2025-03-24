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

# Runs all tests, excluding e2e
test-unit:
	go test `go list ./... | grep -v ./e2e` -parallel 4

# E2E tests using local memstore, leveraging op-e2e framework. Also tests the standard client against the proxy.
test-e2e-local:
	# Add the -v flag to observe logs as the run is happening on CI, given that this test takes ~10 minutes to run.
	# Good to have early feedback when needed.
	BACKEND=memstore go test -v -timeout 20m ./e2e -parallel 4

# E2E tests using holesky testnet backend, leveraging op-e2e framework. Also tests the standard client against the proxy.
# If holesky tests are failing, consider checking https://dora.holesky.ethpandaops.io/epochs for block production status.
test-e2e-testnet:
	# Add the -v flag to observe logs as the run is happening on CI, given that this test takes ~20 minutes to run.
	# Good to have early feedback when needed.
	BACKEND=testnet go test -v -timeout 30m ./e2e -parallel 4

test-e2e-preprod:
	# Add the -v flag to observe logs as the run is happening on CI, given that this test takes ~20 minutes to run.
	# Good to have early feedback when needed.
	BACKEND=preprod go test -v -timeout 30m ./e2e -parallel 4

# Very simple fuzzer which generates random bytes arrays and sends them to the proxy using the standard client.
# To clean the cached corpus, run `go clean -fuzzcache` before running this.
test-fuzz:
	go test ./fuzz -fuzz=FuzzProxyClientServerV1 -fuzztime=1m
	go test ./fuzz -fuzz=FuzzProxyClientServerV2 -fuzztime=1m

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

## calls --help on binary and routes output to file while ignoring dynamic fields specific
## to indivdual builds (e.g, version)
gen-static-help-output: eigenda-proxy
	@echo "Storing binary output to docs/help_out.txt"
	@./bin/eigenda-proxy --help | sed '/^VERSION:/ {N;d;}' > docs/help_out.txt

go-gen-mocks:
	@echo "generating go mocks..."
	@GO111MODULE=on go generate --run "mockgen*" ./...

op-devnet-allocs:
	@echo "Generating devnet allocs..."
	@./scripts/op-devnet-allocs.sh

benchmark:
	go test -benchmem -run=^$ -bench . ./benchmark -test.parallel 4

.PHONY: \
	clean \
	test
