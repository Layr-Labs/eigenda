build: clean
	go mod tidy
	go build -o ./bin/semverscan ./cmd

clean:
	rm -rf ./bin

run: build 
	./bin/semverscan --help
