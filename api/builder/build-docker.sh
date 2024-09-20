#!/usr/bin/env bash

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# Add the --no-cache flag to force a rebuild.
# Add the --progress=plain flag to show verbose output during the build.

docker build \
  -f "${SCRIPT_DIR}/Dockerfile" \
  --tag pbuf-compiler:latest \
  .