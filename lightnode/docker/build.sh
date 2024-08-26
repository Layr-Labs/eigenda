#!/usr/bin/env bash

source docker/default-args.sh
# Create a file called docker/args.sh to override the default values of GIT_URL and BRANCH_OR_COMMIT.
source docker/args.sh

echo "git url: ${GIT_URL}"
echo "branch or commit: ${BRANCH_OR_COMMIT}"
echo "go url: ${GO_URL}"

# Create a file with information about this build. This file will be copied into the docker image.
rm docker/build-info.txt 2> /dev/null || true
touch docker/build-info.txt
echo "git URL: ${GIT_URL}" >> docker/build-info.txt

# This will return an empty string if a git sha is provided. Will return sha and branch name if a branch is provided.
COMMIT_SHA=$(git ls-remote $GIT_URL $BRANCH_OR_COMMIT)
if [ -z "$COMMIT_SHA" ]; then
  echo "target: ${BRANCH_OR_COMMIT}" >> docker/build-info.txt
else
  echo "target: ${COMMIT_SHA}" >> docker/build-info.txt
fi

echo "docker build commit: $(git rev-parse HEAD)" >> docker/build-info.txt

# Docker image is split into three parts:
#  - the base image with OS level packages installed
#  - the git image with the code cloned and go modules downloaded
#  - final image with the code built.
#
# The purpose for this split is to prevent docker from deleting intermediate layers which are time consuming to build.

# Add the --no-cache flag to force a rebuild.
# Add the --progress=plain flag to show verbose output during the build.

docker build \
  --build-arg="GO_URL=${GO_URL}" \
  -f docker/lnode-base.dockerfile \
  --tag lnode-base:latest \
  .
if [ $? -ne 0 ]; then
  echo "Failed to build lnode-base"
  exit 1
fi

docker build \
  --build-arg="GIT_URL=${GIT_URL}" \
  --build-arg="BRANCH_OR_COMMIT=${BRANCH_OR_COMMIT}" \
  -f docker/lnode-git.dockerfile \
  --tag lnode-git:latest \
  .
if [ $? -ne 0 ]; then
  echo "Failed to build lnode-git"
  exit 1
fi

docker build \
  --build-arg="GIT_URL=${GIT_URL}" \
  --build-arg="BRANCH_OR_COMMIT=${BRANCH_OR_COMMIT}" \
  -f docker/lnode.dockerfile \
  --tag lnode:latest \
  .
if [ $? -ne 0 ]; then
  echo "Failed to build lnode"
  exit 1
fi

# Don't leave trash on the filesystem.
rm docker/build-info.txt 2> /dev/null || true
