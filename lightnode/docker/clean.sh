#!/usr/bin/env bash

# Cleans the docker image and all cached steps.

docker image rm lnode-base || true
docker image rm lnode-git || true
docker image rm lnode || true
docker builder prune -f
