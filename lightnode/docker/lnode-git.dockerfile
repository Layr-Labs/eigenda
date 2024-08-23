# lnode-git builds on top of lnode-base by cloning the eigenda repository. A useful caching step to avoid
# re-downloading the repository every time the image is built.

FROM lnode-base

# Clone eigenda repository
# TODO this should use environment variables or something
RUN git clone https://github.com/cody-littley/eigenda.git
