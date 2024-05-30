build: clean
	go mod tidy
	go build -o ./bin/batchgen ./cmd

clean:
	rm -rf ./bin

proto: 
	protoc --go_out=. store_chunks.proto

run: build
	./bin/batchgen 10000 localhost:32005 localhost:32005
