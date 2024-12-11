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

# Generate unified HTML doc
echo "Generating unified HTML documentation..."
docker run --rm \
  -v "${DOCS_DIR}":/out \
  -v "${PROTO_DIR}":/protos \
  pseudomuto/protoc-gen-doc \
  "${PROTO_FILES[@]}" \
  --doc_opt=html,eigenda-protos.html 2>/dev/null

if [ $? -ne 0 ]; then
  echo "Failed to generate unified HTML documentation."
  exit 1
fi

# Generate unified markdown doc
echo "Generating unified markdown documentation..."
docker run --rm \
  -v "${DOCS_DIR}":/out \
  -v "${PROTO_DIR}":/protos \
  pseudomuto/protoc-gen-doc \
  "${PROTO_FILES[@]}" \
  --doc_opt=markdown,eigenda-protos.md 2>/dev/null

if [ $? -ne 0 ]; then
  echo "Failed to generate unified markdown documentation."
  exit 1
fi

# Generate individual markdown/HTML docs
for PROTO_FILE in "${PROTO_FILES[@]}"; do
  PROTO_NAME=$(basename "${PROTO_FILE}" .proto)

  echo "Generating markdown documentation for ${PROTO_NAME}..."
  docker run --rm \
    -v "${DOCS_DIR}":/out \
    -v "${PROTO_DIR}":/protos \
    pseudomuto/protoc-gen-doc \
    "${PROTO_FILE}" \
    --doc_opt=markdown,"${PROTO_NAME}.md" 2>/dev/null

  if [ $? -ne 0 ]; then
    echo "Failed to generate documentation for ${PROTO_NAME}."
    exit 1
  fi

  echo "Generating HTML documentation for ${PROTO_NAME}..."
  docker run --rm \
    -v "${DOCS_DIR}":/out \
    -v "${PROTO_DIR}":/protos \
    pseudomuto/protoc-gen-doc \
    "${PROTO_FILE}" \
    --doc_opt=html,"${PROTO_NAME}.html" 2>/dev/null

  if [ $? -ne 0 ]; then
    echo "Failed to generate documentation for ${PROTO_NAME}."
    exit 1
  fi
done