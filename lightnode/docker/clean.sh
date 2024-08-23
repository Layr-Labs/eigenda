#!/usr/bin/env bash

# TODO remove warnings if images are not found
# TODO figure out why lnode-git won't go away by itself

docker image rm lnode
docker image rm lnode-git
docker image rm lnode-base
