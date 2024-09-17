#!/usr/bin/env bash

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "${SCRIPT_DIR}"
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "${REPO_ROOT}"

docker build \
  --build-arg="REPO_ROOT=${REPO_ROOT}" \
  -f "lightnode/docker/Dockerfile" \
  --tag lnode:latest \
  .
