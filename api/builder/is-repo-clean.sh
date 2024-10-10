#!/usr/bin/env bash

# This script exits with error code 0 if the git repository is clean, and error code 1 if it is not.
# This is utilized by the github workflow that checks to see if the repo is clean after recompiling
# protobufs.

if output=$(git status --porcelain) && [ -z "$output" ]; then
  echo "Repository is clean."
  exit 0
else
  echo "Repository is dirty:"
  git status
  exit 1
fi