#!/usr/bin/env bash

# Sets up the data directory for the light node container.

# Default arguments.
source docker/default-args.sh
# Local overrides for arguments.
source docker/args.sh

# Create the data directory if it doesn't exist.
mkdir -p $DATA_PATH

echo "Using data directory $DATA_PATH"
