#!/bin/bash
set -e  # Exit on any error

##### This script is meant to be run in ci #####
# It tests that the env vars defined in .env.example.holesky are correct.
# It starts the eigenda-proxy with those env vars, waits 5 seconds, and then kills the proxy.
# If any deprecated flags are still being used in .env.example.holesky, the script will fail.

# build the eigenda-proxy binary
make

# Start the eigenda-proxy with the env vars defined in .env.example.holesky
set -a; source .env.example.holesky; set +a
./bin/eigenda-proxy &
PID=$!

# Ensure we kill the process on script exit
trap "kill $PID" EXIT

# Wait 5 seconds for startup to happen
echo "sleeping 5 seconds to let the proxy start up"
sleep 5

echo "Pinging the proxy's health endpoint"
curl 'http://localhost:3100/health'

# Script will automatically kill process due to trap
# If eigenda-proxy has failed, trap will error out and script will exit with an error code