build: clean
	go mod tidy
	go build -o ./bin/ejections ./cmd

clean:
	rm -rf ./bin

run: build
	./bin/ejections --help
