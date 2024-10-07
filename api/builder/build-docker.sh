#!/usr/bin/env bash

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

ARCH=$(uname -m)
if [ "${ARCH}" == "arm64" ]; then
  PROTOC_URL='https://github.com/protocolbuffers/protobuf/releases/download/v23.4/protoc-23.4-linux-aarch_64.zip'
elif [ "${ARCH}" == "x86_64" ]; then
  PROTOC_URL='https://github.com/protocolbuffers/protobuf/releases/download/v23.4/protoc-23.4-linux-x86_64.zip'
else
  echo "Unsupported architecture: ${ARCH}"
  exit 1
fi

# Add the --no-cache flag to force a rebuild.
# Add the --progress=plain flag to show verbose output during the build.

docker build \
  -f "${SCRIPT_DIR}/Dockerfile" \
  --tag pbuf-compiler:latest \
  --build-arg PROTOC_URL="${PROTOC_URL}" \
  --build-arg UID=$(id -u) \
  .

if [ $? -ne 0 ]; then
  exit 1
fi