SHELL := /bin/bash

clean:
	rm -rf ./bin

build: clean
	go mod tidy
	go build -o ./bin/litt ./cli

debug-build: clean
	go mod tidy
	go build -gcflags "all=-N -l"  -o ./bin/litt ./cli
