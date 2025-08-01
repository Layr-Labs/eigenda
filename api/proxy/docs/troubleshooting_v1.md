# Troubleshooting V1

List of common bugs and issues that you may encounter while using eigenda-proxy.

## "batch hash mismatch" error

> Note: Nitro L2s that settle to an Ethereum network do not suffer from this error because the sequencer inbox is an actual contract which verifies the cert against EigenDA service manager state, as opposed to op's batcher inbox which is just an EOA. Arbitrum Layer 3s that use EigenDA are also susceptible to the mismatch error since certificates can't be verified on the Layer 2 sequencer inbox.
> 
```bash
t=2025-02-10T07:25:24+0000 lvl=info msg=request role=eigenda_proxy method=GET url=/get/0x010000f901d8f852f842a02da35fdbe4fed916eb80fafb6e12711a35327e1e7055b78fc9085812e9d6e28fa021c98859127f075fc90785a730ce1b3894586badee45a698a5cce9c2d6d36269820200cac480213701c401213708f901818302270581c3f873eba0eea9df513a56cbdc8a8c045ddaa61a2a7f86097ad0b99b5b08baa7119065effe82000182636283328734a0c7c898eb39076fbc6b367c51c1908a68bbcac3d81f739057aec168aa3a78e1b20083328792a0c3169896044de8cee3270fdd643157665353a6963149cb158b63137ac405251ab90100c3904bddb63b2be1f2e779b9ecf289fe40693eb943cb260a7d5c95ac6021f35883ade0e4d82e071b86158510fc83cd2c2d61f4a72e6bc6d5c533bacd44eb6a6b97f742d7518f6475f65fb0656fe58d8ad2cc8961c770680c9538593c5a8bc3861fc11e38477076f06911e6de93d9b641b535115700bcdcc863b58b021a7a177b3343e80ab645867061216e191740e4558754e7d13fc3f235743f8d865c0f976a91a3c7912d03e56bb955ddec854c2d661edfc43d5c366040296eac657960a8b40ec6e01b17624d18f32e930b8712cc1a00a2f696229ce10b5b03b3c7c992b601c8ef46390b2fd016687cd1ed4cd70310018dfa97418c7819a0cda7322a7eefa4820001 status=500 duration=585.735929ms err="Error: get request failed with commitment f901d8f852f842a02da35fdbe4fed916eb80fafb6e12711a35327e1e7055b78fc9085812e9d6e28fa021c98859127f075fc90785a730ce1b3894586badee45a698a5cce9c2d6d36269820200cac480213701c401213708f901818302270581c3f873eba0eea9df513a56cbdc8a8c045ddaa61a2a7f86097ad0b99b5b08baa7119065effe82000182636283328734a0c7c898eb39076fbc6b367c51c1908a68bbcac3d81f739057aec168aa3a78e1b20083328792a0c3169896044de8cee3270fdd643157665353a6963149cb158b63137ac405251ab90100c3904bddb63b2be1f2e779b9ecf289fe40693eb943cb260a7d5c95ac6021f35883ade0e4d82e071b86158510fc83cd2c2d61f4a72e6bc6d5c533bacd44eb6a6b97f742d7518f6475f65fb0656fe58d8ad2cc8961c770680c9538593c5a8bc3861fc11e38477076f06911e6de93d9b641b535115700bcdcc863b58b021a7a177b3343e80ab645867061216e191740e4558754e7d13fc3f235743f8d865c0f976a91a3c7912d03e56bb955ddec854c2d661edfc43d5c366040296eac657960a8b40ec6e01b17624d18f32e930b8712cc1a00a2f696229ce10b5b03b3c7c992b601c8ef46390b2fd016687cd1ed4cd70310018dfa97418c7819a0cda7322a7eefa4820001: failed to verify batch: batch hash mismatch, onchain: 3f5b8c23001297ae107370b73956650a20d6cdb2097bddde70de590ef35ef5a9, computed: e7c00bb840cbfaa80b9f505d8df5e5fdbdfe87852908de0628d811b415616702 (Mode: optimism_generic, CertVersion: 0)" commitment_mode=optimism_generic cert_version=0
```

