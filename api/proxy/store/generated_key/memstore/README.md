# Memstore Backend

The Memstore backend is a simple in-memory key-value store that is meant to replace a real EigenDA backend (talking to the disperser) for testing and development purposes. It is **never** recommended for production use.

## Usage

```bash
./bin/eigenda-proxy --memstore.enabled
```

## Configuration

See [memconfig/config.go](./memconfig/config.go) for the configuration options.
These can all be set via their respective flags or environment variables. Run `./bin/eigenda-proxy --help | grep memstore` to see these.

## Config REST API

The Memstore backend also provides a REST API for changing the configuration at runtime. This is useful for testing different configurations without restarting the proxy.

The API consists of GET and PATCH methods on the `/memstore/config` resource.

### Get the current configuration

```bash
$ curl http://localhost:3100/memstore/config | jq
{
  "MaxBlobSizeBytes": 16777216,
  "BlobExpiration": "25m0s",
  "PutLatency": "0s",
  "GetLatency": "0s",
  "PutReturnsFailoverError": false
}
```

### Set a configuration option

The PATCH request allows to patch the configuration. This allows only sending a subset of the configuration options. The other fields will be left intact.

```bash
$ curl -X PATCH http://localhost:3100/memstore/config -d '{"PutReturnsFailoverError": true}'
{"MaxBlobSizeBytes":16777216,"BlobExpiration":"25m0s","PutLatency":"0s","GetLatency":"0s","PutReturnsFailoverError":true}
```

One can of course still build a jq pipe to produce the same result (although still using PATCH instead of PUT since that is the only method available):
```bash
$ curl http://localhost:3100/memstore/config | \
  jq '.PutLatency = "5s" | .GetLatency = "2s"' | \
  curl -X PATCH http://localhost:3100/memstore/config -d @-
```

#### PUT with GET returning derivation error
The configuration allows users to configure memstore to return specific derivation error responses during `/get` payload retrievals while still allowing `/put` request operations to succeed normally with the payload persisted to ephemeral db. This enables fast iteration testing of a rollup client's handling of derivation errors without requiring a complex setup.

The `PutWithGetReturnsDerivationError` field supports three states:
1. **Field omitted**: No change to current configuration
2. **Set an error**: `{"PutWithGetReturnsDerivationError": {"StatusCode": 3, "Msg": "test error", "Reset": false}}`
3. **Reset to nil (disabled)**: `{"PutWithGetReturnsDerivationError": {"Reset": true}}`

##### Setting a derivation error
Configure memstore to return a specific derivation error for all subsequent `/get` requests:

```bash
curl -X PATCH http://localhost:3100/memstore/config \
  -d '{"PutWithGetReturnsDerivationError": {"StatusCode": 3, "Msg": "Invalid cert", "Reset": false}}'
```

This will cause all future `/get` requests to return an HTTP 418 error with the specified derivation error, while `/put` requests continue to succeed normally.

##### Resetting derivation error behavior
To disable the derivation error behavior and return to normal operation:

```bash
curl -X PATCH http://localhost:3100/memstore/config \
  -d '{"PutWithGetReturnsDerivationError": {"Reset": true}}'
```

After this, `/get` requests will return stored data normally instead of errors.

A very important invariant is that no key can ever be overwritten.

### Golang client
A simple HTTP client implementation lives in `/clients/memconfig_client/` and can be imported for manipulating the config using more structured types.