#!/usr/bin/env bash

# The directory where this script is located.
SCRIPT_DIR=$(dirname "$0")

export BUILDKIT_PROGRESS=plain

docker build -t lnode-base -f $SCRIPT_DIR/lnode-base.dockerfile .
docker build -t lnode-git -f $SCRIPT_DIR/lnode-git.dockerfile .
docker build --no-cache -t lnode -f $SCRIPT_DIR/lnode.dockerfile .