#!/bin/bash

cd optimism

echo "Initializing monorepo..."
make install-geth &&
git submodule update --init --recursive &&
make devnet-allocs &&
mv .devnet/* ../e2e/resources/optimism/ &&
mv packages/contracts-bedrock/deploy-config/devnetL1.json ../e2e/resources/optimism/devnetL1.json
