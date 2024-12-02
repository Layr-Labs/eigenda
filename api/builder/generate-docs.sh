#!/usr/bin/env bash

# This script generates protobuf documentation.

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

API_DIR="${SCRIPT_DIR}/.."
PROTO_DIR="${API_DIR}/proto"
DOCS_DIR="${API_DIR}/docs"

# Function to get the relative path of file in argument 1 with respect directory in argument 2.
# Doesn't use the convenient 'realpath --relative-to' because it's not available on macOS.
relativePath() {
 python3 -c 'import os.path, sys; print(os.path.relpath(sys.argv[1],sys.argv[2]))' "${1}" "${2}"
}

# Find all .proto files.
PROTO_FILES=( $(find "${PROTO_DIR}" -name '*.proto') )

# Make the proto files relative to the proto directory.
for i in "${!PROTO_FILES[@]}"; do
    PROTO_FILES[$i]=$(relativePath "${PROTO_FILES[$i]}" "${PROTO_DIR}")
done

# Sort the proto files alphabetically. Required for deterministic output.
IFS=$'\n' PROTO_FILES=($(sort <<<"${PROTO_FILES[*]}")); unset IFS

# Generate HTML doc
docker run --rm \
  -v "${DOCS_DIR}":/out \
  -v "${PROTO_DIR}":/protos \
  pseudomuto/protoc-gen-doc \
  "${PROTO_FILES[@]}" \
  --doc_opt=html,eigenda-protos.html

if [ $? -ne 0 ]; then
  exit 1
fi

# Generate markdown doc
docker run --rm \
  -v "${DOCS_DIR}":/out \
  -v "${PROTO_DIR}":/protos \
  pseudomuto/protoc-gen-doc \
  "${PROTO_FILES[@]}" \
  --doc_opt=markdown,eigenda-protos.md

if [ $? -ne 0 ]; then
  exit 1
fi
