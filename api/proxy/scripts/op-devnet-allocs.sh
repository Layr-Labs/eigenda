#!/bin/bash

VERSION=$(cat go.mod | grep -m 1 ethereum-optimism/optimism | awk '{print $2}' | sed 's/\/v//g')
VERSION=$(echo ${VERSION} | sed 's/v//g')

REPO_NAME=optimism-$(echo ${VERSION} | sed 's/v//g')

echo "Downloading ${REPO_NAME} ..."
git clone --branch v${VERSION} https://github.com/ethereum-optimism/optimism.git ${REPO_NAME}
cd ${REPO_NAME}

echo "Initializing monorepo..."
make install-geth &&
    git submodule update --init --recursive &&
    # make devnet-allocs creates allocs for different configurations: .devnet-alt-da, .devnet-l200, .devnet-mt-cannon, etc.
    # we only need the .devnet-alt-da (although we could use .devnet-alt-da-generic once the PR that adds it is merged)
    make devnet-allocs &&
    cp -R .devnet-alt-da ../. &&
    # op's e2e test suite hardcodes this path so we need to copy it to this awkward location
    mv packages/contracts-bedrock/deploy-config/devnetL1.json ../packages/contracts-bedrock/deploy-config/devnetL1.json

STATUS=$?

## Force cleanup of monorepo
echo "${STATUS} Cleaning up ${REPO_NAME} repo ..."
cd ../ &&
    rm -rf ${REPO_NAME}

if [ $? -eq 0 ]; then
    echo "Successfully cleaned up ${REPO_NAME} repo"
    exit ${STATUS}
else
    echo "Failed to clean up ${REPO_NAME} repo"
    exit ${STATUS}
fi
