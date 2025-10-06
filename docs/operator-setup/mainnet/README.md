## Mainnet

Head over to our [EigenDA operator guides](https://docs.eigencloud.xyz/products/eigenda/operator-guides/overview) for installation instructions and more details.

## Blazar (EigenDA V2) Migration
Operators running node version `<=v0.8.6` will need to define new v2 specific environment variables, expose 2 new ports, and update their socket registration as part of the migration to v2.

## Migration Steps
### 1. Update `.env` with v2 specific environment variables
```
NODE_V2_RUNTIME_MODE=v1-and-v2

NODE_V2_DISPERSAL_PORT=32006
NODE_V2_RETRIEVAL_PORT=32007

# Internal ports for Nginx reverse proxy
NODE_INTERNAL_V2_DISPERSAL_PORT=${NODE_V2_DISPERSAL_PORT}
NODE_INTERNAL_V2_RETRIEVAL_PORT=${NODE_V2_RETRIEVAL_PORT}
```

### 2. Update `MAIN_SERVICE_IMAGE`
```
MAIN_SERVICE_IMAGE=ghcr.io/layr-labs/eigenda/opr-node:2.2.0
```

### 3. Update socket registration
EigenDA v2 adds new ports to the socket registration. Socket registration update is required to receive v2 traffic.

Ensure that you are using the latest version of the [eigenda-operator-setup](https://github.com/Layr-Labs/eigenda-operator-setup/releases) before updating the socket.
```
(eigenda-operator-setup) > ./run.sh update-socket
You are about to update your socket to: 23.93.87.155:32005;32004;32006;32007
Confirm? [Y/n]
```

### 4. Restart the node and monitor for reachability checks
The node will check reachability of v1 & v2 sockets. If reachability checks are failing, check that the new ports are open and accessible.
```
Feb 20 19:47:07.861 INF node/node.go:743 Reachability check v1 - dispersal socket ONLINE component=Node status="node.Dispersal is available" socket=operator.eigenda.xyz:32001
Feb 20 19:47:07.861 INF node/node.go:750 Reachability check v1 - retrieval socket ONLINE component=Node status="node.Retrieval is available" socket=operator.eigenda.xyz:32002
Feb 20 19:47:07.867 INF node/node.go:743 Reachability check v2 - dispersal socket ONLINE component=Node status="validator.Dispersal is available" socket=operator.eigenda.xyz:32003
Feb 20 19:47:07.867 INF node/node.go:750 Reachability check v2 - retrieval socket ONLINE component=Node status="validator.Retrieval is available" socket=operator.eigenda.xyz:32005
```

### 5. Confirm v2 StoreChunks requests are being served
```
Feb 20 19:50:36.741 INF grpc/server_v2.go:140 new StoreChunks request batchHeaderHash=873ac1c7faeec0f1e5c886142d0b364a94b3e906f1b4b4f1b0466a5f79cecefb numBlobs=14 referenceBlockNumber=3393054
Feb 20 19:50:41.765 INF grpc/server_v2.go:140 new StoreChunks request batchHeaderHash=76873d64609d50aaf90e1c435c9278c588f1a174a4c0b4a721438a7d44bb2f1e numBlobs=18 referenceBlockNumber=3393054
Feb 20 19:50:46.760 INF grpc/server_v2.go:140 new StoreChunks request batchHeaderHash=8182f31c9b58e04f0a09dfbf1634a73e47a660b441f65c7a35ef9e7afd064493 numBlobs=16 referenceBlockNumber=3393054

```

## V2 Remote BLS Signer API Migration
__Breaking Change:__ With the release of EigenDA Blazar (v2), the the node software has been upgraded to use the latest [cerberus](https://github.com/Layr-Labs/cerberus) remote BLS signer release.

This change requires the operator to define the `NODE_BLS_SIGNER_API_KEY` environment variable.

Follow the steps from the [cerberus setup guide](https://github.com/Layr-Labs/cerberus?tab=readme-ov-file#remote-signer-implementation-of-cerberus-api) to create an API key.


## V2 Environment Variable Reference

### `NODE_RUNTIME_MODE`
This environment variable will be used to determine the runtime mode of the EigenDA node.

- `v1-and-v2`: The node will serve both v1 and v2 traffic (default)
- `v2-only`: The node will serve v2 traffic only
- `v1-only`: The node will serve v1 traffic only

The `v1-only` & `v2-only` modes are intended for isolating traffic to separate validator instances - where 1 instance serves v1 traffic and a second instance serves v2 traffic.

### `NODE_V2_DISPERSAL_PORT`
<ins>Operators must publically expose this port</ins>. This port will be used to listen for dispersal requests from the EigenDA v2 API. IP whitelisting is no longer required with v2.

### `NODE_V2_RETRIEVAL_PORT`
<ins>Operators must publically expose this port</ins>. This port will be used to listen for retrieval requests from the EigenDA v2 API.

### `NODE_INTERNAL_V2_DISPERSAL_PORT`
This port is intended for Nginx reverse proxy use.

### `NODE_INTERNAL_V2_RETRIEVAL_PORT`
This port is intended for Nginx reverse proxy use.


## Advanced - Multi drive support for V2 LittDB

LittDB is capable of partitioning the chunks DB across multiple drives. See https://github.com/Layr-Labs/eigenda/blob/master/node/database-paths.md for more details.

### Example .env that defines 2 littDB partitions
```
NODE_LITT_DB_STORAGE_HOST_PATH_1=${NODE_DB_PATH_HOST}/chunk_v2_litt_1
NODE_LITT_DB_STORAGE_HOST_PATH_2=${NODE_DB_PATH_HOST}/chunk_v2_litt_2
NODE_LITT_DB_STORAGE_PATHS=/data/operator/db/chunk_v2_litt_1,/data/operator/db/chunk_v2_litt_2
```
The `NODE_LITT_DB_STORAGE_HOST_PATH_X` should map to a separately mounted drive folder for each partition.
The `NODE_LITT_DB_STORAGE_PATHS` needs to map to the docker volume mounts defined in `docker-compose.yaml`
#### Example docker-compose volume mounts
```
   volumes:
      - "${NODE_ECDSA_KEY_FILE_HOST}:/app/operator_keys/ecdsa_key.json:readonly"
      - "${NODE_BLS_KEY_FILE_HOST}:/app/operator_keys/bls_key.json:readonly"
      - "${NODE_G1_PATH_HOST}:/app/g1.point:readonly"
      - "${NODE_G2_PATH_HOST}:/app/g2.point.powerOf2:readonly"
      - "${NODE_CACHE_PATH_HOST}:/app/cache:rw"
      - "${NODE_LOG_PATH_HOST}:/app/logs:rw"
      - "${NODE_DB_PATH_HOST}:/data/operator/db:rw"
      - "${NODE_LITT_DB_STORAGE_HOST_PATH_1}:/data/operator/db/chunk_v2_litt_1:rw"
      - "${NODE_LITT_DB_STORAGE_HOST_PATH_2}:/data/operator/db/chunk_v2_litt_2:rw"
```
#### Example node output after loading both partitions
```
eigenda-native-node  | Jun  2 18:11:47.430 INF node/validator_store.go:108 Using littDB at paths /data/operator/db/chunk_v2_litt_1, /data/operator/db/chunk_v2_litt_2
eigenda-native-node  | Jun  2 18:11:47.430 INF littbuilder/db_impl.go:118 LittDB started, current data size: 0
eigenda-native-node  | Jun  2 18:11:47.430 INF littbuilder/db_impl.go:169 creating table chunks
