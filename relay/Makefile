SHELL := /bin/bash

# Clean the light node build files.
clean:
	rm -rf ./bin

# Build the light node.
build: clean
	go mod tidy
	go build -tags=viper_bind_struct -o ./bin/relay ./main

# Run the light node.
run: build
	./bin/relay
