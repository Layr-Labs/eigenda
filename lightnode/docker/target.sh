#!/usr/bin/env bash

# This script is used to specify which branch or commit to build into the docker image.

# TODO can this be made more automatic?

# The url where the git repository is located. Updating this URL requires "make clean-docker" to be run
# in order to be picked up (otherwise the previously cached value may be used).
export GIT_URL=https://github.com/cody-littley/eigenda.git

# The branch or commit to checkout. This value can be updated by a "make build-docker" without calling
# "make clean-docker" first.
export BRANCH_OR_COMMIT=lightnode-docker
