build: clean
	go mod tidy
	go build -o ./bin/relayload ./cmd

clean:
	rm -rf ./bin

lint: 
	golangci-lint run ./...

run: build 
	./bin/relayload --help
