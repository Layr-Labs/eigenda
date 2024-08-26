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
# Add the --no-cache flag to force a rebuild.
# Add the --progress=plain flag to show verbose output during the build.
build-docker:
	docker build \
		--progress=plain \
		--build-arg='GIT_URL=https://github.com/cody-littley/eigenda.git' \
		--build-arg='BRANCH_OR_COMMIT=lightnode-docker' \
		-f docker/Dockerfile \
		-t lnode .

# Run the docker image for the light node.
run-docker:
	docker run --rm lnode
