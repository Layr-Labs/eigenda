SHELL := /bin/bash

# Clean the light node build files.
clean:
	rm -rf ./bin

# Build the light node.
build: clean
	go mod tidy
	go build -o ./bin/lnode ./main
	ls -alh ./bin

# Run the light node.
run: build
	./bin/lnode

# Delete the docker images for the light node.
clean-docker:
	docker image rm lnode || true # don't fail if the image doesn't exist
	docker builder prune -f

# Build the docker images for the light node.
build-docker:
	./docker/build.sh

# Run the docker image for the light node.
run-docker:
	docker run --rm lnode

# Open an interactive bash shell inside the light node docker container. Useful for debugging issues with the image.
debug-docker:
	docker container run --rm -it lnode bash
