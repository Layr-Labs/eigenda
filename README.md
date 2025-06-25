![Unit Tests](https://github.com/Layr-Labs/eigenda/actions/workflows/unit-tests.yml/badge.svg)
![Integration Tests](https://github.com/Layr-Labs/eigenda/actions/workflows/integration-tests.yml/badge.svg)
![Linter](https://github.com/Layr-Labs/eigenda/actions/workflows/golangci-lint.yml/badge.svg)
![Contracts](https://github.com/Layr-Labs/eigenda/actions/workflows/test-contracts.yml/badge.svg)
![Go Coverage](https://github.com/Layr-Labs/eigenda/wiki/coverage.svg)

# EigenDA

## Overview

EigenDA is a secure, high-throughput, and decentralized data availability (DA) service built on top of Ethereum using the [EigenLayer](https://github.com/Layr-Labs/eigenlayer-contracts) restaking primitives.

To understand more about how EigenDA works and how it transforms the modern landscape of data availability, continue reading [EigenDA introduction](https://www.blog.eigenlayer.xyz/intro-to-eigenda-hyperscale-data-availability-for-rollups/).

To dive deep into the technical details, continue reading [EigenDA protocol spec](https://layr-labs.github.io/eigenda/) in mdBook.

If you're interested in integrating your rollup with EigenDA, please fill out the [EigenDA Partner Registration](https://docs.google.com/forms/d/e/1FAIpQLSdXvfxgRfIHWYu90FqN-2yyhgrYm9oExr0jSy7ERzbMUimJew/viewform).

## API Documentation

The EigenDA public API is documented [here](https://github.com/Layr-Labs/eigenda/tree/master/api/docs).

## Operating EigenDA Node

If you want to be an EigenDA operator and run a node, please clone [Operator Setup Guide](https://github.com/Layr-Labs/eigenda-operator-setup) GitHub repo and follow the instructions there.

## Contributing
We welcome all contributions! There are many ways to contribute to the project, including but not limited to:

- Opening a PR
- [Submitting feature requests or bugs](https://github.com/Layr-Labs/eigenda/issues/new/choose)
- Improving our product or contribution documentation
- Voting on [open issues](https://github.com/Layr-Labs/eigenda/issues) or
  contributing use cases to a feature request

### Dependency Management

We use [mise](https://mise.jdx.dev/) to manage dependencies in EigenDA. This is still a work in progress, as it currently only manages go and golangci-lint dependencies.
The goal is to eventually get exact parity and reproducibility between our CI and local environments, so that we can reproduce and debug failing CI issues locally.
To install dependencies using mise, first [install and activate mise](https://mise.jdx.dev/getting-started.html), and then run `mise install` in the root of the repository.

## Contact

- [Open an Issue](https://github.com/Layr-Labs/eigenda/issues/new/choose)
- [EigenLayer/EigenDA forum](https://forum.eigenlayer.xyz)
- [Email](mailto:eigenda-support@eigenlabs.org)
- [Follow us on X](https://x.com/eigen_da)
