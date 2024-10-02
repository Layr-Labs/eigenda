#!/usr/bin/env bash

# Starts the container with an interactive bash shell. This is useful for debugging the container.

# Do setup for the data directory. This is a directory where data that needs
# to persist in-between container runs is stored.
source ./docker/setup-data-dir.sh

docker container run \
  --mount "type=bind,source=${DATA_PATH},target=/home/lnode/data" \
  --rm \
  -it \
  lnode bash
