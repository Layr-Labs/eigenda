#!/usr/bin/env bash

# Cleans the docker image and all cached steps.
docker image rm lnode 2> /dev/null || true
docker builder prune -f
