build: clean
	go mod tidy
	go build -o ./bin/batchgen ./cmd

clean:
	rm -rf ./bin

proto: 
	protoc --go_out=. store_chunks.proto

run: build
	./bin/batchgen localhost:32005
