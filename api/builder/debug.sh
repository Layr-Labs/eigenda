#!/usr/bin/env bash

# This is a handy little script for debugging the docker container. Attaches a bash shell to the container.

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT="${SCRIPT_DIR}/../.."

docker container run \
  --rm \
  --mount "type=bind,source=${ROOT},target=/home/user/eigenda" \
  -it \
  pbuf-compiler bash

