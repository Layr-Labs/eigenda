SHELL := /bin/bash

clean:
	rm -rf ./bin

build:
	go mod tidy
	go build -o ./bin/benchmark ./benchmark/cmd

