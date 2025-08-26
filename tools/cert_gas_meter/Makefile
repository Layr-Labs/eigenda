build: clean
	go mod tidy
	go build -o ./bin/cert-gas-meter ./cmd

clean:
	rm -rf ./bin

lint: 
	golangci-lint run ./...

run: build 
	./bin/quorumscan --help
