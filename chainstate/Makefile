.PHONY: build clean test

build:
	go build -o ./bin/chainstate-indexer ./main

clean:
	rm -rf ./bin

test:
	go test ./... -v
