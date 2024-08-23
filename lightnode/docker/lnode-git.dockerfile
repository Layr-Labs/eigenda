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
RUN git fetch
RUN git checkout $BRANCH_OR_COMMIT
RUN git pull

# Reset the workdir to the home directory
WORKDIR /home/lnode
