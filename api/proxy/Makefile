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
	rm -rf bin/eigenda-proxy

docker-build:
	# we only use this to build the docker image locally, so we give it the dev tag as a reminder
	cd ../.. && SEMVER=$(GIT_TAG) GIT_SHORT_SHA=$(GIT_COMMIT) GITDATE=$(BUILD_TIME) docker buildx bake proxy --load

run-memstore-server: build
	./bin/eigenda-proxy --memstore.enabled --metrics.enabled  --storage.backends-to-enable v2

disperse-test-blob:
	curl -X POST -d my-blob-content http://127.0.0.1:3100/put/ | xxd -p | tr -d '\n'

# Runs all tests, excluding e2e
test-unit:
	gotestsum --format pkgname-and-test-fails -- -short -parallel 4 ./...

# TODO: Add support for E2E network tests with a `hoodi` testnet backend.

# E2E tests using local memstore. Also tests the standard client against the proxy.
test-e2e-local:
	BACKEND=memstore gotestsum --format testname -- -v -timeout 10m ./test/e2e -parallel 8

# E2E tests using hoodi testnet backend.
test-e2e-hoodi:
	BACKEND=hoodi-testnet gotestsum --format testname -- -v -timeout 20m ./test/e2e -parallel 32

# E2E tests using hoodi preprod backend. this is expected to fail since retrieval ingress is turned off.
# this is useful for testing CertVerifier deployments in the hood-preprod and serves utility as a validation
# tool which is why it will never be ran in monorepo CI.
test-e2e-hoodi-preprod:
	BACKEND=hoodi-preprod gotestsum --format testname -- -v -timeout 20m ./test/e2e -parallel 32

# E2E tests using sepolia testnet backend. Also tests the standard client against the proxy.
# If sepolia tests are failing, consider checking https://sepolia.etherscan.io/ for block production status.
test-e2e-sepolia:
	BACKEND=sepolia gotestsum --format testname -- -v -timeout 20m ./test/e2e -parallel 32

# To clean the cached corpus, run `go clean -fuzzcache` before running this.
test-fuzz:
	go test ./test/fuzz -fuzz=FuzzProxyClientServerV1 -fuzztime=1m
	go test ./test/fuzz -fuzz=FuzzProxyClientServerV2 -fuzztime=1m

benchmark:
	go test -benchmem -run=^$ -bench . ./test/benchmark -test.parallel 4


.PHONY: format
format:
	# We also format line lengths. The length here should match that in the lll linter in .golangci.yml
	go fmt ./...
	golines --write-output --shorten-comments --max-len 120 .

## calls --help on binary and routes output to file while ignoring dynamic fields specific
## to indivdual builds (e.g, version)
gen-static-help-output: build
	@echo "Storing proxy help output to docs/help_out.txt"
# removes the VERSION line which makes the output non-deterministic (changes with each commit)
	@./bin/eigenda-proxy --help | sed '/^VERSION:/ {N;d;}' > docs/help_out.txt
	@echo "Storing proxy metrics output to docs/metrics_out.txt"
	@./bin/eigenda-proxy doc metrics > docs/metrics_out.txt

mocks:
	@echo "generating go mocks..."
	@GO111MODULE=on go generate --run "mockgen*" ./...

op-devnet-allocs:
	@echo "Generating devnet allocs..."
	@./scripts/op-devnet-allocs.sh

deps:
	mise install

.PHONY: build clean docker-build test format benchmark deps mocks
