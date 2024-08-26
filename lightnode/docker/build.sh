#!/usr/bin/env bash

source docker/default-args.sh
# Create a file called docker/args.sh to override the default values of GIT_URL and BRANCH_OR_COMMIT.
source docker/args.sh

echo "Building lnode docker image with GIT_URL=${GIT_URL} and BRANCH_OR_COMMIT=${BRANCH_OR_COMMIT}"

# Create a file with information about this build. This file will be copied into the docker image.
rm docker/build-info.txt
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

# Add the --no-cache flag to force a rebuild.
# Add the --progress=plain flag to show verbose output during the build.
docker build \
  --build-arg="GIT_URL=${GIT_URL}" \
  --build-arg="BRANCH_OR_COMMIT=${BRANCH_OR_COMMIT}" \
  -f docker/Dockerfile \
  -t lnode .
