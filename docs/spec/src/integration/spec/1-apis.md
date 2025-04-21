# APIs

Below we give a summary of the APIs relevant to understanding the EigenDA high-level diagram

![](../../assets/integration/high-level-diagram.png)

### Proxy

See our gorilla/mux [routes](https://github.com/Layr-Labs/eigenda-proxy/blob/main/server/routing.go) for full detail, but the gist is that proxy presents a REST endpoint based off of the [op da-server spec](https://specs.optimism.io/experimental/alt-da.html#da-server) to rollup batchers:

```
# OP
POST /put body: <preimage_bytes> → <hex_encoded_commitment>
GET /get/{hex_encoded_commitment} → <preimage_bytes>
# NITRO
Same as OP but add a `?commitment_mode=standard` query param 
to both POST and GET methods.
```

### Disperser

The disperser presents a [grpc v2 service](https://github.com/Layr-Labs/eigenda/blob/ce89dab18d2f8f55004002e17dd3a18529277845/api/proto/disperser/v2/disperser_v2.proto#L10) endpoint

```bash
$ EIGENDA_DISPERSER_PREPROD=disperser-preprod-holesky.eigenda.xyz:443
$ grpcurl $EIGENDA_DISPERSER_PREPROD list disperser.v2.Disperser
disperser.v2.Disperser.DisperseBlob
disperser.v2.Disperser.GetBlobCommitment
disperser.v2.Disperser.GetBlobStatus
disperser.v2.Disperser.GetPaymentState
```

### Relay

Relays similarly present a [grpc service](https://github.com/Layr-Labs/eigenda/blob/ce89dab18d2f8f55004002e17dd3a18529277845/api/proto/relay/relay.proto#L10) endpoint

```bash
$ EIGENDA_RELAY_PREPROD=relay-1-preprod-holesky.eigenda.xyz:443
$ grpcurl $EIGENDA_RELAY_PREPROD list relay.Relay
relay.Relay.GetBlob
relay.Relay.GetChunks
```

### Contracts

The most important contract for rollups integrations is the EigenDACertVerifier, which presents a [function](https://github.com/Layr-Labs/eigenda/blob/98e21397e3471d170f3131549cdbc7113c0cdfaf/contracts/src/core/EigenDACertVerifier.sol#L86) to validate Certs:

```solidity
/**
 * @notice Verifies a blob cert for the specified quorums with the default security thresholds
 * @param batchHeader The batch header of the blob 
 * @param blobInclusionInfo The inclusion proof for the blob cert
 * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
 */
function verifyDACertV2(
    BatchHeaderV2 calldata batchHeader,
    BlobInclusionInfo calldata blobInclusionInfo,
    NonSignerStakesAndSignature calldata nonSignerStakesAndSignature
) external view
```