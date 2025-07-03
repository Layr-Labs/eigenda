# Find all Go files under load/
# The build command will rebuild the binary if any Go files change.
GO_SOURCES_LOAD := $(shell find load -name "*.go" -type f)

build: $(GO_SOURCES_LOAD)
	go build -o bin/load load/main/load_main.go

# Makefile doesn't allow forwarding of arguments, so we use ARGS. Call this as:
# make generate-load ARGS="config/environment/preprod.json config/load/100kb_s-1mb-3x.json"
generate-load: build
	./bin/load $(ARGS)

clean:
	rm -rf bin 2>/dev/null || true

test:
	cd correctness && go test
