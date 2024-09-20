#!/usr/bin/env bash

# This script fully deletes the pbuf-compiler docker image and all cached steps.

# Cleans the docker image and all cached steps.
docker image rm pbuf-compiler 2> /dev/null || true
docker builder prune -f