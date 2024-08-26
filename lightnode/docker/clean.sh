#!/usr/bin/env bash

# Cleans the docker image and all cached steps.

docker image rm lnode || true # don't fail if the image doesn't exist
docker builder prune -f
