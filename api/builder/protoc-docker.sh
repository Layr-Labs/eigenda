#!/usr/bin/env bash

# This script builds the eigenDA protobufs. It does this by running protoc.sh inside of the pbuf-compiler container.

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT="${SCRIPT_DIR}/../.."

if [ -z "$(docker images -q pbuf-compiler:latest 2> /dev/null)" ]; then
  echo "Docker image pbuf-compiler:latest does not exist. Building it now..."
  "${SCRIPT_DIR}"/build-docker.sh
fi

if [ $? -ne 0 ]; then
  exit 1
fi

docker container run \
  --rm \
  --mount "type=bind,source=${ROOT},target=/home/user/eigenda" \
  pbuf-compiler bash -c "source ~/.bashrc && eigenda/api/builder/protoc.sh"

if [ $? -ne 0 ]; then
  exit 1
fi