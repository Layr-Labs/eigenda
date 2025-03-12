## Cancel Ethereum Transaction Lambda Function

This AWS Lambda function is designed to cancel a pending Ethereum transaction by sending a replacement “do-nothing” transaction. The replacement transaction sends 0 ETH from your own address back to itself, effectively canceling the pending transaction when mined. The function dynamically calculates the required gas fees by querying an Ethereum JSON‑RPC endpoint, signs the transaction using a key stored in AWS KMS, and then submits the signed transaction to the Ethereum network.

### How It Works
1. Dynamic Gas Fee Calculation:
- It reads the maximum priority fee (in gwei) from the environment variable PRIORITY_FEE_GWEI and converts it to wei.
- The maximum fee per gas is computed as the sum of the base fee and the priority fee.

2.	Transaction Construction:
- The transaction is built with a nonce supplied by the environment variable NONCE.
- The transaction sends 0 ETH from your account (calculated using your AWS KMS public key) back to your own address.
- It uses Ethereum mainnet settings (chain ID hardcoded to 1) and is built as an EIP-1559 type 2 transaction.


### Required Environment Variables

Ensure the following environment variables are set in your Lambda configuration:
- `ETH_RPC` : The URL of the Ethereum JSON‑RPC endpoint
- `NONCE` : The transaction nonce to be used for the cancellation transaction.
- `PRIORITY_FEE_GWEI` : The maximum priority fee per gas (in gwei) for the transaction. The value will be converted to wei.
- `KMS_KEY_ID` : The AWS KMS Key ID used to retrieve the public key for calculating the Ethereum address.


### Lambda Packaging

This project contains a Dockerfile and Makefile to create AWS Lambda deployment packages using Docker with Python 3.12.

### Prerequisites

- Docker
- Make


### Create the Lambda Package

```bash
make
```


### Clean Up

```bash
make clean
```

