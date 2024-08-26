#!/usr/bin/env bash

# Starts the container with an interactive bash shell. This is useful for debugging the container.

# Do setup for the data directory. This is a directory where data that needs
# to persist in-between container runs is stored.
./docker/setup-data-dir.sh

docker container run --rm -it lnode bash