This error returned [here](https://github.com/Layr-Labs/eigenda/blob/86e27fa0342f4638a356ba9738cf998374889ee3/api/proxy/store/generated_key/eigenda/verify/cert.go#L163) happens when verifying a v1 cert, where the batch hash in the cert does not match the hash that was computed when the cert was bridged onchain by the EigenDA disperser.

This typically results from an L1 reorg while running the eigenda-proxy with the `--eigenda.confirmation-depth` flag set to 0. In this setting, the disperser bridges the cert onchain via the confirmBatch function, then the proxy (while polling the GetBlobStatus) endpoint receives the cert, and without waiting for the cert to have been onchain for a few blocks (to prevent reorgs) immediately returns it to the batcher, which then sends it onchain. At this point the L1 reorgs, forcing the disperser to resent a new confirmBatch transaction, which lands in a new block. The cert for the blob then gets updated with a different [confirmation_block_number](https://github.com/Layr-Labs/eigenda/blob/f305e046ae3e611e19c15e134571cc2ec83062b4/api/proto/disperser/disperser.proto#L242). This makes the cert in the batcher inbox outdated and no longer valid. Hence, the rollup nodes, after reading the cert from the batcher-inbox, submits it to the proxy via a GET request, which hashes the batchHeader contained in the cert and compares it to the batch hash in the onchain. Because their `confirmation_block_number` differ, the hash is naturally different, hence the observed `batch hash mismatch` error.

To troubleshoot this error, you will first need the [request_id](https://github.com/Layr-Labs/eigenda/blob/f305e046ae3e611e19c15e134571cc2ec83062b4/api/proto/disperser/disperser.proto#L112) of the blob. This can be found in the proxy's logs, and is also used as the URL slug for our blob explorer: in https://blobs-holesky.eigenda.xyz/blobs/05cb90531098964cc73f0e6b782f52f3b4387e3aed3ef6b89cd9a3c118dd781e-313733393032373134323137353330363738362f302f33332f312f33332fe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855, the request_id is `05cb90531098964cc73f0e6b782f52f3b4387e3aed3ef6b89cd9a3c118dd781e-313733393032373134323137353330363738362f302f33332f312f33332fe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`. From this request_id, you can query for the cert from the disperser:
```bash
REQUEST_ID=05cb90531098964cc73f0e6b782f52f3b4387e3aed3ef6b89cd9a3c118dd781e-313733393032373134323137353330363738362f302f33332f312f33332fe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
REQUEST_ID_BASE64=$(echo -n $REQUEST_ID | base64)
grpcurl -d @ disperser-holesky.eigenda.xyz:443 disperser.Disperser.GetBlobStatus <<EOM
{
  "request_id": "$REQUEST_ID_BASE64"
}
EOM
```
And then you can compare the confirmation_block_number returned from the disperser to the one in the cert that was posted to the rollup's batcher inbox:
```bash
OP_ALTDA_COMMITMENT=010000f901d8f852f842a02da35fdbe4fed916eb80fafb6e12711a35327e1e7055b78fc9085812e9d6e28fa021c98859127f075fc90785a730ce1b3894586badee45a698a5cce9c2d6d36269820200cac480213701c401213708f901818302270581c3f873eba0eea9df513a56cbdc8a8c045ddaa61a2a7f86097ad0b99b5b08baa7119065effe82000182636283328734a0c7c898eb39076fbc6b367c51c1908a68bbcac3d81f739057aec168aa3a78e1b20083328792a0c3169896044de8cee3270fdd643157665353a6963149cb158b63137ac405251ab90100c3904bddb63b2be1f2e779b9ecf289fe40693eb943cb260a7d5c95ac6021f35883ade0e4d82e071b86158510fc83cd2c2d61f4a72e6bc6d5c533bacd44eb6a6b97f742d7518f6475f65fb0656fe58d8ad2cc8961c770680c9538593c5a8bc3861fc11e38477076f06911e6de93d9b641b535115700bcdcc863b58b021a7a177b3343e80ab645867061216e191740e4558754e7d13fc3f235743f8d865c0f976a91a3c7912d03e56bb955ddec854c2d661edfc43d5c366040296eac657960a8b40ec6e01b17624d18f32e930b8712cc1a00a2f696229ce10b5b03b3c7c992b601c8ef46390b2fd016687cd1ed4cd70310018dfa97418c7819a0cda7322a7eefa4820001
RLP_ENCODED_CERT=${OP_ALTDA_COMMITMENT:6}
cast --from-rlp $RLP_ENCODED_CERT | jq
```
Unfortunately `cast --from-rlp` doesn't currently support schemas, so you will have to compare the arrays-of-hex-strings output to the [BlobInfo](https://github.com/Layr-Labs/eigenda/blob/f305e046ae3e611e19c15e134571cc2ec83062b4/api/proto/disperser/disperser.proto#L178) schema (note that [CertV1 === BlobInfo](https://github.com/Layr-Labs/eigenda/blob/86e27fa0342f4638a356ba9738cf998374889ee3/api/proxy/store/generated_key/eigenda/verify/certificate.go#L31)).

And here's a full script to compare both (replace `REQUEST_ID`, `OP_ALTDA_COMMITMENT`, and `DISPERSER_ENDPOINT` with your values):
```bash
DISPERSER_ENDPOINT=disperser-holesky.eigenda.xyz:443
# Get the request_id from your batcher's proxy logs
REQUEST_ID=05cb90531098964cc73f0e6b782f52f3b4387e3aed3ef6b89cd9a3c118dd781e-313733393032373134323137353330363738362f302f33332f312f33332fe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
# Get the cert from the rollup's batcher inbox
OP_ALTDA_COMMITMENT=010000f901d8f852f842a02da35fdbe4fed916eb80fafb6e12711a35327e1e7055b78fc9085812e9d6e28fa021c98859127f075fc90785a730ce1b3894586badee45a698a5cce9c2d6d36269820200cac480213701c401213708f901818302270581c3f873eba0eea9df513a56cbdc8a8c045ddaa61a2a7f86097ad0b99b5b08baa7119065effe82000182636283328734a0c7c898eb39076fbc6b367c51c1908a68bbcac3d81f739057aec168aa3a78e1b20083328792a0c3169896044de8cee3270fdd643157665353a6963149cb158b63137ac405251ab90100c3904bddb63b2be1f2e779b9ecf289fe40693eb943cb260a7d5c95ac6021f35883ade0e4d82e071b86158510fc83cd2c2d61f4a72e6bc6d5c533bacd44eb6a6b97f742d7518f6475f65fb0656fe58d8ad2cc8961c770680c9538593c5a8bc3861fc11e38477076f06911e6de93d9b641b535115700bcdcc863b58b021a7a177b3343e80ab645867061216e191740e4558754e7d13fc3f235743f8d865c0f976a91a3c7912d03e56bb955ddec854c2d661edfc43d5c366040296eac657960a8b40ec6e01b17624d18f32e930b8712cc1a00a2f696229ce10b5b03b3c7c992b601c8ef46390b2fd016687cd1ed4cd70310018dfa97418c7819a0cda7322a7eefa4820001

REQUEST_ID_BASE64=$(echo -n $REQUEST_ID | base64)
DISPERSER_CERT_CONFIRMATION_BLOCK_NUM=$((
    grpcurl -d @ $DISPERSER_ENDPOINT disperser.Disperser.GetBlobStatus <<EOM
{
  "request_id": "$REQUEST_ID_BASE64"
}
EOM
) | jq .info.blobVerificationProof.batchMetadata.confirmationBlockNumber)

RLP_ENCODED_CERT=${OP_ALTDA_COMMITMENT:6}
BATCHER_INBOX_CERT_CONFIRMATION_BLOCK_NUM=$(cast --from-rlp $RLP_ENCODED_CERT | jq -r '.[1][2][3]' | cast --to-dec)

echo "Disperser cert confirmation block number: $DISPERSER_CERT_CONFIRMATION_BLOCK_NUM"
echo "Batcher inbox cert confirmation block number: $BATCHER_INBOX_CERT_CONFIRMATION_BLOCK_NUM"
```