build: clean
	go mod tidy
	go build -o ./bin/opscan ./cmd

clean:
	rm -rf ./bin

run: build 
	./bin/opscan --help
