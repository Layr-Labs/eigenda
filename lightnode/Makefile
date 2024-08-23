# Clean the light node build files.
clean:
	rm -rf ./bin

# Build the light node.
build: clean
	go mod tidy
	go build -o ./bin/lnode ./main

# Run the light node.
run: build
	./bin/lnode

# Delete the docker images for the light node.
clean-docker:
	./docker/clean.sh

# Build the docker images for the light node.
build-docker:
	./docker/build.sh
