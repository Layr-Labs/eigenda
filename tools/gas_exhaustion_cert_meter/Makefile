build: clean
	go mod tidy
	go build -o ./bin/gas-exhaustion-cert-meter ./cmd

clean:
	rm -rf ./bin

lint: 
	golangci-lint run ./...

run: build
	./bin/gas-exhaustion-cert-meter --help
