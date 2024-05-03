# EigenDA Plasma DA Server

## Introduction

This simple DA server implementation supports local storage via file based storage and remote via S3.
LevelDB is only recommended for usage in local devnets where connecting to S3 is not convenient.
See the [S3 doc](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/) for more information
on how to configure the S3 client.

## S3 Configuration

Depending on your cloud provider a wide array of configurations are available. The S3 client will
load configurations from the environment, shared credentials and shared config files.
Sample environment variables are provided below:

```bash
export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN=YOUR_SESSION_TOKEN
export AWS_REGION=YOUR_REGION
```

You can find out more about AWS authentication [here](https://docs.aws.amazon.com/sdkref/latest/guide/creds-config-files.html).

Additionally, these variables can be used with a google cloud S3 endpoint as well, i.e:

```bash
export AWS_ENDPOINT_URL="https://storage.googleapis.com"
export AWS_ACCESS_KEY_ID=YOUR_GOOGLE_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=YOUR_GOOGLE_ACCESS_KEY_SECRET
```

## EigenDA Configuration
Additional cli args are provided for targeting an EigenDA network backend:
- `--eigenda-rpc`: RPC host of disperser service. (e.g, on holesky this is `disperser-holesky.eigenda.xyz:443`)
- `--eigenda-status-query-timeout`: (default: 1m) Duration for which a client will wait for a blob to finalize after being sent for dispersal.
- `--eigenda-status-query-retry-interval`: (default: 5s) How often a client will attempt a retry when awaiting network blob finalization. 
