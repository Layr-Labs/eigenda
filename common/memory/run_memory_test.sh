#!/bin/bash

# Set the memory limit (2GB by default, but can be overridden)
MEMORY_LIMIT=${1:-2g}

# Directory containing the Dockerfile and where the command should be executed
cd "$(dirname "$0")/../.."

# Build the Docker image
echo "Building Docker image..."
docker build -t eigenda-memory-test -f common/memory/Dockerfile.memtest .

# Run the container with the specified memory limit
echo "Running test with ${MEMORY_LIMIT} memory limit..."
docker run --rm -m "${MEMORY_LIMIT}" eigenda-memory-test

echo "Test completed."