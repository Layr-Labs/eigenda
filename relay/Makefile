SHELL := /bin/bash

# Clean the light node build files.
clean:
	rm -rf ./bin

# Build the light node.
build: clean
	go mod tidy
	go build -o ./bin/relay ./cmd

# Run the light node.
run: build
	./bin/relay
