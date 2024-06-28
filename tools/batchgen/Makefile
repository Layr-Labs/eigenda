build: clean
	go mod tidy
	go build -o ./bin/batchgen ./cmd

clean:
	rm -rf ./bin

run: build
	./bin/batchgen --host localhost:32005

