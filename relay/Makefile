SHELL := /bin/bash

# Clean the light node build files.
clean:
	rm -rf ./bin

# Build the light node.
build: clean
	go mod tidy
	go build -o ./bin/relay ./cmd

build-wasm: clean
	go mod tidy
	GOOS=js GOARCH=wasm go build -o ./bin/relay.wasm ./cmd

# Run the light node.
run: build
	./bin/relay
