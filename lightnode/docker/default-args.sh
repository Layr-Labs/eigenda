# Default arguments for building the docker image.
# To override these values locally, create a file named `args.sh` in the same directory. 'args.sh' is ignored by git.

# The location where the go binary will be downloaded from.
export GO_URL=https://go.dev/dl/go1.21.12.linux-arm64.tar.gz

# The location of the code to clone.
export GIT_URL=https://github.com/Layr-Labs/eigenda.git

# The name of the branch or the commit sha to clone.
export BRANCH_OR_COMMIT=master

# The location on the host file system where light node data will be stored.
export DATA_PATH=~/.lnode-data
