# Integration Utils

A unified command-line tool for EigenDA integration utilities.

## Commands

### `parse-altdacommitment`
Parse and display EigenDA certificates from hex-encoded RLP strings. Hex strings can be obtained from eigenda-proxy output or rollup inbox data. For OP rollups, remove the '1' prefix byte from calldata before parsing.

### `gas-exhaustion-cert-meter` 
Estimates gas costs for verifying EigenDA certificates when all operators are non-signers (worst case scenario).

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
```