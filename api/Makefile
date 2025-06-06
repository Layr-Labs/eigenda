# Buf commands to lint/format proto files
# All of these commands are run by the github action in `.github/workflows/buf-proto.yaml`
proto-lint:
	buf lint

proto-format:
	buf format -w

# Builds the protobuf files inside a docker container.
protoc: clean proto-format
	./builder/protoc-docker.sh

# Builds the protobuf files locally (i.e. without docker).
protoc-local: clean
	./builder/protoc.sh

clean:
	./builder/clean.sh