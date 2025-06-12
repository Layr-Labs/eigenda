SHELL := /bin/bash

clean:
	rm -rf ./bin

build: clean
	go mod tidy
	go build -o ./bin/litt ./cli
