#!/usr/bin/env bash

# Cleans the docker image and all cached steps.
docker image rm pbuf-compiler 2> /dev/null || true
docker builder prune -f