build: clean
	mkdir -p ./bin
	go build -o ./bin/integration_utils ./cmd

clean:
	rm -rf ./bin

lint: 
	golangci-lint run ./...

run: build
	./bin/integration_utils --help