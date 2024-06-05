build: clean
	go mod tidy
	go build -o ./bin/batchgen ./cmd

clean:
	rm -rf ./bin

run: build
	./bin/batchgen 1000 10 localhost:32005

