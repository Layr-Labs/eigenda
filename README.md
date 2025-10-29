![Unit Tests](https://github.com/Layr-Labs/eigenda/actions/workflows/unit-tests.yml/badge.svg)
![Integration Tests](https://github.com/Layr-Labs/eigenda/actions/workflows/integration-tests.yml/badge.svg)
![Linter](https://github.com/Layr-Labs/eigenda/actions/workflows/golangci-lint.yml/badge.svg)
![Contracts](https://github.com/Layr-Labs/eigenda/actions/workflows/test-contracts.yml/badge.svg)
[![codecov](https://codecov.io/github/Layr-Labs/eigenda/graph/badge.svg?token=EKLGVKW1VN)](https://codecov.io/github/Layr-Labs/eigenda)

# EigenDA

## Overview

EigenDA is a secure, high-throughput, and decentralized data availability (DA) service built on top of Ethereum using the [EigenLayer](https://github.com/Layr-Labs/eigenlayer-contracts) restaking primitives.

To understand more about how EigenDA works and how it transforms the modern landscape of data availability, continue reading [EigenDA introduction](https://www.blog.eigenlayer.xyz/intro-to-eigenda-hyperscale-data-availability-for-rollups/).

To dive deep into the technical details, continue reading [EigenDA protocol spec](https://layr-labs.github.io/eigenda/) in mdBook.

If you're interested in integrating your rollup with EigenDA, follow the rollup guides [here](https://docs.eigencloud.xyz/products/eigenda/api/disperser-v2-API/overview)

## API Documentation

The EigenDA public API is documented [here](https://docs.eigencloud.xyz/products/eigenda/api/disperser-v2-API/overview).

## Operating EigenDA Node

If you want to be an EigenDA operator and run a node, please clone [Operator Setup Guide](https://github.com/Layr-Labs/eigenda-operator-setup) GitHub repo and follow the instructions there.

## Repository Structure

- **`./rust`** - Sovereign SDK EigenDA adapter: A data availability adapter implementation for [Sovereign SDK](https://github.com/Sovereign-Labs/sovereign-sdk) rollups that enables them to use EigenDA as their data availability layer.

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

To set up your development environment, first [install and activate mise](https://mise.jdx.dev/getting-started.html), then run:

```bash
mise install              # Install all development tools
mise run install-hooks    # Install git pre-commit hooks
```

### Pre-commit Hooks

We provide pre-commit hooks to automatically check your code before committing. These hooks run linting and formatting checks to catch issues early.

The hooks are installed automatically when you run `mise run install-hooks` (see Dependency Management above).

The pre-commit hook will run the following checks:
- **Linting**: Runs `golangci-lint` to check code quality
- **Go mod tidy check**: Ensures `go.mod` and `go.sum` are up to date
- **Format checking**: Verifies Go and Solidity code formatting

If any checks fail, the commit will be blocked. You can:
- Fix the issues by running `make fmt` to auto-format code and `go mod tidy` if needed
- Bypass the hooks (not recommended) using `git commit --no-verify`

**Note**: You can also manually install/update hooks by running `./scripts/install-hooks.sh`

## Contact

- [Open an Issue](https://github.com/Layr-Labs/eigenda/issues/new/choose)
- [EigenDA forum](https://forum.eigenlayer.xyz/c/eigenda-research/36)
- [Email](mailto:eigenda-support@eigenlabs.org)
- [Follow us on X](https://x.com/eigen_da)
