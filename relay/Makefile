SHELL := /bin/bash

# Build the light node.
build:
	go build -o ./bin/relay ./cmd

# Clean the light node build files.
clean:
	rm -rf ./bin

# Run the light node.
run: build
	./bin/relay
