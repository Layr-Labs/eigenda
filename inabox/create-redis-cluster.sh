#!/bin/bash
set -e

REDIS_CLUSTER_ID="test-eigenda-redis-cluster"
REDIS_PORT="6379"
AWS_REGION="us-east-1"

# Check if the Redis cluster already exists
function redis_cluster_exists() {
    aws elasticache describe-cache-clusters --region $AWS_REGION | grep -q $REDIS_CLUSTER_ID
    return $?
}

# Start Redis service using LocalStack
function create_redis_cluster() {
    aws elasticache create-cache-cluster \
        --cache-cluster-id $REDIS_CLUSTER_ID \
        --engine redis \
        --cache-node-type cache.t2.micro \
        --num-cache-nodes 1 \
        --port $REDIS_PORT \
        --region $AWS_REGION
}

# Check if Redis cluster exists and create it if it does not
if redis_cluster_exists; then
    echo "Redis cluster $REDIS_CLUSTER_ID already exists."
else
    echo "Creating Redis cluster $REDIS_CLUSTER_ID."
    create_redis_cluster
fi
