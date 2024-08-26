# syntax=docker/dockerfile:1

# lnode-git builds on top of lnode-base by cloning the eigenda repository. A useful caching step to avoid
# re-downloading the entire repository every time the image is built.

FROM lnode-base

# The url of the git repository to clone.
ARG GIT_URL
# The branch or commit to check out.
ARG BRANCH_OR_COMMIT

# Clone eigenda repository
RUN git clone $GIT_URL eigenda
WORKDIR /home/lnode/eigenda
RUN git checkout $BRANCH_OR_COMMIT
WORKDIR /home/lnode

# Download all go dependencies.
# This is done as a separate step to avoid repeat work every time the latest commit in the branch is updated.
WORKDIR /home/lnode/eigenda/lightnode
RUN bash -c 'source ~/.bashrc && go mod download'
WORKDIR /home/lnode
