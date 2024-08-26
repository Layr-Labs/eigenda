# syntax=docker/dockerfile:1

# lnode pulls the lates commit in the target branch and builds the light node.

FROM lnode-git

# The url of the git repository to clone.
ARG GIT_URL
# The branch or commit to check out.
ARG BRANCH_OR_COMMIT

# Copy a file containing build information. Useful for detective work on an otherwise unlabelled image.
# This is also useful for forcing docker to invalidate caches when the build target is the latest commit in a branch.
WORKDIR /home/lnode/eigenda/lightnode
COPY --chown=lnode docker/build-info.txt /home/lnode
WORKDIR /home/lnode

# Just in case we are tracking the latest commit in a branch, pull again. This is a no-op if the target is a commit sha.
WORKDIR /home/lnode/eigenda
RUN git pull
WORKDIR /home/lnode

# Build the light node
WORKDIR /home/lnode/eigenda/lightnode
RUN bash -c 'source ~/.bashrc && make build'
RUN ln -s /home/lnode/eigenda/lightnode/bin/lnode ~/lnode
WORKDIR /home/lnode

# Data that should persist between container restarts will be stored in this mount point.
RUN mkdir ~/data

# Data that shouldn't persist between container restarts will be stored in this location.
RUN mkdir ~/tmp

# Make everything except ~/tmp and ~/data read only to enforce good file system hygiene.
RUN chmod -R -w ~
RUN chmod -R +w ~/tmp
RUN chmod -R +w ~/data

# Run the light node when the container starts.
CMD ["/home/lnode/lnode"]
