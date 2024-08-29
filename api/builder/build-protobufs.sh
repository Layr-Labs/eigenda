#!/usr/bin/env bash

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT="${SCRIPT_DIR}/../.."

docker container run \
  --rm \
  --mount "type=bind,source=${ROOT},target=/home/user/eigenda" \
  pbuf-compiler bash -c 'source .bashrc && cd eigenda && make protoc'
