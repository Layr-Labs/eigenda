clean:
	rm -rf ./bin

build: clean
	# cd .. && make protoc
	go mod tidy
	go build -o ./bin/server ./cmd