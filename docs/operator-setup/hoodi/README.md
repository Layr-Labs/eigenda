# Hoodi

Head over to our [EigenDA operator guides](https://docs.eigencloud.xyz/products/eigenda/operator-guides/overview) for installation instructions and more details.

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
