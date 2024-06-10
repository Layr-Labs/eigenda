APP_NAME = eigenda-proxy
LINTER_VERSION = v1.52.1
LINTER_URL = https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
GET_LINT_CMD = "curl -sSfL $(LINTER_URL) | sh -s -- -b $(go env GOPATH)/bin $(LINTER_VERSION)"

GITCOMMIT ?= $(shell git rev-parse HEAD)
GITDATE ?= $(shell git show -s --format='%ct')
VERSION := v0.0.0

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.Version=$(VERSION)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

.PHONY: eigenda-proxy
eigenda-proxy:
	env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/eigenda-proxy ./cmd/server

.PHONY: docker-build
docker-build:
	@docker build -t $(APP_NAME) .

run-server:
	./bin/eigenda-proxy

clean:
	rm bin/eigenda-proxy

test: optimism-test
	go test -v ./...

optimism-test:
	OPTIMISM=true go test -timeout 50m -v ./e2e/... -deploy-config ../.devnet/devnetL1.json

holesky-test:
	TESTNET=true go test -timeout 50m -v ./e2e
	
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

gosec:
	@echo "Running security scan with gosec..."
	gosec ./...

submodules:
	git submodule update --init --recursive

srs:
	if ! test -f /operator-setup/resources/g1.point; then \
		cd operator-setup && ./srs_setup.sh; \
	fi

op-devnet-allocs: submodules
	@echo "Generating devnet allocs..."
	@./scripts/op-devnet-allocs.sh

.PHONY: \
	clean \
	test
