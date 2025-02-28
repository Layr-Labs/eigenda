#!/usr/bin/env bash

# Starts the container and runs the light node.

# Do setup for the data directory. This is a directory where data that needs
# to persist in-between container runs is stored.
source ./docker/setup-data-dir.sh

docker run \
  --rm \
  --mount "type=bind,source=${DATA_PATH},target=/home/lnode/data" \
  lnode
