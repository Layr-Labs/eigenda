![Unit Tests](https://github.com/Layr-Labs/eigenda/actions/workflows/unit-tests.yml/badge.svg)
![Integration Tests](https://github.com/Layr-Labs/eigenda/actions/workflows/integration-tests.yml/badge.svg)
![Linter](https://github.com/Layr-Labs/eigenda/actions/workflows/golangci-lint.yml/badge.svg)
![Contracts](https://github.com/Layr-Labs/eigenda/actions/workflows/test-contracts.yml/badge.svg)
![Go Coverage](https://github.com/Layr-Labs/eigenda/wiki/coverage.svg)

# EigenDA

## Overview

EigenDA is a secure, high-throughput, and decentralized data availability (DA) service built on top of Ethereum using the [EigenLayer](https://github.com/Layr-Labs/eigenlayer-contracts) restaking primitives.

To understand more how EigenDA works and how it transforms the modern landscape of data availability, continue reading [EigenDA introduction](https://www.blog.eigenlayer.xyz/intro-to-eigenda-hyperscale-data-availability-for-rollups/).

To dive deep into the technical details, continue reading [EigenDA protocol spec](https://github.com/Layr-Labs/eigenda/blob/master/docs/spec/overview.md).

If you're interested in integrating your rollup with EigenDA, please fill out the [EigenDA questionnaire](https://docs.google.com/forms/d/e/1FAIpQLSez6PG-BL6C6Mc4QY1M--vbV219OGL_0Euv2zhJ1HmcUiU7cw/viewform).

## Why EigenDA?
As the [first actively validated service (AVS)](https://www.blog.eigenlayer.xyz/twelve-early-projects-building-on-eigenlayer/) built on EigenLayer, EigenDA transforms the data availability by providing high throughput and low cost serivice with security derived from Ethereum.

- Aligning with Ethereum ecosystem and building toward the Ethereum scaling [endgame](https://vitalik.ca/general/2021/12/06/endgame.html)
- A standard for high throughput and low cost data availability to enable growth of new on-chain use cases
- Horizontally scaling both security and throughput with the amount of restake and operators in the network, and therefore protecting decentralization (less work needed from each operator as network scaling)
- Innovative features such as Dual Quorum (two sperate quorums can be required to attest to the availability of data, for example ETH quorum and rollup's native token), customizable safety and liveness.

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

## Contact

- [Open an Issue](https://github.com/Layr-Labs/eigenda/issues/new/choose)
- [EigenLayer/EigenDA forum](https://forum.eigenlayer.xyz/c/eigenda/9)
- [Email](mailto:eigenda-support@eigenlabs.org)
- [Follow us on Twitter](https://twitter.com/eigenlayer)
