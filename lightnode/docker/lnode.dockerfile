# lnode is the top level dockerfile for the light node. It builds on top of lnode-git by checking out and building the
# target branch of the eigenda repository.

FROM lnode-git

# Check out the target branch/commit
WORKDIR /home/lnode/eigenda
# TODO this should use environment variables or something
RUN git fetch
RUN git checkout master
RUN git pull

RUN pwd
RUN ls -alh

# Build the light node
WORKDIR /home/lnode/eigenda/lightnode
RUN pwd
RUN ls -alh
RUN make build

# Reset the workdir to the home directory
WORKDIR /home/lnode

# Run the light node
CMD eigenda/lightnode/bin/lnode
