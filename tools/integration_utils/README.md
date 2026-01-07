# Integration Utils

A unified command-line tool for EigenDA integration utilities.

## Commands

### `parse-altdacommitment`
Parse and display EigenDA certificates from hex-encoded RLP strings. Hex strings can be obtained from eigenda-proxy output or rollup inbox data. For OP rollups, remove the '1' prefix byte from calldata before parsing.

### `gas-exhaustion-cert-meter` 
Estimates gas costs for verifying EigenDA certificates when all operators are non-signers (worst case scenario).

### `validate-cert-verifier`
Validates the CertVerifier contract by dispersing a test blob to EigenDA, constructing a `DA Cert` from the disperser's reply, and verifying that the CertVerifier contract correctly verifies the returned certificate using `checkDACert`. This is useful for integration testing and validating CertVerifier deployments.

## Usage

```bash
# Build the tool
make build

# Run with help
./bin/integration_utils --help

# Parse a certificate
./bin/integration_utils parse-altdacommitment --hex <hex_string>

# Estimate gas costs
./bin/integration_utils gas-exhaustion-cert-meter --help

# Validate CertVerifier contract
./bin/integration_utils validate-cert-verifier \
  --eigenda-network hoodi_testnet \
  --json-rpc-url <RPC_URL> \
  --signer-auth-key <PRIVATE_KEY> \
  --cert-verifier-address <CONTRACT_ADDRESS>
```