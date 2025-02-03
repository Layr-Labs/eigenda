# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common/v2/common_v2.proto](#common_v2_common_v2-proto)
    - [Batch](#common-v2-Batch)
    - [BatchHeader](#common-v2-BatchHeader)
    - [BlobCertificate](#common-v2-BlobCertificate)
    - [BlobHeader](#common-v2-BlobHeader)
    - [PaymentHeader](#common-v2-PaymentHeader)
  
- [Scalar Value Types](#scalar-value-types)



<a name="common_v2_common_v2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common/v2/common_v2.proto



<a name="common-v2-Batch"></a>

### Batch
Batch is a batch of blob certificates


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [BatchHeader](#common-v2-BatchHeader) |  | header contains metadata about the batch |
| blob_certificates | [BlobCertificate](#common-v2-BlobCertificate) | repeated | blob_certificates is the list of blob certificates in the batch |






<a name="common-v2-BatchHeader"></a>

### BatchHeader
BatchHeader is the header of a batch of blobs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_root | [bytes](#bytes) |  | batch_root is the root of the merkle tree of the hashes of blob certificates in the batch |
| reference_block_number | [uint64](#uint64) |  | reference_block_number is the block number that the state of the batch is based on for attestation |






<a name="common-v2-BlobCertificate"></a>

### BlobCertificate
BlobCertificate contains a full description of a blob and how it is dispersed. Part of the certificate
is provided by the blob submitter (i.e. the blob header), and part is provided by the disperser (i.e. the relays).
Validator nodes eventually sign the blob certificate once they are in custody of the required chunks
(note that the signature is indirect; validators sign the hash of a Batch, which contains the blob certificate).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| blob_header | [BlobHeader](#common-v2-BlobHeader) |  | blob_header contains data about the blob. |
| signature | [bytes](#bytes) |  | signature is an ECDSA signature signed by the blob request signer&#39;s account ID over the BlobHeader&#39;s blobKey, which is a keccak hash of the serialized BlobHeader, and used to verify against blob dispersal request&#39;s account ID |
| relay_keys | [uint32](#uint32) | repeated | relay_keys is the list of relay keys that are in custody of the blob. The relays custodying the data are chosen by the Disperser to which the DisperseBlob request was submitted. It needs to contain at least 1 relay number. To retrieve a blob from the relay, one can find that relay&#39;s URL in the EigenDARelayRegistry contract: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDARelayRegistry.sol |






<a name="common-v2-BlobHeader"></a>

### BlobHeader
BlobHeader contains the information describing a blob and the way it is to be dispersed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [uint32](#uint32) |  | The blob version. Blob versions are pushed onchain by EigenDA governance in an append only fashion and store the maximum number of operators, number of chunks, and coding rate for a blob. On blob verification, these values are checked against supplied or default security thresholds to validate the security assumptions of the blob&#39;s availability. |
| quorum_numbers | [uint32](#uint32) | repeated | quorum_numbers is the list of quorum numbers that the blob is part of. All quorums must be specified (including required quorums).

The following quorums are currently required: - 0: ETH - 1: EIGEN |
| commitment | [common.BlobCommitment](#common-BlobCommitment) |  | commitment is the KZG commitment to the blob |
| payment_header | [PaymentHeader](#common-v2-PaymentHeader) |  | payment_header contains payment information for the blob |
| salt | [uint32](#uint32) |  | salt is used to ensure that the dispersal request is intentionally unique. This is currently only useful for reserved payments when the same blob is submitted multiple times within the same reservation period. On-demand payments already have unique cumulative_payment values for intentionally unique dispersal requests. |






<a name="common-v2-PaymentHeader"></a>

### PaymentHeader
PaymentHeader contains payment information for a blob.
At least one of reservation_period or cumulative_payment must be set, and reservation_period 
is always considered before cumulative_payment. If reservation_period is set but not valid, 
the server will reject the request and not proceed with dispersal. If reservation_period is not set 
and cumulative_payment is set but not valid, the server will reject the request and not proceed with dispersal.
Once the server has accepted the payment header, a client cannot cancel or rollback the payment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [string](#string) |  | The account ID of the disperser client. This account ID is an eth wallet address of the user, corresponding to the key used by the client to sign the BlobHeader. |
| reservation_period | [uint32](#uint32) |  | The reservation period of the dispersal request is used for rate-limiting the user&#39;s account against their dedicated bandwidth. This method requires users to set up reservation accounts with EigenDA team, and the team will set up an on-chain record of reserved bandwidth for the user for some period of time. The dispersal client&#39;s accountant will set this value to the current timestamp divided by the on-chain configured reservation period interval, mapping each request to a time-based window and is serialized and parsed as a uint32. The disperser server then validates that it matches either the current or the previous period.

Example Usage Flow: 1. The user sets up a reservation with the EigenDA team, including throughput (symbolsPerSecond), startTimestamp, endTimestamp, and reservationPeriodInterval. 2. When sending a dispersal request at time t, the client computes reservation_period = floor(t / reservationPeriodInterval). 3. The request includes this reservation_period index. The disperser checks: - If the reservation is active (t &gt;= startTimestamp and t &lt; endTimestamp). - Whether the user still has bandwidth capacity (hasn’t exceeded symbolsPerSecond * reservationPeriodInterval). 4. Server always go ahead with recording the received request, and then categorize the scenarios - If the remaining bandwidth is sufficient for the request, the dispersal request proceeds. - If the remaining bandwidth is not enough for the request, server fills up the current bin and overflowing the extra to a future bin. - If the bandwidth has already been exhausted, the request is rejected. 5. Once the dispersal request signature has been verified, the server will not roll back the payment or the usage records. Users should be aware of this when planning their usage. The dispersal client written by EigenDA team takes account of this. 6. When the reservation ends or usage is exhausted, the client must wait for the next reservation period or switch to on-demand. |
| cumulative_payment | [bytes](#bytes) |  | Cumulative payment is the total amount of tokens paid by the requesting account, including the current request. This value is serialized as an uint256 and parsed as a big integer, and must match the user’s on-chain deposit limits as well as the recorded payments for all previous requests. Because it is a cumulative (not incremental) total, requests can arrive out of order and still unambiguously declare how much of the on-chain deposit can be deducted.

Example Decision Flow: 1. In the set up phase, the user must deposit tokens into the EigenDA PaymentVault contract. The payment vault contract specifies the minimum number of symbols charged per dispersal, the pricing per symbol, and the maximum global rate for on-demand dispersals. The user should calculate the amount of tokens they would like to deposit based on their usage. The first time a user make a request, server will immediate read the contract for the on-chain balance. When user runs out of on-chain balance, the server will reject the request and not proceed with dispersal. When a user top up on-chain, the server will only refresh every few minutes for the top-up to take effect. 2. The disperser client accounts how many tokens they’ve already paid (previousCumPmt). 3. They calculate the incremental amount of tokens needed for the current request needs based on protocol defined pricing. 4. They take the sum of previousCumPmt &#43; new incremental payment and place it in the “cumulative_payment” field. 5. The disperser checks this new cumulative total against on-chain deposits and prior records (largest previous payment and smallest later payment if exists). 6. If the payment number is valid, the request is confirmed and disperser proceeds with dispersal; otherwise it’s rejected. |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

