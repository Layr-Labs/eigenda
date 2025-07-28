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
The configuration allows users to configure memstore to inject specific derivation error responses during GET operations while allowing PUT operations to succeed normally. This enables testing client handling of derivation errors without requiring complex setup.
Specifically, users send a PATCH request that sets the desired derivation error for all subsequent GET requests. After that, when the user sends data to the proxy, the PUT operation succeeds as normal—the error injection only affects the GET path. Behind the scenes, upon a GET request, the proxy returns either the stored data or the specified derivation error depending on its configuration.
The PATCH request is sticky, meaning it will take effect on multiple GET requests unless reset.

```bash
 curl -X PATCH http://localhost:3100/memstore/config -d '{"PutWithGetReturnsDerivationError": {"StatusCode": 3}}'
 {"MaxBlobSizeBytes":2048,"BlobExpiration":"45m0s","PutLatency":"0s","GetLatency":"0s","PutReturnsFailoverError":false,"PutWithGetReturnsDerivationError": {"StatusCode": 3}}
```

A very important invariant is that no key can ever be overwritten.

### Golang client
A simple HTTP client implementation lives in `/clients/memconfig_client/` and can be imported for manipulating the config using more structured types.