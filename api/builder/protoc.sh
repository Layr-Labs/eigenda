#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail

# This script builds the eigenDA protobufs.

# The location where this script can be found.
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

# Build protobufs in the api/proto directory.

API_DIR="${SCRIPT_DIR}/.."
PROTO_DIR="${API_DIR}/proto"
GRPC_DIR="${API_DIR}/grpc"
mkdir -p "${GRPC_DIR}"

if [ $? -ne 0 ]; then
	exit 1
fi

PROTO_FILES=($(find "${PROTO_DIR}" -name '*.proto'))

protoc -I "${PROTO_DIR}" \
	--go_out="${GRPC_DIR}" \
	--go_opt=paths=source_relative \
	--go-grpc_out="${GRPC_DIR}" \
	--go-grpc_opt=paths=source_relative \
	${PROTO_FILES[@]}

if [ $? -ne 0 ]; then
	exit 1
fi
