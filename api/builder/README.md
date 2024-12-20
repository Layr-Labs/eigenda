This directory contains scripts for building a docker image capable of compiling the EigenDA protobufs. I found
it difficult to control the exact build version of the protobufs, since the version depends on whatever is installed
locally when they are built. This is an attempt to standardize the protobuf build process.

# Usage

To build the docker image, run the following command:

```bash
./api/builder/build-docker.sh
```

Once the docker image is built, you can build the protobufs via the following command:

```bash
./api/builder/protoc-docker.sh
```

# Caveats

I've tested this on my m3 macbook. It's possible that the docker image may have trouble on other architectures.
Please report any issues you encounter with this build process to the EigenDA team. The goal is to be architecturally
agnostic, but that isn't a priority in the very short term.
