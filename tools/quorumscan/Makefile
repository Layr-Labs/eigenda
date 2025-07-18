build:
	go mod tidy
	go build -o ./bin/quorumscan ./cmd

clean:
	rm -rf ./bin

lint: 
	golangci-lint run ./...

run: build 
	./bin/quorumscan --help
