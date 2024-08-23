# lnode is the top level dockerfile for the light node. It builds on top of lnode-git by checking out and building the
# target branch of the eigenda repository.

FROM lnode-git

# The branch or commit to check out.
ARG GIT_BRANCH_OR_COMMIT

# Check out the target branch/commit
WORKDIR /home/lnode/eigenda
RUN git fetch
RUN echo "Checking out $GIT_BRANCH_OR_COMMIT"
RUN git checkout $GIT_BRANCH_OR_COMMIT
RUN git pull

# Build the light node
WORKDIR /home/lnode/eigenda/lightnode
RUN pwd
RUN ls -alh
RUN make build

# Reset the workdir to the home directory
WORKDIR /home/lnode

# Run the light node
CMD eigenda/lightnode/bin/lnode
