#!/bin/bash

REDIS_CONTAINER_NAME="my-redis"
REDIS_PORT="6379"

echo "Pulling the latest Redis image from Docker Hub..."
docker pull redis:latest

echo "Checking if Redis container ($REDIS_CONTAINER_NAME) already exists..."

# Check if the Redis container is already running
if [ "$(docker ps -q -f name=$REDIS_CONTAINER_NAME)" ]; then
    echo "Redis container ($REDIS_CONTAINER_NAME) is already running."
    exit 0
fi

# Check if the Redis container exists but is stopped
if [ "$(docker ps -aq -f name=$REDIS_CONTAINER_NAME)" ]; then
    echo "Starting existing Redis container ($REDIS_CONTAINER_NAME)..."
    docker start "$REDIS_CONTAINER_NAME"
else
    echo "Starting a new Redis container ($REDIS_CONTAINER_NAME)..."
    docker run --name "$REDIS_CONTAINER_NAME" -p "$REDIS_PORT:$REDIS_PORT" -d redis
fi

echo "Redis setup complete."