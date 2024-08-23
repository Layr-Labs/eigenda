#!/usr/bin/env bash

# The directory where this script is located.
SCRIPT_DIR=$(dirname "$0")

source $SCRIPT_DIR/target.sh

# Uncomment this to enable more verbose output.
export BUILDKIT_PROGRESS=plain

# TODO flag for caching or no caching

docker build \
  -t lnode-base \
  -f $SCRIPT_DIR/lnode-base.dockerfile .

docker build --no-cache \
  -t lnode-git \
  --build-arg GIT_URL=$GIT_URL \
  --build-arg BRANCH_OR_COMMIT=$BRANCH_OR_COMMIT \
  -f $SCRIPT_DIR/lnode-git.dockerfile .

docker build --no-cache \
  -t lnode \
  --build-arg BRANCH_OR_COMMIT=$BRANCH_OR_COMMIT \
  -f $SCRIPT_DIR/lnode.dockerfile .
