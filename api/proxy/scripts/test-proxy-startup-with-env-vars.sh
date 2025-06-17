#!/bin/bash
set -e  # Exit on any error

##### This script is meant to be run in ci #####
# It tests that the env vars defined in the specified environment file are correct.
# It starts the eigenda-proxy with those env vars, waits 5 seconds, and then kills the proxy.
# If any deprecated flags are still being used in the specified environment file, the script will fail.

# Check if an environment file is provided
if [ $# -eq 0 ]; then
    echo "Error: No environment file specified"
    echo "Usage: $0 <environment_file_path>"
    exit 1
fi

ENV_FILE=$1

# Check if the environment file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: Environment file $ENV_FILE does not exist"
    echo "Current working directory: $(pwd)"
    echo "Files in current directory:"
    ls -la

    exit 1
fi

echo "Using environment file: $ENV_FILE"

# build the eigenda-proxy binary
make

# Start the eigenda-proxy with the env vars defined in the specified environment file
set -a; source "$ENV_FILE"; set +a
./bin/eigenda-proxy &
PID=$!

# Ensure we kill the process on script exit
trap "kill $PID" EXIT

# Wait 10 seconds for startup to happen (actual startup takes ~5 seconds with max blob length=1MiB)
echo "sleeping 10 seconds to let the proxy start up"
sleep 10

echo "Pinging the proxy's health endpoint"
curl 'http://localhost:3100/health'

# Script will automatically kill process due to trap
# If eigenda-proxy has failed, trap will error out and script will exit with an error code