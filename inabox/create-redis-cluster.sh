#!/bin/bash

# Specify the Redis Docker image version
REDIS_IMAGE="redis:latest"

# Specify the Redis container name
REDIS_CONTAINER_NAME="test-eigenda-redis-cluster"

# Specify the port to be used for Redis
REDIS_PORT="6379"

# Check if the Redis container is already running
if [ $(docker ps -q -f name="$REDIS_CONTAINER_NAME") ]; then
    echo "Redis container ($REDIS_CONTAINER_NAME) is already running."
else
    # Check if the Redis container exists but is stopped
    if [ $(docker ps -aq -f name="$REDIS_CONTAINER_NAME") ]; then
        # Start the existing Redis container
        docker start "$REDIS_CONTAINER_NAME"
    else
        # Run a new Redis container
        docker run --name "$REDIS_CONTAINER_NAME" -p "$REDIS_PORT:$REDIS_PORT" -d "$REDIS_IMAGE"
    fi
    echo "Redis server started and available on port $REDIS_PORT"
fi
